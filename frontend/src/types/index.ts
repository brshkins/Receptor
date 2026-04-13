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
  id: string;
  title: string;
  description: string;
  ingredients: Ingredient[] | string[];
  instructions: string[];
  cookingTime: number;
  difficulty: 'easy' | 'medium' | 'hard';
  servings?: number;
  imageUrl?: string;
  isFavorite?: boolean;
  createdAt?: string;
  authorId?: string;
  author?: User;
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
  sort?: 'asc' | 'desc';
  sortBy?: 'name' | 'time' | 'difficulty';
  limit?: number;
  page?: number;
}