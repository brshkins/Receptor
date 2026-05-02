package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ingredientClassMap = map[string]string{
	"apple":          "apple",
	"apples":         "apple",
	"banana":         "banana",
	"bananas":        "banana",
	"orange":         "orange",
	"oranges":        "orange",
	"tomato":         "tomato",
	"tomatoes":       "tomato",
	"cucumber":       "cucumber",
	"cucumbers":      "cucumber",
	"potato":         "potato",
	"potatoes":       "potato",
	"onion":          "onion",
	"onions":         "onion",
	"carrot":         "carrot",
	"carrots":        "carrot",
	"cabbage":        "cabbage",
	"cabbages":       "cabbage",
	"mushroom":       "mushroom",
	"mushrooms":      "mushroom",
	"garlic":         "garlic",
	"honey":          "honey",
	"walnut":         "walnut",
	"walnuts":        "walnut",
	"pepper":         "pepper",
	"peppers":        "pepper",
	"milk":           "milk",
	"cheese":         "cheese",
	"yogurt":         "yogurt",
	"butter":         "butter",
	"egg":            "eggs",
	"eggs":           "eggs",
	"bread":          "bread",
	"beef":           "meat",
	"pork":           "meat",
	"chicken":        "meat",
	"turkey":         "meat",
	"meat":           "meat",
	"fish":           "fish",
	"salmon":         "fish",
	"tuna":           "fish",
	"juice":          "juice",
	"water":          "water",
	"instant noodle": "instant_noodle",
	"instant noodles": "instant_noodle",
	"instant_noodle": "instant_noodle",
	"ramen":          "instant_noodle",
	"noodles":        "instant_noodle",
}

func normalizeIngredientName(raw string) string {
	norm := strings.ToLower(strings.TrimSpace(raw))
	if mapped, ok := ingredientClassMap[norm]; ok {
		return mapped
	}
	return raw
}

func ingredientDedupKey(raw string) string {
	norm := strings.ToLower(strings.TrimSpace(raw))
	if mapped, ok := ingredientClassMap[norm]; ok {
		return mapped
	}
	return norm
}

var categoryImagePools = map[string][]string{
	"pasta": {
		"https://images.unsplash.com/photo-1621996346565-e3dbc646d9a9?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1563379926898-05f4575a45d8?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1473093295043-cdd812d0e601?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1551183053-bf91a1d81141?w=1200&auto=format&fit=crop",
	},
	"meat": {
		"https://images.unsplash.com/photo-1544025162-d76694265947?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1608039829572-08f78c7c5e63?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1529692236671-f94f979ee425?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1547592180-2a62e887fd72?w=1200&auto=format&fit=crop",
	},
	"vegetarian": {
		"https://images.unsplash.com/photo-1512621776951-a57141f2eefd?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1540420773420-3366772f4999?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1490645935967-10de6ba17061?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=1200&auto=format&fit=crop",
	},
	"breakfast": {
		"https://images.unsplash.com/photo-1533089860892-a7c6f0a88666?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1525351484163-7529144344d8?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1493770348161-369560ae357d?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1484723091739-30a097e8f929?w=1200&auto=format&fit=crop",
	},
	"dessert": {
		"https://images.unsplash.com/photo-1563729784474-d77dbb933a9e?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1551782450-a2132b4ba21d?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1578985545062-69928b1d9587?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1563805042-7684c019e1cb?w=1200&auto=format&fit=crop",
	},
	"soup": {
		"https://images.unsplash.com/photo-1547592166-23ac45744acd?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1578247329974-62f1fbf53b5d?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1604908177522-8b0ee18b8338?w=1200&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1569718212165-3a8278d5f624?w=1200&auto=format&fit=crop",
	},
}

func seedRecipeImageURL(category, title, explicit string) string {
	if s := strings.TrimSpace(explicit); s != "" {
		return s
	}
	cat := strings.ToLower(strings.TrimSpace(category))
	pool := categoryImagePools[cat]
	if len(pool) == 0 {
		pool = categoryImagePools["vegetarian"]
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(strings.ToLower(strings.TrimSpace(title))))
	idx := int(h.Sum32()) % len(pool)
	if idx < 0 {
		idx = -idx
	}
	return pool[idx]
}

type seedRecipe struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
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
	if err := applyRecipeStockImages(ctx, pool); err != nil {
		fatal(fmt.Errorf("apply stock images: %w", err))
	}
}

// applyRecipeStockImages подставляет URL из recipe_stock_images.json (в т.ч. после правок в 003_recipe_stock_images.sql).
// Без этого сид оставляет Unsplash из categoryImagePools, а отредактированная миграция 003 может не выполниться повторно.
func applyRecipeStockImages(ctx context.Context, pool *pgxpool.Pool) error {
	paths := make([]string, 0, 4)
	if _, thisFile, _, ok := runtime.Caller(0); ok {
		paths = append(paths, filepath.Join(filepath.Dir(thisFile), "recipe_stock_images.json"))
	}
	if wd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(wd, "recipe_stock_images.json"))
	}
	if exe, err := os.Executable(); err == nil {
		paths = append(paths, filepath.Join(filepath.Dir(exe), "recipe_stock_images.json"))
	}

	var raw []byte
	var readErr error
	for _, p := range paths {
		raw, readErr = os.ReadFile(p)
		if readErr == nil {
			break
		}
	}
	if readErr != nil {
		return nil
	}

	var m map[string]string
	if err := json.Unmarshal(raw, &m); err != nil || len(m) == 0 {
		return nil
	}

	titles := make([]string, 0, len(m))
	urls := make([]string, 0, len(m))
	for t, u := range m {
		t = strings.TrimSpace(t)
		u = strings.TrimSpace(u)
		if t == "" || u == "" {
			continue
		}
		titles = append(titles, t)
		urls = append(urls, u)
	}
	if len(titles) == 0 {
		return nil
	}

	_, err := pool.Exec(ctx, `
		UPDATE recipes AS r SET image_url = v.url
		FROM (
			SELECT * FROM unnest($1::text[], $2::text[]) AS v(title, url)
		) v
		WHERE r.title = v.title`,
		titles, urls)
	return err
}

func readRecipesJSON(filename string) ([]seedRecipe, error) {
	paths := make([]string, 0, 3)

	if _, thisFile, _, ok := runtime.Caller(0); ok {
		paths = append(paths, filepath.Join(filepath.Dir(thisFile), filename))
	}
	if wd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(wd, filename))
	}
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

	img := seedRecipeImageURL(r.Category, r.Title, r.ImageURL)
	var recipeID int
	if err := tx.QueryRow(
		ctx,
		`INSERT INTO recipes (title, description, image_url, cooking_time, category)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		r.Title,
		r.Description,
		img,
		r.CookingTime,
		r.Category,
	).Scan(&recipeID); err != nil {
		return err
	}

	seenIngredients := make(map[string]struct{}, len(r.Ingredients))
	for _, ing := range r.Ingredients {
		key := ingredientDedupKey(ing.Name)
		if _, seen := seenIngredients[key]; seen {
			continue
		}
		seenIngredients[key] = struct{}{}

		name := normalizeIngredientName(ing.Name)
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO ingredients (name) VALUES ($1)
			 ON CONFLICT (name) DO NOTHING`,
			name,
		); err != nil {
			return err
		}

		var ingredientID int
		if err := tx.QueryRow(ctx, `SELECT id FROM ingredients WHERE name=$1`, name).Scan(&ingredientID); err != nil {
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

