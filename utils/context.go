package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rafimuhammad01/portofolio-api/internal/jwt"
)

func GetPayloadFromContext(c *gin.Context) (*jwt.Payload, error) {
	payloadContext, _ := c.Get(AuthorizationPayloadKey)

	payload, ok := payloadContext.(*jwt.Payload)
	if !ok {
		return nil, errors.New("failed to assert context to payload")
	}

	return payload, nil
}
