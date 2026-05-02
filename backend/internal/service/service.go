package service

import (
	"context"

	"receptor/backend/internal/dto"
)

// ML клиент для распознавания ингредиентов по изображению.
type MLClient interface {
	DetectIngredients(ctx context.Context, image []byte) ([]string, error)
}

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
	GetDetails(ctx context.Context, id int64) (*dto.RecipeDetailsResponse, error)
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

// Загрузка изображения, распознавание ингредиентов и подбор рецептов.
type UploadService interface {
	UploadAndMatch(ctx context.Context, image []byte) ([]dto.MatchResponse, error)
}
