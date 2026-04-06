package model

// Соответствует шагам приготовления в БД.
type RecipeStep struct {
	RecipeID    int64
	StepNumber  int
	Description string
}
