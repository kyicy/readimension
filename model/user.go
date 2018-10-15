package model

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User is related to `users` table in mysql database
type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"type:varchar(100);unique_index;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	List     List
	ListID   uint
}

// BeforeCreate is a hook function
func (u *User) BeforeCreate() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.Password = string(bytes)
	return err
}

// UpdatePassword converts password string into hash. You need mannual update the instance later!
func (u *User) UpdatePassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	u.Password = string(bytes)
	return err
}

// ValidatePassword compares a password string against stored password hash
func (u *User) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
