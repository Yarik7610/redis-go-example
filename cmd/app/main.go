package main

import (
	"log"

	"github.com/Yarik7610/redis-learning/internal/config"
	"github.com/Yarik7610/redis-learning/internal/controller"
	"github.com/Yarik7610/redis-learning/internal/repository"
	"github.com/Yarik7610/redis-learning/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()
	defer cfg.Close()

	userRepo := repository.NewUserRepository(cfg.DB)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	r := gin.Default()
	userController.RegisterRoutes(r)

	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
