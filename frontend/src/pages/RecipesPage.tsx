import React, { useEffect, useState, useCallback } from 'react';
import { RecipeList } from '../components/Recipes/RecipeList';
import { SearchBar } from '../components/Common/SearchBar';
import { RecipeFilters } from '../components/Recipes/RecipeFilters';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import { recipesService } from '../services/recipesService';
import { Recipe, RecipeFilters as FilterType } from '../types';
import styles from './Pages.module.css';

const RecipesPage: React.FC = () => {
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [filters, setFilters] = useState<FilterType>({
    sort: 'asc',
    sortBy: 'name'
  });

  const loadRecipes = useCallback(async () => {
    try {
      setLoading(true);
      const data = await recipesService.getRecipes({
        search: searchTerm || undefined,
        ...filters
      });
      setRecipes(data.recipes);
    } catch (error) {
      console.error('Failed to load recipes:', error);
    } finally {
      setLoading(false);
    }
  }, [searchTerm, filters]);

  useEffect(() => {
    loadRecipes();
  }, [loadRecipes]);

  const handleSearch = (value: string) => {
    setSearchTerm(value);
  };

  const handleFilterChange = (newFilters: FilterType) => {
    setFilters(newFilters);
  };

  const handleFavoriteToggle = (id: string) => {
    // Обновляем локальное состояние
    setRecipes(prevRecipes => 
      prevRecipes.map(recipe => 
        recipe.id === id 
          ? { ...recipe, isFavorite: !recipe.isFavorite }
          : recipe
      )
    );
  };

  if (loading && recipes.length === 0) {
    return <LoadingSpinner text="Загружаем рецепты..." />;
  }

  return (
    <div className={styles.pageContainer}>
      <div className={styles.pageHeader}>
        <h1 className={styles.pageTitle}>📖 Все рецепты</h1>
        <p className={styles.pageDescription}>
          Найди идеальный рецепт из нашей коллекции
        </p>
      </div>
      
      <SearchBar 
        onSearch={handleSearch} 
        placeholder="Поиск по названию или ингредиентам..."
        initialValue={searchTerm}
      />
      
      <RecipeFilters 
        onFilterChange={handleFilterChange}
        initialFilters={filters}
      />
      
      {loading ? (
        <LoadingSpinner text="Обновляем..." />
      ) : (
        <>
          {recipes.length > 0 ? (
            <>
              <p style={{ marginBottom: '16px', color: '#8B6F5E' }}>
                Найдено рецептов: {recipes.length}
              </p>
              <RecipeList 
                recipes={recipes} 
                onFavoriteToggle={handleFavoriteToggle}
              />
            </>
          ) : (
            <div className={styles.emptyState}>
              <div className={styles.emptyIcon}>🔍</div>
              <h3 className={styles.emptyTitle}>Рецепты не найдены</h3>
              <p className={styles.emptyText}>
                Попробуйте изменить параметры поиска или фильтры
              </p>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default RecipesPage;