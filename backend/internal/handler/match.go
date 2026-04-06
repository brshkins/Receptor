package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/service"
)

// MatchHandler exposes match HTTP endpoints (backend_spec §7 Match).
type MatchHandler struct {
	svc service.MatchService
}

// NewMatchHandler wires match HTTP handlers.
func NewMatchHandler(svc service.MatchService) *MatchHandler {
	return &MatchHandler{svc: svc}
}

// Match handles POST /match/by-ingredients.
func (h *MatchHandler) Match(c *gin.Context) {
	var req dto.MatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	out, err := h.svc.Match(c.Request.Context(), &req)
	if err != nil {
		RespondError(c, err)
		return
	}
	JSONOK(c, out)
}
