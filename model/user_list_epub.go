package model

type UserListEpub struct {
	UserID uint `gorm:"primary_key;auto_increment:false"`
	ListID uint `gorm:"primary_key;auto_increment:false"`
	EpubID uint `gorm:"primary_key;auto_increment:false"`
}
