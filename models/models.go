package models

type User struct {
	ID       int    `json:"id" validate:"required"`
	UserName string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}
