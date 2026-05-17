// src/pages/RecipesPage.tsx
import React, { useEffect, useState, useCallback } from 'react';
import { RecipeList } from '../components/Recipes/RecipeList';
import { SearchBar } from '../components/Common/SearchBar';
import { FilterPanel } from '../components/Recipes/FilterPanel';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import { recipesService } from '../services/recipesService';
import { favoritesService } from '../services/favoritesService';
import { Recipe, RecipeFilters as FilterType } from '../types';
import { translateCategory, translateRecipeTitle } from '../utils/translations';
import styles from './Pages.module.css';

const FILTERS_STORAGE_KEY = 'receptor_recipe_filters';

// Загружаем фильтры из localStorage
const loadSavedFilters = (): FilterType => {
  try {
    const saved = localStorage.getItem(FILTERS_STORAGE_KEY);
    if (saved) return JSON.parse(saved);
  } catch (e) { /* ignore */ }
  return { sort: 'asc', sortBy: 'name' };
};

// Сохраняем фильтры в localStorage
const saveFiltersToStorage = (filters: FilterType) => {
  try {
    localStorage.setItem(FILTERS_STORAGE_KEY, JSON.stringify(filters));
  } catch (e) { /* ignore */ }
};

const RecipesPage: React.FC = () => {
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [filters, setFilters] = useState<FilterType>(loadSavedFilters);
  const [favoriteIds, setFavoriteIds] = useState<Set<string>>(new Set());

  // Сохраняем фильтры при каждом изменении
  useEffect(() => {
    saveFiltersToStorage(filters);
  }, [filters]);

  const loadFavoriteIds = useCallback(async () => {
    try {
      const favorites = await favoritesService.getAll();
      const responseData = favorites as any;
      const data = Array.isArray(responseData) 
        ? responseData 
        : responseData?.data || [];
      const ids = new Set<string>(data.map((f: any) => String(f.id || f.recipe_id)));
      setFavoriteIds(ids);
      return ids;
    } catch (error) {
      console.error('Ошибка загрузки избранного:', error);
    }
    return new Set<string>();
  }, []);

  const loadRecipes = useCallback(async () => {
    try {
      setLoading(true);
      setLoadError(null);
      const recipesList = await recipesService.getAll();

      const favIds = await loadFavoriteIds();
      let filtered: Recipe[] = recipesList.map((r) => ({
        ...r,
        isFavorite: favIds.has(String(r.id)),
      }));

      const q = searchTerm.trim().toLowerCase();
      if (q) {
        filtered = filtered.filter((r) => {
          const titleRu = translateRecipeTitle(r.title).toLowerCase();
          const titleEn = (r.title ?? '').toLowerCase();
          return titleRu.includes(q) || titleEn.includes(q);
        });
      }

      if (filters.category) {
        filtered = filtered.filter(r => r.category === filters.category);
      }

      // Сортировка
      if (filters.sortBy === 'time') {
        filtered.sort((a, b) => {
          const ta = a.cooking_time ?? 0;
          const tb = b.cooking_time ?? 0;
          return filters.sort === 'asc' ? ta - tb : tb - ta;
        });
      } else if (filters.sortBy === 'name') {
        filtered.sort((a, b) => 
          filters.sort === 'asc' 
            ? a.title.localeCompare(b.title) 
            : b.title.localeCompare(a.title)
        );
      } else if (filters.sortBy === 'category') {
        filtered.sort((a, b) => {
          const ca = translateCategory(a.category) || '';
          const cb = translateCategory(b.category) || '';
          return filters.sort === 'asc' ? ca.localeCompare(cb) : cb.localeCompare(ca);
        });
      }

      setRecipes(filtered);
    } catch (error: unknown) {
      console.error('Failed to load recipes:', error);
      setLoadError('Не удалось загрузить рецепты');
      setRecipes([]);
    } finally { setLoading(false); }
  }, [searchTerm, filters, loadFavoriteIds]);

  useEffect(() => { loadRecipes(); }, [loadRecipes]);

  const handleSearch = (value: string) => setSearchTerm(value);
  
  const handleFilterChange = (newFilters: FilterType) => {
    setFilters(newFilters);
  };

  const handleFavoriteToggle = async (id: string) => {
    const recipe = recipes.find(r => String(r.id) === id);
    if (!recipe) return;
    const newState = !recipe.isFavorite;

    setRecipes(prev => prev.map(r => String(r.id) === id ? { ...r, isFavorite: newState } : r));
    setFavoriteIds(prev => { const s = new Set(prev); newState ? s.add(id) : s.delete(id); return s; });

    try {
      if (newState) await favoritesService.add(id);
      else await favoritesService.remove(id);
    } catch (error: any) {
      if (error?.response?.status === 409 && newState) return;
      if (error?.response?.status === 401) alert('Чтобы добавить в избранное, войдите в аккаунт');
      setRecipes(prev => prev.map(r => String(r.id) === id ? { ...r, isFavorite: !newState } : r));
      setFavoriteIds(prev => { const s = new Set(prev); newState ? s.delete(id) : s.add(id); return s; });
    }
  };

  const isFilterActive = filters.category !== undefined || filters.sort !== 'asc' || filters.sortBy !== 'name';

  if (loading && recipes.length === 0 && !loadError) {
    return <LoadingSpinner text="Загружаем рецепты..." />;
  }

  return (
    <div className={styles.pageContainer}>
      <div className={styles.pageHeader}>
        <h1 className={styles.pageTitle}>
          <span className={styles.pageTitleEmoji}>📖</span>
          <span className={styles.pageTitleText}>Все рецепты</span>
        </h1>
        <p className={styles.pageDescription}>Найди идеальный рецепт из нашей коллекции</p>
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

      {loadError ? (
        <div className={styles.emptyState} role="alert">
          <div className={styles.emptyIcon}>⚠️</div>
          <h3 className={styles.emptyTitle}>Не удалось загрузить рецепты</h3>
          <p className={styles.emptyText}>{loadError}</p>
          <button type="button" className={styles.backButtonBig} onClick={() => loadRecipes()}>
            Повторить
          </button>
        </div>
      ) : loading ? (
        <LoadingSpinner text="Обновляем..." />
      ) : recipes.length > 0 ? (
        <>
          <p className={styles.resultsCount}>Найдено рецептов: {recipes.length}</p>
          <RecipeList recipes={recipes} onFavoriteToggle={handleFavoriteToggle} />
        </>
      ) : (
        <div className={styles.emptyState}>
          <div className={styles.emptyIcon}>🔍</div>
          <h3 className={styles.emptyTitle}>Рецепты не найдены</h3>
          <p className={styles.emptyText}>Попробуйте изменить параметры поиска или фильтры</p>
        </div>
      )}
    </div>
  );
};

export default RecipesPage;