import React, { useState } from 'react';
import styles from './Common.module.css';

interface SearchBarProps {
  onSearch: (value: string) => void;
  placeholder?: string;
  initialValue?: string;
}

export const SearchBar: React.FC<SearchBarProps> = ({ 
  onSearch, 
  placeholder = 'Поиск рецептов...',
  initialValue = ''
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
      <input
        type="text"
        className={styles.searchInput}
        placeholder={placeholder}
        value={value}
        onChange={handleChange}
      />
      <button type="submit" className={styles.searchButton}>
        🔍 Найти
      </button>
    </form>
  );
};