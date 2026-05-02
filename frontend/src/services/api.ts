// src/services/api.ts
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

/** Backend wraps success as `{ data: T }`; axios puts that in `response.data`. */
export const unwrap = (res: AxiosResponse) => res.data?.data ?? res.data;

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
      withCredentials: true,
    });

    this.client.interceptors.request.use((config) => {
      // Иначе дефолтный application/json ломает POST multipart (поле file).
      if (config.data instanceof FormData && config.headers) {
        config.headers.delete('Content-Type');
      }
      const token = localStorage.getItem('token');
      if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    this.client.interceptors.response.use(
      (response: AxiosResponse) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('token');
          if (!window.location.pathname.includes('/auth')) {
            window.location.href = '/auth';
          }
        }
        return Promise.reject(error);
      }
    );
  }

  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.get(url, config);
    return unwrap(response) as T;
  }

  async post<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.post(url, data, config);
    return unwrap(response) as T;
  }

  /** Сырой POST (для сервисов с явным unwrap телом ответа). */
  rawPost(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<AxiosResponse> {
    return this.client.post(url, data, config);
  }

  async put<T>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.put(url, data, config);
    return unwrap(response) as T;
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.delete(url, config);
    return unwrap(response) as T;
  }
}

export const apiClient = new ApiClient();
export default apiClient;