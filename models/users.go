package models

import (
	"time"

	"gorm.io/gorm"
)

type NPLStatus string

func (n *NPLStatus) toString() string {
	return string(*n)
}

const (
	NPL_SMOOTH_PAYMENT     NPLStatus = ""
	NPL_DELINQUENT_PAYMENT NPLStatus = "delinquent"
)

type User struct {
	ID        int64          `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
}

type Profile struct {
	ID              int64          `json:"id" gorm:"primarykey"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at"`
	UserID          int64          `json:"user_id"`
	PlatformAccount string         `json:"platform_account"`
	NplStatus       NPLStatus      `json:"npl_status"`
}

type UserWithProfile struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	PlatformAccount string    `json:"platform_account"`
	NplStatus       NPLStatus `json:"npl_status"`
}
