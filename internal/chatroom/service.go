package chatroom

import (
	"errors"

	"github.com/google/uuid"
	"mozho_chat/internal/models"
	"mozho_chat/internal/repository"
	"mozho_chat/internal/chatroom/dto"
)

type Service interface {
	CreateRoom(userID string, input dto.CreateChatRoomRequest) (*dto.ChatRoomResponse, error)
	GetRoom(userID, roomID string) (*dto.ChatRoomResponse, error)
	JoinRoom(userID, roomID string) error
	LeaveRoom(userID, roomID string) error
	ListRooms(userID string) ([]dto.ChatRoomResponse, error)
	DeleteRoom(userID, roomID string) error
}

type chatRoomService struct {
	repo repository.ChatRoomRepository
	userRepo repository.UserRepository
}

func NewService(repo repository.ChatRoomRepository, userRepo repository.UserRepository) Service {
	return &chatRoomService{repo: repo, userRepo: userRepo}
}

func (s *chatRoomService) CreateRoom(userID string, input dto.CreateChatRoomRequest) (*dto.ChatRoomResponse, error) {
	// Check that other user exists
	_, err := s.userRepo.FindByID(input.OtherUserID.String())
	if err != nil {
		return nil, errors.New("other user not found")
	}

	// Check if a chat room already exists between these two users
	existingRoom, err := s.repo.FindRoomBetweenUsers(userID, input.OtherUserID.String())
	if err != nil {
		return nil, err
	}
	if existingRoom != nil {
		// Return existing room
		response := mapChatRoomToDTO(existingRoom)
		return &response, nil
	}

	// Create new chat room
	room := &models.ChatRoom{
		ID:      uuid.New(),
		IsGroup: false, // This is a direct message, not a group
	}
	if err := s.repo.Create(room); err != nil {
		return nil, err
	}

	// Add both users
	if err := s.repo.AddUser(room.ID.String(), userID); err != nil {
		return nil, err
	}
	if err := s.repo.AddUser(room.ID.String(), input.OtherUserID.String()); err != nil {
		return nil, err
	}

	// Load room with users to return
	room, err = s.repo.FindByID(room.ID.String())
	if err != nil {
		return nil, err
	}

	response := mapChatRoomToDTO(room)
	return &response, nil
}

func (s *chatRoomService) GetRoom(userID, roomID string) (*dto.ChatRoomResponse, error) {
	// Check if user is in the room
	inRoom, err := s.repo.IsUserInRoom(roomID, userID)
	if err != nil {
		return nil, err
	}
	if !inRoom {
		return nil, errors.New("user not authorized to view this room")
	}

	// Get the room with users
	room, err := s.repo.FindByID(roomID)
	if err != nil {
		return nil, err
	}

	response := mapChatRoomToDTO(room)
	return &response, nil
}

func (s *chatRoomService) JoinRoom(userID, roomID string) error {
	// Check if user already in room
	inRoom, err := s.repo.IsUserInRoom(roomID, userID)
	if err != nil {
		return err
	}
	if inRoom {
		return errors.New("user already in room")
	}

	// Check room user count
	count, err := s.repo.CountUsers(roomID)
	if err != nil {
		return err
	}
	if count >= 2 {
		return errors.New("room is full")
	}

	return s.repo.AddUser(roomID, userID)
}

func (s *chatRoomService) LeaveRoom(userID, roomID string) error {
	inRoom, err := s.repo.IsUserInRoom(roomID, userID)
	if err != nil {
		return err
	}
	if !inRoom {
		return errors.New("user not in room")
	}
	if err := s.repo.RemoveUser(roomID, userID); err != nil {
		return err
	}

	// Optionally: delete room if no users left
	count, err := s.repo.CountUsers(roomID)
	if err != nil {
		return err
	}
	if count == 0 {
		room, err := s.repo.FindByID(roomID)
		if err != nil {
			return err
		}
		return s.repo.Delete(room)
	}

	return nil
}

func (s *chatRoomService) ListRooms(userID string) ([]dto.ChatRoomResponse, error) {
	rooms, err := s.repo.ListRoomsByUser(userID)
	if err != nil {
		return nil, err
	}
	res := make([]dto.ChatRoomResponse, 0, len(rooms))
	for _, room := range rooms {
		res = append(res, mapChatRoomToDTO(&room))
	}
	return res, nil
}

func (s *chatRoomService) DeleteRoom(userID, roomID string) error {
	// Check if user is in room
	inRoom, err := s.repo.IsUserInRoom(roomID, userID)
	if err != nil {
		return err
	}
	if !inRoom {
		return errors.New("user not authorized to delete this room")
	}

	room, err := s.repo.FindByID(roomID)
	if err != nil {
		return err
	}

	return s.repo.Delete(room)
}

func mapChatRoomToDTO(room *models.ChatRoom) dto.ChatRoomResponse {
	// Map users from the room's Users association
	users := make([]dto.UserBasic, len(room.Users))
	for i, user := range room.Users {
		users[i] = dto.UserBasic{
			ID:       user.ID,
			Username: user.Username,
		}
	}
	return dto.ChatRoomResponse{
		ID:    room.ID,
		Users: users,
	}
}
