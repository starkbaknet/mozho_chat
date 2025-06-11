package models

import (
    "time"

    "github.com/google/uuid"
)

type ChatRoomMember struct {
    ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    ChatRoomID uuid.UUID `gorm:"column:chat_room_id;type:uuid;not null;index"`
    UserID     uuid.UUID `gorm:"column:user_id;type:uuid;not null;index"`
    JoinedAt   time.Time `gorm:"column:joined_at;autoCreateTime"`
}
