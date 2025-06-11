package models

import (
    "time"

    "github.com/google/uuid"
)

type ChatRoom struct {
    ID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Name    string
    Users   []User    `gorm:"many2many:chat_room_members;foreignKey:ID;joinForeignKey:ChatRoomID;References:ID;joinReferences:UserID"`
    IsGroup bool      `gorm:"default:false;not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}
