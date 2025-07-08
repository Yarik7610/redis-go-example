package repository

import (
	"database/sql"
	"errors"

	"github.com/Yarik7610/redis-learning/internal/model"
)

// Тут работает с БД нашей, пишем наши CRUD базовые запросы, которые будет потом переиспользовать в сервисе

type UserRepository interface {
	Save(user *model.User) error
	GetById(id int) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) Save(user *model.User) error {
	query := "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email"
	return ur.db.QueryRow(query, user.Name, user.Email).Scan(&user.ID, &user.Name, &user.Email)
}

func (ur *userRepository) GetById(id int) (*model.User, error) {
	var user model.User
	query := "SELECT id, name, email FROM users WHERE id = $1"
	err := ur.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	query := "SELECT id, name, email FROM users WHERE email = $1"
	err := ur.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
