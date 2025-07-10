package db

import (
	"github.com/tikayere/userservice/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&repository.User{}, &repository.Role{}, repository.Permission{}, repository.RolePermission{})

	return db, nil

}
