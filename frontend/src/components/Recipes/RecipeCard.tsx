import React from 'react';
import { Link } from 'react-router-dom';
import { FavoriteButton } from './FavoriteButton';
import styles from './Recipes.module.css';

interface RecipeCardProps {
  id: number | string;
  title: string;
  description?: string;
  image?: string;        // ← бэкенд возвращает image
  imageUrl?: string;     // ← оставляем для совместимости
  cooking_time?: number; // ← бэкенд возвращает cooking_time
  cookingTime?: number;  // ← оставляем для совместимости
  category?: string;     // ← вместо difficulty
  difficulty?: string;
  isFavorite?: boolean;
  onFavoriteToggle?: (id: string) => void;
}

export const RecipeCard: React.FC<RecipeCardProps> = ({
  id,
  title,
  description,
  image,
  imageUrl,
  cooking_time,
  cookingTime,
  category,
  difficulty,
  isFavorite = false,
  onFavoriteToggle
}) => {
  // Используем то, что пришло
  const img = image || imageUrl;
  const time = cooking_time || cookingTime;
  const diff = difficulty || category;
  
  const getDifficultyText = (diff?: string) => {
    if (!diff) return '';
    const map: Record<string, string> = {
      'easy': '🥚 Легко',
      'medium': '👨‍🍳 Средне',
      'hard': '🔥 Сложно',
      'pasta': '🍝 Паста',
      'soup': '🍲 Суп',
      'salad': '🥗 Салат',
      'dessert': '🍰 Десерт'
    };
    return map[diff] || diff;
  };

  return (
    <div className={styles.recipeCard}>
      <Link to={`/recipes/${id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
        {img ? (
          <img src={img} alt={title} className={styles.recipeImage} />
        ) : (
          <div className={styles.recipeImage} style={{ 
            display: 'flex', 
            alignItems: 'center', 
            justifyContent: 'center',
            fontSize: '48px'
          }}>
            {getCategoryEmoji(category)}
          </div>
        )}
        <div className={styles.recipeContent}>
          <h3 className={styles.recipeTitle}>{title}</h3>
          
          <div className={styles.recipeMeta}>
            {time && <span>⏱️ {time} мин</span>}
            {diff && <span>📊 {getDifficultyText(diff)}</span>}
          </div>
        </div>
      </Link>
      <FavoriteButton 
        recipeId={String(id)}
        isFavorite={isFavorite}
        onToggle={onFavoriteToggle}
        className={styles.favoriteButton}
      />
    </div>
  );
};

// Эмодзи по категории
const getCategoryEmoji = (category?: string): string => {
  const map: Record<string, string> = {
    'pasta': '🍝',
    'soup': '🍲',
    'salad': '🥗',
    'dessert': '🍰',
    'meat': '🥩',
    'fish': '🐟',
    'breakfast': '🍳',
    'bakery': '🥐'
  };
  return map[category || ''] || '🍳';
};