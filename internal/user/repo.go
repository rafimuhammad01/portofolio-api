package user

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"net/http"
)

// NewRepo PostgreSQL
func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

type Repo interface {
	List() (*ListUser, int, string, error)
	Create(CreateUserAPIRequest) (*User, int, string, error)
	GetByID(int) (*User, int, string, error)
	GetByUsername(string) (*User, int, string, error)
	GetUserIDAndPasswordByUsername(string) (*LoginResponse, int, string, error)
}

type repo struct {
	db *sqlx.DB
}

func (r repo) List() (*ListUser, int, string, error) {
	var users ListUser

	err := r.db.Select(&users.Users, "SELECT id, username, full_name FROM users")
	if err != nil {
		return nil, http.StatusInternalServerError, "internal server error", err
	}

	users.Count = len(users.Users)

	return &users, http.StatusOK, "success", nil
}

func (r repo) Create(request CreateUserAPIRequest) (*User, int, string, error) {
	var user User
	err := r.db.Get(&user, "INSERT INTO users (username, full_name, password) VALUES ($1, $2, $3) RETURNING id, username, full_name", request.Username, request.FullName, request.Password)
	if err != nil {
		return nil, http.StatusInternalServerError, "internal server error", err
	}

	return &user, http.StatusCreated, "success", nil
}

func (r repo) GetByID(ID int) (*User, int, string, error) {
	var user User
	err := r.db.Get(&user, "SELECT id, username, full_name FROM users WHERE id=$1", ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, "user not found", nil
		}
		return nil, http.StatusInternalServerError, "internal server error", err
	}

	return &user, http.StatusOK, "success", nil
}

func (r repo) GetByUsername(username string) (*User, int, string, error) {
	var user User
	err := r.db.Get(&user, "SELECT id, username, full_name FROM users WHERE username=$1", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, "user not found", nil
		}
		return nil, http.StatusInternalServerError, "internal server error", err
	}

	return &user, http.StatusOK, "success", nil
}

func (r repo) GetUserIDAndPasswordByUsername(username string) (*LoginResponse, int, string, error) {
	var user LoginResponse
	err := r.db.Get(&user, "SELECT id, password FROM users WHERE username=$1", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, "user not found", nil
		}
		return nil, http.StatusInternalServerError, "internal server error", err
	}

	return &user, http.StatusOK, "success", nil
}
