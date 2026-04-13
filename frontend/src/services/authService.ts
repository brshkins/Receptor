import { apiClient } from './api';

export interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
}

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

export const authService = {
  async register(data: RegisterData): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/register', data);
    if (response.token) {
      localStorage.setItem('token', response.token);
    }
    return response;
  },

  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/login', credentials);
    if (response.token) {
      localStorage.setItem('token', response.token);
    }
    return response;
  },

  async getMe(): Promise<User> {
    return apiClient.get<User>('/auth/me');
  },

  logout(): void {
    localStorage.removeItem('token');
    window.location.href = '/auth';
  },

  isAuthenticated(): boolean {
    return !!localStorage.getItem('token');
  },

  getToken(): string | null {
    return localStorage.getItem('token');
  }
};