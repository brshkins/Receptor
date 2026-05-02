package dto

type MatchResponse struct {
	RecipeID           int64    `json:"recipe_id"`
	Title              string   `json:"title"`
	Image              string   `json:"image"`
	CookingTime        int      `json:"cooking_time"`
	MatchPercent       float64  `json:"match_percent"`
	Ingredients        []string `json:"ingredients"`
	MissingIngredients []string `json:"missing_ingredients"`
}

