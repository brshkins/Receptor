// src/pages/RecipesPage.tsx
import React, { useEffect, useState, useCallback } from 'react';
import { RecipeList } from '../components/Recipes/RecipeList';
import { SearchBar } from '../components/Common/SearchBar';
import { FilterPanel } from '../components/Recipes/FilterPanel';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import { recipesService } from '../services/recipesService';
import { Recipe, RecipeFilters as FilterType } from '../types';
import { favoritesService } from '../services/favoritesService';
import styles from './Pages.module.css';

const RecipesPage: React.FC = () => {
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [filters, setFilters] = useState<FilterType>({
    sort: 'alphabet_asc'
  });

  const loadRecipes = useCallback(async () => {
    try {
      setLoading(true);
      
      const data = await recipesService.getAll({
        search: searchTerm || undefined,
        category: filters.category,
        max_time: filters.max_time,
        sort: filters.sort,
      });

      // If authenticated, overlay favorite state from GET /favorites.
      const token = localStorage.getItem('token');
      if (token) {
        try {
          const favs = await favoritesService.list();
          const favSet =
            Array.isArray(favs) && (favs.length === 0 || typeof favs[0] === 'number')
              ? new Set<number>(favs as number[])
              : new Set<number>((favs as Array<{ id: number }>).map((f) => f.id));
          setRecipes(data.map((r) => ({ ...r, isFavorite: favSet.has(r.id) })));
        } catch {
          setRecipes(data);
        }
      } else {
        setRecipes(data);
      }
    } catch (error) {
      console.error('Failed to load recipes:', error);
      setRecipes([]);
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

  const handleFavoriteToggle = async (id: number) => {
    try {
      const currentlyFavorite = recipes.find((r) => r.id === id)?.isFavorite;
      if (currentlyFavorite) {
        await favoritesService.remove(id);
      } else {
        await favoritesService.add(id);
      }
      setRecipes((prev) => prev.map((r) => (r.id === id ? { ...r, isFavorite: !r.isFavorite } : r)));
    } catch {
      // no-op (UI stays as-is)
    }
  };

  const isFilterActive =
    !!filters.category || typeof filters.max_time === 'number' || (filters.sort ?? 'alphabet_asc') !== 'alphabet_asc';

  if (loading && recipes.length === 0) {
    return <LoadingSpinner text="Загружаем рецепты..." />;
  }

  return (
  <div className={styles.pageContainer}>
    <div className={styles.pageHeader}>
      <h1 className={styles.pageTitle}>
        <span className={styles.pageTitleEmoji}>📖</span>
        <span className={styles.pageTitleText}>Все рецепты</span>
      </h1>
      <p className={styles.pageDescription}>
        Найди идеальный рецепт из нашей коллекции
      </p>
    </div>
      
      <div className={styles.searchSection}>
        <SearchBar 
          onSearch={handleSearch} 
          onFilterClick={() => setIsFilterOpen(true)}
          placeholder="Поиск по названию..."
          initialValue={searchTerm}
          isFilterActive={isFilterActive}
        />
      </div>
      
      <FilterPanel
        filters={filters}
        onFilterChange={handleFilterChange}
        isOpen={isFilterOpen}
        onClose={() => setIsFilterOpen(false)}
      />
      
      {loading ? (
        <LoadingSpinner text="Обновляем..." />
      ) : (
        <>
          {recipes.length > 0 ? (
            <>
              <p className={styles.resultsCount}>
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