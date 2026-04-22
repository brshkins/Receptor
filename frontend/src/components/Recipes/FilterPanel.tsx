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

  const handleCategoryChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.trim();
    setLocalFilters((prev) => ({ ...prev, category: value || undefined }));
  };

  const handleMaxTimeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const raw = e.target.value;
    const value = raw === '' ? undefined : Number(raw);
    setLocalFilters((prev) => ({
      ...prev,
      max_time: typeof value === 'number' && !Number.isNaN(value) ? value : undefined,
    }));
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setLocalFilters((prev) => ({ ...prev, sort: e.target.value as FilterType['sort'] }));
  };

  const handleApply = () => {
    onFilterChange(localFilters);
    onClose();
  };

  const handleReset = () => {
    const defaultFilters: FilterType = {
      sort: 'alphabet_asc'
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
          {/* Category */}
          <div className={styles.filterChipContainer}>
            <div className={styles.filterChipLabel}>
              <span className={styles.labelIcon}>🏷️</span>
              <span>Категория</span>
            </div>
            <input
              className={styles.filterSelect}
              value={localFilters.category || ''}
              onChange={handleCategoryChange}
              placeholder="Например: soup"
            />
          </div>

          {/* Max time */}
          <div className={styles.filterChipContainer}>
            <div className={styles.filterChipLabel}>
              <span className={styles.labelIcon}>⏱️</span>
              <span>Макс. время (мин)</span>
            </div>
            <input
              className={styles.filterSelect}
              type="number"
              min={0}
              value={typeof localFilters.max_time === 'number' ? String(localFilters.max_time) : ''}
              onChange={handleMaxTimeChange}
              placeholder="Например: 30"
            />
          </div>

          {/* Sort */}
          <div className={styles.filterChipContainer}>
            <div className={styles.filterChipLabel}>
              <span className={styles.labelIcon}>🔄</span>
              <span>Сортировка</span>
            </div>
            <select className={styles.filterSelect} value={localFilters.sort || 'alphabet_asc'} onChange={handleSortChange}>
              <option value="alphabet_asc">Название (А→Я)</option>
              <option value="alphabet_desc">Название (Я→А)</option>
              <option value="time_asc">Время (по возрастанию)</option>
              <option value="time_desc">Время (по убыванию)</option>
            </select>
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