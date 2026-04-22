export interface User {
  id: number;
  email: string;
  name: string;
}

export interface Recipe {
  id: number;
  title: string;
  image: string;
  cooking_time: number;
  category: string;
  isFavorite?: boolean;
}

export interface RecipeFilters {
  search?: string;
  category?: string;
  max_time?: number;
  sort?: 'alphabet_asc' | 'alphabet_desc' | 'time_asc' | 'time_desc';
}