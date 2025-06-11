// internal/repository/message_repository.go
package repository

import (
	"context"
	"mozho_chat/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	FindByChatRoom(chatRoomID string, limit, offset int) ([]models.Message, error)
	MarkRead(userID, messageID string) error
	MarkUnread(userID, messageID string) error
	MarkDelivered(userID, messageID string) error
	MarkUndelivered(userID, messageID string) error
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
	status := &models.MessageStatus{
		UserID:    mustParseUUID(userID),
		MessageID: mustParseUUID(messageID),
	}
	return r.db.
		Where("user_id = ? AND message_id = ?", userID, messageID).
		Assign(models.MessageStatus{Read: true}).
		FirstOrCreate(status).Error
}

func (r *messageRepository) MarkUnread(userID, messageID string) error {
	return r.db.Model(&models.MessageStatus{}).
		Where("user_id = ? AND message_id = ?", userID, messageID).
		Update("read", false).Error
}

func (r *messageRepository) MarkDelivered(userID, messageID string) error {
	status := &models.MessageStatus{
		UserID:    mustParseUUID(userID),
		MessageID: mustParseUUID(messageID),
	}
	return r.db.
		Where("user_id = ? AND message_id = ?", userID, messageID).
		Assign(models.MessageStatus{Delivered: true}).
		FirstOrCreate(status).Error
}

func (r *messageRepository) MarkUndelivered(userID, messageID string) error {
	return r.db.Model(&models.MessageStatus{}).
		Where("user_id = ? AND message_id = ?", userID, messageID).
		Update("delivered", false).Error
}

// Helper function to parse UUID from string
func mustParseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		panic("invalid UUID: " + s)
	}
	return id
}
