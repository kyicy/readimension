package model

import (
	"github.com/jinzhu/gorm"
)

type Hero struct {
	gorm.Model
	Users []User
}
