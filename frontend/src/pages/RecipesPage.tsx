// src/pages/RecipesPage.tsx
import React, { useEffect, useState, useCallback } from 'react';
import { RecipeList } from '../components/Recipes/RecipeList';
import { SearchBar } from '../components/Common/SearchBar';
import { FilterPanel } from '../components/Recipes/FilterPanel';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import { recipesService } from '../services/recipesService';
import { favoritesService } from '../services/favoritesService';
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
      
      const response = await recipesService.getAll({
        search: searchTerm || undefined,
      });
      
      console.log('Ответ бэкенда:', response);
      
      let recipesList: Recipe[] = [];
      
      if (Array.isArray(response)) {
        recipesList = response;
      } else if (response && typeof response === 'object') {
        if ('data' in response && Array.isArray((response as any).data)) {
          recipesList = (response as any).data;
        } else if ('recipes' in response && Array.isArray((response as any).recipes)) {
          recipesList = (response as any).recipes;
        } else {
          const values = Object.values(response);
          recipesList = values.filter((item): item is Recipe => {
            return item !== null && 
                   typeof item === 'object' && 
                   'id' in item && 
                   'title' in item;
          });
        }
      }
      
      console.log('Массив рецептов:', recipesList);
      
      let filteredRecipes: Recipe[] = [...recipesList];
      
      // Фильтрация по сложности
      if (filters.difficulty) {
        filteredRecipes = filteredRecipes.filter(recipe => {
          const recipeDiff = recipe.difficulty || recipe.category || '';
          return recipeDiff === filters.difficulty;
        });
      }
      
      // Сортировка
      if (filters.sortBy === 'time') {
        filteredRecipes.sort((a, b) => {
          const timeA = a.cookingTime ?? a.cooking_time ?? 0;
          const timeB = b.cookingTime ?? b.cooking_time ?? 0;
          return filters.sort === 'asc' ? timeA - timeB : timeB - timeA;
        });
      } else if (filters.sortBy === 'difficulty') {
        const difficultyOrder: Record<string, number> = { 
          easy: 1, medium: 2, hard: 3 
        };
        filteredRecipes.sort((a, b) => {
          const diffA = a.difficulty || a.category || '';
          const diffB = b.difficulty || b.category || '';
          const valueA = difficultyOrder[diffA] ?? 0;
          const valueB = difficultyOrder[diffB] ?? 0;
          return filters.sort === 'asc' ? valueA - valueB : valueB - valueA;
        });
      } else if (filters.sortBy === 'name') {
        filteredRecipes.sort((a, b) => {
          const titleA = a.title || '';
          const titleB = b.title || '';
          return filters.sort === 'asc' 
            ? titleA.localeCompare(titleB)
            : titleB.localeCompare(titleA);
        });
      }
      
      setRecipes(filteredRecipes);
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

  const handleFavoriteToggle = async (id: string) => {
    const recipe = recipes.find(r => String(r.id) === id);
    if (!recipe) return;
    
    // Оптимистичное обновление UI
    setRecipes(prevRecipes => 
      prevRecipes.map(r => 
        String(r.id) === id 
          ? { ...r, isFavorite: !r.isFavorite }
          : r
      )
    );
    
    // Отправка на бэкенд
    try {
      if (recipe.isFavorite) {
        await favoritesService.remove(id);
      } else {
        await favoritesService.add(id);
      }
    } catch (error) {
      console.error('Ошибка избранного:', error);
      // Откатываем при ошибке
      setRecipes(prevRecipes => 
        prevRecipes.map(r => 
          String(r.id) === id 
            ? { ...r, isFavorite: !r.isFavorite }
            : r
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