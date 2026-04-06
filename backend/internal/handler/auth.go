package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/middleware"
	"receptor/backend/internal/service"
)

// Предоставляет HTTP-эндпоинты для аутентификации.
type AuthHandler struct {
	svc service.AuthService
}

// Создает новый экземпляр AuthHandler.
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Регистрация нового пользователя.
func (h *AuthHandler) Register(c *gin.Context) {
	var in dto.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.svc.Register(c.Request.Context(), in)
	if err != nil {
		RespondError(c, err)
		return
	}
	JSONOK(c, out)
}

// Вход в систему.
func (h *AuthHandler) Login(c *gin.Context) {
	var in dto.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.svc.Login(c.Request.Context(), in)
	if err != nil {
		RespondError(c, err)
		return
	}
	JSONOK(c, out)
}

// Получение информации о текущем пользователе.
func (h *AuthHandler) Me(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		JSONError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	out, err := h.svc.Me(c.Request.Context(), uid)
	if err != nil {
		RespondError(c, err)
		return
	}
	JSONOK(c, out)
}
