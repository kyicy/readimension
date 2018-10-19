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
	return "/u/covers/" + e.SHA256 + "." + e.CoverFormat
}

func (e *Epub) IsZipped() bool {
	return e.SizeByMB <= 10.0
}

func (e *Epub) StoreName() string {
	if e.IsZipped() {
		return e.SHA256 + ".epub"
	} else {
		return e.SHA256
	}
}
