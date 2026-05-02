/**
 * Карта ввода (RU/EN) → токены для поиска в БД.
 * Значение: одна строка или несколько (общие категории → конкретные имена).
 */
type IngredientMapValue = string | string[];

const ingredientInputToDBName: Record<string, IngredientMapValue> = {
  chicken: 'meat',
  beef: 'meat',
  pork: 'meat',
  turkey: 'meat',
  egg: 'eggs',
  eggs: 'eggs',
  tuna: 'fish',
  salmon: 'fish',
  курица: 'meat',
  куриное: 'meat',
  куриный: 'meat',
  говядина: 'meat',
  свинина: 'meat',
  индейка: 'meat',
  мясо: 'meat',
  рыба: 'fish',
  лосось: 'fish',
  тунец: 'fish',
  яйцо: 'eggs',
  яйца: 'eggs',
  яиц: 'eggs',
  креветки: 'shrimp',
  креветка: 'shrimp',
  shrimp: 'shrimp',
  помидор: 'tomato',
  помидоры: 'tomato',
  томат: 'tomato',
  томаты: 'tomato',
  tomato: 'tomato',
  tomatoes: 'tomato',
  сыр: ['cheese', 'mozzarella', 'cheddar', 'parmesan'],
  cheese: ['cheese', 'mozzarella', 'cheddar', 'parmesan'],
  молоко: 'milk',
  milk: 'milk',
  мука: 'flour',
  flour: 'flour',
  сливочное: 'butter',
  'сливочное масло': 'butter',
  butter: 'butter',
  масло: 'olive oil',
  'оливковое масло': 'olive oil',
  оливковое: 'olive oil',
  'olive oil': 'olive oil',
  oil: 'olive oil',
  паста: 'pasta',
  pasta: 'pasta',
  спагетти: 'pasta',
  spaghetti: 'pasta',
  чеснок: 'garlic',
  garlic: 'garlic',
  лук: 'onion',
  onion: 'onion',
  onions: 'onion',
  картофель: 'potato',
  картошка: 'potato',
  potato: 'potato',
  potatoes: 'potato',
  рис: 'rice',
  rice: 'rice',
  бекон: 'meat',
  bacon: 'meat',
  соль: 'salt',
  salt: 'salt',
  перец: 'pepper',
  pepper: 'pepper',
  peppers: 'pepper',
  сахар: 'sugar',
  sugar: 'sugar',
  лимон: 'lemon',
  lemon: 'lemon',
  базилик: 'basil',
  basil: 'basil',
  петрушка: 'parsley',
  parsley: 'parsley',
  шпинат: 'spinach',
  spinach: 'spinach',
  морковь: 'carrot',
  carrot: 'carrot',
  carrots: 'carrot',
  огурец: 'cucumber',
  огурцы: 'cucumber',
  cucumber: 'cucumber',
  cucumbers: 'cucumber',
  гриб: 'mushroom',
  грибы: 'mushroom',
  mushroom: 'mushroom',
  mushrooms: 'mushroom',
  орегано: 'oregano',
  oregano: 'oregano',
  чили: 'chili',
  'перец чили': 'chili',
  chili: 'chili',
  йогурт: 'yogurt',
  yogurt: 'yogurt',
  сметана: 'cream',
  cream: 'cream',
  банан: 'banana',
  banana: 'banana',
  яблоко: 'apple',
  яблоки: 'apple',
  apple: 'apple',
  apples: 'apple',
  ягоды: 'berry',
  berry: 'berry',
  berries: 'berry',
  фасоль: 'beans',
  beans: 'beans',
  вода: 'water',
  water: 'water',
};

/**
 * trim, lower, ё→е; по словарю — одна строка или список токенов.
 */
export function normalizeIngredient(raw: string): string | string[] {
  const s = raw.trim().toLowerCase().replace(/ё/g, 'е');
  if (!s) return '';
  const canon = ingredientInputToDBName[s];
  if (canon === undefined) return s;
  return canon;
}

/** Токены для одной строки ввода (для flatMap). */
export function normalizeIngredientTokens(raw: string): string[] {
  const v = normalizeIngredient(raw);
  if (v === '') return [];
  return Array.isArray(v) ? v : [v];
}

export function normalizeMatchIngredientNames(input: string[]): string[] {
  const seen = new Set<string>();
  const out: string[] = [];
  for (const raw of input) {
    for (const t of normalizeIngredientTokens(raw)) {
      if (!t) continue;
      if (seen.has(t)) continue;
      seen.add(t);
      out.push(t);
    }
  }
  return out;
}
