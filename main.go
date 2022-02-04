package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rafimuhammad01/portofolio-api/api"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	// Load env var
	err := godotenv.Load()
	if err != nil{
		logrus.Fatal(".env not found, will use default env")
	}

	// Creating router
	router := gin.Default()
	s := api.NewServer(router)
	s.Init()

	// Running server
	s.RunServer(os.Getenv("PORT"))
}