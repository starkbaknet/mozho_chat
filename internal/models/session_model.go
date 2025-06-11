package models

import (
    "time"

    "github.com/google/uuid"
)

type Session struct {
    ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
    RefreshToken string    `gorm:"not null"`
    Device       string
    IPAddress    string
    ExpiresAt    *time.Time
    CreatedAt    time.Time `gorm:"autoCreateTime"`
}
