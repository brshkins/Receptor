package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"receptor/backend/internal/model"
)

type ingredientsRepository struct {
	db *pgxpool.Pool
}

func NewIngredientsRepository(db *pgxpool.Pool) IngredientsRepository {
	return &ingredientsRepository{db: db}
}

func (r *ingredientsRepository) GetByNames(ctx context.Context, names []string) ([]*model.Ingredient, error) {
	if len(names) == 0 {
		return []*model.Ingredient{}, nil
	}

	// Только точное совпадение (case-insensitive) по name и aliases.
	// Подстрочный LIKE давал ложные матчи: например inp "car" → ingredient "carrot".
	const q = `
		SELECT DISTINCT i.id, i.name, COALESCE(i.aliases, ARRAY[]::text[])
		FROM ingredients i
		WHERE EXISTS (
			SELECT 1
			FROM unnest($1::text[]) AS q(inp)
			WHERE TRIM(q.inp) <> ''
			AND LOWER(TRIM(i.name)) = LOWER(TRIM(q.inp))
		)
		OR EXISTS (
			SELECT 1
			FROM unnest(COALESCE(i.aliases, ARRAY[]::text[])) AS a(alias),
			     unnest($1::text[]) AS q(inp)
			WHERE TRIM(q.inp) <> ''
			AND TRIM(a.alias) <> ''
			AND LOWER(TRIM(a.alias)) = LOWER(TRIM(q.inp))
		)`

	rows, err := r.db.Query(ctx, q, names)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model.Ingredient
	for rows.Next() {
		var ing model.Ingredient
		if err := rows.Scan(&ing.ID, &ing.Name, &ing.Aliases); err != nil {
			return nil, err
		}
		out = append(out, &ing)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *ingredientsRepository) Create(ctx context.Context, ingredient *model.Ingredient) error {
	aliases := ingredient.Aliases
	if aliases == nil {
		aliases = []string{}
	}
	const q = `
		INSERT INTO ingredients (name, aliases)
		VALUES ($1, $2)
		RETURNING id`
	return r.db.QueryRow(ctx, q, ingredient.Name, aliases).Scan(&ingredient.ID)
}
