// src/components/Recipes/FavoriteButton.tsx
import React, { useState } from 'react';
import styles from './Recipes.module.css';

interface FavoriteButtonProps {
  recipeId: string;
  isFavorite: boolean;
  onToggle?: (id: string) => void;
  className?: string;
}

export const FavoriteButton: React.FC<FavoriteButtonProps> = ({
  recipeId,
  isFavorite: initialFavorite,
  onToggle,
  className
}) => {
  const [isFavorite, setIsFavorite] = useState(initialFavorite);
  const [isLoading, setIsLoading] = useState(false);

  // Синхронизируем с пропсами
  React.useEffect(() => {
    setIsFavorite(initialFavorite);
  }, [initialFavorite]);

  const handleClick = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (isLoading || !onToggle) return;
    
    setIsLoading(true);
    try {
      await onToggle(recipeId);
      setIsFavorite(!isFavorite);
    } catch (error) {
      console.error('Failed to toggle favorite:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <button
      className={`${className || styles.favoriteButton} ${isFavorite ? styles.favoriteActive : ''}`}
      onClick={handleClick}
      disabled={isLoading}
      aria-label={isFavorite ? 'Удалить из избранного' : 'Добавить в избранное'}
    >
      {isFavorite ? '❤️' : '🤍'}
    </button>
  );
};