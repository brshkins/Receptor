// src/services/authService.ts
import { apiClient } from './api';

export interface User {
  id: number | string;
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
    const response = await apiClient.post<any>('/auth/register', data);
    
    // Бэкенд возвращает { data: { token: "..." } }
    const token = response?.data?.token || response?.token;
    const user = response?.data?.user || response?.user;
    
    if (token) {
      localStorage.setItem('token', token);
    }
    
    return { token, user };
  },

  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const response = await apiClient.post<any>('/auth/login', credentials);
    
    // Бэкенд возвращает { data: { token: "..." } }
    const token = response?.data?.token || response?.token;
    const user = response?.data?.user || response?.user;
    
    if (token) {
      localStorage.setItem('token', token);
    }
    
    return { token, user };
  },

  async getMe(): Promise<User> {
    const response = await apiClient.get<any>('/auth/me');
    // Бэкенд может возвращать { data: { ... } }
    return response?.data || response;
  },

  logout(): void {
    localStorage.removeItem('token');
  },

  isAuthenticated(): boolean {
    return !!localStorage.getItem('token');
  },

  getToken(): string | null {
    return localStorage.getItem('token');
  }
};