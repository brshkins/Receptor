package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type seedRecipe struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CookingTime int    `json:"cooking_time"`
	Category    string `json:"category"`
	Ingredients []struct {
		Name   string `json:"name"`
		Amount string `json:"amount"`
	} `json:"ingredients"`
	Steps []string `json:"steps"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fatal(errors.New("DB_URL is empty"))
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		fatal(err)
	}

	recipes, err := readRecipesJSON("recipes.json")
	if err != nil {
		fatal(err)
	}

	for i := range recipes {
		if err := seedOneRecipe(ctx, pool, &recipes[i]); err != nil {
			fatal(fmt.Errorf("seed recipe %q: %w", recipes[i].Title, err))
		}
	}
}

func readRecipesJSON(filename string) ([]seedRecipe, error) {
	paths := make([]string, 0, 3)

	// Prefer the directory of this source file (works with `go run`).
	if _, thisFile, _, ok := runtime.Caller(0); ok {
		paths = append(paths, filepath.Join(filepath.Dir(thisFile), filename))
	}
	// Fallback to current working directory (works with `go run` from backend/).
	if wd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(wd, filename))
	}
	// As a last resort, try next to the executable (works with `go build`).
	if exe, err := os.Executable(); err == nil {
		paths = append(paths, filepath.Join(filepath.Dir(exe), filename))
	}

	var b []byte
	var readErr error
	for _, p := range paths {
		b, readErr = os.ReadFile(p)
		if readErr == nil {
			break
		}
	}
	if readErr != nil {
		return nil, fmt.Errorf("read %s: %w", filename, readErr)
	}

	var recipes []seedRecipe
	if err := json.Unmarshal(b, &recipes); err != nil {
		return nil, err
	}
	return recipes, nil
}

func seedOneRecipe(ctx context.Context, pool *pgxpool.Pool, r *seedRecipe) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var recipeID int
	if err := tx.QueryRow(
		ctx,
		`INSERT INTO recipes (title, description, image_url, cooking_time, category)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		r.Title,
		r.Description,
		"",
		r.CookingTime,
		r.Category,
	).Scan(&recipeID); err != nil {
		return err
	}

	for _, ing := range r.Ingredients {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO ingredients (name) VALUES ($1)
			 ON CONFLICT (name) DO NOTHING`,
			ing.Name,
		); err != nil {
			return err
		}

		var ingredientID int
		if err := tx.QueryRow(ctx, `SELECT id FROM ingredients WHERE name=$1`, ing.Name).Scan(&ingredientID); err != nil {
			return err
		}

		if _, err := tx.Exec(
			ctx,
			`INSERT INTO recipe_ingredients (recipe_id, ingredient_id, amount)
			 VALUES ($1, $2, $3)`,
			recipeID,
			ingredientID,
			ing.Amount,
		); err != nil {
			return err
		}
	}

	for idx, step := range r.Steps {
		stepNumber := idx + 1
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO recipe_steps (recipe_id, step_number, description)
			 VALUES ($1, $2, $3)`,
			recipeID,
			stepNumber,
			step,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

