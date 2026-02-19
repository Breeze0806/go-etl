package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Breeze0806/go-etl/cmd/web/middleware"
	"github.com/Breeze0806/go-etl/cmd/web/models"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	db  *sql.DB
	jwt *middleware.JWTMiddleware
}

func NewAuthHandler(db *sql.DB, jwt *middleware.JWTMiddleware) *AuthHandler {
	return &AuthHandler{
		db:  db,
		jwt: jwt,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
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
		"INSERT INTO users (username, password, email, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		req.Username, hashedPassword, req.Email, role, now, now,
	)
	if err != nil {
		fmt.Printf("Registration error: %v\n", err)
		if isDuplicateError(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user: " + err.Error()})
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

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := h.db.QueryRow(
		"SELECT id, username, password, role, created_at, updated_at FROM users WHERE username = ?",
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if !models.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := h.jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	err := h.db.QueryRow(
		"SELECT id, username, password, role, created_at, updated_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)

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

func isDuplicateError(err error) bool {
	return err.Error() == "UNIQUE constraint failed: users.username"
}
