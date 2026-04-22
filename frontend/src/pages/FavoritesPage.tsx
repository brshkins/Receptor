// src/pages/FavoritesPage.tsx
import React, { useCallback, useEffect, useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { favoritesService } from '../services/favoritesService';
import { recipesService } from '../services/recipesService';
import { RecipeList } from '../components/Recipes/RecipeList';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import { Recipe } from '../types';
import styles from './Pages.module.css';

const FavoritesPage: React.FC = () => {
  const { isLoading: authLoading, isAuthenticated } = useAuth();
  const [loading, setLoading] = useState(true);
  const [favorites, setFavorites] = useState<Recipe[]>([]);

  const load = useCallback(async () => {
    if (!isAuthenticated) {
      setFavorites([]);
      setLoading(false);
      return;
    }
    try {
      setLoading(true);
      const data = await favoritesService.list();
      if (Array.isArray(data) && (data.length === 0 || typeof data[0] === 'number')) {
        const ids = data as number[];
        const recipes = await Promise.all(ids.map((id) => recipesService.getById(String(id))));
        setFavorites(recipes.map((r) => ({ ...r, isFavorite: true })));
      } else {
        const recipes = data as Array<{
          id: number;
          title: string;
          image: string;
          cooking_time: number;
          category: string;
        }>;
        setFavorites(recipes.map((r) => ({ ...r, isFavorite: true })));
      }
    } finally {
      setLoading(false);
    }
  }, [isAuthenticated]);

  useEffect(() => {
    if (authLoading) return;
    load();
  }, [authLoading, load]);

  const handleFavoriteToggle = async (id: number) => {
    try {
      await favoritesService.remove(id);
      setFavorites((prev) => prev.filter((r) => r.id !== id));
    } catch (e) {
      // ignore
    }
  };

  if (authLoading) return <LoadingSpinner text="Загружаем..." />;
  if (!isAuthenticated) {
    return (
      <div className={styles.pageContainer}>
        <div className={styles.pageHeader}>
          <h1 className={styles.pageTitle}>
            <span className={styles.pageTitleEmoji}>❤️</span>
            <span className={styles.pageTitleText}>Избранное</span>
          </h1>
          <p className={styles.pageDescription}>Войдите, чтобы увидеть избранные рецепты</p>
        </div>
      </div>
    );
  }

  if (loading) return <LoadingSpinner text="Загружаем избранное..." />;

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

      {favorites.length > 0 ? (
        <RecipeList recipes={favorites} onFavoriteToggle={handleFavoriteToggle} />
      ) : (
        <div className={styles.emptyFavorites}>
          <div className={styles.emptyIcon}>📖</div>
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