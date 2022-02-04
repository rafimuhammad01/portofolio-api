package utils

import "os"

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func GetAccessTokenDuration() string {
	return os.Getenv("JWT_ACCESS_TOKEN_DURATION")
}

func GetRefreshTokenDuration() string {
	return os.Getenv("JWT_REFRESH_TOKEN_DURATION")
}
