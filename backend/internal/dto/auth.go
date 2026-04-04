package dto

// RegisterInput — тело запроса регистрации (см. backend_spec.md).
type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginInput — тело запроса входа (см. backend_spec.md).
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse — ответ с токеном (см. backend_spec.md).
type AuthResponse struct {
	Token string `json:"token"`
}
