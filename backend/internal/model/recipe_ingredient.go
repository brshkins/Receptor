package model

// Соответствует связи рецепт–ингредиент в БД.
type RecipeIngredient struct {
	RecipeID      int64
	IngredientID  int64
	Amount        string
}
