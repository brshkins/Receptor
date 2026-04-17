import React from 'react';
import { Recipe } from '../../types';
import { RecipeCard } from './RecipeCard';
import styles from './Recipes.module.css';

interface RecipeListProps {
  recipes: Recipe[];
  onFavoriteToggle?: (id: string) => void;
}

export const RecipeList: React.FC<RecipeListProps> = ({ 
  recipes, 
  onFavoriteToggle 
}) => {
  if (!recipes || recipes.length === 0) {
    return (
      <div className={styles.emptyState}>
        <p>Рецепты не найдены</p>
      </div>
    );
  }

  return (
    <div className={styles.recipeGrid}>
      {recipes.map((recipe) => (
        <RecipeCard
          key={recipe.id}
          id={recipe.id}
          title={recipe.title}
          imageUrl={recipe.imageUrl}
          cookingTime={recipe.cookingTime}
          difficulty={recipe.difficulty}
          isFavorite={recipe.isFavorite}
          onFavoriteToggle={onFavoriteToggle}
        />
      ))}
    </div>
  );
};