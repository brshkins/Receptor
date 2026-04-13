// src/components/Recipes/FilterPanel.tsx
import React, { useState, useEffect } from 'react';
import { RecipeFilters as FilterType } from '../../types';
import styles from './Recipes.module.css';

interface FilterPanelProps {
  filters: FilterType;
  onFilterChange: (filters: FilterType) => void;
  isOpen: boolean;
  onClose: () => void;
}

export const FilterPanel: React.FC<FilterPanelProps> = ({
  filters,
  onFilterChange,
  isOpen,
  onClose
}) => {
  const [localFilters, setLocalFilters] = useState<FilterType>(filters);

  // Сбрасываем локальные фильтры при открытии
  useEffect(() => {
    if (isOpen) {
      setLocalFilters(filters);
    }
  }, [isOpen, filters]);

  const handleDifficultyChange = (difficulty: FilterType['difficulty']) => {
    setLocalFilters(prev => ({
      ...prev,
      difficulty: prev.difficulty === difficulty ? undefined : difficulty
    }));
  };

  const handleSortChange = (sort: FilterType['sort']) => {
    setLocalFilters(prev => ({ ...prev, sort }));
  };

  const handleSortByChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setLocalFilters(prev => ({ 
      ...prev, 
      sortBy: e.target.value as FilterType['sortBy']
    }));
  };

  const handleApply = () => {
    onFilterChange(localFilters);
    onClose();
  };

  const handleReset = () => {
    const defaultFilters: FilterType = {
      sort: 'asc',
      sortBy: 'name'
    };
    setLocalFilters(defaultFilters);
  };

  const handleClose = () => {
    // Сбрасываем локальные изменения при закрытии без применения
    setLocalFilters(filters);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <>
      <div className={styles.filterOverlay} onClick={handleClose} />
      <div className={styles.filterPanel}>
        <div className={styles.filterHeader}>
          <h3>
            <span className={styles.filterIcon}>🪄</span>
            Фильтры
          </h3>
          <button className={styles.closeButton} onClick={handleClose}>
            <span>✕</span>
          </button>
        </div>
        
        <div className={styles.filterContent}>
          {/* Плашка с фильтрами сложности */}
          <div className={styles.filterChipContainer}>
            <div className={styles.filterChipLabel}>
              <span className={styles.labelIcon}>📊</span>
              <span>Сложность</span>
            </div>
            <div className={styles.chipGroup}>
              <button 
                className={`${styles.chip} ${styles.chipSmall} ${localFilters.difficulty === 'easy' ? styles.chipActive : ''}`}
                onClick={() => handleDifficultyChange('easy')}
              >
                <span className={styles.chipIcon}>🥚</span>
                <span className={styles.chipText}>Легко</span>
              </button>
              
              <button 
                className={`${styles.chip} ${styles.chipSmall} ${localFilters.difficulty === 'medium' ? styles.chipActive : ''}`}
                onClick={() => handleDifficultyChange('medium')}
              >
                <span className={styles.chipIcon}>👨‍🍳</span>
                <span className={styles.chipText}>Средне</span>
              </button>
              
              <button 
                className={`${styles.chip} ${styles.chipSmall} ${localFilters.difficulty === 'hard' ? styles.chipActive : ''}`}
                onClick={() => handleDifficultyChange('hard')}
              >
                <span className={styles.chipIcon}>🔥</span>
                <span className={styles.chipText}>Сложно</span>
              </button>
            </div>
          </div>

          {/* Плашка с сортировкой */}
          <div className={styles.filterChipContainer}>
            <div className={styles.filterChipLabel}>
              <span className={styles.labelIcon}>📝</span>
              <span>Сортировка</span>
            </div>
            <select 
              className={styles.filterSelect}
              value={localFilters.sortBy} 
              onChange={handleSortByChange}
            >
              <option value="name">По названию</option>
              <option value="time">По времени</option>
              <option value="difficulty">По сложности</option>
            </select>
          </div>

          {/* Плашка с порядком сортировки */}
          <div className={styles.filterChipContainer}>
            <div className={styles.filterChipLabel}>
              <span className={styles.labelIcon}>🔄</span>
              <span>Порядок</span>
            </div>
            <div className={styles.chipGroup}>
              <button 
                className={`${styles.chip} ${styles.chipSmall} ${localFilters.sort === 'asc' ? styles.chipActive : ''}`}
                onClick={() => handleSortChange('asc')}
              >
                <span className={styles.chipIcon}>🔼</span>
                <span className={styles.chipText}>Возр.</span>
              </button>
              
              <button 
                className={`${styles.chip} ${styles.chipSmall} ${localFilters.sort === 'desc' ? styles.chipActive : ''}`}
                onClick={() => handleSortChange('desc')}
              >
                <span className={styles.chipIcon}>🔽</span>
                <span className={styles.chipText}>Убыв.</span>
              </button>
            </div>
          </div>
        </div>
        
        <div className={styles.filterFooter}>
          <button className={styles.resetButton} onClick={handleReset}>
            <span className={styles.buttonIcon}>🔄</span>
            <span>Сбросить</span>
          </button>
          <button className={styles.applyButton} onClick={handleApply}>
            <span className={styles.buttonIcon}>✓</span>
            <span>Применить</span>
          </button>
        </div>
      </div>
    </>
  );
};