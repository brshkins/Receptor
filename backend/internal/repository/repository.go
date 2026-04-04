package repository

import (
	"context"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/model"
)

// UsersRepository — контракт доступа к пользователям (см. backend_spec.md §5).
type UsersRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

// RecipesRepository — контракт доступа к рецептам (см. backend_spec.md §5).
type RecipesRepository interface {
	GetAll(ctx context.Context, filter dto.RecipeFilter) ([]*model.Recipe, error)
	GetByID(ctx context.Context, id int64) (*model.Recipe, error)
}

// IngredientsRepository — контракт доступа к ингредиентам (см. backend_spec.md §5).
type IngredientsRepository interface {
	GetByNames(ctx context.Context, names []string) ([]*model.Ingredient, error)
	Create(ctx context.Context, ingredient *model.Ingredient) error
}

// FavoritesRepository — контракт доступа к избранному (см. backend_spec.md §5).
type FavoritesRepository interface {
	Add(ctx context.Context, userID int64, recipeID int64) error
	Remove(ctx context.Context, userID int64, recipeID int64) error
	GetByUser(ctx context.Context, userID int64) ([]*model.Recipe, error)
}
