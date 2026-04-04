package model

// DetectedIngredient соответствует распознанным ингредиентам по загрузке в БД.
type DetectedIngredient struct {
	UploadID      int64
	IngredientID  int64
	Confidence    float64
}
