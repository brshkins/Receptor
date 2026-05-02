package dto

type MatchResponse struct {
	RecipeID           int64    `json:"recipe_id"`
	Title              string   `json:"title"`
	Image              string   `json:"image"`
	MatchPercent       float64  `json:"match_percent"`
	Ingredients        []string `json:"ingredients"`
	MissingIngredients []string `json:"missing_ingredients"`
}

