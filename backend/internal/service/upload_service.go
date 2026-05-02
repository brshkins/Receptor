package service

import (
	"context"

	"receptor/backend/internal/dto"
)

type uploadService struct {
	ml    MLClient
	match MatchService
}

func NewUploadService(ml MLClient, match MatchService) UploadService {
	return &uploadService{ml: ml, match: match}
}

func (s *uploadService) UploadAndMatch(ctx context.Context, image []byte) ([]dto.MatchResponse, error) {
	ingredients, err := s.ml.DetectIngredients(ctx, image)
	if err != nil {
		return nil, err
	}

	return s.match.Match(ctx, &dto.MatchRequest{Ingredients: ingredients})
}

