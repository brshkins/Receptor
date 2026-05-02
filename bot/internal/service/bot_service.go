package service

import (
	"context"
	"fmt"
	"strings"

	"receptor/bot/internal/client"
	"receptor/bot/internal/dto"
)

type BotService interface {
	Register(ctx context.Context, email, password, name string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	Me(ctx context.Context, token string) (*dto.User, error)

	GetRecipes(ctx context.Context) ([]dto.RecipeResponse, error)
	// GetRecipesPaged загружает полный список с API и отдаёт срез для отображения в Telegram (бэкенд без offset/limit).
	GetRecipesPaged(ctx context.Context, page, limit int) (recipes []dto.RecipeResponse, hasNext bool, err error)
	GetRecipeByID(ctx context.Context, id int64) (*dto.RecipeResponse, error)
	GetRecipeDetails(ctx context.Context, id int64) (*dto.RecipeDetailsResponse, error)

	GetFavorites(ctx context.Context, token string) ([]dto.RecipeResponse, error)
	AddFavorite(ctx context.Context, token string, recipeID int64) error
	RemoveFavorite(ctx context.Context, token string, recipeID int64) error

	Match(ctx context.Context, ingredients []string) ([]dto.MatchResponse, error)
	UploadImage(ctx context.Context, image []byte) ([]dto.MatchResponse, error)
}

type botService struct {
	backend client.BackendClient
}

func NewBotService(backend client.BackendClient) (BotService, error) {
	if backend == nil {
		return nil, fmt.Errorf("backend client is nil")
	}
	return &botService{backend: backend}, nil
}

func (s *botService) Register(ctx context.Context, email, password, name string) (string, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)
	name = strings.TrimSpace(name)
	if email == "" || password == "" || name == "" {
		return "", fmt.Errorf("email, password and name are required")
	}
	return s.backend.Register(ctx, email, password, name)
}

func (s *botService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)
	if email == "" || password == "" {
		return "", fmt.Errorf("email and password are required")
	}
	return s.backend.Login(ctx, email, password)
}

func (s *botService) Me(ctx context.Context, token string) (*dto.User, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	return s.backend.Me(ctx, token)
}

func (s *botService) GetRecipes(ctx context.Context) ([]dto.RecipeResponse, error) {
	return s.backend.GetRecipes(ctx)
}

func (s *botService) GetRecipesPaged(ctx context.Context, page, limit int) ([]dto.RecipeResponse, bool, error) {
	all, err := s.backend.GetRecipes(ctx)
	if err != nil {
		return nil, false, err
	}
	if page < 0 {
		page = 0
	}
	if limit <= 0 {
		limit = 5
	}
	start := page * limit
	if start >= len(all) {
		return []dto.RecipeResponse{}, false, nil
	}
	end := start + limit
	if end > len(all) {
		end = len(all)
	}
	slice := append([]dto.RecipeResponse(nil), all[start:end]...)
	hasNext := end < len(all)
	return slice, hasNext, nil
}

func (s *botService) GetRecipeByID(ctx context.Context, id int64) (*dto.RecipeResponse, error) {
	return s.backend.GetRecipeByID(ctx, id)
}

func (s *botService) GetRecipeDetails(ctx context.Context, id int64) (*dto.RecipeDetailsResponse, error) {
	return s.backend.GetRecipeDetails(ctx, id)
}

func (s *botService) GetFavorites(ctx context.Context, token string) ([]dto.RecipeResponse, error) {
	return s.backend.GetFavorites(ctx, token)
}

func (s *botService) AddFavorite(ctx context.Context, token string, recipeID int64) error {
	return s.backend.AddFavorite(ctx, token, recipeID)
}

func (s *botService) RemoveFavorite(ctx context.Context, token string, recipeID int64) error {
	return s.backend.RemoveFavorite(ctx, token, recipeID)
}

func (s *botService) Match(ctx context.Context, ingredients []string) ([]dto.MatchResponse, error) {
	return s.backend.Match(ctx, ingredients)
}

func (s *botService) UploadImage(ctx context.Context, image []byte) ([]dto.MatchResponse, error) {
	return s.backend.UploadImage(ctx, image)
}

