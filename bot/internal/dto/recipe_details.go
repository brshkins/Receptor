package dto

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

