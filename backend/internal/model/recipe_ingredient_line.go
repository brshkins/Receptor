package model

// Строка списка ингредиентов (имя из БД + количество из recipe_ingredients.amount).
type RecipeIngredientLine struct {
	Name   string
	Amount string
}
