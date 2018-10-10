package model

import (
	"github.com/jinzhu/gorm"
)

type Epub struct {
	gorm.Model
	Title      string `gorm:"type:varchar(255)"`
	SHA256     string `gorm:"type:varchar(255)"`
	FilePath   string `gorm:"type:varchar(255)"`
	SizeByMB   float64
	OwnerCount uint `gorm:"default:0"`
	ViewCount  uint `gorm:"default:0"`
}
