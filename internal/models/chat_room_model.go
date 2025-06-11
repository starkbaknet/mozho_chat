package models

import (
    "time"

    "github.com/google/uuid"
)

type ChatRoom struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Name      string
    IsGroup   bool      `gorm:"default:false;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}
