package dto

import "github.com/google/uuid"

type CreateChatRoomRequest struct {
	OtherUserID uuid.UUID `json:"other_user_id" binding:"required"`
}

type ChatRoomResponse struct {
	ID    uuid.UUID   `json:"id"`
	Users []UserBasic `json:"users"`
}

type UserBasic struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}
