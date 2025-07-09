package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Yarik7610/redis-go-example/internal/model"
	"github.com/Yarik7610/redis-go-example/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Бизнес-логика прлиожения, отделяет работу от протокола (контроллер) и от БД (репо)
// Координация работы между разными репозиториями, валидируем данные входные

type UserService interface {
	CreateUser(ctx *gin.Context) (*model.User, error)
	GetUser(id int) (*model.User, error)
}

type userService struct {
	repo        repository.UserRepository
	redisClient *redis.Client
}

func NewUserService(repo repository.UserRepository, redisClient *redis.Client) UserService {
	return &userService{repo: repo, redisClient: redisClient}
}

func (s *userService) CreateUser(ctx *gin.Context) (*model.User, error) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		return nil, errors.New("invalid input: " + err.Error())
	}

	existingUser, err := s.repo.GetByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	if err := s.repo.Save(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// В идеале тут сделать структуру с одним ID и провалидровать его также через ShouldBindJSON
func (s *userService) GetUser(id int) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	cacheKey := fmt.Sprintf("user:%d", id)
	var cachedUser model.User
	err := s.redisClient.HGetAll(context.Background(), cacheKey).Scan(&cachedUser)
	if err != nil {
		log.Printf("Redis HGetAll error: %v\n", err)
	} else {
		emptyUser := model.User{}
		if cachedUser != emptyUser {
			return &cachedUser, nil
		}
	}

	user, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	err = s.redisClient.HSet(context.Background(), cacheKey, user).Err()
	if err != nil {
		log.Printf("Redis HSet error: %v\n", err)
	}
	err = s.redisClient.Expire(context.Background(), cacheKey, time.Minute).Err()
	if err != nil {
		log.Printf("Redis Expire error: %v\n", err)
	}

	return user, nil
}
