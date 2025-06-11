package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/datatypes"
)

type Message struct {
    ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    RoomID    uuid.UUID      `gorm:"type:uuid;not null;index"`
    SenderID  uuid.UUID      `gorm:"type:uuid;not null;index"`
    Content   string
    Media     datatypes.JSON `gorm:"type:jsonb"`
    CreatedAt time.Time      `gorm:"autoCreateTime"`
}
