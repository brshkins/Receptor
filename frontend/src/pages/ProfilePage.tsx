// src/pages/ProfilePage.tsx
import React, { useEffect, useState, useCallback } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { favoritesService } from '../services/favoritesService';
import styles from './Pages.module.css';

const ProfilePage: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [favoritesCount, setFavoritesCount] = useState<number | null>(null);

  const loadFavoritesCount = useCallback(async () => {
    if (!user) return;
    
    try {
      const favorites = await favoritesService.getAll();
      const count = Array.isArray(favorites) ? favorites.length : 0;
      setFavoritesCount(count);
    } catch (error) {
      console.error('Ошибка загрузки избранного:', error);
      setFavoritesCount(0);
    }
  }, [user]);

  // Загружаем при монтировании
  useEffect(() => {
    loadFavoritesCount();
  }, [loadFavoritesCount]);

  // Обновляем счётчик при возврате на страницу
  useEffect(() => {
    const handleFocus = () => {
      loadFavoritesCount();
    };
    
    window.addEventListener('focus', handleFocus);
    return () => window.removeEventListener('focus', handleFocus);
  }, [loadFavoritesCount]);

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  if (!user) {
    return (
      <div className={styles.pageContainer}>
        <div className={styles.emptyState}>
          <div className={styles.emptyIcon}>🔐</div>
          <h3 className={styles.emptyTitle}>Вы не авторизованы</h3>
          <p className={styles.emptyText}>
            <Link to="/auth">Войдите</Link> или <Link to="/auth">зарегистрируйтесь</Link>
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.pageContainer}>
      <div className={styles.pageHeader}>
        <h1 className={styles.pageTitle}>
          <span className={styles.pageTitleEmoji}>👤</span>
          <span className={styles.pageTitleText}>Личный кабинет</span>
        </h1>
      </div>
      
      <div className={styles.profileCard}>
        {/* Аватар и информация */}
        <div className={styles.profileHeader}>
          <div className={styles.profileAvatar}>
            {user.name?.charAt(0).toUpperCase() || '👤'}
          </div>
          <div className={styles.profileInfo}>
            <h2>{user.name}</h2>
            <p>{user.email}</p>
          </div>
        </div>
        
        {/* Карточки действий */}
        <div className={styles.profileActions}>
          <Link to="/favorites" className={styles.actionCard}>
            <span className={styles.actionIcon}>❤️</span>
            <div className={styles.actionInfo}>
              <h4>Избранные рецепты</h4>
              <p>
                {favoritesCount === null 
                  ? 'Загрузка...' 
                  : favoritesCount > 0 
                    ? `Сохранено ${favoritesCount} рецептов` 
                    : 'Пока нет избранных рецептов'}
              </p>
            </div>
            <span className={styles.actionArrow}>→</span>
          </Link>
          
          <Link to="/recipes" className={styles.actionCard}>
            <span className={styles.actionIcon}>📖</span>
            <div className={styles.actionInfo}>
              <h4>Все рецепты</h4>
              <p>Просмотр всех доступных рецептов</p>
            </div>
            <span className={styles.actionArrow}>→</span>
          </Link>
          
          <Link to="/match" className={styles.actionCard}>
            <span className={styles.actionIcon}>🔍</span>
            <div className={styles.actionInfo}>
              <h4>Подбор рецептов</h4>
              <p>Найти рецепт по продуктам</p>
            </div>
            <span className={styles.actionArrow}>→</span>
          </Link>
        </div>
        
        {/* Кнопка выхода */}
        <button onClick={handleLogout} className={styles.logoutButton}>
          <span className={styles.logoutIcon}>🚪</span>
          <span>Выйти из аккаунта</span>
        </button>
      </div>
    </div>
  );
};

export default ProfilePage;