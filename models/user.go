package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id        uint   `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Password  []byte `lson:"-"`
	Phone     string `json:"phone" db:"phone"`
}

func (user *User) SetPassword(password string) {
	hashedRassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user.Password = hashedRassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
