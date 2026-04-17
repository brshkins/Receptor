// src/components/Common/SearchBar.tsx
import React, { useState } from 'react';
import styles from './Common.module.css';

interface SearchBarProps {
  onSearch: (value: string) => void;
  onFilterClick?: () => void;
  placeholder?: string;
  initialValue?: string;
  isFilterActive?: boolean;
}

export const SearchBar: React.FC<SearchBarProps> = ({ 
  onSearch, 
  onFilterClick,
  placeholder = 'Поиск рецептов...',
  initialValue = '',
  isFilterActive = false
}) => {
  const [value, setValue] = useState(initialValue);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch(value);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
  };

  return (
    <form className={styles.searchBar} onSubmit={handleSubmit}>
      {/* Кнопка фильтра слева */}
      {onFilterClick && (
        <button 
          type="button" 
          className={`${styles.filterButton} ${isFilterActive ? styles.filterActive : ''}`}
          onClick={onFilterClick}
          title="Фильтры"
        >
          <span className={styles.filterIcon}>🪄</span>
        </button>
      )}
      
      {/* Поле поиска */}
      <input
        type="text"
        className={styles.searchInput}
        placeholder={placeholder}
        value={value}
        onChange={handleChange}
      />
      
      {/* Кнопка поиска справа */}
      <button type="submit" className={styles.searchButton}>
        <span className={styles.searchIcon}>🔍</span>
      </button>
    </form>
  );
};