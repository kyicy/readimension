package model

import (
	"github.com/jinzhu/gorm"
)

type Book struct {
	gorm.Model

	Epub   Epub
	EpubID uint

	UserID uint
}
