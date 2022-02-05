package user

import "github.com/rafimuhammad01/portofolio-api/internal/jwt"

// User entity represent users table in database
type User struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	FullName string `json:"full_name" db:"full_name"`
}

type ListUser struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

// ListUserAPIResponse API response for List
type ListUserAPIResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    *ListUser `json:"data,omitempty"`
	Errors  []string  `json:"errors,omitempty"`
}

// CreateUserAPIRequest create user request body from client
type CreateUserAPIRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name" db:"full_name"`
	Password string `json:"password"`
}

type CreateUserAPIResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    *User    `json:"data,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

type GetUserByIDAPIResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    *User    `json:"data,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// LoginAPIRequest login request body
type LoginAPIRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type IDAndPassword struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
}

type LoginAPIResponse struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    *jwt.JWTAPIResponse `json:"data,omitempty"`
	Errors  []string            `json:"errors,omitempty"`
}

type RefreshTokenAPIRequest struct {
	RefreshToken string `json:"refresh_token"`
}
