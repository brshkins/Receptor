import React, { useState } from 'react';
import { favoritesService } from '../../services/favoritesService';
import { useAuth } from '../../contexts/AuthContext';
import styles from './Recipes.module.css';

interface FavoriteButtonProps {
  recipeId: number;
  isFavorite: boolean;
  onToggle?: (id: number) => void;
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
  const { isAuthenticated } = useAuth();

  const handleClick = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (isLoading) return;
    if (!isAuthenticated) {
      window.location.href = '/auth';
      return;
    }
    
    setIsLoading(true);
    try {
      if (isFavorite) {
        await favoritesService.remove(recipeId);
        setIsFavorite(false);
      } else {
        await favoritesService.add(recipeId);
        setIsFavorite(true);
      }
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