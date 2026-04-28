// src/pages/HomePage.tsx
import React from 'react';
import { Link } from 'react-router-dom';
import styles from './Pages.module.css';

const HomePage: React.FC = () => {
  return (
    <div className={styles.pageContainer}>
      {/* Hero секция с информацией */}
      <div className={styles.hero}>
        <h1 className={styles.heroTitle}>
          <span className={styles.heroTitleEmoji}>🍳</span>
          <span className={styles.heroTitleText}>Receptor AI</span>
        </h1>
        <p className={styles.heroSubtitle}>
          Ваш умный помощник в мире кулинарии
        </p>
        
        {/* Информационный блок */}
        <div className={styles.infoBadges}>
          <div className={styles.infoBadge}>
            <span className={styles.badgeIcon}>📖</span>
            <span className={styles.badgeText}>129+ рецептов</span>
          </div>
          <div className={styles.infoBadge}>
            <span className={styles.badgeIcon}>🔍</span>
            <span className={styles.badgeText}>Умный подбор</span>
          </div>
          <div className={styles.infoBadge}>
            <span className={styles.badgeIcon}>📸</span>
            <span className={styles.badgeText}>Поиск по фото</span>
          </div>
        </div>
        
        <p className={styles.heroDescription}>
          Сфотографируйте продукты, выберите ингредиенты или просто найдите 
          рецепт — мы поможем приготовить вкусное блюдо из того, что есть под рукой
        </p>
      </div>

      {/* Кликабельные карточки-ссылки */}
      <div className={styles.features}>
        <Link to="/match" className={styles.featureCard}>
          <div className={styles.featureIcon}>📸🤖</div>
          <h3 className={styles.featureTitle}>Подбор рецептов</h3>
          <p className={styles.featureDescription}>
            Сфотографируйте продукты или выберите вручную — искусственный интеллект 
            подберёт подходящие рецепты из нашей базы
          </p>
          <span className={styles.featureArrow}>→</span>
        </Link>
        
        <Link to="/recipes" className={styles.featureCard}>
          <div className={styles.featureIcon}>📚🍝</div>
          <h3 className={styles.featureTitle}>Все рецепты</h3>
          <p className={styles.featureDescription}>
            129 проверенных рецептов с подробными инструкциями, 
            ингредиентами и временем приготовления
          </p>
          <span className={styles.featureArrow}>→</span>
        </Link>
        
        <Link to="/favorites" className={styles.featureCard}>
          <div className={styles.featureIcon}>❤️📌</div>
          <h3 className={styles.featureTitle}>Избранное</h3>
          <p className={styles.featureDescription}>
            Сохраняйте любимые рецепты, чтобы они всегда были под рукой. 
            Ваша личная кулинарная книга
          </p>
          <span className={styles.featureArrow}>→</span>
        </Link>
      </div>
    </div>
  );
};

export default HomePage;