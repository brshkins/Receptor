package service

import (
	"context"

	"receptor/backend/internal/dto"
	"receptor/backend/internal/repository"
)

type favoriteService struct {
	favorites repository.FavoritesRepository
}

func NewFavoriteService(favorites repository.FavoritesRepository) FavoriteService {
	return &favoriteService{favorites: favorites}
}

func (s *favoriteService) Add(ctx context.Context, userID int64, recipeID int64) error {
	return s.favorites.Add(ctx, userID, recipeID)
}

func (s *favoriteService) Remove(ctx context.Context, userID int64, recipeID int64) error {
	return s.favorites.Remove(ctx, userID, recipeID)
}

func (s *favoriteService) GetAll(ctx context.Context, userID int64) ([]dto.RecipeResponse, error) {
	list, err := s.favorites.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.RecipeResponse, 0, len(list))
	for _, r := range list {
		out = append(out, recipeModelToResponse(r))
	}
	return out, nil
}
