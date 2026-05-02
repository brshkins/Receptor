package dto

// Детальная карточка рецепта для GET /recipes/:id.
// Не заменяет RecipeResponse (list item), а расширяет его для details-view.
type RecipeDetailsResponse struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Image       string   `json:"image"`
	CookingTime int      `json:"cooking_time"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	Ingredients []string `json:"ingredients"`
}

