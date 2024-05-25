package models

type Post struct {
	Id          uint   `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Image       string `json:"image" db:"image"`
	UserId      string `lson:"userid"`
	User        User   `json:"user";gorm:"foreignkey:UserID"`
}