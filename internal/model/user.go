package model

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name" validate:"required,min=2"`
	Email string `json:"email" validate:"required,email"`
}
