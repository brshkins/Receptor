package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/service"
)

// RecipeHandler exposes recipe list/detail HTTP endpoints (backend_spec §7 Recipes).
type RecipeHandler struct {
	svc service.RecipeService
}

// NewRecipeHandler wires recipe HTTP handlers.
func NewRecipeHandler(svc service.RecipeService) *RecipeHandler {
	return &RecipeHandler{svc: svc}
}

// List handles GET /recipes with optional search, category, max_time, sort.
func (h *RecipeHandler) List(c *gin.Context) {
	filter, aborted := parseRecipeFilter(c)
	if aborted {
		return
	}
	list, err := h.svc.GetAll(c.Request.Context(), filter)
	if err != nil {
		RespondError(c, err)
		return
	}
	JSONOK(c, list)
}

func parseRecipeFilter(c *gin.Context) (dto.RecipeFilter, bool) {
	var f dto.RecipeFilter
	if s := c.Query("search"); s != "" {
		f.Search = &s
	}
	if cat := c.Query("category"); cat != "" {
		f.Category = &cat
	}
	if mt := c.Query("max_time"); mt != "" {
		v, err := strconv.Atoi(mt)
		if err != nil {
			JSONError(c, http.StatusBadRequest, "invalid max_time")
			return f, true
		}
		f.MaxTime = &v
	}
	if sortStr := c.Query("sort"); sortStr != "" {
		s := dto.RecipeSort(sortStr)
		switch s {
		case dto.RecipeSortAlphabetAsc, dto.RecipeSortAlphabetDesc, dto.RecipeSortTimeAsc, dto.RecipeSortTimeDesc:
			f.Sort = s
		default:
			JSONError(c, http.StatusBadRequest, "invalid sort")
			return f, true
		}
	}
	return f, false
}

// GetByID handles GET /recipes/:id.
func (h *RecipeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	out, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		RespondError(c, err)
		return
	}
	JSONOK(c, out)
}
