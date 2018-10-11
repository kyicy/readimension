package model

import (
	"github.com/jinzhu/gorm"
)

type List struct {
	gorm.Model

	Name string `gorm:"type:varchar(255)"`

	UserID uint
	Epubs  []Epub `gorm:"many2many:list_epubs;"`

	UpVote   uint `gorm:"default:0"`
	DownVote uint `gorm:"default:0"`

	Raindrops []Raindrop
}
