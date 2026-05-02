// src/pages/RecipesPage.tsx
import React, { useEffect, useState, useCallback } from 'react';
import { RecipeList } from '../components/Recipes/RecipeList';
import { SearchBar } from '../components/Common/SearchBar';
import { FilterPanel } from '../components/Recipes/FilterPanel';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import { recipesService } from '../services/recipesService';
import { favoritesService } from '../services/favoritesService';
import { Recipe, RecipeFilters as FilterType } from '../types';
import { translateCategory } from '../utils/translations';
import styles from './Pages.module.css';

const RecipesPage: React.FC = () => {
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [filters, setFilters] = useState<FilterType>({ sort: 'asc', sortBy: 'name' });
  const [favoriteIds, setFavoriteIds] = useState<Set<string>>(new Set());

  const loadFavoriteIds = useCallback(async () => {
  try {
    const response = await favoritesService.getAll();
    console.log('Ответ избранного:', response);
    
    // response может быть { data: [...] } или просто массивом
    const favorites = Array.isArray(response) 
      ? response 
      : (response as any)?.data || [];
    
    if (Array.isArray(favorites)) {
      const ids = new Set(favorites.map((f: any) => String(f.id || f.recipe_id)));
      console.log('ID избранных:', Array.from(ids));
      setFavoriteIds(ids);
      return ids;
    }
  } catch (error) {
    console.error('Ошибка загрузки избранного:', error);
  }
  return new Set<string>();
}, []);

  const loadRecipes = useCallback(async () => {
    try {
      setLoading(true);
      const response = await recipesService.getAll({ search: searchTerm || undefined });

      let recipesList: Recipe[] = [];
      if (Array.isArray(response)) recipesList = response;
      else if (response && typeof response === 'object' && 'data' in response && Array.isArray((response as any).data)) {
        recipesList = (response as any).data;
      }

      const favIds = await loadFavoriteIds();
      let filtered = recipesList.map(r => ({ ...r, isFavorite: favIds.has(String(r.id)) }));

      if (filters.category) filtered = filtered.filter(r => r.category === filters.category);

      if (filters.sortBy === 'time') {
        filtered.sort((a, b) => {
          const ta = a.cookingTime ?? a.cooking_time ?? 0;
          const tb = b.cookingTime ?? b.cooking_time ?? 0;
          return filters.sort === 'asc' ? ta - tb : tb - ta;
        });
      } else if (filters.sortBy === 'name') {
        filtered.sort((a, b) => filters.sort === 'asc' ? a.title.localeCompare(b.title) : b.title.localeCompare(a.title));
      } else if (filters.sortBy === 'category') {
        filtered.sort((a, b) => {
          const ca = translateCategory(a.category) || '';
          const cb = translateCategory(b.category) || '';
          return filters.sort === 'asc' ? ca.localeCompare(cb) : cb.localeCompare(ca);
        });
      }

      setRecipes(filtered);
    } catch (error) {
      console.error('Failed to load recipes:', error);
      setRecipes([]);
    } finally { setLoading(false); }
  }, [searchTerm, filters, loadFavoriteIds]);

  useEffect(() => { loadRecipes(); }, [loadRecipes]);

  const handleSearch = (value: string) => setSearchTerm(value);
  const handleFilterChange = (newFilters: FilterType) => setFilters(newFilters);

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

  if (loading && recipes.length === 0) return <LoadingSpinner text="Загружаем рецепты..." />;

  return (
    <div className={styles.pageContainer}>
      <div className={styles.pageHeader}>
        <h1 className={styles.pageTitle}><span className={styles.pageTitleEmoji}>📖</span><span className={styles.pageTitleText}>Все рецепты</span></h1>
        <p className={styles.pageDescription}>Найди идеальный рецепт из нашей коллекции</p>
      </div>
      <div className={styles.searchSection}>
        <SearchBar onSearch={handleSearch} onFilterClick={() => setIsFilterOpen(true)} placeholder="Поиск по названию..." initialValue={searchTerm} isFilterActive={isFilterActive} />
      </div>
      <FilterPanel filters={filters} onFilterChange={handleFilterChange} isOpen={isFilterOpen} onClose={() => setIsFilterOpen(false)} />
      {loading ? <LoadingSpinner text="Обновляем..." /> : (
        recipes.length > 0 ? (
          <>
            <p className={styles.resultsCount}>Найдено рецептов: {recipes.length}</p>
            <RecipeList recipes={recipes} onFavoriteToggle={handleFavoriteToggle} />
          </>
        ) : (
          <div className={styles.emptyState}><div className={styles.emptyIcon}>🔍</div><h3 className={styles.emptyTitle}>Рецепты не найдены</h3><p className={styles.emptyText}>Попробуйте изменить параметры поиска или фильтры</p></div>
        )
      )}
    </div>
  );
};

export default RecipesPage;