package model

import (
	"github.com/jinzhu/gorm"
)

type Epub struct {
	gorm.Model
	Title       string `gorm:"type:varchar(255)"`
	SHA256      string `gorm:"type:varchar(255);unique;not null"`
	SizeByMB    float64
	Author      string `gorm:"type:varchar(255)"`
	HasCover    bool
	CoverFormat string
}

func (e *Epub) CoverPath() string {
	return "/covers/" + e.SHA256 + "." + e.CoverFormat
}
