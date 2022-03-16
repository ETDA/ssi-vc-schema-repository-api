package models

import (
	"gorm.io/gorm"
	"time"
)

type Token struct {
	ID        string          `json:"id" gorm:"id"`
	Name      string          `json:"name" gorm:"name"`
	Token     string          `json:"token" gorm:"token"`
	Role      string          `json:"role" gorm:"role"`
	CreatedAt *time.Time      `json:"created_at" gorm:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"deleted_at"`
}

func (m Token) TableName() string {
	return "tokens"
}
