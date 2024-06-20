package models

import (
	"time"
)

type ChatRoom struct {
	Id        uint       `json:"id" gorm:"primaryKey"`
	Posts     []*Post    `gorm:"many2many:chats_posts;"`
	CreatedAt *time.Time `json:"created_at" gorm:"not null;default:now()"`
	Messages  []Message  `json:"messages"`
}
