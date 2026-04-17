import React, { useState } from 'react';
import { RecipeFilters as FilterType } from '../../types';
import styles from './Recipes.module.css';

interface RecipeFiltersProps {
  onFilterChange: (filters: FilterType) => void;
  initialFilters?: FilterType;
}

export const RecipeFilters: React.FC<RecipeFiltersProps> = ({ 
  onFilterChange, 
  initialFilters = {} 
}) => {
  const [filters, setFilters] = useState<FilterType>({
    difficulty: initialFilters.difficulty,
    sort: initialFilters.sort || 'asc',
    sortBy: initialFilters.sortBy || 'name'
  });

  const handleDifficultyChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newFilters = { 
      ...filters, 
      difficulty: e.target.value as FilterType['difficulty'] 
    };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newFilters = { 
      ...filters, 
      sort: e.target.value as FilterType['sort']
    };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  const handleSortByChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const newFilters = { 
      ...filters, 
      sortBy: e.target.value as FilterType['sortBy']
    };
    setFilters(newFilters);
    onFilterChange(newFilters);
  };

  return (
    <div className={styles.filtersContainer}>
      <select 
        className={styles.sortSelect}
        value={filters.difficulty || ''} 
        onChange={handleDifficultyChange}
      >
        <option value="">Все сложности</option>
        <option value="easy">Легко</option>
        <option value="medium">Средне</option>
        <option value="hard">Сложно</option>
      </select>

      <select 
        className={styles.sortSelect}
        value={filters.sortBy} 
        onChange={handleSortByChange}
      >
        <option value="name">По названию</option>
        <option value="time">По времени</option>
        <option value="difficulty">По сложности</option>
      </select>

      <select 
        className={styles.sortSelect}
        value={filters.sort} 
        onChange={handleSortChange}
      >
        <option value="asc">По возрастанию</option>
        <option value="desc">По убыванию</option>
      </select>
    </div>
  );
};