package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/service"
)

// Предоставляет HTTP-эндпоинты для списка и карточки рецепта.
type RecipeHandler struct {
	svc service.RecipeService
}

// Создает новый экземпляр RecipeHandler.
func NewRecipeHandler(svc service.RecipeService) *RecipeHandler {
	return &RecipeHandler{svc: svc}
}

// Получение списка рецептов с фильтрами.
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

// Получение карточки рецепта по ID.
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

// Получение детальной карточки рецепта (описание и шаги).
func (h *RecipeHandler) GetDetails(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		JSONError(c, http.StatusBadRequest, "invalid id")
		return
	}
	out, err := h.svc.GetDetails(c.Request.Context(), id)
	if err != nil {
		RespondError(c, err)
		return
	}
	JSONOK(c, out)
}
