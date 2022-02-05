package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rafimuhammad01/portofolio-api/internal/jwt"
)

var (
	ErrAssertion = errors.New("failed to assert context to payload")
)

func GetPayloadFromContext(c *gin.Context) (*jwt.Payload, error) {
	payloadContext, _ := c.Get(AuthorizationPayloadKey)

	payload, ok := payloadContext.(*jwt.Payload)
	if !ok {
		return nil, errors.Wrap(ErrAssertion, "failed to assert context to payload")
	}

	return payload, nil
}
