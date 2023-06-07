package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Password string
	Email    string
	Role     string
	Approved int
	Deleted  int `gorm:"default:0"`
}

type Topic struct {
	gorm.Model
	Topic  string
	UserID int `json:"user_id"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
