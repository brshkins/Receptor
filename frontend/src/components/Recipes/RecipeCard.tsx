import React from 'react';
import { Link } from 'react-router-dom';
import { FavoriteButton } from './FavoriteButton';
import styles from './Recipes.module.css';

interface RecipeCardProps {
  id: number;
  title: string;
  image?: string;
  cooking_time?: number;
  category?: string;
  isFavorite?: boolean;
  onFavoriteToggle?: (id: number) => void;
}

export const RecipeCard: React.FC<RecipeCardProps> = ({
  id,
  title,
  image,
  cooking_time,
  category,
  isFavorite = false,
  onFavoriteToggle
}) => {
  return (
    <div className={styles.recipeCard}>
      <Link to={`/recipes/${id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
        {image ? (
          <img src={image} alt={title} className={styles.recipeImage} />
        ) : (
          <div className={styles.recipeImage} style={{ 
            display: 'flex', 
            alignItems: 'center', 
            justifyContent: 'center',
            fontSize: '48px'
          }}>
            🍳
          </div>
        )}
        <div className={styles.recipeContent}>
          <h3 className={styles.recipeTitle}>{title}</h3>
          <div className={styles.recipeMeta}>
            {typeof cooking_time === 'number' && <span>⏱️ {cooking_time} мин</span>}
            {category && <span>🏷️ {category}</span>}
          </div>
        </div>
      </Link>
      <FavoriteButton 
        recipeId={id}
        isFavorite={isFavorite}
        onToggle={onFavoriteToggle}
        className={styles.favoriteButton}
      />
    </div>
  );
};