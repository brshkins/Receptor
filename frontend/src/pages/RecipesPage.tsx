// src/pages/RecipesPage.tsx
import React, { useEffect, useState, useCallback } from 'react';
import { RecipeList } from '../components/Recipes/RecipeList';
import { SearchBar } from '../components/Common/SearchBar';
import { FilterPanel } from '../components/Recipes/FilterPanel';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import { recipesService } from '../services/recipesService';
import { Recipe, RecipeFilters as FilterType } from '../types';
import styles from './Pages.module.css';

const RecipesPage: React.FC = () => {
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [filters, setFilters] = useState<FilterType>({
    sort: 'asc',
    sortBy: 'name'
  });

  const loadRecipes = useCallback(async () => {
    try {
      setLoading(true);
      
      // Используем правильный метод из сервиса
      const data = await recipesService.getAll({
        search: searchTerm || undefined,
        sort: filters.sort,
        limit: 50 // Добавляем лимит для производительности
      });
      
      // Применяем фильтрацию на клиенте (если бэкенд не поддерживает)
      let filteredRecipes = data;
      
      // Фильтрация по сложности (если бэкенд не фильтрует)
      if (filters.difficulty) {
        filteredRecipes = filteredRecipes.filter(
          recipe => recipe.difficulty === filters.difficulty
        );
      }
      
      // Сортировка (если бэкенд не сортирует)
      if (filters.sortBy === 'time') {
        filteredRecipes.sort((a, b) => {
          return filters.sort === 'asc' 
            ? a.cookingTime - b.cookingTime 
            : b.cookingTime - a.cookingTime;
        });
      } else if (filters.sortBy === 'difficulty') {
        const difficultyOrder = { easy: 1, medium: 2, hard: 3 };
        filteredRecipes.sort((a, b) => {
          const diffA = difficultyOrder[a.difficulty as keyof typeof difficultyOrder] || 0;
          const diffB = difficultyOrder[b.difficulty as keyof typeof difficultyOrder] || 0;
          return filters.sort === 'asc' ? diffA - diffB : diffB - diffA;
        });
      }
      
      setRecipes(filteredRecipes);
    } catch (error) {
      console.error('Failed to load recipes:', error);
      // Показываем заглушку при ошибке
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

  const handleFavoriteToggle = async (id: string) => {
    // Оптимистичное обновление UI
    setRecipes(prevRecipes => 
      prevRecipes.map(recipe => 
        recipe.id === id 
          ? { ...recipe, isFavorite: !recipe.isFavorite }
          : recipe
      )
    );
    
    // TODO: Добавить API call для сохранения в избранное
    try {
      // await favoritesService.toggle(id);
    } catch (error) {
      // Откатываем при ошибке
      setRecipes(prevRecipes => 
        prevRecipes.map(recipe => 
          recipe.id === id 
            ? { ...recipe, isFavorite: !recipe.isFavorite }
            : recipe
        )
      );
    }
  };

  // Проверяем, активны ли фильтры
  const isFilterActive = filters.difficulty !== undefined || 
                         filters.sort !== 'asc' || 
                         filters.sortBy !== 'name';

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