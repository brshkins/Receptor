package dto

// RecipeResponse — элемент списка/карточки рецепта в API (см. backend_spec.md).
type RecipeResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Image       string `json:"image"`
	CookingTime int    `json:"cooking_time"`
	Category    string `json:"category"`
}

// RecipeSort — значение query sort для списка рецептов (см. backend_spec.md, GET /recipes).
type RecipeSort string

const (
	RecipeSortAlphabetAsc  RecipeSort = "alphabet_asc"
	RecipeSortAlphabetDesc RecipeSort = "alphabet_desc"
	RecipeSortTimeAsc      RecipeSort = "time_asc"
	RecipeSortTimeDesc     RecipeSort = "time_desc"
)

// RecipeFilter — фильтры списка рецептов (search, category, max_time, sort; см. backend_spec.md).
type RecipeFilter struct {
	Search   *string
	Category *string
	MaxTime  *int
	Sort     RecipeSort
}
