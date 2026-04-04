package model

// Favorite соответствует таблице избранного в БД.
type Favorite struct {
	UserID   int64
	RecipeID int64
}
