// src/pages/FavoritesPage.tsx
import React, { useEffect, useState } from 'react';
import { RecipeList } from '../components/Recipes/RecipeList';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import { favoritesService } from '../services/favoritesService';
import { Recipe } from '../types';
import styles from './Pages.module.css';

const FavoritesPage: React.FC = () => {
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadFavorites();
  }, []);

  const loadFavorites = async () => {
    try {
      setLoading(true);
      const data = await favoritesService.getAll();
      console.log('Избранное с бэкенда:', data);
      
      const favorites = Array.isArray(data) ? data : (data as any)?.data || [];
      console.log('Массив избранного:', favorites);
      
      const recipesWithFavorite = favorites.map((r: any) => ({
        ...r,
        isFavorite: true
      }));
      
      setRecipes(recipesWithFavorite);
    } catch (error) {
      console.error('Ошибка загрузки избранного:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleFavoriteToggle = async (id: string) => {
    setRecipes(prev => prev.filter(r => String(r.id) !== id));
    try {
      await favoritesService.remove(id);
    } catch (error) {
      console.error('Ошибка удаления:', error);
      loadFavorites();
    }
  };

  if (loading) {
    return <LoadingSpinner text="Загружаем избранное..." />;
  }

  return (
    <div className={styles.pageContainer}>
      <div className={styles.pageHeader}>
        <h1 className={styles.pageTitle}>
          <span className={styles.pageTitleEmoji}>❤️</span>
          <span className={styles.pageTitleText}>Избранное</span>
        </h1>
        <p className={styles.pageDescription}>
          Ваши любимые рецепты
        </p>
      </div>

      {recipes.length > 0 ? (
        <RecipeList recipes={recipes} onFavoriteToggle={handleFavoriteToggle} />
      ) : (
        <div className={styles.emptyFavorites}>
          <div className={styles.emptyIcon}>❤️</div>
          <h3 className={styles.emptyTitle}>Пока пусто</h3>
          <p className={styles.emptyText}>
            Добавляйте рецепты в избранное, чтобы они всегда были под рукой
          </p>
        </div>
      )}
    </div>
  );
};

export default FavoritesPage;