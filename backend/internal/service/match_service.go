package service

import (
	"context"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/repository"
)

type matchService struct {
	db          *pgxpool.Pool
	ingredients repository.IngredientsRepository
}

func NewMatchService(db *pgxpool.Pool, ingredients repository.IngredientsRepository) MatchService {
	return &matchService{db: db, ingredients: ingredients}
}

func (s *matchService) Match(ctx context.Context, req *dto.MatchRequest) ([]dto.MatchResponse, error) {
	if req == nil {
		return []dto.MatchResponse{}, nil
	}

	normalized := normalizeIngredientNames(req.Ingredients)
	if len(normalized) == 0 {
		return []dto.MatchResponse{}, nil
	}

	found, err := s.ingredients.GetByNames(ctx, normalized)
	if err != nil {
		return nil, err
	}

	userIngredientIDs := make(map[int64]struct{}, len(found))
	for _, ing := range found {
		userIngredientIDs[ing.ID] = struct{}{}
	}

	rows, err := s.db.Query(ctx, `
		SELECT r.id, r.title, r.image_url, ri.ingredient_id, i.name
		FROM recipes r
		INNER JOIN recipe_ingredients ri ON r.id = ri.recipe_id
		INNER JOIN ingredients i ON ri.ingredient_id = i.id
		ORDER BY r.id, ri.ingredient_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type bundle struct {
		title   string
		image   string
		byID    map[int64]string
	}
	recipes := make(map[int64]*bundle)

	for rows.Next() {
		var recipeID int64
		var title, imageURL, ingName string
		var ingredientID int64
		if err := rows.Scan(&recipeID, &title, &imageURL, &ingredientID, &ingName); err != nil {
			return nil, err
		}
		b, ok := recipes[recipeID]
		if !ok {
			b = &bundle{byID: make(map[int64]string)}
			recipes[recipeID] = b
		}
		b.title = title
		b.image = imageURL
		b.byID[ingredientID] = ingName
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := make([]dto.MatchResponse, 0, len(recipes))
	for recipeID, b := range recipes {
		total := len(b.byID)
		if total == 0 {
			continue
		}
		matched := 0
		missing := make([]string, 0)
		for id, name := range b.byID {
			if _, ok := userIngredientIDs[id]; ok {
				matched++
			} else {
				missing = append(missing, name)
			}
		}
		sort.Strings(missing)

		pct := float64(matched) / float64(total)
		out = append(out, dto.MatchResponse{
			RecipeID:           recipeID,
			Title:              b.title,
			Image:              b.image,
			MatchPercent:       pct,
			MissingIngredients: missing,
		})
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].MatchPercent != out[j].MatchPercent {
			return out[i].MatchPercent > out[j].MatchPercent
		}
		return out[i].RecipeID < out[j].RecipeID
	})

	return out, nil
}

func normalizeIngredientNames(in []string) []string {
	seen := make(map[string]struct{})
	var out []string
	for _, s := range in {
		s = strings.ToLower(strings.TrimSpace(s))
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
