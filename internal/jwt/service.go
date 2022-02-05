package jwt

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

// Service is an interface for managing tokens
type Service interface {
	CreateToken(username string, userID int, duration time.Duration) (string, time.Time, error)
	VerifyToken(token string) (*Payload, error)
	CreateRefreshToken(accessToken string, duration time.Duration, ctx context.Context) (string, error)
	GetUserInformation(refreshToken string, ctx context.Context) (userID int, username string, err error)
	LogoutUser(userID string) (string, error)
}

type service struct {
	secretKey string
	repo      Repo
}

// NewService creates a new JWTMaker
func NewService(secretKey string, repo Repo) Service {
	return &service{
		secretKey: secretKey,
		repo:      repo,
	}
}

// VerifyToken checks if the token is valid or not
func (s *service) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.Wrap(ErrInvalidToken, ErrInvalidToken.Error())
		}
		return []byte(s.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, errors.Wrap(ErrExpiredToken, ErrExpiredToken.Error())
		}
		return nil, errors.Wrap(ErrInvalidToken, ErrInvalidToken.Error())
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errors.Wrap(ErrInvalidToken, ErrInvalidToken.Error())
	}

	return payload, nil
}

func (s *service) CreateToken(username string, userID int, duration time.Duration) (string, time.Time, error) {
	payload, err := s.repo.NewPayload(username, userID, duration)
	if err != nil {
		return "", time.Time{}, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", time.Time{}, errors.Wrap(ErrIntervalServer, err.Error())
	}
	return token, payload.ExpiredAt, nil
}

func (s *service) CreateRefreshToken(accessToken string, duration time.Duration, ctx context.Context) (string, error) {
	payload, err := s.VerifyToken(accessToken)
	if err != nil {
		return "", err
	}

	refreshToken, err := uuid.NewUUID()
	if err != nil {
		return "", errors.Wrap(ErrIntervalServer, err.Error())
	}

	err = s.repo.StoreRefreshToken(refreshToken.String(), payload.UserID, payload.Username, duration, ctx)
	if err != nil {
		return "", err
	}

	return refreshToken.String(), nil
}

func (s *service) GetUserInformation(refreshToken string, ctx context.Context) (userID int, username string, err error) {
	userID, username, err = s.repo.GetRefreshToken(refreshToken, ctx)
	if err != nil {
		return 0, "", err
	}

	return userID, username, nil
}

func (s *service) LogoutUser(userID string) (string, error) {
	panic("implement me")
}
