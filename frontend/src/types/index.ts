// User types
export interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
  createdAt?: string;
}

// Ingredient types
export interface Ingredient {
  id?: string;
  name: string;
  amount?: string;
  unit?: string;
}

// Recipe types
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
  isFavorite?: boolean;
}

// API types
export interface ApiError {
  message: string;
  statusCode: number;
  errors?: Record<string, string[]>;
}

export interface PaginationParams {
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}

// Auth types
export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  name: string;
  email: string;
  password: string;
}

/** POST /auth/register, /auth/login — payload после unwrap */
export interface AuthTokenDto {
  token: string;
}

/** GET /auth/me — payload после unwrap */
export interface MeDto {
  id: number;
  email: string;
  name: string;
}

/** GET /recipes, GET /recipes/:id, GET /favorites — элемент списка */
export interface RecipeListItemDto {
  id: number;
  title: string;
  image: string;
  cooking_time: number;
  category: string;
}

/** GET /recipes/:id/details */
export interface RecipeDetailsDto extends RecipeListItemDto {
  description: string;
  steps: string[];
  /** Имена ингредиентов из backend (`ingredients`); при отсутствии ключа — пусто в UI */
  ingredients?: string[];
}

/** POST /match/by-ingredients, POST /upload — элемент результата (строго по backend). */
export type MatchResponse = {
  recipe_id: number;
  title: string;
  image: string;
  match_percent: number;
  ingredients?: string[];
  missing_ingredients: string[];
};

/** @deprecated используйте AuthTokenDto + MeDto */
export interface AuthResponse {
  token: string;
  user: User;
}

// Component props types
export interface RecipeCardProps {
  id: string;
  title: string;
  imageUrl?: string;
  cookingTime?: number;
  difficulty?: string;
  isFavorite?: boolean;
  onFavoriteToggle?: (id: string) => void;
}

// Filter types
export interface RecipeFilters {
  search?: string;
  ingredients?: string[];
  difficulty?: 'easy' | 'medium' | 'hard';
  category?: string;  // ← ДОБАВЬ ЭТУ СТРОКУ
  sort?: 'asc' | 'desc';
  sortBy?: 'name' | 'time' | 'difficulty' | 'category';  // ← Добавь 'category'
  limit?: number;
  page?: number;
}