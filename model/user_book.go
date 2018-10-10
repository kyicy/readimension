package model

import (
	"github.com/jinzhu/gorm"
)

type UserBook struct {
	gorm.Model

	Epub   Epub
	EpubID uint

	UserID uint
}
