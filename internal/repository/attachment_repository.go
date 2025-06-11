package repository

import (
	"context"
	"mozho_chat/internal/models"

	"gorm.io/gorm"
)

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *models.Attachment) error
	FindByMessageID(messageID string) ([]models.Attachment, error)
}

type attachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) Create(ctx context.Context, attachment *models.Attachment) error {
	return r.db.WithContext(ctx).Create(attachment).Error
}

func (r *attachmentRepository) FindByMessageID(messageID string) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := r.db.Where("message_id = ?", messageID).Find(&attachments).Error
	return attachments, err
}

