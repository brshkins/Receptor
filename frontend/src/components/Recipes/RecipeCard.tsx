import React from 'react';
import { Link } from 'react-router-dom';
import { FavoriteButton } from './FavoriteButton';
import styles from './Recipes.module.css';

interface RecipeCardProps {
  id: string;
  title: string;
  imageUrl?: string;
  cookingTime?: number;
  difficulty?: string;
  isFavorite?: boolean;
  onFavoriteToggle?: (id: string) => void;
}

export const RecipeCard: React.FC<RecipeCardProps> = ({
  id,
  title,
  imageUrl,
  cookingTime,
  difficulty,
  isFavorite = false,
  onFavoriteToggle
}) => {
  return (
    <div className={styles.recipeCard}>
      <Link to={`/recipes/${id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
        {imageUrl ? (
          <img src={imageUrl} alt={title} className={styles.recipeImage} />
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
            {cookingTime && <span>⏱️ {cookingTime} мин</span>}
            {difficulty && <span>📊 {difficulty}</span>}
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