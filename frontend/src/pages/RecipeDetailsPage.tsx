import React, { useEffect, useMemo, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import { recipesService } from '../services/recipesService';
import { Recipe } from '../types';
import pageStyles from './Pages.module.css';
import recipeStyles from '../components/Recipes/Recipes.module.css';

const RecipeDetailsPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const recipeId = useMemo(() => {
    if (!id) return null;
    const n = Number(id);
    if (!Number.isInteger(n) || n <= 0) return null;
    return n;
  }, [id]);

  const [loading, setLoading] = useState(true);
  const [recipe, setRecipe] = useState<Recipe | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    const load = async () => {
      if (!recipeId) {
        setLoading(false);
        setRecipe(null);
        setError(null);
        return;
      }
      setLoading(true);
      setError(null);
      try {
        const data = await recipesService.getById(String(recipeId));
        if (cancelled) return;
        setRecipe(data);
      } catch (e: any) {
        if (cancelled) return;
        setRecipe(null);
        setError(e?.response?.data?.error || 'Не удалось загрузить рецепт');
      } finally {
        if (!cancelled) setLoading(false);
      }
    };
    load();
    return () => {
      cancelled = true;
    };
  }, [recipeId]);

  if (loading) return <LoadingSpinner text="Загружаем рецепт..." />;

  if (!recipeId) {
    return (
      <div className={pageStyles.pageContainer}>
        <div className={pageStyles.emptyState}>
          <div className={pageStyles.emptyIcon}>⚠️</div>
          <h3 className={pageStyles.emptyTitle}>Некорректный рецепт</h3>
          <p className={pageStyles.emptyText}>
            <Link to="/recipes">Вернуться к списку рецептов</Link>
          </p>
        </div>
      </div>
    );
  }

  if (!recipe) {
    return (
      <div className={pageStyles.pageContainer}>
        <div className={pageStyles.emptyState}>
          <div className={pageStyles.emptyIcon}>{error ? '⚠️' : '🔍'}</div>
          <h3 className={pageStyles.emptyTitle}>{error ? 'Ошибка' : 'Рецепт не найден'}</h3>
          <p className={pageStyles.emptyText}>{error || 'Попробуйте открыть другой рецепт'}</p>
          <p className={pageStyles.emptyText}>
            <Link to="/recipes">Вернуться к списку рецептов</Link>
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className={pageStyles.pageContainer}>
      <div className={pageStyles.pageHeader}>
        <h1 className={pageStyles.pageTitle}>
          <span className={pageStyles.pageTitleEmoji}>📖</span>
          <span className={pageStyles.pageTitleText}>{recipe.title}</span>
        </h1>
        <p className={pageStyles.pageDescription}>
          <Link to="/recipes">← Все рецепты</Link>
        </p>
      </div>

      <div className={recipeStyles.recipeCard} style={{ maxWidth: 720, margin: '0 auto' }}>
        {recipe.image ? (
          <img src={recipe.image} alt={recipe.title} className={recipeStyles.recipeImage} style={{ height: 360 }} />
        ) : (
          <div
            className={recipeStyles.recipeImage}
            style={{ height: 360, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 64 }}
          >
            🍳
          </div>
        )}
        <div className={recipeStyles.recipeContent}>
          <div className={recipeStyles.recipeMeta}>
            <span>⏱️ {recipe.cooking_time} мин</span>
            <span>🏷️ {recipe.category}</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default RecipeDetailsPage;

