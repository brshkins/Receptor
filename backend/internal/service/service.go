package service

import (
	"context"

	"receptor/backend/internal/dto"
)

// Регистрация, вход, текущий пользователь.
type AuthService interface {
	Register(ctx context.Context, in dto.RegisterInput) (*dto.AuthResponse, error)
	Login(ctx context.Context, in dto.LoginInput) (*dto.AuthResponse, error)
	Me(ctx context.Context, userID int64) (*dto.MeResponse, error)
}

// Список и карточка рецепта с фильтрами.
type RecipeService interface {
	GetAll(ctx context.Context, filter dto.RecipeFilter) ([]dto.RecipeResponse, error)
	GetByID(ctx context.Context, id int64) (*dto.RecipeResponse, error)
}

// Избранное.
type FavoriteService interface {
	Add(ctx context.Context, userID int64, recipeID int64) error
	Remove(ctx context.Context, userID int64, recipeID int64) error
	GetAll(ctx context.Context, userID int64) ([]dto.RecipeResponse, error)
}

// Подбор рецептов по ингредиентам.
type MatchService interface {
	Match(ctx context.Context, req *dto.MatchRequest) ([]dto.MatchResponse, error)
}
