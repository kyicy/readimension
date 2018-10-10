package model

import (
	"github.com/jinzhu/gorm"
)

type BookList struct {
	gorm.Model

	Name string `gorm:"type:varchar(255)"`

	UserID uint
	Epubs  []Epub `gorm:"many2many:book_list_epubs;"`
}
