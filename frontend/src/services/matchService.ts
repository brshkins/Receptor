import { apiClient } from './api';

export interface MatchResponseItem {
  recipe_id: number;
  title: string;
  image: string;
  match_percent: number;
  missing_ingredients: string[];
}

export const matchService = {
  async matchByIngredients(ingredients: string[]): Promise<MatchResponseItem[]> {
    return apiClient.post<MatchResponseItem[]>('/match/by-ingredients', { ingredients });
  }
};