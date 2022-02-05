package user

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rafimuhammad01/portofolio-api/internal/jwt"
	"github.com/rafimuhammad01/portofolio-api/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func NewService(repo Repo, jwtService jwt.Service) Service {
	return &service{
		repo:       repo,
		jwtService: jwtService,
	}
}

type Service interface {
	List() (*ListUser, error)
	Create(username, fullName, password string) (*User, error)
	Get(ID int) (*User, error)
	Login(username, password string, ctx context.Context) (string, string, time.Time, error)
	RefreshToken(refreshToken string, ctx context.Context) (string, string, time.Time, error)
}

type service struct {
	repo       Repo
	jwtService jwt.Service
}

func (s service) List() (*ListUser, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s service) Create(username, fullName, password string) (*User, error) {
	// Check Username Uniqueness
	userByUsername, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	// Bad request if username already exist
	if userByUsername != nil {
		return nil, errors.Wrap(ErrUsernameAlreadyExist, "username already taken")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(ErrInternalServer, err.Error())
	}

	password = string(hashedPassword)

	// Insert data to DB
	user, err := s.repo.Create(username, fullName, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s service) Get(ID int) (*User, error) {
	user, err := s.repo.GetByID(ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s service) Login(username, password string, ctx context.Context) (string, string, time.Time, error) {
	user, err := s.repo.GetUserIDAndPasswordByUsername(username)
	if err != nil {
		return "", "", time.Time{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return "", "", time.Time{}, errors.Wrap(ErrInvalidUsernameOrPassword, err.Error())
		}
		return "", "", time.Time{}, err
	}

	duration, _ := time.ParseDuration(utils.GetAccessTokenDuration())
	accessToken, expAt, err := s.jwtService.CreateToken(
		username,
		user.ID,
		duration,
	)
	if err != nil {
		return "", "", time.Time{}, err
	}

	duration, _ = time.ParseDuration(utils.GetRefreshTokenDuration())
	refreshToken, err := s.jwtService.CreateRefreshToken(accessToken, duration, ctx)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken, refreshToken, expAt, nil
}

func (s service) RefreshToken(refreshToken string, ctx context.Context) (string, string, time.Time, error) {
	userID, username, err := s.jwtService.GetUserInformation(refreshToken, ctx)
	if err != nil {
		return "", "", time.Time{}, err
	}

	duration, _ := time.ParseDuration(utils.GetAccessTokenDuration())
	accessToken, expAt, err := s.jwtService.CreateToken(
		username,
		userID,
		duration,
	)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken, refreshToken, expAt, nil
}
