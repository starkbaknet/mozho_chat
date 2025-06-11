package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mozho_chat/internal/repository"
	"mozho_chat/internal/user"
	"mozho_chat/internal/chatroom" 
	"mozho_chat/pkg/middleware"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	v1 := r.Group("/api/v1")

	userRepo := repository.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(v1)

	chatRoomRepo := repository.NewChatRoomRepository(db)
	chatRoomService := chatroom.NewService(chatRoomRepo, userRepo)
	chatRoomHandler := chatroom.NewHandler(chatRoomService)
	chatRoomHandler.RegisterRoutes(v1)

	return r
}
