import { apiClient } from './api';

export interface UploadResponse {
  success: boolean;
  fileUrl?: string;
  matches?: Array<{
    id: string;
    title: string;
    confidence: number;
  }>;
  message?: string;
}

export const uploadService = {
  async uploadAndMatch(file: File): Promise<UploadResponse> {
    const formData = new FormData();
    formData.append('image', file);

    return apiClient.post<UploadResponse>('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },

  async uploadImage(file: File): Promise<{ url: string }> {
    const formData = new FormData();
    formData.append('image', file);

    return apiClient.post<{ url: string }>('/upload/image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
};