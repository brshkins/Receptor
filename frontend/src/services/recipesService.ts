// src/services/recipesService.ts
import { apiClient } from './api';
import { translateCategory, translateDifficulty, translateIngredient } from '../utils/translations';

export interface Recipe {
  id: number | string;
  title: string;
  description?: string;
  image?: string;
  imageUrl?: string;
  cooking_time?: number;
  cookingTime?: number;
  category?: string;
  difficulty?: string;
  ingredients?: string[];
  isFavorite?: boolean;
  instructions?: string[];
  servings?: number;
}

export interface RecipesResponse {
  data: Recipe[];
  total?: number;
}

export interface RecipeFilters {
  search?: string;
  ingredients?: string[];
  difficulty?: string;
  category?: string;
  sort?: 'asc' | 'desc';
  sortBy?: 'name' | 'time' | 'difficulty' | 'category';
  limit?: number;
  page?: number;
}

export const recipesService = {
  async getAll(params?: {
    search?: string;
  }): Promise<Recipe[]> {
    const queryParams = new URLSearchParams();
    if (params?.search) queryParams.append('search', params.search);
    
    const query = queryParams.toString();
    const response = await apiClient.get<RecipesResponse>(`/recipes${query ? `?${query}` : ''}`);
    
    // Извлекаем массив рецептов из ответа
    const recipes = response?.data || response;
    return Array.isArray(recipes) ? recipes : [];
  },

  async getById(id: string | number): Promise<any> {
    const response = await apiClient.get<any>(`/recipes/${id}`);
    // Бэкенд может вернуть { data: {...} } или сразу объект
    return response?.data || response;
  },

  async getRecipes(filters?: RecipeFilters): Promise<Recipe[]> {
    const queryParams = new URLSearchParams();
    
    if (filters?.search) queryParams.append('search', filters.search);
    if (filters?.category) queryParams.append('category', filters.category);
    if (filters?.difficulty) queryParams.append('difficulty', filters.difficulty);
    
    const query = queryParams.toString();
    const response = await apiClient.get<RecipesResponse>(`/recipes${query ? `?${query}` : ''}`);
    
    const recipes = response?.data || response;
    return Array.isArray(recipes) ? recipes : [];
  },
};
