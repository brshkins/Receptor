package service

import (
	"context"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/model"
)

// AuthService — регистрация, вход, текущий пользователь (см. backend_spec.md §6).
type AuthService interface {
	Register(ctx context.Context, in dto.RegisterInput) (*dto.AuthResponse, error)
	Login(ctx context.Context, in dto.LoginInput) (*dto.AuthResponse, error)
	Me(ctx context.Context, userID int64) (*model.User, error)
}

// RecipeService — список и карточка рецепта с фильтрами (см. backend_spec.md §6).
type RecipeService interface {
	GetAll(ctx context.Context, filter dto.RecipeFilter) ([]dto.RecipeResponse, error)
	GetByID(ctx context.Context, id int64) (*dto.RecipeResponse, error)
}

// FavoriteService — избранное (см. backend_spec.md §6).
type FavoriteService interface {
	Add(ctx context.Context, userID int64, recipeID int64) error
	Remove(ctx context.Context, userID int64, recipeID int64) error
	GetAll(ctx context.Context, userID int64) ([]dto.RecipeResponse, error)
}

// MatchService — подбор рецептов по ингредиентам (см. backend_spec.md §6).
type MatchService interface {
	Match(ctx context.Context, req *dto.MatchRequest) ([]dto.MatchResponse, error)
}
