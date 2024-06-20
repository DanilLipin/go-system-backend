package models

import "time"

type Post struct {
	Id             uint        `json:"id" gorm:"primaryKey"`
	Title          string      `json:"title" gorm:"not null"`
	Description    string      `json:"description" gorm:"not null"`
	Images         string      `json:"images" gorm:"type:text"`
	IsUsed         bool        `json:"is_used"`
	EstimatedPrice uint        `json:"estimated_price"`
	CreatedAt      *time.Time  `json:"created_at" gorm:"not null;default:now()"`
	IsSubmit       bool        `json:"is_submit"`
	IsArchive      bool        `json:"is_archive"`
	UserId         string      `json:"user_id" gorm:"not null"`
	User           User        `json:"user" gorm:"foreignKey:UserId"`
	CategoryId     string      `json:"category_id" gorm:"not null"`
	Category       Category    `json:"category" gorm:"foreignKey:CategoryId"`
	Chats          []*ChatRoom `gorm:"many2many:chats_posts;"`
	DeletedAt      *time.Time  `json:"deleted_at" gorm:"index"`
}
