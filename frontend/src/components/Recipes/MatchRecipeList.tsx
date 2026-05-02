import React from 'react';
import { Link } from 'react-router-dom';
import type { MatchResponse } from '../../types';
import { translateIngredientDisplayName } from '../../utils/recipeTextRu';
import { translateRecipeTitle, formatTime } from '../../utils/translations';
import styles from './Recipes.module.css';

/** API отдаёт `match_percent` как долю 0–1 (см. docs/API.md). */
function matchPercentForDisplay(v: unknown): number {
  if (typeof v !== 'number' || Number.isNaN(v)) return 0;
  if (v >= 0 && v <= 1) return Math.round(v * 100);
  return Math.round(v);
}

interface MatchRecipeListProps {
  items: MatchResponse[];
}

/** Список результатов матча: поля только из MatchResponse, переход на /recipes/:recipe_id. */
export const MatchRecipeList: React.FC<MatchRecipeListProps> = ({ items }) => {
  if (!items?.length) {
    return (
      <div className={styles.emptyState}>
        <p>Рецепты не найдены</p>
      </div>
    );
  }

  return (
    <div className={styles.recipeGrid}>
      {items.map((m) => {
        const titleRu = translateRecipeTitle(m.title);
        const pct = matchPercentForDisplay(m.match_percent);
        const timeLabel =
          typeof m.cooking_time === 'number' && m.cooking_time > 0
            ? formatTime(m.cooking_time)
            : '';
        return (
          <div key={m.recipe_id} className={styles.recipeCard}>
            <Link to={`/recipes/${m.recipe_id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
              {m.image ? (
                <img src={m.image} alt={titleRu} className={styles.recipeImage} />
              ) : (
                <div className={styles.recipeImagePlaceholder}>
                  <span className={styles.placeholderEmoji}>🍽️</span>
                </div>
              )}
              <div className={styles.recipeContent}>
                <h3 className={styles.recipeTitle}>{titleRu}</h3>
                <div className={styles.recipeMeta}>
                  {timeLabel ? (
                    <span className={styles.metaItem}>
                      <span className={styles.metaIcon}>⏱️</span>
                      <span>{timeLabel}</span>
                    </span>
                  ) : null}
                  <span className={styles.metaItem}>
                    <span className={styles.metaIcon}>📊</span>
                    <span>Совпадение: {pct}%</span>
                  </span>
                </div>
                {(m.missing_ingredients?.length ?? 0) > 0 && (
                  <p className={styles.matchMissing}>
                    Не хватает:{' '}
                    {(m.missing_ingredients ?? []).map(translateIngredientDisplayName).join(', ')}
                  </p>
                )}
              </div>
            </Link>
          </div>
        );
      })}
    </div>
  );
};
