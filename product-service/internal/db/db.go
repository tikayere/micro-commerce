package db

import (
	"github.com/tikayere/productservice/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&repository.Product{})

	return db, nil

}
