// src/pages/MatchPage.tsx
import React, { useState } from 'react';
import { matchService, MatchResponseItem } from '../services/matchService';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import styles from './Pages.module.css';
import recipeStyles from '../components/Recipes/Recipes.module.css';

const MatchPage: React.FC = () => {
  const [ingredients, setIngredients] = useState<string[]>([]);
  const [currentIngredient, setCurrentIngredient] = useState('');
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<MatchResponseItem[]>([]);
  const [error, setError] = useState<string | null>(null);

  const addIngredient = () => {
    const v = currentIngredient.trim();
    if (!v) return;
    if (ingredients.includes(v)) {
      setCurrentIngredient('');
      return;
    }
    setIngredients((prev) => [...prev, v]);
    setCurrentIngredient('');
  };

  const removeIngredient = (v: string) => {
    setIngredients((prev) => prev.filter((x) => x !== v));
  };

  const onSubmit = async () => {
    if (ingredients.length === 0) return;
    setLoading(true);
    setError(null);
    try {
      const data = await matchService.matchByIngredients(ingredients);
      setResults(data);
    } catch (e: any) {
      setResults([]);
      setError(e?.response?.data?.error || 'Не удалось подобрать рецепты');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.pageContainer}>
      <div className={styles.pageHeader}>
        <h1 className={styles.pageTitle}>
          <span className={styles.pageTitleEmoji}>🤖</span>
          <span className={styles.pageTitleText}>Подбор рецептов по ингредиентам</span>
        </h1>
        <p className={styles.pageDescription}>Введите ингредиенты и получите подходящие рецепты</p>
      </div>

      <div className={styles.ingredientsPanel}>
        <div className={styles.ingredientInput}>
          <input
            type="text"
            value={currentIngredient}
            onChange={(e) => setCurrentIngredient(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                e.preventDefault();
                addIngredient();
              }
            }}
            placeholder="Например: tomato, cheese..."
            className={styles.input}
          />
          <button onClick={addIngredient} className={styles.addButton} disabled={!currentIngredient.trim()}>
            <span>+</span> Добавить
          </button>
        </div>

        {ingredients.length > 0 && (
          <div className={styles.selectedSection}>
            <h4>Выбранные продукты ({ingredients.length})</h4>
            <div className={styles.selectedGrid}>
              {ingredients.map((ing) => (
                <div key={ing} className={styles.selectedChip}>
                  {ing}
                  <button onClick={() => removeIngredient(ing)} className={styles.removeButton}>
                    ✕
                  </button>
                </div>
              ))}
            </div>
          </div>
        )}

        <button className={styles.findButton} onClick={onSubmit} disabled={ingredients.length === 0 || loading}>
          {loading ? (
            <>
              <span className={styles.spinnerSmall}></span>
              Подбираем...
            </>
          ) : (
            <>
              <span>🔍</span> Подобрать рецепты
            </>
          )}
        </button>
      </div>

      {loading && results.length === 0 && <LoadingSpinner text="Подбираем рецепты..." />}

      {error && (
        <div className={styles.emptyState}>
          <div className={styles.emptyIcon}>⚠️</div>
          <h3 className={styles.emptyTitle}>Ошибка</h3>
          <p className={styles.emptyText}>{error}</p>
        </div>
      )}

      {!loading && !error && results.length > 0 && (
        <div style={{ marginTop: 24 }}>
          <p className={styles.resultsCount}>Найдено: {results.length}</p>
          <div className={recipeStyles.recipeGrid}>
            {results.map((m) => (
              <div key={m.recipe_id} className={recipeStyles.recipeCard}>
                {m.image ? (
                  <img src={m.image} alt={m.title} className={recipeStyles.recipeImage} />
                ) : (
                  <div
                    className={recipeStyles.recipeImage}
                    style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: '48px' }}
                  >
                    🍳
                  </div>
                )}
                <div className={recipeStyles.recipeContent}>
                  <h3 className={recipeStyles.recipeTitle}>{m.title}</h3>
                  <div className={recipeStyles.recipeMeta}>
                    <span>✅ {Math.round(m.match_percent)}%</span>
                    {m.missing_ingredients?.length ? <span>➖ нет: {m.missing_ingredients.join(', ')}</span> : <span>🎯 всё есть</span>}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default MatchPage;