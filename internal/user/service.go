package user

import (
	"github.com/rafimuhammad01/portofolio-api/internal/jwt"
	"github.com/rafimuhammad01/portofolio-api/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func NewService(repo Repo, jwtService jwt.Service) Service {
	return &service{
		repo:       repo,
		jwtService: jwtService,
	}
}

type Service interface {
	List() (*ListUserAPIResponse, error)
	Create(CreateUserAPIRequest) (*CreateUserAPIResponse, error)
	Get(int) (*GetUserByIDAPIResponse, error)
	Login(LoginAPIRequest) (*LoginAPIResponse, error)
	RefreshToken(accessToken string) (*LoginAPIResponse, error)
}

type service struct {
	repo       Repo
	jwtService jwt.Service
}

func (s service) Login(request LoginAPIRequest) (*LoginAPIResponse, error) {
	user, status, message, err := s.repo.GetUserIDAndPasswordByUsername(request.Username)
	if err != nil {
		return &LoginAPIResponse{
			Status:  status,
			Message: "internal server error",
		}, err
	}

	if user == nil {
		return &LoginAPIResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
			Errors: []string{
				"wrong username/password",
			},
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return &LoginAPIResponse{
			Status:  http.StatusUnauthorized,
			Message: "unauthorized",
			Errors: []string{
				"wrong username/password",
			},
		}, nil
	}

	duration, _ := time.ParseDuration(utils.GetAccessTokenDuration())
	accessToken, expAt, err := s.jwtService.CreateToken(
		request.Username,
		user.ID,
		duration,
	)
	if err != nil {
		return &LoginAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}, err
	}

	duration, _ = time.ParseDuration(utils.GetRefreshTokenDuration())
	refreshToken, err := s.jwtService.CreateRefreshToken(accessToken, duration)
	if err != nil {
		return &LoginAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}, err
	}

	return &LoginAPIResponse{
		Status:  http.StatusOK,
		Message: message,
		Data: &jwt.JWTAPIResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiredAt:    expAt,
		},
	}, nil
}

func (s service) Get(ID int) (*GetUserByIDAPIResponse, error) {
	user, status, message, err := s.repo.GetByID(ID)
	if err != nil {
		return &GetUserByIDAPIResponse{
			Status:  status,
			Message: "internal server error",
		}, err
	}

	if user == nil {
		return &GetUserByIDAPIResponse{
			Status:  http.StatusNotFound,
			Message: "not found",
			Errors:  []string{message},
		}, nil
	}

	return &GetUserByIDAPIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    user,
	}, nil
}

func (s service) List() (*ListUserAPIResponse, error) {
	users, status, message, err := s.repo.List()
	if err != nil {
		if status == http.StatusNotFound {
			return &ListUserAPIResponse{
				Message: "not found",
				Status:  status,
				Errors:  []string{message},
			}, nil
		} else {
			return &ListUserAPIResponse{
				Message: "internal server error",
				Status:  status,
			}, err
		}
	}

	return &ListUserAPIResponse{
		Status:  status,
		Message: "success get list users",
		Data:    users,
	}, err
}

func (s service) Create(request CreateUserAPIRequest) (*CreateUserAPIResponse, error) {
	// Check Username Uniqueness
	userByUsername, status, message, err := s.repo.GetByUsername(request.Username)
	if err != nil {
		return &CreateUserAPIResponse{
			Message: "internal server error",
			Status:  status,
		}, err
	}

	// Bad request if username already exist
	if userByUsername != nil && status != http.StatusNotFound && err == nil {
		return &CreateUserAPIResponse{
			Message: "bad request",
			Status:  http.StatusBadRequest,
			Errors:  []string{"Username already exist"},
		}, err
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return &CreateUserAPIResponse{
			Message: "internal server error",
			Status:  http.StatusInternalServerError,
		}, err
	}

	request.Password = string(hashedPassword)

	// Insert data to DB
	user, status, message, err := s.repo.Create(request)
	if err != nil {
		return &CreateUserAPIResponse{
			Message: "internal server error",
			Status:  status,
		}, err
	}

	return &CreateUserAPIResponse{
		Status:  status,
		Message: message,
		Data:    user,
	}, err
}

func (s service) RefreshToken(accessToken string) (*LoginAPIResponse, error) {
	payload, refreshToken, err := s.jwtService.GetRefreshToken(accessToken)
	if err != nil {
		if err == jwt.ErrInvalidToken {
			return &LoginAPIResponse{
				Message: "bad request",
				Status:  http.StatusBadRequest,
				Errors:  []string{err.Error()},
			}, nil
		}
		return &LoginAPIResponse{
			Message: "internal server error",
			Status:  http.StatusInternalServerError,
		}, err
	}

	duration, _ := time.ParseDuration(utils.GetAccessTokenDuration())
	accessToken, expAt, err := s.jwtService.CreateToken(
		payload.Username,
		payload.UserID,
		duration,
	)
	if err != nil {
		return &LoginAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		}, err
	}

	return &LoginAPIResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data: &jwt.JWTAPIResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiredAt:    expAt,
		},
	}, nil
}
