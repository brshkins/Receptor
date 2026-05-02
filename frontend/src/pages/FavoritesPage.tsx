// src/pages/FavoritesPage.tsx
import React, { useEffect, useState } from 'react';
import { RecipeList } from '../components/Recipes/RecipeList';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import { favoritesService } from '../services/favoritesService';
import { Recipe } from '../types';
import { getApiErrorMessage } from '../utils/apiError';
import toast from 'react-hot-toast';
import styles from './Pages.module.css';

const FavoritesPage: React.FC = () => {
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);

  useEffect(() => {
    loadFavorites();
  }, []);

  const loadFavorites = async () => {
    try {
      setLoading(true);
      setLoadError(null);
      const favorites = await favoritesService.getAll();
      const recipesWithFavorite = favorites.map((r) => ({
        ...r,
        isFavorite: true,
      }));

      setRecipes(recipesWithFavorite);
    } catch (error) {
      console.error('Ошибка загрузки избранного:', error);
      setLoadError(getApiErrorMessage(error, 'Не удалось загрузить избранное'));
      setRecipes([]);
    } finally {
      setLoading(false);
    }
  };

  const handleFavoriteToggle = async (id: string) => {
    const previous = recipes;
    const removed = previous.find((r) => String(r.id) === id);
    if (!removed) return;

    setRecipes((prev) => prev.filter((r) => String(r.id) !== id));
    try {
      await favoritesService.remove(id);
    } catch (error) {
      console.error('Ошибка удаления:', error);
      setRecipes(previous);
      toast.error(getApiErrorMessage(error, 'Не удалось удалить из избранного'));
    }
  };

  if (loading) {
    return <LoadingSpinner text="Загружаем избранное..." />;
  }

  if (loadError) {
    return (
      <div className={styles.pageContainer}>
        <div className={styles.pageHeader}>
          <h1 className={styles.pageTitle}>
            <span className={styles.pageTitleEmoji}>❤️</span>
            <span className={styles.pageTitleText}>Избранное</span>
          </h1>
        </div>
        <div className={styles.emptyState} role="alert">
          <div className={styles.emptyIcon}>⚠️</div>
          <h3 className={styles.emptyTitle}>Не удалось загрузить</h3>
          <p className={styles.emptyText}>{loadError}</p>
          <button type="button" className={styles.backButtonBig} onClick={() => loadFavorites()}>
            Повторить
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.pageContainer}>
      <div className={styles.pageHeader}>
        <h1 className={styles.pageTitle}>
          <span className={styles.pageTitleEmoji}>❤️</span>
          <span className={styles.pageTitleText}>Избранное</span>
        </h1>
        <p className={styles.pageDescription}>Ваши любимые рецепты</p>
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
