package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mozho_chat/internal/repository"
	"mozho_chat/internal/user"
	"mozho_chat/internal/chatroom"
	"mozho_chat/internal/message"
	"mozho_chat/pkg/middleware"
	"mozho_chat/pkg/s3"
	"mozho_chat/pkg/encryption"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	v1 := r.Group("/api/v1")

	// User
	userRepo := repository.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(v1)

	// Chat Room
	chatRoomRepo := repository.NewChatRoomRepository(db)
	chatRoomService := chatroom.NewService(chatRoomRepo, userRepo)
	chatRoomHandler := chatroom.NewHandler(chatRoomService)
	chatRoomHandler.RegisterRoutes(v1)

	// Message ✅
	messageRepo := repository.NewMessageRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)
	s3Service, err := s3.NewS3Service()
	if err != nil {
		panic("Failed to initialize S3 service: " + err.Error())
	}
	encryptionService := encryption.NewEncryptionService()
	messageService := message.NewMessageService(messageRepo, chatRoomRepo, userRepo, attachmentRepo, s3Service, encryptionService)
	messageHandler := message.NewHandler(messageService)
	messageHandler.RegisterRoutes(v1)

	return r
}
