package dto

// Тело POST /match/by-ingredients.
type MatchRequest struct {
	Ingredients []string `json:"ingredients"`
}

// Элемент результата подбора рецептов.
type MatchResponse struct {
	RecipeID           int64    `json:"recipe_id"`
	Title              string   `json:"title"`
	Image              string   `json:"image"`
	MatchPercent       float64  `json:"match_percent"`
	MissingIngredients []string `json:"missing_ingredients"`
}
