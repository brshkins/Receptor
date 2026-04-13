import React from 'react';
import styles from './Pages.module.css';

const MatchPage: React.FC = () => {
  return (
    <div className={styles.pageContainer}>
      <div className={styles.pageHeader}>
        <h1 className={styles.pageTitle}>🔍 Подбор рецептов</h1>
        <p className={styles.pageDescription}>
          Выберите ингредиенты или загрузите фото продуктов
        </p>
      </div>
      
      <div style={{ 
        background: 'white', 
        padding: '32px', 
        borderRadius: 'var(--radius-lg)',
        textAlign: 'center' 
      }}>
        <div style={{ fontSize: '64px', marginBottom: '20px' }}>🤖</div>
        <h3>ИИ-подбор рецептов</h3>
        <p style={{ color: '#8B6F5E', marginTop: '12px' }}>
          Функционал в разработке
        </p>
      </div>
    </div>
  );
};

export default MatchPage;