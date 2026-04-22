// src/services/favoritesService.ts
import { apiClient } from './api';
import { Recipe } from '../types';

export const favoritesService = {
  async getAll(): Promise<Recipe[]> {
    return apiClient.get<Recipe[]>('/favorites');
  },

  async add(recipeId: string | number): Promise<void> {
    return apiClient.post<void>(`/favorites/${recipeId}`);
  },

  async remove(recipeId: string | number): Promise<void> {
    return apiClient.delete<void>(`/favorites/${recipeId}`);
  },
};