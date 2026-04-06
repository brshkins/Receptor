package model

// Соответствует таблице распознанных ингредиентов в БД.
type DetectedIngredient struct {
	UploadID      int64
	IngredientID  int64
	Confidence    float64
}
