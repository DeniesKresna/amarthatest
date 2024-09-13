// database/connection.go
package config

import (
	"errors"

	"github.com/DeniesKresna/amarthatest/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func (c *Config) InitDatabase() (err error) {
	dsn := "denies:deniespassword@tcp(localhost:3999)/amarthatest?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&models.User{}, &models.Profile{}, &models.Loan{}, &models.Repayment{})
	if err != nil {
		return
	}

	err = initDBRecords(db)

	c.DB = db
	if err != nil {
		return
	}
	return
}

func initDBRecords(db *gorm.DB) (err error) {
	if db.Migrator().HasTable(&models.User{}) {
		err = db.First(&models.User{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Create([]models.User{
				{
					Name:  "Onic Albert",
					Email: "albert@onic.com",
				}, {
					Name:  "RRQ Skylar",
					Email: "skylar@rrq.com",
				}, {
					Name:  "Evos Tazz",
					Email: "tazz@evos.id",
				},
			}).Error
			if err != nil {
				return
			}
		}
	}

	if db.Migrator().HasTable(&models.Profile{}) {
		err = db.First(&models.Profile{}).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Create([]models.Profile{
				{
					UserID:          1,
					PlatformAccount: "1234",
					NplStatus:       "",
				}, {
					UserID:          2,
					PlatformAccount: "9283",
					NplStatus:       "",
				}, {
					UserID:          3,
					PlatformAccount: "2222",
					NplStatus:       "",
				},
			}).Error
			if err != nil {
				return
			}
		}
	}
	return
}
