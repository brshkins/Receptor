package repository

import (
	"context"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/model"
)

// Контракт доступа к пользователям.
type UsersRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

// Контракт доступа к рецептам.
type RecipesRepository interface {
	GetAll(ctx context.Context, filter dto.RecipeFilter) ([]*model.Recipe, error)
	GetByID(ctx context.Context, id int64) (*model.Recipe, error)
}

// Контракт доступа к ингредиентам.
type IngredientsRepository interface {
	GetByNames(ctx context.Context, names []string) ([]*model.Ingredient, error)
	Create(ctx context.Context, ingredient *model.Ingredient) error
}

// Контракт доступа к избранному.
type FavoritesRepository interface {
	Add(ctx context.Context, userID int64, recipeID int64) error
	Remove(ctx context.Context, userID int64, recipeID int64) error
	GetByUser(ctx context.Context, userID int64) ([]*model.Recipe, error)
}
