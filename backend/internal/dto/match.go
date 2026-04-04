package dto

// MatchRequest — тело POST /match/by-ingredients (см. backend_spec.md).
type MatchRequest struct {
	Ingredients []string `json:"ingredients"`
}

// MatchResponse — элемент результата подбора рецептов (см. backend_spec.md).
type MatchResponse struct {
	RecipeID           int64    `json:"recipe_id"`
	Title              string   `json:"title"`
	Image              string   `json:"image"`
	MatchPercent       float64  `json:"match_percent"`
	MissingIngredients []string `json:"missing_ingredients"`
}
