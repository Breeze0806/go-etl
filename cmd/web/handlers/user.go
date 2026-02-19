package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/Breeze0806/go-etl/cmd/web/middleware"
	"github.com/Breeze0806/go-etl/cmd/web/models"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	db  *sql.DB
	jwt *middleware.JWTMiddleware
}

func NewUserHandler(db *sql.DB, jwt *middleware.JWTMiddleware) *UserHandler {
	return &UserHandler{
		db:  db,
		jwt: jwt,
	}
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	username := c.Query("username")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var query string
	var args []interface{}

	if username != "" {
		query = `
			SELECT id, username, role, created_at, updated_at 
			FROM users 
			WHERE username LIKE ? 
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		args = []interface{}{"%" + username + "%", limit, offset}
	} else {
		query = `
			SELECT id, username, role, created_at, updated_at 
			FROM users 
			ORDER BY id DESC 
			LIMIT ? OFFSET ?
		`
		args = []interface{}{limit, offset}
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query users"})
		return
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to scan user"})
			return
		}
		users = append(users, user.ToResponse())
	}

	if users == nil {
		users = []models.UserResponse{}
	}

	var total int
	if username != "" {
		err = h.db.QueryRow("SELECT COUNT(*) FROM users WHERE username LIKE ?", "%"+username+"%").Scan(&total)
	} else {
		err = h.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count users"})
		return
	}

	c.JSON(http.StatusOK, models.UserListResponse{
		Users: users,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func (h *UserHandler) Create(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := models.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	role := models.ValidateRole(req.Role)
	now := time.Now()

	result, err := h.db.Exec(
		"INSERT INTO users (username, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		req.Username, hashedPassword, role, now, now,
	)
	if err != nil {
		if isDuplicateError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	id, _ := result.LastInsertId()
	user := &models.User{
		ID:        id,
		Username:  req.Username,
		Password:  hashedPassword,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	c.JSON(http.StatusCreated, user.ToResponse())
}

func (h *UserHandler) Update(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	currentUserID, _ := c.Get("user_id")
	currentUserRole, _ := c.Get("role")

	currentUserIDInt := currentUserID.(int64)
	roleStr := currentUserRole.(string)

	if currentUserIDInt != userID && roleStr != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	err = h.db.QueryRow(
		"SELECT id, username, password, role, created_at, updated_at FROM users WHERE id = ?",
		userID,
	).Scan(&existingUser.ID, &existingUser.Username, &existingUser.Password, &existingUser.Role, &existingUser.CreatedAt, &existingUser.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	username := existingUser.Username
	role := existingUser.Role

	if req.Username != "" {
		username = req.Username
	}
	if req.Role != "" {
		role = models.ValidateRole(req.Role)
	}

	if roleStr != "admin" && role != existingUser.Role {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions to change role"})
		return
	}

	now := time.Now()

	_, err = h.db.Exec(
		"UPDATE users SET username = ?, role = ?, updated_at = ? WHERE id = ?",
		username, role, now, userID,
	)
	if err != nil {
		if isDuplicateError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	var updatedUser models.User
	err = h.db.QueryRow(
		"SELECT id, username, role, created_at, updated_at FROM users WHERE id = ?",
		userID,
	).Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.Role, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser.ToResponse())
}

func (h *UserHandler) Delete(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	currentUserID, _ := c.Get("user_id")
	currentUserIDInt := currentUserID.(int64)

	if currentUserIDInt == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete yourself"})
		return
	}

	result, err := h.db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var user models.User
	err = h.db.QueryRow(
		"SELECT id, username, role, created_at, updated_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}
