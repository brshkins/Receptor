import { apiClient, unwrap } from './api';
import type { MatchResponse } from '../types';
import { normalizeIngredient } from '../utils/matchIngredientNormalize';

export const matchService = {
  async matchByIngredients(ingredients: string[]): Promise<MatchResponse[]> {
    console.log('INGREDIENTS BEFORE:', ingredients);
    const normalized = ingredients
      .flatMap((raw) => {
        const v = normalizeIngredient(raw);
        if (v === '') return [];
        return Array.isArray(v) ? v : [v];
      })
      .filter(Boolean);
    console.log('INGREDIENTS AFTER:', normalized);

    const response = await apiClient.rawPost('/match/by-ingredients', {
      ingredients: normalized,
    });
    const data = unwrap(response);
    return Array.isArray(data) ? (data as MatchResponse[]) : [];
  },
};
