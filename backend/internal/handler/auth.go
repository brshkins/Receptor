package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/middleware"
	"receptor/backend/internal/service"
)

// AuthHandler exposes auth HTTP endpoints (backend_spec §7 Auth).
type AuthHandler struct {
	svc service.AuthService
}

// NewAuthHandler wires auth HTTP handlers.
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Register handles POST /auth/register.
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

// Login handles POST /auth/login.
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

// Me handles GET /auth/me (requires JWT middleware).
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
