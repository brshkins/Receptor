package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"receptor/backend/internal/middleware"
	"receptor/backend/internal/service"
)

// FavoriteHandler exposes favorites HTTP endpoints (backend_spec §7 Favorites).
type FavoriteHandler struct {
	svc service.FavoriteService
}

// NewFavoriteHandler wires favorites HTTP handlers.
func NewFavoriteHandler(svc service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{svc: svc}
}

// Add handles POST /favorites/:id (requires JWT middleware).
func (h *FavoriteHandler) Add(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		JSONError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	rid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || rid <= 0 {
		JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Add(c.Request.Context(), uid, rid); err != nil {
		RespondError(c, err)
		return
	}
	JSONData(c, http.StatusOK, nil)
}

// Remove handles DELETE /favorites/:id (requires JWT middleware).
func (h *FavoriteHandler) Remove(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		JSONError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	rid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || rid <= 0 {
		JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Remove(c.Request.Context(), uid, rid); err != nil {
		RespondError(c, err)
		return
	}
	JSONData(c, http.StatusOK, nil)
}

// List handles GET /favorites (requires JWT middleware).
func (h *FavoriteHandler) List(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		JSONError(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := h.svc.GetAll(c.Request.Context(), uid)
	if err != nil {
		RespondError(c, err)
		return
	}
	JSONOK(c, list)
}
