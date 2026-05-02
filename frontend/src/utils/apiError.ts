import { AxiosError } from 'axios';

/** Типовые тексты ошибок backend → RU (тело `{ "error": "..." }` без смены API). */
const API_ERROR_RU: Record<string, string> = {
  'invalid request body': 'Неверный запрос',
  'invalid credentials': 'Неверный email или пароль',
  unauthorized: 'Требуется вход',
  'not found': 'Не найдено',
  'invalid id': 'Неверный идентификатор',
  'email already exists': 'Этот email уже зарегистрирован',
  conflict: 'Операция невозможна: ресурс уже существует',
  'internal server error': 'Ошибка сервера',
  'invalid max_time': 'Некорректное время приготовления',
  'invalid sort': 'Некорректная сортировка',
  'file is required': 'Нужно выбрать файл',
  'file too large': 'Файл слишком большой',
  'invalid file': 'Не удалось прочитать файл',
  'file must be an image': 'Файл должен быть изображением',
};

function translateKnownMessage(msg: string): string {
  const key = msg.trim().toLowerCase();
  if (API_ERROR_RU[key]) return API_ERROR_RU[key];
  if (key === 'network error' || key.includes('network error')) {
    return 'Нет соединения с сервером';
  }
  return msg;
}

/** Сообщение из тела ответа backend `{ "error": "..." }` или fallback (всё на русском в UI). */
export function getApiErrorMessage(error: unknown, fallback: string): string {
  if (error instanceof AxiosError) {
    const data = error.response?.data as { error?: string } | undefined;
    if (data && typeof data.error === 'string' && data.error.trim()) {
      return translateKnownMessage(data.error);
    }
    if (error.message) {
      return translateKnownMessage(error.message);
    }
  }
  if (error instanceof Error && error.message) {
    return translateKnownMessage(error.message);
  }
  return fallback;
}
