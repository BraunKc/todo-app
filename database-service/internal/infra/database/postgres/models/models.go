package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID           string
	Username     string
	PasswordHash string
}

type Task struct {
	gorm.Model
	ID          string
	Title       string
	Description string
	Status      int32
	Priority    int32
	DueDate     int64
	CreatedAt   int64
}
