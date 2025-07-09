package model

type User struct {
	ID    int    `json:"id" redis:"id"`
	Name  string `json:"name" validate:"required,min=2" redis:"name"`
	Email string `json:"email" validate:"required,email" redis:"email"`
}
