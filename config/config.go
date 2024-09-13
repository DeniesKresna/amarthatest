package config

import "gorm.io/gorm"

type Config struct {
	DB *gorm.DB
}
