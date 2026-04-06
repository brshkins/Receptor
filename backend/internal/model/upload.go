package model

import "time"

// Соответствует таблице загрузок в БД.
type Upload struct {
	ID        int64
	UserID    int64
	ImageURL  string
	CreatedAt time.Time
}
