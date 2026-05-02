// src/pages/RecipeDetailPage.tsx
import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { recipesService } from '../services/recipesService';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import {
  translateRecipeTitle,
  translateCategory,
  getCategoryEmoji,
  formatTime,
} from '../utils/translations';
import {
  translateIngredientDisplayName,
  translateRecipeEnglishText,
} from '../utils/recipeTextRu';
import { formatRecipeAmountRu } from '../utils/formatRecipeAmountRu';
import { recipeBodyRuByTitleLower } from '../utils/recipeBodyRu.generated';
import type { RecipeDetailsDto } from '../types';
import styles from './Pages.module.css';

const RecipeDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [recipe, setRecipe] = useState<RecipeDetailsDto | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadRecipe();
  }, [id]);

  const loadRecipe = async () => {
    if (!id) return;

    try {
      setLoading(true);
      setError('');
      const recipeData = await recipesService.getRecipeDetails(id);
      setRecipe(recipeData);
    } catch (e) {
      console.error('Ошибка загрузки рецепта:', e);
      setError('Не удалось загрузить рецепт');
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <LoadingSpinner text="Загружаем рецепт..." />;

  if (error || !recipe) {
    return (
      <div className={styles.pageContainer}>
        <div className={styles.emptyState}>
          <div className={styles.emptyIcon}>😕</div>
          <h3 className={styles.emptyTitle}>{error || 'Рецепт не найден'}</h3>
          <button className={styles.backButtonBig} onClick={() => navigate('/recipes')}>
            ← К рецептам
          </button>
        </div>
      </div>
    );
  }

  const titleRu = translateRecipeTitle(recipe.title);
  const categoryRu = translateCategory(recipe.category);
  const categoryEmoji = getCategoryEmoji(recipe.category);
  const timeFormatted = formatTime(recipe.cooking_time);
  const bodyRu = recipeBodyRuByTitleLower[recipe.title.trim().toLowerCase()];
  const steps =
    bodyRu?.steps_ru != null && bodyRu.steps_ru.length > 0
      ? bodyRu.steps_ru
      : (recipe.steps ?? []).map((t) => translateRecipeEnglishText(t));
  const descriptionText =
    (bodyRu?.description_ru && bodyRu.description_ru.trim()) ||
    translateRecipeEnglishText(recipe.description ?? '');
  const ingredients = recipe.ingredients ?? [];

  return (
    <div className={styles.pageContainer}>
      <button type="button" onClick={() => navigate(-1)} className={styles.backButtonBig}>
        ← Назад
      </button>

      <div className={styles.recipeDetail}>
        <div className={styles.detailHeader}>
          {recipe.image ? (
            <img src={recipe.image} alt={titleRu} className={styles.detailHeroImage} />
          ) : (
            <div className={styles.detailImagePlaceholder}>
              <span className={styles.detailEmoji}>{categoryEmoji}</span>
            </div>
          )}
          <div className={styles.detailHeaderInfo}>
            <h1 className={styles.detailTitle}>{titleRu}</h1>
            <div className={styles.detailMeta}>
              <span className={styles.detailMetaItem}>
                <span className={styles.metaIcon}>⏱️</span>
                <span>{timeFormatted}</span>
              </span>
              {categoryRu && (
                <span className={styles.detailMetaItem}>
                  <span className={styles.metaIcon}>{categoryEmoji}</span>
                  <span>{categoryRu}</span>
                </span>
              )}
            </div>
            {descriptionText.trim() ? (
              <p className={styles.detailDescription}>{descriptionText}</p>
            ) : null}
          </div>
        </div>

        {ingredients.length > 0 && (
          <div className={styles.detailSection}>
            <h2 className={styles.detailSectionTitle}>
              <span>🥕</span> Ингредиенты
            </h2>
            <div className={styles.ingredientsList}>
              {ingredients.map((line, index) => {
                const nameRu = translateIngredientDisplayName(line.name);
                const amt = formatRecipeAmountRu(line.amount ?? '');
                return (
                  <div key={`${index}-${line.name}`} className={styles.ingredientItem}>
                    <span className={styles.ingredientName}>{nameRu}</span>
                    {amt ? <span className={styles.ingredientAmount}>{amt}</span> : null}
                  </div>
                );
              })}
            </div>
          </div>
        )}

        <div className={styles.detailSection}>
          <h2 className={styles.detailSectionTitle}>
            <span>📝</span> Инструкция
          </h2>
          {steps.length > 0 ? (
            <ol className={styles.stepsList}>
              {steps.map((text, index) => (
                <li key={index} className={styles.stepItem}>
                  <span className={styles.stepNumber}>{index + 1}</span>
                  <span className={styles.stepText}>{text}</span>
                </li>
              ))}
            </ol>
          ) : (
            <p className={styles.noData}>Инструкция не указана</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default RecipeDetailPage;
