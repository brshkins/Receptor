import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import styles from './Pages.module.css';

const ProfilePage: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

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
        <div className={styles.profileHeader}>
          <div className={styles.profileAvatar}>
            {user.name?.charAt(0).toUpperCase() || '👤'}
          </div>
          <div className={styles.profileInfo}>
            <h2>{user.name}</h2>
            <p>{user.email}</p>
          </div>
        </div>
        
        <div className={styles.profileStats}>
          <div className={styles.statCard}>
            <div className={styles.statValue}>—</div>
            <div className={styles.statLabel}>Рецептов</div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statValue}>—</div>
            <div className={styles.statLabel}>В избранном</div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statValue}>—</div>
            <div className={styles.statLabel}>Подборок</div>
          </div>
        </div>
        
        <button onClick={handleLogout} className={styles.logoutProfileButton}>
          🚪 Выйти из аккаунта
        </button>
      </div>
    </div>
  );
};

export default ProfilePage;