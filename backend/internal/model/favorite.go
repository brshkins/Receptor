package model

// Соответствует таблице избранного в БД.
type Favorite struct {
	UserID   int64
	RecipeID int64
}
