package domain

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id" form:"id"`
	Name     string    `json:"name" form:"name" validate:"required,min=5,max=100"`
	Email    string    `json:"email" form:"email" validate:"required,email"`
	Password string    `json:"password" form:"password" validate:"required,min=6,max=30"`
}
type LoginUser struct {
	ID       uuid.UUID `json:"id" form:"id"`
	Email    string    `json:"email" form:"email" validate:"required,email"`
	Password string    `json:"password" form:"password" validate:"required,min=6,max=30"`
}
