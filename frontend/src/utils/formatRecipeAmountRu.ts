/** Количества из recipes.json (EN) → отображение рядом с русским названием продукта. */
export function formatRecipeAmountRu(amount: string): string {
  const s0 = amount.trim();
  if (!s0) return '';
  const low = s0.toLowerCase();
  if (low === 'to taste') return 'по вкусу';
  if (low === 'pinch') return 'щепотка';

  let s = s0;
  const pairs: [RegExp, string][] = [
    [/\bto taste\b/gi, 'по вкусу'],
    [/\bpinch\b/gi, 'щепотка'],
    [/\btbsp\b/gi, 'ст. л.'],
    [/\btsp\b/gi, 'ч. л.'],
    [/\bcloves?\b/gi, 'зубч.'],
    [/\bpcs?\b/gi, 'шт.'],
    [/\bml\b/gi, 'мл'],
    [/\bg\b/gi, 'г'],
  ];
  for (const [re, ru] of pairs) {
    s = s.replace(re, ru);
  }
  return s.trim();
}
