package service

import (
	"errors"
	"net/mail"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ValidateRegisterName проверяет имя при регистрации в Telegram.
func ValidateRegisterName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("Введите имя.")
	}
	n := utf8.RuneCountInString(name)
	if n < 2 {
		return errors.New("Имя слишком короткое — минимум 2 символа.")
	}
	if n > 80 {
		return errors.New("Имя слишком длинное — не более 80 символов.")
	}
	hasLetter := false
	for _, r := range name {
		if unicode.IsLetter(r) {
			hasLetter = true
			break
		}
	}
	if !hasLetter {
		return errors.New("В имени должна быть хотя бы одна буква.")
	}
	return nil
}

// ValidateEmail проверяет формат email (регистрация и вход).
func ValidateEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return errors.New("Введите email.")
	}
	if len(email) > 254 {
		return errors.New("Email слишком длинный.")
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || strings.TrimSpace(addr.Address) == "" {
		return errors.New("Некорректный email. Пример: user@example.com")
	}
	return nil
}

// ValidateRegisterPassword проверяет пароль при регистрации (bcrypt учитывает не более 72 байт).
func ValidateRegisterPassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return errors.New("Введите пароль.")
	}
	if len(password) > 72 {
		return errors.New("Пароль слишком длинный — не более 72 символов (ограничение безопасного хранения).")
	}
	if utf8.RuneCountInString(password) < 8 {
		return errors.New("Пароль слишком короткий — минимум 8 символов.")
	}
	var hasLetter, hasDigit bool
	for _, r := range password {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsNumber(r) {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return errors.New("Пароль должен содержать и буквы, и цифры.")
	}
	return nil
}
