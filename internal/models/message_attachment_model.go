package models

import (
    "time"

    "github.com/google/uuid"
)

type Attachment struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
    MessageID uuid.UUID `gorm:"type:uuid;index;not null"`
    Key       string    `gorm:"type:text;not null"`
    FileName  string    `gorm:"type:text;not null"`
    MimeType  string    `gorm:"type:text;not null"`
    Size      int64     `gorm:"not null"`
    CreatedAt time.Time
}
