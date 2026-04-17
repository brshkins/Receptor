// src/pages/FavoritesPage.tsx
import React from 'react';
import styles from './Pages.module.css';

const FavoritesPage: React.FC = () => {
  return (
    <div className={styles.pageContainer}>
      <div className={styles.pageHeader}>
        <h1 className={styles.pageTitle}>
          <span className={styles.pageTitleEmoji}>❤️</span>
          <span className={styles.pageTitleText}>Избранное</span>
        </h1>
        <p className={styles.pageDescription}>
          Ваши любимые рецепты
        </p>
      </div>
      <div className={styles.emptyFavorites}>
        <div className={styles.emptyIcon}>📖</div>
        <h3 className={styles.emptyTitle}>Пока пусто</h3>
        <p className={styles.emptyText}>
          Добавляйте рецепты в избранное, чтобы они всегда были под рукой
        </p>
      </div>
    </div>
  );
};

export default FavoritesPage;