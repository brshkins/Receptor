import { apiClient } from './api';

export interface Recipe {
  id: number;
  title: string;
  image: string;
  cooking_time: number;
  category: string;
}

export interface RecipeFilters {
  search?: string;
  category?: string;
  max_time?: number;
  sort?: 'alphabet_asc' | 'alphabet_desc' | 'time_asc' | 'time_desc';
}

export const recipesService = {
  async getAll(params?: RecipeFilters): Promise<Recipe[]> {
    const queryParams = new URLSearchParams();
    if (params?.search) queryParams.append('search', params.search);
    if (params?.category) queryParams.append('category', params.category);
    if (typeof params?.max_time === 'number') queryParams.append('max_time', String(params.max_time));
    if (params?.sort) queryParams.append('sort', params.sort);

    const query = queryParams.toString();
    return apiClient.get<Recipe[]>(`/recipes${query ? `?${query}` : ''}`);
  },

  async getById(id: string): Promise<Recipe> {
    return apiClient.get<Recipe>(`/recipes/${id}`);
  },
};