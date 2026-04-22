import { apiClient } from './api';

export interface FavoriteRecipe {
  id: number;
  title: string;
  image: string;
  cooking_time: number;
  category: string;
}

export type FavoritesListResponse = FavoriteRecipe[] | number[];

export const favoritesService = {
  async list(): Promise<FavoritesListResponse> {
    return apiClient.get<FavoritesListResponse>('/favorites');
  },

  async add(recipeId: number): Promise<null> {
    return apiClient.post<null>(`/favorites/${recipeId}`);
  },

  async remove(recipeId: number): Promise<null> {
    return apiClient.delete<null>(`/favorites/${recipeId}`);
  },
};
