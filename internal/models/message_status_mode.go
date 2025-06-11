package models

import (
    "time"

    "github.com/google/uuid"
)

type MessageStatus struct {
    ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    MessageID  uuid.UUID `gorm:"type:uuid;not null;index"`
    UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
    Delivered  bool      `gorm:"default:false;not null"`
    Read       bool      `gorm:"default:false;not null"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
