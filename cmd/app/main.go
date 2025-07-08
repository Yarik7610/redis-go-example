package main

import (
	"log"

	"github.com/Yarik7610/redis-go-example/internal/config"
	"github.com/Yarik7610/redis-go-example/internal/controller"
	"github.com/Yarik7610/redis-go-example/internal/repository"
	"github.com/Yarik7610/redis-go-example/internal/service"
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
