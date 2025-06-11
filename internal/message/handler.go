// internal/message/handler/handler.go
package message

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"mozho_chat/internal/message/dto"
	"mozho_chat/pkg/middleware"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	messages := rg.Group("/messages")
	messages.Use(middleware.AuthMiddleware())
	{
		messages.POST("/send", h.SendMessage)
		messages.GET("/:chat_room_id", h.GetMessages)
		messages.POST("/:message_id/read", h.MarkRead)
		messages.POST("/:message_id/unread", h.MarkUnread)
		messages.POST("/generate-key", h.GenerateKey)
	}
}

func (h *Handler) SendMessage(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	var req dto.SendMessageRequest
	req.ReceiverID = c.PostForm("receiver_id")
	req.Content = c.PostForm("content")
	req.Algorithm = c.PostForm("algorithm")
	req.EncryptionKey = c.PostForm("encryption_key")
	files := c.Request.MultipartForm.File["attachments"]

	message, err := h.service.SendMessage(userID, req, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

func (h *Handler) GetMessages(c *gin.Context) {
	userID := c.GetString("user_id")
	chatRoomID := c.Param("chat_room_id")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
		return
	}

	messages, err := h.service.GetMessages(chatRoomID, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *Handler) MarkRead(c *gin.Context) {
	userID := c.GetString("user_id")
	messageID := c.Param("message_id")

	if err := h.service.MarkMessageRead(userID, messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *Handler) MarkUnread(c *gin.Context) {
	userID := c.GetString("user_id")
	messageID := c.Param("message_id")

	if err := h.service.MarkMessageUnread(userID, messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as unread"})
}

func (h *Handler) GenerateKey(c *gin.Context) {
	key, err := h.service.GenerateAESKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": key})
}
