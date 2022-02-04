package jwt

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"strconv"
	"time"
)

// Service is an interface for managing tokens
type Service interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(username string, userID int, duration time.Duration) (string, time.Time, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)

	CreateRefreshToken(token string, duration time.Duration) (string, error)

	GetRefreshToken(token string) (*Payload, string, error)

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
			return nil, ErrInvalidToken
		}
		return []byte(s.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
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
		return "", time.Time{}, err
	}
	return token, payload.ExpiredAt, nil
}

func (s *service) CreateRefreshToken(token string, duration time.Duration) (string, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if !ok || !errors.Is(verr.Inner, ErrExpiredToken) {
			return "", ErrInvalidToken
		}
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return "", ErrInvalidToken
	}

	refreshToken, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	err = s.repo.StoreRefreshToken(ctx, strconv.Itoa(payload.UserID), refreshToken.String(), duration)
	if err != nil {
		return "", err
	}

	return refreshToken.String(), nil
}

func (s *service) GetRefreshToken(token string) (*Payload, string, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if !ok || !errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, "", ErrInvalidToken
		}
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, "", ErrInvalidToken
	}

	ctx := context.Background()
	refreshToken, err := s.repo.GetRefreshToken(ctx, strconv.Itoa(payload.UserID))
	if err != nil {
		return nil, "", err
	}

	if refreshToken == "" {
		return nil, "", ErrInvalidToken
	}

	return payload, refreshToken, nil
}

func (s *service) LogoutUser(userID string) (string, error) {
	panic("implement me")
}
