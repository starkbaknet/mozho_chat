package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mozho_chat/internal/repository"
	"mozho_chat/internal/user"
	"mozho_chat/pkg/middleware"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	userRepo := repository.NewUserRepository(db)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(v1)


	r.Use(middleware.CORSMiddleware())

	return r
}
