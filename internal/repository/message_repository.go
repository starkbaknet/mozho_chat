// internal/repository/message_repository.go
package repository

import (
	"context"
	"mozho_chat/internal/models"

	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	FindByChatRoom(chatRoomID string, limit, offset int) ([]models.Message, error)
	MarkRead(userID, messageID string) error
	MarkUnread(userID, messageID string) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *messageRepository) FindByChatRoom(chatRoomID string, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.
		Where("chat_room_id = ?", chatRoomID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *messageRepository) MarkRead(userID, messageID string) error {
	return r.db.Model(&models.MessageStatus{}).
		Where("user_id = ? AND message_id = ?", userID, messageID).
		Assign(models.MessageStatus{Read: true}).
		FirstOrCreate(&models.MessageStatus{}).Error
}

func (r *messageRepository) MarkUnread(userID, messageID string) error {
	return r.db.Model(&models.MessageStatus{}).
		Where("user_id = ? AND message_id = ?", userID, messageID).
		Update("read", false).Error
}
