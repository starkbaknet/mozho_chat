package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"mozho_chat/internal/models"
)

type ChatRoomRepository interface {
	Create(room *models.ChatRoom) error
	AddUser(roomID, userID string) error
	RemoveUser(roomID, userID string) error
	FindByID(id string) (*models.ChatRoom, error)
	ListRoomsByUser(userID string) ([]models.ChatRoom, error)
	Delete(room *models.ChatRoom) error
	CountUsers(roomID string) (int64, error)
	IsUserInRoom(roomID, userID string) (bool, error)
	FindRoomBetweenUsers(userID1, userID2 string) (*models.ChatRoom, error)
}

type chatRoomRepo struct {
	db *gorm.DB
}

func NewChatRoomRepository(db *gorm.DB) ChatRoomRepository {
	return &chatRoomRepo{db: db}
}

func (r *chatRoomRepo) Create(room *models.ChatRoom) error {
	return r.db.Create(room).Error
}

func (r *chatRoomRepo) AddUser(roomID, userID string) error {
	cru := models.ChatRoomMember{
		ChatRoomID: uuidFromString(roomID),
		UserID:     uuidFromString(userID),
		JoinedAt:   time.Now(),
	}
	return r.db.Create(&cru).Error
}

func (r *chatRoomRepo) RemoveUser(roomID, userID string) error {
	return r.db.Delete(&models.ChatRoomMember{}, "chat_room_id = ? AND user_id = ?", roomID, userID).Error
}

func (r *chatRoomRepo) FindByID(id string) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := r.db.Preload("Users").First(&room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *chatRoomRepo) ListRoomsByUser(userID string) ([]models.ChatRoom, error) {
	var rooms []models.ChatRoom
	err := r.db.Joins("JOIN chat_room_members crm ON crm.chat_room_id = chat_rooms.id").
		Where("crm.user_id = ?", userID).
		Preload("Users").
		Find(&rooms).Error
	return rooms, err
}

func (r *chatRoomRepo) Delete(room *models.ChatRoom) error {
	return r.db.Delete(room).Error
}

func (r *chatRoomRepo) CountUsers(roomID string) (int64, error) {
	var count int64
	err := r.db.Model(&models.ChatRoomMember{}).Where("chat_room_id = ?", roomID).Count(&count).Error
	return count, err
}

func (r *chatRoomRepo) IsUserInRoom(roomID, userID string) (bool, error) {
	var crm models.ChatRoomMember
	err := r.db.First(&crm, "chat_room_id = ? AND user_id = ?", roomID, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (r *chatRoomRepo) FindRoomBetweenUsers(userID1, userID2 string) (*models.ChatRoom, error) {
	var room models.ChatRoom
	// Find rooms where both users are members and it's not a group chat
	err := r.db.Table("chat_rooms").
		Joins("JOIN chat_room_members crm1 ON crm1.chat_room_id = chat_rooms.id").
		Joins("JOIN chat_room_members crm2 ON crm2.chat_room_id = chat_rooms.id").
		Where("crm1.user_id = ? AND crm2.user_id = ? AND crm1.user_id != crm2.user_id AND chat_rooms.is_group = false", userID1, userID2).
		Preload("Users").
		First(&room).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // No room found, return nil without error
	}
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func uuidFromString(id string) (u uuid.UUID) {
	u, _ = uuid.Parse(id)
	return
}
