package model

import "time"

// User соответствует таблице пользователей в БД.
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Name         string
	CreatedAt    time.Time
}
