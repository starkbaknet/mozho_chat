package dto

import (
	"mozho_chat/internal/models"
	"time"
)

type SendMessageRequest struct {
	ReceiverID        string `json:"receiver_id" form:"receiver_id" binding:"required,uuid"`
	Content           string `json:"content" form:"content" binding:"required"`
	Algorithm         string `json:"algorithm" form:"algorithm" binding:"required"`
	EncryptionKey     string `json:"encryption_key" form:"encryption_key"`
}

type MessageResponse struct {
	ID          string                   `json:"id"`
	ChatRoomID  string                   `json:"chat_room_id"`
	SenderID    string                   `json:"sender_id"`
	Content     string                   `json:"content,omitempty"`
	Encrypted   bool                     `json:"encrypted"`
	Encryption  *EncryptionMetadata      `json:"encryption,omitempty"`
	Attachments []string                 `json:"attachments,omitempty"`
	CreatedAt   string                   `json:"created_at"`
}

type EncryptionMetadata struct {
	Algorithm string `json:"algorithm"`
	Key       string `json:"key,omitempty"`
}

// NewMessageResponse creates a MessageResponse from a Message model
func NewMessageResponse(msg *models.Message) *MessageResponse {
	response := &MessageResponse{
		ID:         msg.ID.String(),
		ChatRoomID: msg.ChatRoomID.String(),
		SenderID:   msg.SenderID.String(),
		Content:    msg.Content,
		CreatedAt:  msg.CreatedAt.Format(time.RFC3339),
	}

	// Add encryption metadata if present
	if msg.EncryptionMetadata.Algorithm != "" {
		response.Encrypted = true
		response.Encryption = &EncryptionMetadata{
			Algorithm: msg.EncryptionMetadata.Algorithm,
			// Don't expose the key in responses for security
		}
	}

	return response
}

// ToMessageResponses converts a slice of Message models to MessageResponse DTOs
func ToMessageResponses(messages []models.Message) []MessageResponse {
	responses := make([]MessageResponse, len(messages))
	for i, msg := range messages {
		responses[i] = *NewMessageResponse(&msg)
	}
	return responses
}

type GetMessagesQuery struct {
	ChatRoomID string `form:"chat_room_id" binding:"required,uuid"`
	Limit      int    `form:"limit,default=50"`
	Offset     int    `form:"offset,default=0"`
}

type MarkReadRequest struct {
	MessageID string `json:"message_id" binding:"required,uuid"`
}
