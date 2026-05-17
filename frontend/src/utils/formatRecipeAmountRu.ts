/** Количества из recipes.json (EN) → отображение рядом с русским названием продукта. */
export function formatRecipeAmountRu(amount: string): string {
  const s0 = amount.trim();
  if (!s0) return '';
  const low = s0.toLowerCase();
  if (low === 'to taste') return 'по вкусу';
  if (low === 'pinch') return 'щепотка';

  let s = s0;
  const pairs: [RegExp, string][] = [
    // Вкус/щепотка
    [/\bto taste\b/gi, 'по вкусу'],
    [/\bpinch\b/gi, 'щепотка'],
    
    // Столовые/чайные ложки
    [/\btbsp\b/gi, 'ст. л.'],
    [/\btsp\b/gi, 'ч. л.'],
    
    // Штуки/зубчики/ломтики
    [/\bcloves?\b/gi, 'зубч.'],
    [/\bpcs?\b/gi, 'шт.'],
    [/\bpieces?\b/gi, 'шт.'],
    [/\bslices?\b/gi, 'ломт.'],
    
    // Объём
    [/\bml\b/gi, 'мл'],
    [/\bl\b/gi, 'л'],
    
    // Вес
    [/\bg\b/gi, 'г'],
    [/\bkg\b/gi, 'кг'],
    
    // Чашки
    [/\bcups?\b/gi, 'чаш.'],
    [/\bcup\b/gi, 'чаш.'],
  ];
  
  for (const [re, ru] of pairs) {
    s = s.replace(re, ru);
  }
  
  return s.trim();
}