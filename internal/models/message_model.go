package models

import (
    "time"

    "github.com/google/uuid"
)

type Message struct {
    ID                uuid.UUID           `gorm:"type:uuid;primaryKey"`
    ChatRoomID        uuid.UUID           `gorm:"type:uuid;not null;index"`
    SenderID          uuid.UUID           `gorm:"type:uuid;not null;index"`
    Content           string              `gorm:"type:text;not null"`
    EncryptionMetadata EncryptionMetadata `gorm:"embedded"`
    CreatedAt         time.Time
    UpdatedAt         time.Time

    Sender      User         `gorm:"foreignKey:SenderID"`
    ChatRoom    ChatRoom     `gorm:"foreignKey:ChatRoomID"`
    Attachments []Attachment `gorm:"foreignKey:MessageID"`
}

type EncryptionMetadata struct {
    Algorithm string `gorm:"type:varchar(20);default:'RSA'"`
    Key       string `gorm:"type:text"`
}
