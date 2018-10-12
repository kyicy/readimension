package model

import (
	"github.com/jinzhu/gorm"
)

type Book struct {
	gorm.Model

	Epub   Epub
	EpubID uint `gorm:"unique_index:user_epub"`
	UserID uint `gorm:"unique_index:user_epub"`
}
