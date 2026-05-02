package service

import (
	"context"
	"log"
	"strings"

	"receptor/backend/internal/dto"
)

type uploadService struct {
	ml    MLClient
	match MatchService
}

func NewUploadService(ml MLClient, match MatchService) UploadService {
	return &uploadService{ml: ml, match: match}
}

func (s *uploadService) UploadAndMatch(ctx context.Context, images [][]byte) ([]dto.MatchResponse, error) {
	if len(images) == 0 {
		return []dto.MatchResponse{}, nil
	}
	var ingredients []string
	for _, image := range images {
		detected, err := s.ml.DetectIngredients(ctx, image)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, detected...)
		log.Println("ML INGREDIENTS:", ingredients)
	}
	ingredients = demoteAppleToPastaWhenGrocerySet(ingredients)
	return s.match.Match(ctx, &dto.MatchRequest{Ingredients: ingredients})
}

// Несколько фото: оранжевые макароны часто дают apple при верных potato/cucumber и без instant_noodle.
func demoteAppleToPastaWhenGrocerySet(ingredients []string) []string {
	lower := make(map[string]struct{}, len(ingredients))
	for _, raw := range ingredients {
		lower[strings.ToLower(strings.TrimSpace(raw))] = struct{}{}
	}
	_, apple := lower["apple"]
	_, noodle := lower["instant_noodle"]
	_, potato := lower["potato"]
	_, cucumber := lower["cucumber"]
	if !apple || noodle || !potato || !cucumber {
		return ingredients
	}
	out := make([]string, 0, len(ingredients))
	for _, raw := range ingredients {
		if strings.ToLower(strings.TrimSpace(raw)) == "apple" {
			out = append(out, "instant_noodle")
			continue
		}
		out = append(out, raw)
	}
	return out
}

