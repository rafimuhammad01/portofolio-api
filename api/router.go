package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rafimuhammad01/portofolio-api/internal/jwt"
	userpkg "github.com/rafimuhammad01/portofolio-api/internal/user"
	"github.com/rafimuhammad01/portofolio-api/middleware"
)

type Routes struct {
	Router      *gin.Engine
	userHandler *userpkg.Handler
	jwtHandler  *jwt.Handler
}

func NewRoutes(router *gin.Engine, userHandler *userpkg.Handler, jwtHandler *jwt.Handler) *Routes {
	return &Routes{
		Router:      router,
		userHandler: userHandler,
		jwtHandler:  jwtHandler,
	}
}

func (r *Routes) Init() {
	v1 := r.Router.Group("/api/v1")

	// User Routing
	user := v1.Group("/user")
	user.POST("/register", r.userHandler.RegisterUser)
	user.GET("/me", middleware.AuthMiddleware(r.jwtHandler), r.userHandler.GetUserByID)
	user.POST("/login", r.userHandler.Login)
	user.POST("/refresh-token", r.userHandler.RefreshToken)
}
