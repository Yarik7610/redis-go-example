package controller

import (
	"net/http"
	"strconv"

	"github.com/Yarik7610/redis-learning/internal/service"
	"github.com/gin-gonic/gin"
)

// Передача данных в следующие слои, никакой логики, кроме логики протокола
// по типу запись ответа, указания заголовков

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{service: service}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	user, err := uc.service.CreateUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := uc.service.GetUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) RegisterRoutes(r *gin.Engine) {
	r.POST("/users/create", uc.CreateUser)
	r.GET("/users/:id", uc.GetUser)
}
