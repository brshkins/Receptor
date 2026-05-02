import { apiClient } from './api';

export interface Recipe {
  id: string;
  title: string;
  description: string;
  ingredients: any[];
  instructions: string[];
  cookingTime: number;
  difficulty: string;
  imageUrl?: string;
}

export interface MatchResult {
  recipe: Recipe;
  matchScore: number;
  matchedIngredients: string[];
}

export interface MatchResponse {
  matches: MatchResult[];
  totalCount: number;
}

export const matchService = {
  async matchByIngredients(ingredients: string[]): Promise<MatchResponse> {
    return apiClient.post<MatchResponse>('/match/by-ingredients', { ingredients });
  },

  async matchByImage(imageFile: File): Promise<MatchResponse> {
    const formData = new FormData();
    formData.append('image', imageFile);
    
    return apiClient.post<MatchResponse>('/match/by-image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
};