package service

import (
	"context"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/model"
	"receptor/backend/internal/repository"
)

type recipeService struct {
	recipes repository.RecipesRepository
}

func NewRecipeService(recipes repository.RecipesRepository) RecipeService {
	return &recipeService{recipes: recipes}
}

func (s *recipeService) GetAll(ctx context.Context, filter dto.RecipeFilter) ([]dto.RecipeResponse, error) {
	list, err := s.recipes.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}
	out := make([]dto.RecipeResponse, 0, len(list))
	for _, r := range list {
		out = append(out, recipeModelToResponse(r))
	}
	return out, nil
}

func (s *recipeService) GetByID(ctx context.Context, id int64) (*dto.RecipeResponse, error) {
	r, err := s.recipes.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := recipeModelToResponse(r)
	return &resp, nil
}

func (s *recipeService) GetDetails(ctx context.Context, id int64) (*dto.RecipeDetailsResponse, error) {
	r, err := s.recipes.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	steps, err := s.recipes.GetStepsByRecipeID(ctx, id)
	if err != nil {
		return nil, err
	}
	outSteps := make([]string, 0, len(steps))
	for _, st := range steps {
		outSteps = append(outSteps, st.Description)
	}
	return &dto.RecipeDetailsResponse{
		ID:          r.ID,
		Title:       r.Title,
		Image:       r.ImageURL,
		CookingTime: r.CookingTime,
		Category:    r.Category,
		Description: r.Description,
		Steps:       outSteps,
	}, nil
}

func recipeModelToResponse(r *model.Recipe) dto.RecipeResponse {
	return dto.RecipeResponse{
		ID:          r.ID,
		Title:       r.Title,
		Image:       r.ImageURL,
		CookingTime: r.CookingTime,
		Category:    r.Category,
	}
}
