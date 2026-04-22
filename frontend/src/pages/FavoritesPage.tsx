import React, { useEffect, useState } from 'react';
import { favoritesService } from '../services/favoritesService';
import { RecipeList } from '../components/Recipes/RecipeList';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
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
      // Помечаем все как избранные
      const favorites = (Array.isArray(data) ? data : []).map(r => ({
        ...r,
        isFavorite: true
      }));
      setRecipes(favorites);
    } catch (error) {
      console.error('Ошибка загрузки избранного:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleFavoriteToggle = async (id: string) => {
    // Удаляем из списка при снятии лайка
    setRecipes(prev => prev.filter(r => String(r.id) !== id));
    
    try {
      await favoritesService.remove(id);
    } catch (error) {
      console.error('Ошибка удаления из избранного:', error);
      // Возвращаем при ошибке
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
        <>
          <p className={styles.resultsCount}>
            В избранном: {recipes.length} рецептов
          </p>
          <RecipeList 
            recipes={recipes} 
            onFavoriteToggle={handleFavoriteToggle}
          />
        </>
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