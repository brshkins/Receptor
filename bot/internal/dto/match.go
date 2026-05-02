package dto

// MatchResponse соответствует backend DTO (см. backend_spec.md).
type MatchResponse struct {
	RecipeID           int64    `json:"recipe_id"`
	Title              string   `json:"title"`
	Image              string   `json:"image"`
	MatchPercent       float64  `json:"match_percent"`
	MissingIngredients []string `json:"missing_ingredients"`
}

