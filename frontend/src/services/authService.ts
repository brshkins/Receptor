// src/services/authService.ts
import { apiClient } from './api';
import type { AuthTokenDto, MeDto } from '../types';

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  name: string;
  email: string;
  password: string;
}

export const authService = {
  async register(data: RegisterData): Promise<AuthTokenDto> {
    const out = await apiClient.post<AuthTokenDto>('/auth/register', data);
    if (out?.token) {
      localStorage.setItem('token', out.token);
    }
    return out;
  },

  async login(credentials: LoginCredentials): Promise<AuthTokenDto> {
    const out = await apiClient.post<AuthTokenDto>('/auth/login', credentials);
    if (out?.token) {
      localStorage.setItem('token', out.token);
    }
    return out;
  },

  async getMe(): Promise<MeDto> {
    return apiClient.get<MeDto>('/auth/me');
  },

  logout(): void {
    localStorage.removeItem('token');
  },

  isAuthenticated(): boolean {
    return !!localStorage.getItem('token');
  },

  getToken(): string | null {
    return localStorage.getItem('token');
  },
};
