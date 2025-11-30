package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primarykey;not null;index"`
	Username     string    `gorm:"type:varchar(64);unique;not null"`
	PasswordHash []byte    `gorm:"not null"`
}

type Task struct {
	ID          uuid.UUID `gorm:"type:uuid;primarykey;not null;index"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Title       string    `gorm:"type:varchar(128);not null"`
	Description string    `gorm:"type:text"`
	Status      uint8     `gorm:"not null"`
	Priority    uint8     `gorm:"not null"`
	DueDate     int64
	CreatedAt   int64 `gorm:"not null"`
}
