// src/utils/translations.ts

import { recipeTitleEnToRuLower } from './recipeTitleRu.generated';

/** Категории: ключ как в API (нижний регистр). categoryMap[key] ?? исходное значение */
export const categoryMap: Record<string, string> = {
  dinner: 'Ужин',
  breakfast: 'Завтрак',
  lunch: 'Обед',
  pasta: 'Паста',
  meat: 'Мясо',
  vegetarian: 'Вегетарианское',
  dessert: 'Десерт',
  soup: 'Суп',
  salad: 'Салат',
};

/** Короткие названия (slug) */
export const titleMap: Record<string, string> = {
  pasta: 'Паста',
  sandwich: 'Сэндвич',
};

/** Частые EN-слова в названиях → RU (если нет в recipeTitleEnToRuLower), по аналогии с ботом */
const recipeTitleWordEnToRu: Record<string, string> = {
  chicken: 'курица',
  beef: 'говядина',
  pork: 'свинина',
  turkey: 'индейка',
  soup: 'суп',
  pasta: 'паста',
  salad: 'салат',
  toast: 'тост',
  bread: 'хлеб',
  tomato: 'помидор',
  potato: 'картофель',
  rice: 'рис',
  egg: 'яйцо',
  eggs: 'яйца',
  cheese: 'сыр',
  milk: 'молоко',
  butter: 'масло',
  garlic: 'чеснок',
  onion: 'лук',
  fish: 'рыба',
  tuna: 'тунец',
  salmon: 'лосось',
  shrimp: 'креветки',
  with: 'с',
  and: 'и',
  style: '',
  bowl: 'боул',
  berry: 'ягоды',
  berries: 'ягоды',
  cream: 'сливки',
  honey: 'мёд',
  apple: 'яблоко',
  banana: 'банан',
  noodle: 'лапша',
  noodles: 'лапша',
  chili: 'чили',
  bean: 'фасоль',
  beans: 'фасоль',
  spaghetti: 'спагетти',
  penne: 'пенне',
  farfalle: 'фарфалле',
  basil: 'базилик',
  omelette: 'омлет',
  pancakes: 'блины',
  cake: 'торт',
};

function recipeTitleWordFallback(title: string): string {
  const words = title.trim().split(/\s+/);
  const out: string[] = [];
  for (const raw of words) {
    const w = raw.replace(/^[^a-zA-Z0-9а-яА-ЯёЁ]+|[^a-zA-Z0-9а-яА-ЯёЁ]+$/g, '');
    if (!w) continue;
    const key = w.toLowerCase();
    const ru = recipeTitleWordEnToRu[key];
    if (ru) {
      if (ru) out.push(ru);
      continue;
    }
    const allLatin = /^[a-z]+$/i.test(w);
    if (allLatin) continue;
    out.push(w.charAt(0).toUpperCase() + w.slice(1).toLowerCase());
  }
  if (out.length === 0) return 'Блюдо';
  return out.join(' ');
}

/** Название рецепта для пользователя: RU из общего словаря с ботом (ключ — lower). */
export const translateRecipeTitle = (title?: string): string => {
  if (!title) return '';
  const t = title.trim();
  if (!t) return '';
  // Уже по-русски — не трогаем. Не используем «любой не-ASCII»: é в «sauté» ломал бы lookup.
  if (/[а-яА-ЯёЁ]/.test(t)) return t;
  const lower = t.toLowerCase();
  const fromSeed = recipeTitleEnToRuLower[lower];
  if (fromSeed) return fromSeed;
  const slug = titleMap[lower];
  if (slug) return slug;
  return recipeTitleWordFallback(t);
};

// ============================================
// ЭМОДЗИ ПО КАТЕГОРИЯМ
// ============================================
export const getCategoryEmoji = (category?: string): string => {
  const emojiMap: Record<string, string> = {
    pasta: '🍝',
    meat: '🥩',
    vegetarian: '🥬',
    breakfast: '🍳',
    dessert: '🍰',
    soup: '🍲',
    salad: '🥗',
    dinner: '🌙',
    lunch: '🍴',
  };
  return emojiMap[category?.toLowerCase() || ''] || '🍽️';
};

// ============================================
// СЛОВАРЬ ЕДИНИЦ ИЗМЕРЕНИЯ
// ============================================
export const unitTranslations: Record<string, string> = {
  'g': 'г',
  'kg': 'кг',
  'ml': 'мл',
  'l': 'л',
  'tbsp': 'ст.л.',
  'tsp': 'ч.л.',
  'cup': 'чашка',
  'cups': 'чашки',
  'piece': 'шт.',
  'pieces': 'шт.',
  'clove': 'зубчик',
  'cloves': 'зубчика',
  'pinch': 'щепотка',
  'to taste': 'по вкусу',
};

// ============================================
// ФУНКЦИИ ПЕРЕВОДА
// ============================================

// Перевод категории
export const translateCategory = (category?: string): string => {
  if (!category) return '';
  return categoryMap[category.toLowerCase()] ?? category;
};

// Перевод единицы измерения
export const translateUnit = (unit?: string): string => {
  if (!unit) return '';
  return unitTranslations[unit.toLowerCase()] || unit;
};

// Перевод сложности
export const translateDifficulty = (difficulty?: string): string => {
  if (!difficulty) return '';
  const map: Record<string, string> = {
    'easy': '🥚 Легко',
    'medium': '👨‍🍳 Средне',
    'hard': '🔥 Сложно',
  };
  return map[difficulty.toLowerCase()] || difficulty;
};

// Форматирование времени
export const formatTime = (minutes?: number): string => {
  if (!minutes) return '';
  if (minutes < 60) return `${minutes} мин`;
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  return mins > 0 ? `${hours} ч ${mins} мин` : `${hours} ч`;
};