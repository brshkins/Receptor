import { apiClient, unwrap } from './api';
import type { MatchResponse } from '../types';

export async function normalizeImage(input: Blob | File): Promise<File> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    const url = URL.createObjectURL(input);

    img.onload = async () => {
      URL.revokeObjectURL(url);

      const canvas = document.createElement('canvas');
      canvas.width = Math.max(img.naturalWidth || img.width, 256);
      canvas.height = Math.max(img.naturalHeight || img.height, 256);

      const ctx = canvas.getContext('2d');
      if (!ctx) return reject(new Error('Canvas error'));

      ctx.drawImage(img, 0, 0, canvas.width, canvas.height);

      try {
        const dataUrl = canvas.toDataURL('image/jpeg', 0.9);
        const blob = await fetch(dataUrl).then((r) => r.blob());

        if (!blob || blob.size === 0) {
          return reject(new Error('Empty blob after normalization'));
        }

        const file = new File([blob], 'photo.jpg', {
          type: 'image/jpeg',
        });

        resolve(file);
      } catch (e) {
        reject(e);
      }
    };

    img.onerror = () => {
      URL.revokeObjectURL(url);
      reject(new Error('Invalid image'));
    };

    img.src = url;
  });
}

export const uploadService = {
  normalizeImage,
  /** Одно или несколько фото — бэкенд склеивает список ингредиентов и один раз матчит рецепты. */
  async uploadAndMatch(input: File | Blob | Array<File | Blob>): Promise<MatchResponse[]> {
    const list = Array.isArray(input) ? input : [input];
    const normalized = await Promise.all(list.map((x) => normalizeImage(x)));

    const formData = new FormData();
    // Одно имя поля `file` несколько раз — так стабильнее парсится в Go (multipart.File["file"]).
    for (const f of normalized) {
      formData.append('file', f, f.name);
    }

    const response = await apiClient.rawPost('/upload', formData);
    const data = unwrap(response);
    return Array.isArray(data) ? (data as MatchResponse[]) : [];
  },
};
