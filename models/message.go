package models

import "time"

type Message struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	Content    string     `json:"content" gorm:"not null"`
	SenderID   uint       `json:"sender_id"`
	Sender     User       `json:"sender" gorm:"foreignKey:SenderID"`
	ChatRoomID uint       `json:"chat_room_id"`
	ChatRoom   ChatRoom   `json:"chat_room" gorm:"foreignKey:ChatRoomID"`
	CreatedAt  *time.Time `json:"created_at" gorm:"not null;default:now()"`
}
