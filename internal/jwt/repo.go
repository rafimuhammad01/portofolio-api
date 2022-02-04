package jwt

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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
	GetRefreshToken(ctx context.Context, userID string) (string, error)
	StoreRefreshToken(ctx context.Context, userID string, refreshToken string, expiresIn time.Duration) error
	DestroyRefreshToken(ctx context.Context, userID string) error
}

type repo struct {
	rdb *redis.Client
}

func (r repo) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	val, err := r.rdb.Get(ctx, userID).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r repo) StoreRefreshToken(ctx context.Context, userID string, refreshToken string, expiresIn time.Duration) error {
	err := r.rdb.Set(ctx, userID, refreshToken, expiresIn).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r repo) DestroyRefreshToken(ctx context.Context, userID string) error {
	err := r.rdb.Del(ctx, userID).Err()
	if err != nil {
		return err
	}
	return nil
}

// NewPayload creates a new token payload with a specific username and duration
func (r repo) NewPayload(username string, userID int, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
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
