package dto

// Элемент списка/карточки рецепта в API.
type RecipeResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Image       string `json:"image"`
	CookingTime int    `json:"cooking_time"`
	Category    string `json:"category"`
}

// Значение query sort для списка рецептов.
type RecipeSort string

const (
	RecipeSortAlphabetAsc  RecipeSort = "alphabet_asc"
	RecipeSortAlphabetDesc RecipeSort = "alphabet_desc"
	RecipeSortTimeAsc      RecipeSort = "time_asc"
	RecipeSortTimeDesc     RecipeSort = "time_desc"
)

// Фильтры списка рецептов.
type RecipeFilter struct {
	Search   *string
	Category *string
	MaxTime  *int
	Sort     RecipeSort
}
