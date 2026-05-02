import { apiClient, unwrap } from './api';
import type { MatchResponse } from '../types';

export const matchService = {
  async matchByIngredients(ingredients: string[]): Promise<MatchResponse[]> {
    const payload = ingredients.map((s) => s.trim()).filter(Boolean);

    const response = await apiClient.rawPost('/match/by-ingredients', {
      ingredients: payload,
    });
    const data = unwrap(response);
    return Array.isArray(data) ? (data as MatchResponse[]) : [];
  },
};
