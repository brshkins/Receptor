import { apiClient } from './api';

export interface Recipe {
  id: string;
  title: string;
  description: string;
  ingredients: any[];
  instructions: string[];
  cookingTime: number;
  difficulty: 'easy' | 'medium' | 'hard';
  servings?: number;
  imageUrl?: string;
  isFavorite?: boolean;
}

export interface RecipesResponse {
  recipes: Recipe[];
  total: number;
  page: number;
  totalPages: number;
}

export interface RecipeFilters {
  search?: string;
  ingredients?: string[];
  difficulty?: 'easy' | 'medium' | 'hard';
  sort?: 'asc' | 'desc';
  sortBy?: 'name' | 'time' | 'difficulty';
  limit?: number;
  page?: number;
}

export const recipesService = {
  // Временно убираем sort и limit — бэкенд их не поддерживает
  async getAll(params?: {
    search?: string;
  }): Promise<Recipe[]> {
    const queryParams = new URLSearchParams();
    if (params?.search) queryParams.append('search', params.search);
    
    const query = queryParams.toString();
    return apiClient.get<Recipe[]>(`/recipes${query ? `?${query}` : ''}`);
  },

  async getById(id: string): Promise<Recipe> {
    return apiClient.get<Recipe>(`/recipes/${id}`);
  },

  async getRecipes(filters?: RecipeFilters): Promise<RecipesResponse> {
    const queryParams = new URLSearchParams();
    
    if (filters?.search) queryParams.append('search', filters.search);
    if (filters?.ingredients?.length) {
      queryParams.append('ingredients', filters.ingredients.join(','));
    }
    if (filters?.difficulty) queryParams.append('difficulty', filters.difficulty);
    // Убираем sort и limit — они не поддерживаются бэкендом
    // if (filters?.sort) queryParams.append('sort', filters.sort);
    // if (filters?.sortBy) queryParams.append('sortBy', filters.sortBy);
    // if (filters?.limit) queryParams.append('limit', String(filters.limit));
    if (filters?.page) queryParams.append('page', String(filters.page));

    const query = queryParams.toString();
    return apiClient.get<RecipesResponse>(`/recipes${query ? `?${query}` : ''}`);
  },
};