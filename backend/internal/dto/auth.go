package dto

// Тело запроса регистрации.
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Тело запроса входа.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Ответ с токеном.
type AuthResponse struct {
	Token string `json:"token"`
}

// Ответ с информацией о пользователе.
type MeResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
