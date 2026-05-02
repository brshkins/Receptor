package dto

type RecipeResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Image       string `json:"image"`
	CookingTime int    `json:"cooking_time"`
	Category    string `json:"category"`
}

