// src/components/Layout/Header.tsx
import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import styles from './Layout.module.css';

export const Header: React.FC = () => {
  const location = useLocation();
  const { user, isAuthenticated, logout } = useAuth();

  const isActive = (path: string) => location.pathname === path;
  const isHomePage = location.pathname === '/';

  return (
    <header className={styles.header}>
      <div className={styles.headerContent}>
        <Link to="/" className={styles.logo}>
          <span className={styles.logoIcon}>🍳</span>
          <span>Receptor</span>
        </Link>
        
        <nav className={styles.nav}>
          <Link to="/" className={`${styles.navLink} ${isActive('/') ? styles.active : ''}`}>
            <span>🏠</span> Главная
          </Link>
          <Link to="/recipes" className={`${styles.navLink} ${isActive('/recipes') ? styles.active : ''}`}>
            <span>📖</span> Рецепты
          </Link>
          <Link to="/match" className={`${styles.navLink} ${isActive('/match') ? styles.active : ''}`}>
            <span>🔍</span> Подбор
          </Link>
          <Link to="/favorites" className={`${styles.navLink} ${isActive('/favorites') ? styles.active : ''}`}>
            <span>❤️</span> Избранное
          </Link>
        </nav>
        
        <div className={styles.userMenu}>
          {isAuthenticated && user ? (
            <div className={styles.userProfile}>
              <Link to="/profile" className={styles.userAvatar}>
                {user.name?.charAt(0).toUpperCase() || '👤'}
              </Link>
              <button onClick={logout} className={styles.logoutButton} title="Выйти">
                🚪
              </button>
            </div>
          ) : (
            <Link to="/auth" className={styles.authButton}>
              <span>🔐</span> Войти
            </Link>
          )}
        </div>
      </div>
    </header>
  );
};