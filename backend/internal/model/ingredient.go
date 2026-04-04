package model

// Ingredient соответствует таблице ингредиентов в БД.
type Ingredient struct {
	ID      int64
	Name    string
	Aliases []string
}
