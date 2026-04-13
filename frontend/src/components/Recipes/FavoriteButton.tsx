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

  const handleClick = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (isLoading) return;
    
    setIsLoading(true);
    try {
      // API call будет здесь
      setIsFavorite(!isFavorite);
      onToggle?.(recipeId);
    } catch (error) {
      console.error('Failed to toggle favorite:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <button
      className={`${className} ${isFavorite ? styles.active : ''}`}
      onClick={handleClick}
      disabled={isLoading}
      aria-label={isFavorite ? 'Remove from favorites' : 'Add to favorites'}
    >
      {isFavorite ? '❤️' : '🤍'}
    </button>
  );
};