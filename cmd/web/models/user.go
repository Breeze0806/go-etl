package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"` // Never expose in JSON
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserResponse is the API response for user data (without password)
type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidateRole checks if the role is valid
func ValidateRole(role string) string {
	if role == "" {
		return "user" // Default role
	}
	// Only allow specific roles
	switch role {
	case "admin", "user", "readonly":
		return role
	default:
		return "user"
	}
}

// Scan scans a row from the database into the User struct
func (u *User) Scan(row *sql.Row) error {
	return row.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
}

// ScanRows scans multiple rows from the database into the User struct
func (u *User) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
}

// ErrUserNotFound is returned when a user is not found
var ErrUserNotFound = errors.New("user not found")

// ErrUserExists is returned when a user already exists
var ErrUserExists = errors.New("user already exists")

// ErrInvalidCredentials is returned when login credentials are invalid
var ErrInvalidCredentials = errors.New("invalid credentials")

// CreateUserRequest represents the request to create a user (admin only)
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role"`
}

// UpdateUserRequest represents the request to update a user (self or admin)
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Role     string `json:"role"`
}

// UserListResponse represents the response for listing users
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}
