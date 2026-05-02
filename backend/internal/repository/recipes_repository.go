package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/model"
)

type recipesRepository struct {
	db *pgxpool.Pool
}

func NewRecipesRepository(db *pgxpool.Pool) RecipesRepository {
	return &recipesRepository{db: db}
}

func (r *recipesRepository) GetAll(ctx context.Context, filter dto.RecipeFilter) ([]*model.Recipe, error) {
	query := `
		SELECT r.id, r.title, r.description, r.image_url, r.cooking_time, r.category, r.created_at
		FROM recipes r
		WHERE 1 = 1`
	args := make([]any, 0, 4)
	n := 1

	if filter.Search != nil && *filter.Search != "" {
		pat := "%" + escapeLikePattern(*filter.Search) + "%"
		query += fmt.Sprintf(" AND (r.title ILIKE $%d ESCAPE '!' OR r.description ILIKE $%d ESCAPE '!')", n, n)
		args = append(args, pat)
		n++
	}
	if filter.Category != nil && *filter.Category != "" {
		query += fmt.Sprintf(" AND r.category = $%d", n)
		args = append(args, *filter.Category)
		n++
	}
	if filter.MaxTime != nil {
		query += fmt.Sprintf(" AND r.cooking_time <= $%d", n)
		args = append(args, *filter.MaxTime)
		n++
	}

	query += " ORDER BY " + recipeOrderByClause(filter.Sort)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRecipes(rows)
}

func (r *recipesRepository) GetByID(ctx context.Context, id int64) (*model.Recipe, error) {
	const q = `
		SELECT id, title, description, image_url, cooking_time, category, created_at
		FROM recipes
		WHERE id = $1`
	row := r.db.QueryRow(ctx, q, id)
	var rec model.Recipe
	if err := row.Scan(recipeScanDest(&rec)...); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *recipesRepository) GetStepsByRecipeID(ctx context.Context, recipeID int64) ([]model.RecipeStep, error) {
	const q = `
		SELECT recipe_id, step_number, description
		FROM recipe_steps
		WHERE recipe_id = $1
		ORDER BY step_number ASC`
	rows, err := r.db.Query(ctx, q, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.RecipeStep, 0)
	for rows.Next() {
		var s model.RecipeStep
		if err := rows.Scan(&s.RecipeID, &s.StepNumber, &s.Description); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func recipeScanDest(rec *model.Recipe) []any {
	return []any{
		&rec.ID,
		&rec.Title,
		&rec.Description,
		&rec.ImageURL,
		&rec.CookingTime,
		&rec.Category,
		&rec.CreatedAt,
	}
}

func scanRecipes(rows pgx.Rows) ([]*model.Recipe, error) {
	var out []*model.Recipe
	for rows.Next() {
		var rec model.Recipe
		if err := rows.Scan(recipeScanDest(&rec)...); err != nil {
			return nil, err
		}
		out = append(out, &rec)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func recipeOrderByClause(sort dto.RecipeSort) string {
	switch sort {
	case dto.RecipeSortAlphabetAsc:
		return "r.title ASC"
	case dto.RecipeSortAlphabetDesc:
		return "r.title DESC"
	case dto.RecipeSortTimeAsc:
		return "r.cooking_time ASC"
	case dto.RecipeSortTimeDesc:
		return "r.cooking_time DESC"
	default:
		return "r.id ASC"
	}
}

func escapeLikePattern(s string) string {
	s = strings.ReplaceAll(s, "!", "!!")
	s = strings.ReplaceAll(s, "%", "!%")
	s = strings.ReplaceAll(s, "_", "!_")
	return s
}
