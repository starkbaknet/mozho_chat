package models

import (
    "time"

    "github.com/google/uuid"
)

type ChatRoomMember struct {
    ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    RoomID   uuid.UUID `gorm:"type:uuid;not null;index"`
    UserID   uuid.UUID `gorm:"type:uuid;not null;index"`
    JoinedAt time.Time `gorm:"autoCreateTime"`
}
