package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rafimuhammad01/portofolio-api/internal/jwt"
	"github.com/rafimuhammad01/portofolio-api/utils"
	"net/http"
	"strings"
)

// AuthMiddleware creates a gin middleware for authorization
func AuthMiddleware(tokenMaker *jwt.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(utils.AuthorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.HTTPError{Status: http.StatusUnauthorized, Message: "unauthorized", Errors: []string{err.Error()}})
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.HTTPError{Status: http.StatusUnauthorized, Message: "unauthorized", Errors: []string{err.Error()}})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != utils.AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.HTTPError{Status: http.StatusUnauthorized, Message: "unauthorized", Errors: []string{err.Error()}})
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.HTTPError{Status: http.StatusUnauthorized, Message: "unauthorized", Errors: []string{err.Error()}})
			return
		}

		ctx.Set(utils.AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
