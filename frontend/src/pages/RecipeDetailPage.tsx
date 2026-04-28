// src/pages/RecipeDetailPage.tsx
import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { recipesService } from '../services/recipesService';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import { 
  translateRecipeTitle, 
  translateCategory, 
  translateIngredient, 
  getCategoryEmoji, 
  formatTime 
} from '../utils/translations';
import styles from './Pages.module.css';

interface RecipeDetail {
  id: number;
  title: string;
  description?: string;
  image?: string;
  image_url?: string;
  cooking_time: number;
  category?: string;
  // Разные возможные форматы ингредиентов
  ingredients?: any[];
  recipe_ingredients?: any[];
  // Разные возможные форматы шагов
  steps?: any[];
  recipe_steps?: any[];
  instructions?: any[];
}

const RecipeDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [recipe, setRecipe] = useState<RecipeDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadRecipe();
  }, [id]);

  const loadRecipe = async () => {
    if (!id) return;
    
    try {
      setLoading(true);
      const response = await recipesService.getById(id);
      console.log('🔥 Ответ бэкенда:', JSON.stringify(response, null, 2));
      
      // Извлекаем рецепт из ответа
      const recipeData = response?.data || response;
      console.log('📋 Данные рецепта:', recipeData);
      console.log('🛒 Ингредиенты:', recipeData?.ingredients);
      console.log('📝 Шаги:', recipeData?.steps);
      console.log('📝 recipe_steps:', recipeData?.recipe_steps);
      console.log('📝 instructions:', recipeData?.instructions);
      
      setRecipe(recipeData);
    } catch (error) {
      console.error('Ошибка загрузки рецепта:', error);
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

  // Собираем ингредиенты из всех возможных полей
  const allIngredients = recipe.recipe_ingredients || recipe.ingredients || [];
  
  // Собираем шаги из всех возможных полей
  const allSteps = recipe.recipe_steps || recipe.steps || recipe.instructions || [];

  return (
    <div className={styles.pageContainer}>
      <button onClick={() => navigate(-1)} className={styles.backButtonBig}>
        ← Назад
      </button>

      <div className={styles.recipeDetail}>
        <div className={styles.detailHeader}>
          <div className={styles.detailImagePlaceholder}>
            <span className={styles.detailEmoji}>{categoryEmoji}</span>
          </div>
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
            {recipe.description && (
              <p className={styles.detailDescription}>{recipe.description}</p>
            )}
          </div>
        </div>

        {/* Ингредиенты */}
        <div className={styles.detailSection}>
          <h2 className={styles.detailSectionTitle}>
            <span>🛒</span> Ингредиенты
          </h2>
          {allIngredients.length > 0 ? (
            <div className={styles.ingredientsList}>
              {allIngredients.map((ing: any, index: number) => {
                const name = ing.name || ing.ingredient_name || ing.ingredient?.name || 'Неизвестно';
                const amount = ing.amount || ing.quantity || '';
                
                return (
                  <div key={index} className={styles.ingredientItem}>
                    <span className={styles.ingredientName}>
                      {translateIngredient(name)}
                    </span>
                    {amount && (
                      <span className={styles.ingredientAmount}>
                        {amount}
                      </span>
                    )}
                  </div>
                );
              })}
            </div>
          ) : (
            <p className={styles.noData}>Ингредиенты не указаны</p>
          )}
        </div>

        {/* Инструкция */}
        <div className={styles.detailSection}>
          <h2 className={styles.detailSectionTitle}>
            <span>📝</span> Инструкция
          </h2>
          {allSteps.length > 0 ? (
            <ol className={styles.stepsList}>
              {allSteps.map((step: any, index: number) => {
                const text = typeof step === 'string' 
                  ? step 
                  : step.description || step.step_description || step.instruction || '';
                const stepNum = step.step_number || step.number || index + 1;
                
                return (
                  <li key={index} className={styles.stepItem}>
                    <span className={styles.stepNumber}>{stepNum}</span>
                    <span className={styles.stepText}>{text}</span>
                  </li>
                );
              })}
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