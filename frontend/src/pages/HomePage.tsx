// src/pages/HomePage.tsx
import React from 'react';
import { Link } from 'react-router-dom';
import styles from './Pages.module.css';

const HomePage: React.FC = () => {
  return (
    <div className={styles.pageContainer}>
      <div className={styles.hero}>
        <h1 className={styles.heroTitle}>
          <span className={styles.heroTitleEmoji}>🍳</span>
          <span className={styles.heroTitleText}>Добро пожаловать в Receptor AI</span>
        </h1>
        <p className={styles.heroSubtitle}>
          Придумаем рецепт из того, что есть в холодильнике
        </p>
        <div className={styles.heroButtons}>
          <Link to="/match" className={`${styles.heroButton} ${styles.heroButtonPrimary}`}>
            🔍 Подобрать рецепт
          </Link>
          <Link to="/recipes" className={`${styles.heroButton} ${styles.heroButtonSecondary}`}>
            📖 Все рецепты
          </Link>
        </div>
      </div>

      <div className={styles.features}>
        <div className={styles.featureCard}>
          <div className={styles.featureIcon}>🤖</div>
          <h3 className={styles.featureTitle}>ИИ-подбор</h3>
          <p className={styles.featureDescription}>
            Загрузите фото продуктов или выберите ингредиенты — мы найдём идеальный рецепт
          </p>
        </div>
        <div className={styles.featureCard}>
          <div className={styles.featureIcon}>📚</div>
          <h3 className={styles.featureTitle}>Сотни рецептов</h3>
          <p className={styles.featureDescription}>
            Огромная база проверенных рецептов с подробными инструкциями
          </p>
        </div>
        <div className={styles.featureCard}>
          <div className={styles.featureIcon}>❤️</div>
          <h3 className={styles.featureTitle}>Избранное</h3>
          <p className={styles.featureDescription}>
            Сохраняйте любимые рецепты и возвращайтесь к ним в любое время
          </p>
        </div>
      </div>
    </div>
  );
};

export default HomePage;