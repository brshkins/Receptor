/**
 * Отображение описаний, шагов и ингредиентов на русском.
 * В API/БД тексты на английском — здесь словари и эвристики для UI.
 */

/** Имена ингредиентов как в БД (lower) → подпись для пользователя. Синхронизировано с bot/internal/handler/recipe_text_ru.go; в seed egg → eggs и т.п. */
export const ingredientNameEnToRu: Record<string, string> = {
  apple: 'Яблоко',
  apples: 'Яблоки',
  banana: 'Банан',
  basil: 'Базилик',
  beans: 'Фасоль',
  beef: 'Говядина',
  berries: 'Ягоды',
  berry: 'Ягоды',
  bread: 'Хлеб',
  broccoli: 'Брокколи',
  butter: 'Масло сливочное',
  cabbage: 'Капуста',
  carrot: 'Морковь',
  carrots: 'Морковь',
  cheddar: 'Чеддер',
  cheese: 'Сыр',
  chicken: 'Курица',
  chili: 'Чили',
  cinnamon: 'Корица',
  cream: 'Сливки',
  cucumber: 'Огурец',
  egg: 'Яйцо',
  eggs: 'Яйца',
  fish: 'Рыба',
  flour: 'Мука',
  garlic: 'Чеснок',
  honey: 'Мёд',
  instant_noodle: 'Лапша быстрого приготовления',
  'instant noodle': 'Лапша быстрого приготовления',
  lemon: 'Лимон',
  meat: 'Мясо',
  milk: 'Молоко',
  mozzarella: 'Моцарелла',
  mushroom: 'Грибы',
  mushrooms: 'Грибы',
  'olive oil': 'Оливковое масло',
  onion: 'Лук',
  onions: 'Лук',
  oregano: 'Орегано',
  parmesan: 'Пармезан',
  parsley: 'Петрушка',
  pasta: 'Паста',
  pepper: 'Перец',
  pork: 'Свинина',
  potato: 'Картофель',
  potatoes: 'Картофель',
  rice: 'Рис',
  salt: 'Соль',
  shrimp: 'Креветки',
  spinach: 'Шпинат',
  sugar: 'Сахар',
  tomato: 'Помидор',
  tomatoes: 'Помидоры',
  tuna: 'Тунец',
  walnut: 'Грецкие орехи',
  walnuts: 'Грецкие орехи',
  water: 'Вода',
  yogurt: 'Йогурт',
};

/** Фразы в инструкциях (длинные первыми), замена без учёта регистра */
const INSTRUCTION_PHRASES: [string, string][] = [
  [
    'warm apple bake with cinnamon and buttery crumble.',
    'Тёплая яблочная запеканка с корицей и масляным крамблом.',
  ],
  ['until al dente', 'до состояния аль денте'],
  ['until cooked through', 'пока не будет готово'],
  ['until golden and bubbling', 'пока не зарумянится и не забулькает'],
  ['until golden', 'пока не зарумянится'],
  ['until crisp', 'пока не станет хрустящим'],
  ['until bubbling', 'пока не забулькает'],
  ['until tender', 'пока не станет мягким'],
  ['preheat oven to', 'Разогрейте духовку до'],
  ['bring to a boil', 'Доведите до кипения'],
  ['bring to boil', 'Доведите до кипения'],
  ['simmer for', 'Томите на слабом огне'],
  ['simmer ', 'Томите на слабом огне '],
  ['boil in salted water', 'варите в подсоленной воде'],
  ['boil ', 'Варите '],
  ['bake until', 'Выпекайте, пока'],
  ['bake for', 'Выпекайте'],
  ['bake ', 'Выпекайте '],
  ['serve with', 'Подавайте с'],
  ['serve hot', 'Подавайте горячим'],
  ['season with salt and pepper', 'Посолите и поперчите'],
  ['season with', 'Приправьте'],
  ['heat olive oil', 'Разогрейте оливковое масло'],
  ['heat oil', 'Разогрейте масло'],
  ['heat ', 'Нагрейте '],
  ['add salt', 'Добавьте соль'],
  ['add ', 'Добавьте '],
  ['stir in', 'Вмешайте'],
  ['stir ', 'Помешайте '],
  ['mix well', 'Хорошо перемешайте'],
  ['mix ', 'Смешайте '],
  ['toss with', 'Смешайте с'],
  ['slice ', 'Нарежьте '],
  ['chop ', 'Мелко нарежьте '],
  ['drain ', 'Слейте '],
  ['top with', 'Сверху выложите'],
  ['rub butter into flour', 'Вотрите масло в муку'],
  ['fold in', 'Аккуратно вмешайте'],
  ['cool before serving', 'Остудите перед подачей'],
  ['warm apple bake with cinnamon and buttery crumble', 'Тёплая яблочная запеканка с корицей и масляным крамблом'],
  ['classic spaghetti with tomato, garlic, and basil', 'Классические спагетти с томатом, чесноком и базиликом'],
];

const DESCRIPTION_WORDS: Record<string, string> = {
  warm: 'Тёплая',
  comforting: 'Сытный',
  classic: 'Классический',
  crispy: 'Хрустящий',
  creamy: 'Сливочный',
  moist: 'Нежный',
  simple: 'Простой',
  fresh: 'Свежий',
  hearty: 'Сытный',
  spicy: 'Острый',
  sweet: 'Сладкий',
  savory: 'Сытный',
  light: 'Лёгкий',
  apple: 'яблочная',
  bake: 'запеканка',
  with: 'с',
  and: 'и',
  buttery: 'масляным',
  crumble: 'крамблом',
  cinnamon: 'корицей',
  tomato: 'томатом',
  garlic: 'чесноком',
  basil: 'базиликом',
  chicken: 'курицей',
  soup: 'суп',
  noodles: 'лапшой',
  carrot: 'морковью',
  beans: 'фасолью',
  onion: 'луком',
  parsley: 'петрушкой',
  bacon: 'беконом',
  eggs: 'яйцами',
  honey: 'мёдом',
  soy: 'соевым',
  sauce: 'соусом',
  spaghetti: 'спагетти',
  style: '',
  flavored: 'с ароматом',
  salad: 'салат',
  plate: 'тарелка',
  breakfast: 'завтрак',
  fried: 'жареные',
  glazed: 'в глазури',
  breast: 'грудка',
};

function escapeRegExp(s: string): string {
  return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

/** Подпись ингредиента для карточки рецепта */
export function translateIngredientDisplayName(name: string): string {
  const key = name.trim().toLowerCase();
  if (ingredientNameEnToRu[key]) return ingredientNameEnToRu[key];
  // Запасной вариант: канонические имена в БД часто во мн. ч. (eggs), в словаре есть обе формы;
  // для редких имён — пробуем убрать хвост мн. ч.
  if (key.length > 3 && key.endsWith('ies')) {
    const y = `${key.slice(0, -3)}y`;
    if (ingredientNameEnToRu[y]) return ingredientNameEnToRu[y];
  }
  if (key.length > 3 && key.endsWith('es')) {
    const woEs = key.slice(0, -2);
    if (ingredientNameEnToRu[woEs]) return ingredientNameEnToRu[woEs];
  }
  if (key.length > 2 && key.endsWith('s')) {
    const woS = key.slice(0, -1);
    if (ingredientNameEnToRu[woS]) return ingredientNameEnToRu[woS];
  }
  return name;
}

/**
 * Грубый перевод описания/шага (EN → понятные русские фразы; не дословный).
 * Если уже есть кириллица — возвращаем как есть.
 */
export function translateRecipeEnglishText(text: string): string {
  if (!text?.trim()) return text;
  if (/[а-яА-ЯёЁ]/.test(text) && !/[a-zA-Z]{5,}/.test(text)) return text;

  let s = text
    .replace(/\u00b0C/g, '°C')
    .replace(/190В°C/gi, '190 °C')
    .replace(/(\d+)В°C/gi, '$1 °C')
    .trim();

  for (const [en, ru] of INSTRUCTION_PHRASES) {
    const re = new RegExp(escapeRegExp(en), 'gi');
    s = s.replace(re, ru);
  }

  // Короткие описания: по словам
  if (s.length < 120 && s.split(/\s+/).length < 25 && !/[а-яА-ЯёЁ]/.test(s)) {
    const words = s.split(/(\s+)/);
    const out = words.map((chunk) => {
      if (/^\s+$/.test(chunk)) return chunk;
      const punct = chunk.match(/^([^a-zA-Z]*)([a-zA-Z'-]+)([^a-zA-Z]*)$/i);
      if (!punct) return chunk;
      const [, pre, w, post] = punct;
      const ru = DESCRIPTION_WORDS[w.toLowerCase()];
      if (ru) return `${pre}${ru}${post}`;
      return chunk;
    });
    s = out.join('');
    s = s.replace(/\s+/g, ' ').trim();
    if (/^[a-z]/i.test(s)) {
      s = s.charAt(0).toUpperCase() + s.slice(1);
    }
  }

  return s;
}
