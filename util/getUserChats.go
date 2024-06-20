package util

import (
	"blogbackend/database"
	"blogbackend/models"
)

func GetChatsByUserIDWithPosts(userID string) ([]models.ChatRoom, error) {
	var chats []models.ChatRoom
	result := database.DB.Preload("Posts").Joins("JOIN chats_posts ON chat_rooms.id = chats_posts.chat_room_id").
		Joins("JOIN posts ON posts.id = chats_posts.post_id").
		Where("posts.user_id = ?", userID).
		Find(&chats)
	return chats, result.Error
}

func PostBelongsToChat(postId uint) (bool, error) {
	db := database.DB

	var chatRooms []models.ChatRoom
	err := db.Joins("JOIN chats_posts ON chats_posts.chat_room_id = chat_rooms.id").
		Where("chats_posts.post_id = ?", postId).
		Find(&chatRooms).Error
	if err != nil {
		return false, err
	}
	return len(chatRooms) > 0, nil
}
