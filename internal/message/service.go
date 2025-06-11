package message

import (
	"context"
	"encoding/base64"
	"errors"
	"mime/multipart"
	"mozho_chat/internal/message/dto"
	"mozho_chat/internal/models"
	"mozho_chat/internal/repository"
	"mozho_chat/pkg/encryption"
	s3upload "mozho_chat/pkg/s3"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	SendMessage(senderID string, input dto.SendMessageRequest, files []*multipart.FileHeader) (*dto.MessageResponse, error)
	GetMessages(chatRoomID, userID string, limit, offset int) ([]dto.MessageResponse, error)
	MarkAsRead(userID, messageID string) error
	MarkMessageRead(userID, messageID string) error
	MarkMessageUnread(userID, messageID string) error
	GenerateAESKey() (string, error)
}

type messageService struct {
	repo         repository.MessageRepository
	roomRepo     repository.ChatRoomRepository
	userRepo     repository.UserRepository
	attachmentRepo repository.AttachmentRepository
	s3Service    s3upload.Service
	encryption   encryption.EncryptionService
}

func NewMessageService(
	repo repository.MessageRepository,
	roomRepo repository.ChatRoomRepository,
	userRepo repository.UserRepository,
	attachmentRepo repository.AttachmentRepository,
	s3Service s3upload.Service,
	encryption encryption.EncryptionService,
) Service {
	return &messageService{repo, roomRepo, userRepo, attachmentRepo, s3Service, encryption}
}

func (s *messageService) SendMessage(senderID string, input dto.SendMessageRequest, files []*multipart.FileHeader) (*dto.MessageResponse, error) {
	// For now, we'll use the ReceiverID as the ChatRoomID
	// In a real implementation, you'd look up or create a private room between sender and receiver
	chatRoomID := mustParseUUID(input.ReceiverID)

	// Encrypt the content
	encryptedContent, encryptionKey, err := s.encryptMessage(input.Content, input.Algorithm, input.EncryptionKey)
	if err != nil {
		return nil, err
	}

	msg := &models.Message{
		ID:         uuid.New(),
		ChatRoomID: chatRoomID,
		SenderID:   mustParseUUID(senderID),
		Content:    encryptedContent,
		CreatedAt:  time.Now(),
		EncryptionMetadata: models.EncryptionMetadata{
			Algorithm: input.Algorithm,
			Key:       encryptionKey,
		},
	}

	if err := s.repo.Create(context.TODO(), msg); err != nil {
		return nil, err
	}

	for _, file := range files {
		uploaded, err := s.s3Service.UploadFile(file, "messages", msg.ID.String(), "attachment", false)
		if err != nil {
			return nil, err
		}
		attachment := &models.Attachment{
			ID:        uuid.New(),
			MessageID: msg.ID,
			Key:       uploaded.Key,
			FileName:  file.Filename,
			Size:      uploaded.Size,
			MimeType:  uploaded.MimeType,
			CreatedAt: time.Now(),
		}
		s.attachmentRepo.Create(context.TODO(), attachment)
	}

	return dto.NewMessageResponse(msg), nil
}

func (s *messageService) GetMessages(chatRoomID, userID string, limit, offset int) ([]dto.MessageResponse, error) {
	msgs, err := s.repo.FindByChatRoom(chatRoomID, limit, offset)
	if err != nil {
		return nil, err
	}
	return dto.ToMessageResponses(msgs), nil
}

func (s *messageService) MarkAsRead(userID, messageID string) error {
	return s.repo.MarkRead(userID, messageID)
}

func (s *messageService) MarkMessageRead(userID, messageID string) error {
	return s.repo.MarkRead(userID, messageID)
}

func (s *messageService) MarkMessageUnread(userID, messageID string) error {
	return s.repo.MarkUnread(userID, messageID)
}

func (s *messageService) GenerateAESKey() (string, error) {
	return s.encryption.GenerateAESKey()
}

// encryptMessage encrypts the plaintext content using the specified algorithm
func (s *messageService) encryptMessage(plaintext, algorithm, providedKey string) (string, string, error) {
	switch algorithm {
	case "AES":
		// Use provided key or generate a new one
		var keyBytes []byte
		var keyForStorage string
		
		if providedKey == "" {
			// Generate a new key
			generatedKey, err := s.encryption.GenerateAESKey()
			if err != nil {
				return "", "", err
			}
			// GenerateAESKey returns base64 encoded key
			keyBytes, err = base64.StdEncoding.DecodeString(generatedKey)
			if err != nil {
				return "", "", err
			}
			keyForStorage = generatedKey
		} else {
			// Decode the provided key if it's base64 encoded
			decodedKey, err := base64.StdEncoding.DecodeString(providedKey)
			if err != nil {
				return "", "", errors.New("invalid base64 encoded key")
			}
			if len(decodedKey) != 32 {
				return "", "", errors.New("AES key must be 32 bytes")
			}
			keyBytes = decodedKey
			keyForStorage = providedKey
		}
		
		encrypted, err := s.encryption.EncryptWithAES(plaintext, string(keyBytes))
		if err != nil {
			return "", "", err
		}
		return encrypted, keyForStorage, nil
	
	case "RSA":
		if providedKey == "" {
			return "", "", errors.New("RSA public key is required")
		}
		encrypted, err := s.encryption.EncryptWithRSA(plaintext, providedKey)
		if err != nil {
			return "", "", err
		}
		return encrypted, providedKey, nil
	
	default:
		return "", "", errors.New("unsupported encryption algorithm")
	}
}

// Helper function to parse UUID from string
func mustParseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		panic("invalid UUID: " + s)
	}
	return id
}
