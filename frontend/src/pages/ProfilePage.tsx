import React from 'react';
import { useAuth } from '../contexts/AuthContext';
import styles from './Pages.module.css';

const ProfilePage: React.FC = () => {
  const { user } = useAuth();

  return (
    <div className={styles.pageContainer}>
      <div className={styles.pageHeader}>
        <h1 className={styles.pageTitle}>👤 Профиль</h1>
        <p className={styles.pageDescription}>
          Ваши данные и настройки
        </p>
      </div>
      
      <div className={styles.profileCard}>
        <div className={styles.profileHeader}>
          <div className={styles.profileAvatar}>
            {user?.name?.charAt(0).toUpperCase() || '👤'}
          </div>
          <div className={styles.profileInfo}>
            <h2>{user?.name || 'Гость'}</h2>
            <p>{user?.email || 'Не авторизован'}</p>
          </div>
        </div>
        
        <div className={styles.profileStats}>
          <div className={styles.statCard}>
            <div className={styles.statValue}>0</div>
            <div className={styles.statLabel}>Рецептов</div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statValue}>0</div>
            <div className={styles.statLabel}>В избранном</div>
          </div>
          <div className={styles.statCard}>
            <div className={styles.statValue}>0</div>
            <div className={styles.statLabel}>Подборок</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProfilePage;