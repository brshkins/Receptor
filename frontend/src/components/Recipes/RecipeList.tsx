// src/components/Recipes/RecipeList.tsx
import React from 'react';
import { Recipe } from '../../types';
import { RecipeCard } from './RecipeCard';
import styles from './Recipes.module.css';

interface RecipeListProps {
  recipes: Recipe[];
  onFavoriteToggle?: (id: string) => void;
}

export const RecipeList: React.FC<RecipeListProps> = ({ recipes, onFavoriteToggle }) => {
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
          description={recipe.description}
          image={recipe.image}
          imageUrl={recipe.imageUrl}
          cooking_time={recipe.cooking_time}   // ← ПРОВЕРЬ ЭТУ СТРОКУ
          cookingTime={recipe.cookingTime}
          category={recipe.category}           // ← ПРОВЕРЬ ЭТУ СТРОКУ
          difficulty={recipe.difficulty}
          isFavorite={recipe.isFavorite}
          onFavoriteToggle={onFavoriteToggle}
        />
      ))}
    </div>
  );
};