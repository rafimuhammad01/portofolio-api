package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rafimuhammad01/portofolio-api/db/postgre"
	"github.com/rafimuhammad01/portofolio-api/db/redis"
	jwt2 "github.com/rafimuhammad01/portofolio-api/internal/jwt"
	user2 "github.com/rafimuhammad01/portofolio-api/internal/user"
	"os"
)

type Server struct {
	Router *gin.Engine
}

func NewServer(router *gin.Engine) *Server {
	return &Server{
		Router: router,
	}
}

var (
	// Handler
	userHandler *user2.Handler
	jwtHandler  *jwt2.Handler

	// Service
	userService user2.Service
	jwtService  jwt2.Service

	// Repo
	userRepo user2.Repo
	jwtRepo  jwt2.Repo
)

func (s Server) Init() {
	// Init DB
	db := postgre.Init()
	rdb := redis.Init()

	// Init internal package
	// JWT
	jwtRepo = jwt2.NewRepo(rdb)
	jwtService = jwt2.NewService(os.Getenv("JWT_SECRET"), jwtRepo)
	jwtHandler = jwt2.NewHandler(jwtService)

	// User
	userRepo = user2.NewRepo(db)
	userService = user2.NewService(userRepo, jwtService)
	userHandler = user2.NewHandler(userService)

	// Start routing
	r := NewRoutes(s.Router, userHandler, jwtHandler)
	r.Init()
}

func (s Server) RunServer(port string) {
	s.Router.Run(":" + port)
}
