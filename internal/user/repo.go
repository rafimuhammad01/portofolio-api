package user

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// NewRepo PostgreSQL
func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

type Repo interface {
	List() (*ListUser, error)
	Create(username, fullName, password string) (*User, error)
	GetByID(ID int) (*User, error)
	GetByUsername(username string) (*User, error)
	GetUserIDAndPasswordByUsername(username string) (*IDAndPassword, error)
}

type repo struct {
	db *sqlx.DB
}

func (r repo) List() (*ListUser, error) {
	var users ListUser

	err := r.db.Select(&users.Users, "SELECT id, username, full_name FROM users")
	if err != nil {
		if err == sql.ErrNoRows {
			return &users, nil
		}
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	users.Count = len(users.Users)

	return &users, nil
}

func (r repo) Create(username, fullName, password string) (*User, error) {
	var user User
	err := r.db.Get(&user, "INSERT INTO users (username, full_name, password) VALUES ($1, $2, $3) RETURNING id, username, full_name", username, fullName, password)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	return &user, nil
}

func (r repo) GetByID(ID int) (*User, error) {
	var user User
	err := r.db.Get(&user, "SELECT id, username, full_name FROM users WHERE id=$1", ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(ErrUserNotFound, err.Error())
		}
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	return &user, nil
}

func (r repo) GetByUsername(username string) (*User, error) {
	var user User
	err := r.db.Get(&user, "SELECT id, username, full_name FROM users WHERE username=$1", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(ErrUserNotFound, err.Error())
		}
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	return &user, nil
}

func (r repo) GetUserIDAndPasswordByUsername(username string) (*IDAndPassword, error) {
	var user IDAndPassword
	err := r.db.Get(&user, "SELECT id, password FROM users WHERE username=$1", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(ErrUserNotFound, err.Error())
		}
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	return &user, nil
}
