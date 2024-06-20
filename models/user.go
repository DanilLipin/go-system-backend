package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           uint       `json:"id" gorm:"primaryKey"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Email        string     `json:"email" gorm:"unique"`
	Password     []byte     `json:"-"`
	Phone        string     `json:"phone"`
	ProfilePhoto string     `json:"profile_photo"`
	IsAdmin      bool       `json:"is_admin"`
	CreatedAt    *time.Time `json:"created_at" gorm:"not null;default:now()"`
	Posts        []Post     `json:"posts"`
}

func (user *User) SetPassword(password string) {
	hashedRassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedRassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
