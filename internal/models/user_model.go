package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/datatypes"
)

type User struct {
    ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Username  string     `gorm:"unique;not null"`
    ChatRooms []ChatRoom `gorm:"many2many:chat_room_members;joinForeignKey:UserID;joinReferences:RoomID"`
    Email     string     `gorm:"unique;not null"`
    PasswordHash string   `gorm:"not null"`
    Profile   datatypes.JSON `gorm:"type:jsonb"`
    CreatedAt time.Time  `gorm:"autoCreateTime"`
}
