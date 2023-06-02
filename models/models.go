package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Password string
	Role     string
	Approved int
	Deleted  int `gorm:"default:0"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
