package model

import (
	"github.com/jinzhu/gorm"
)

type Raindrop struct {
	gorm.Model

	Description string `gorm:"type:varchar(255)"`

	UserID uint

	ListID uint
}
