package model

import "time"

// Соответствует таблице рецептов в БД.
type Recipe struct {
	ID          int64
	Title       string
	Description string
	ImageURL    string
	CookingTime int
	Category    string
	CreatedAt   time.Time
}
