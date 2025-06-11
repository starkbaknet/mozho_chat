package models

import (
    "time"

    "github.com/google/uuid"
)

type UserPublicKey struct {
    ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
    UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
    PublicKey  string    `gorm:"type:text;not null"`
    KeyVersion string    `gorm:"type:varchar(50);default:'v1'"`
    CreatedAt  time.Time
    ExpiresAt  *time.Time `gorm:"default:null"`
    
    User User `gorm:"foreignKey:UserID"`
}
