package model

// RecipeStep соответствует шагам приготовления в БД.
type RecipeStep struct {
	RecipeID    int64
	StepNumber  int
	Description string
}
