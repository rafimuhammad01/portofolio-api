package jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

// NewRepo for Payload
func NewRepo(rdb *redis.Client) Repo {
	return &repo{
		rdb: rdb,
	}
}

type Repo interface {
	NewPayload(username string, userID int, duration time.Duration) (*Payload, error)
	StoreRefreshToken(refreshToken string, userID int, username string, duration time.Duration, ctx context.Context) error
	GetRefreshToken(refreshToken string, ctx context.Context) (userID int, username string, err error)
}

type repo struct {
	rdb *redis.Client
}

// NewPayload creates a new token payload with a specific username and duration
func (r repo) NewPayload(username string, userID int, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(ErrIntervalServer, err.Error())
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (r repo) StoreRefreshToken(refreshToken string, userID int, username string, duration time.Duration, ctx context.Context) error {
	val := map[string]interface{}{
		"user_id":  userID,
		"username": username,
	}

	content, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(ErrIntervalServer, err.Error())
	}

	err = r.rdb.Set(ctx, refreshToken, content, duration).Err()
	if err != nil {
		return errors.Wrap(ErrIntervalServer, err.Error())
	}

	return nil
}

func (r repo) GetRefreshToken(refreshToken string, ctx context.Context) (userID int, username string, err error) {
	var val map[string]interface{}

	result, err := r.rdb.Get(ctx, refreshToken).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, "", errors.Wrap(ErrInvalidToken, err.Error())
		}
		return 0, "", errors.Wrap(ErrIntervalServer, err.Error())
	}

	err = json.Unmarshal([]byte(result), &val)
	if err != nil {
		return 0, "", errors.Wrap(ErrIntervalServer, err.Error())
	}

	ID, ok := val["user_id"].(float64)
	if !ok {
		return 0, "", errors.Wrap(ErrIntervalServer, fmt.Sprintf("[failed to assert interface to int] should be %T", val["user_id"]))
	}

	username = fmt.Sprint(val["username"])

	return int(ID), username, nil

}
