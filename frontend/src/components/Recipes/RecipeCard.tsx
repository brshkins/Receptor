// src/components/Recipes/RecipeCard.tsx
import React from 'react';
import { Link } from 'react-router-dom';
import { FavoriteButton } from './FavoriteButton';
import { translateCategory, getCategoryEmoji, translateRecipeTitle, formatTime } from '../../utils/translations';
import styles from './Recipes.module.css';

interface RecipeCardProps {
  id: number | string;
  title: string;
  description?: string;
  image?: string;
  imageUrl?: string;
  cooking_time?: number;
  cookingTime?: number;
  category?: string;
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
  const img = image || imageUrl;
  const time = cooking_time || cookingTime;
  const titleRu = translateRecipeTitle(title);
  const categoryRu = translateCategory(category);
  const emoji = getCategoryEmoji(category);
  const timeFormatted = formatTime(time); // "25 мин" или "1 ч 20 мин"

  return (
    <div className={styles.recipeCard}>
      <Link to={`/recipes/${id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
        {img ? (
          <img src={img} alt={titleRu} className={styles.recipeImage} />
        ) : (
          <div className={styles.recipeImagePlaceholder}>
            <span className={styles.placeholderEmoji}>{emoji}</span>
          </div>
        )}
        <div className={styles.recipeContent}>
          <h3 className={styles.recipeTitle}>{titleRu}</h3>
          
          {/* Время и сложность */}
          <div className={styles.recipeMeta}>
            {time && (
              <span className={styles.metaItem}>
                <span className={styles.metaIcon}>⏱️</span>
                <span>{timeFormatted}</span>
              </span>
            )}
            {categoryRu && (
              <span className={styles.metaItem}>
                <span className={styles.metaIcon}>{emoji}</span>
                <span>{categoryRu}</span>
              </span>
            )}
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