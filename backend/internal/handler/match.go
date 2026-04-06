package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/service"
)

// Предоставляет HTTP-эндпоинты для подбора рецептов по ингредиентам.
type MatchHandler struct {
	svc service.MatchService
}

// Создает новый экземпляр MatchHandler.
func NewMatchHandler(svc service.MatchService) *MatchHandler {
	return &MatchHandler{svc: svc}
}

// Подбор рецептов по ингредиентам.
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
