package service

import (
	"errors"

	"github.com/Yarik7610/redis-go-example/internal/model"
	"github.com/Yarik7610/redis-go-example/internal/repository"
	"github.com/gin-gonic/gin"
)

// Бизнес-логика прлиожения, отделяет работу от протокола (контроллер) и от БД (репо)
// Координация работы между разными репозиториями, валидируем данные входные

type UserService interface {
	CreateUser(ctx *gin.Context) (*model.User, error)
	GetUser(id int) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
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
	return s.repo.GetById(id)
}
