package model

import (
	"github.com/jinzhu/gorm"
)

// DB is gorm pointer
var DB *gorm.DB

// LoadModel setup all stuff about gorm model
func LoadModel(db *gorm.DB) {
	DB = db
	DB.AutoMigrate(
		&Book{},
		&Epub{},
		&List{},
		&Raindrop{},
		&User{})
}
