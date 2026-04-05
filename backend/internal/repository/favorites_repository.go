package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"receptor/backend/internal/model"
)

type favoritesRepository struct {
	db *pgxpool.Pool
}

func NewFavoritesRepository(db *pgxpool.Pool) FavoritesRepository {
	return &favoritesRepository{db: db}
}

func (r *favoritesRepository) Add(ctx context.Context, userID int64, recipeID int64) error {
	const q = `
		INSERT INTO favorites (user_id, recipe_id)
		VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, q, userID, recipeID)
	return err
}

func (r *favoritesRepository) Remove(ctx context.Context, userID int64, recipeID int64) error {
	const q = `
		DELETE FROM favorites
		WHERE user_id = $1 AND recipe_id = $2`
	_, err := r.db.Exec(ctx, q, userID, recipeID)
	return err
}

func (r *favoritesRepository) GetByUser(ctx context.Context, userID int64) ([]*model.Recipe, error) {
	const q = `
		SELECT r.id, r.title, r.description, r.image_url, r.cooking_time, r.category, r.created_at
		FROM favorites f
		INNER JOIN recipes r ON r.id = f.recipe_id
		WHERE f.user_id = $1
		ORDER BY r.id ASC`

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRecipes(rows)
}
