package model

import (
	"gorm.io/gorm"
)

// DB is gorm pointer
var DB *gorm.DB

// LoadModel setup all stuff about gorm model
func LoadModel(db *gorm.DB) {
	DB = db
	DB.AutoMigrate(&Epub{})
	DB.AutoMigrate(&List{})
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&UserListEpub{})
}
