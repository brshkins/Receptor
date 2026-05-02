// src/services/favoritesService.ts
import { apiClient } from './api';
import type { RecipeListItemDto } from '../types';

export const favoritesService = {
  async getAll(): Promise<RecipeListItemDto[]> {
    const list = await apiClient.get<RecipeListItemDto[]>('/favorites');
    return Array.isArray(list) ? list : [];
  },

  async add(recipeId: string | number): Promise<void> {
    await apiClient.post<null>(`/favorites/${recipeId}`);
  },

  async remove(recipeId: string | number): Promise<void> {
    await apiClient.delete<null>(`/favorites/${recipeId}`);
  },
};
