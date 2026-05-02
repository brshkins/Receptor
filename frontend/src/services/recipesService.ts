// src/services/recipesService.ts
import { apiClient } from './api';
import type { RecipeListItemDto, RecipeDetailsDto } from '../types';

export interface RecipeFilters {
  search?: string;
  category?: string;
  max_time?: number;
  /** backend: alphabet_asc | alphabet_desc | time_asc | time_desc */
  sort?: string;
}

export const recipesService = {
  async getAll(params?: { search?: string }): Promise<RecipeListItemDto[]> {
    const queryParams = new URLSearchParams();
    if (params?.search) queryParams.append('search', params.search);
    const query = queryParams.toString();
    const list = await apiClient.get<RecipeListItemDto[]>(`/recipes${query ? `?${query}` : ''}`);
    return Array.isArray(list) ? list : [];
  },

  async getById(id: string | number): Promise<RecipeListItemDto> {
    return apiClient.get<RecipeListItemDto>(`/recipes/${id}`);
  },

  /** GET /recipes/:id/details — описание, шаги, ингредиенты */
  async getRecipeDetails(id: string | number): Promise<RecipeDetailsDto> {
    return apiClient.get<RecipeDetailsDto>(`/recipes/${id}/details`);
  },

  async getRecipes(filters?: RecipeFilters): Promise<RecipeListItemDto[]> {
    const queryParams = new URLSearchParams();
    if (filters?.search) queryParams.append('search', filters.search);
    if (filters?.category) queryParams.append('category', filters.category);
    if (filters?.max_time !== undefined) queryParams.append('max_time', String(filters.max_time));
    if (filters?.sort) queryParams.append('sort', filters.sort);
    const query = queryParams.toString();
    const list = await apiClient.get<RecipeListItemDto[]>(`/recipes${query ? `?${query}` : ''}`);
    return Array.isArray(list) ? list : [];
  },
};
