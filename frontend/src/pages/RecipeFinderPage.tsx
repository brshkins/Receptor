// src/pages/RecipeFinderPage.tsx
import React, { useState, useRef, useEffect, useCallback } from 'react';
import { MatchRecipeList } from '../components/Recipes/MatchRecipeList';
import { LoadingSpinner } from '../components/Common/LoadingSpinner';
import { matchService } from '../services/matchService';
import { uploadService } from '../services/uploadService';
import type { MatchResponse } from '../types';
import { getApiErrorMessage } from '../utils/apiError';
import styles from './Pages.module.css';

const SEARCH_STATE_KEY = 'receptor_finder_search';

// Загружаем сохранённое состояние
const loadSavedState = () => {
  try {
    const saved = sessionStorage.getItem(SEARCH_STATE_KEY);
    if (saved) return JSON.parse(saved);
  } catch (e) { /* ignore */ }
  return null;
};

// Сохраняем состояние
const saveState = (hasSearched: boolean, recipes: MatchResponse[]) => {
  try {
    sessionStorage.setItem(SEARCH_STATE_KEY, JSON.stringify({ hasSearched, recipes }));
  } catch (e) { /* ignore */ }
};

// Очищаем состояние
const clearState = () => {
  sessionStorage.removeItem(SEARCH_STATE_KEY);
};

const RecipeFinderPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'ingredients' | 'camera'>('ingredients');
  const [ingredients, setIngredients] = useState<string[]>([]);
  const [currentIngredient, setCurrentIngredient] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  
  // Загружаем сохранённые результаты при старте
  const savedState = loadSavedState();
  const [recipes, setRecipes] = useState<MatchResponse[]>(savedState?.recipes || []);
  const [searchError, setSearchError] = useState<string | null>(null);
  const [hasSearched, setHasSearched] = useState<boolean>(savedState?.hasSearched || false);

  // Сохраняем результаты при изменении
  useEffect(() => {
    if (hasSearched && recipes.length > 0) {
      saveState(hasSearched, recipes);
    }
  }, [hasSearched, recipes]);

  const resetSearchResults = () => {
    setHasSearched(false);
    setSearchError(null);
    setRecipes([]);
    clearState();
  };
  
  // Для камеры
  const [stream, setStream] = useState<MediaStream | null>(null);
  const [capturedImages, setCapturedImages] = useState<string[]>([]);
  const [isMultiMode, setIsMultiMode] = useState(false);
  const videoRef = useRef<HTMLVideoElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  /** Русские подписи = то, что уходит в API; нормализация на backend */
  const popularIngredients = [
    'помидоры',
    'сыр',
    'курица',
    'паста',
    'лук',
    'чеснок',
    'сметана',
    'яйцо',
    'картофель',
    'морковь',
    'рис',
    'грибы',
  ];

  // ========== СЦЕНАРИЙ 1: РУЧНОЙ ВВОД ==========
  const handleAddIngredient = () => {
    if (currentIngredient.trim() && !ingredients.includes(currentIngredient.trim())) {
      setIngredients([...ingredients, currentIngredient.trim()]);
      setCurrentIngredient('');
    }
  };

  const handleRemoveIngredient = (ingredient: string) => {
    setIngredients(ingredients.filter(i => i !== ingredient));
  };

  const handleTogglePopular = (ingredient: string) => {
    if (ingredients.includes(ingredient)) {
      setIngredients(ingredients.filter(i => i !== ingredient));
    } else {
      setIngredients([...ingredients, ingredient]);
    }
  };

  const handleFindByIngredients = async () => {
    if (ingredients.length === 0) return;

    setSearchError(null);
    setIsLoading(true);
    setHasSearched(true);

    try {
      const recipes = await matchService.matchByIngredients(ingredients);
      setRecipes(recipes);
      saveState(true, recipes); // ← Сохраняем сразу
    } catch (error) {
      console.error('Ошибка поиска:', error);
      setSearchError(getApiErrorMessage(error, 'Не удалось выполнить подбор по ингредиентам'));
      setRecipes([]);
    } finally {
      setIsLoading(false);
    }
  };

  // ========== СЦЕНАРИЙ 2: КАМЕРА ==========
  const startCamera = async () => {
    try {
      if (stream) {
        stream.getTracks().forEach(track => track.stop());
      }

      const mediaStream = await navigator.mediaDevices.getUserMedia({
        video: {
          facingMode: 'environment',
          width: { ideal: 1920 },
          height: { ideal: 1080 }
        }
      });

      setStream(mediaStream);

      setTimeout(() => {
        if (videoRef.current) {
          videoRef.current.srcObject = mediaStream;
        }
      }, 100);
    } catch (error) {
      console.error('Ошибка камеры:', error);
      alert('Не удалось получить доступ к камере');
    }
  };

  const stopCamera = () => {
    if (stream) {
      stream.getTracks().forEach(track => track.stop());
      setStream(null);
    }
  };

  const capturePhoto = () => {
    if (videoRef.current && canvasRef.current) {
      const video = videoRef.current;
      const canvas = canvasRef.current;
      
      canvas.width = video.videoWidth || 640;
      canvas.height = video.videoHeight || 480;
      
      const context = canvas.getContext('2d');
      if (context) {
        context.drawImage(video, 0, 0, canvas.width, canvas.height);
        const imageData = canvas.toDataURL('image/jpeg', 0.9);
        
        if (isMultiMode) {
          setCapturedImages(prev => [...prev, imageData]);
        } else {
          setCapturedImages([imageData]);
          stopCamera();
        }
      }
    }
  };

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);
    if (files.length === 0) return;

    const readers = files.map(file => {
      return new Promise<string>((resolve) => {
        const reader = new FileReader();
        reader.onload = (event) => {
          resolve(event.target?.result as string);
        };
        reader.readAsDataURL(file);
      });
    });

    Promise.all(readers).then(images => {
      if (isMultiMode) {
        setCapturedImages(prev => [...prev, ...images]);
      } else {
        setCapturedImages(images);
      }
      stopCamera();
    });
  };

  const removeImage = (index: number) => {
    setCapturedImages(prev => prev.filter((_, i) => i !== index));
  };

  const clearAllImages = () => {
    setCapturedImages([]);
    stopCamera();
  };

  const retakePhoto = () => {
    setCapturedImages([]);
    stopCamera();
  };

  const handleFindByImages = async () => {
    if (capturedImages.length === 0) return;

    setSearchError(null);
    setIsLoading(true);
    setHasSearched(true);

    try {
      const blobs = await Promise.all(
        capturedImages.map((dataUrl) => fetch(dataUrl).then((r) => r.blob()))
      );
      const recipes = await uploadService.uploadAndMatch(blobs);
      setRecipes(recipes);
      saveState(true, recipes); // ← Сохраняем сразу
    } catch (error) {
      console.error('Ошибка распознавания:', error);
      setSearchError(getApiErrorMessage(error, 'Не удалось распознать фото или выполнить подбор'));
      setRecipes([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleTabChange = (tab: 'ingredients' | 'camera') => {
    setActiveTab(tab);
    resetSearchResults();
    if (tab === 'ingredients') {
      stopCamera();
      setCapturedImages([]);
    }
  };

  useEffect(() => {
    return () => {
      stopCamera();
    };
  }, [stream]);

  return (
    <div className={styles.pageContainer}>
      <div className={styles.pageHeader}>
        <h1 className={styles.pageTitle}>
          <span className={styles.pageTitleEmoji}>🔍</span>
          <span className={styles.pageTitleText}>Подбор рецептов</span>
        </h1>
        <p className={styles.pageDescription}>
          Найдите рецепт по вашим продуктам — вручную или с помощью камеры
        </p>
      </div>

      {/* Вкладки */}
      <div className={styles.matchTabs}>
        <button 
          className={`${styles.matchTab} ${activeTab === 'ingredients' ? styles.active : ''}`}
          onClick={() => handleTabChange('ingredients')}
        >
          <span className={styles.tabIcon}>📝</span>
          <span className={styles.tabText}>Список продуктов</span>
        </button>
        <button 
          className={`${styles.matchTab} ${activeTab === 'camera' ? styles.active : ''}`}
          onClick={() => handleTabChange('camera')}
        >
          <span className={styles.tabIcon}>📸</span>
          <span className={styles.tabText}>Камера</span>
          <span className={styles.tabBadge}>AI</span>
        </button>
      </div>

      <div className={styles.matchContent}>
        {/* ВКЛАДКА: РУЧНОЙ ВВОД */}
        {activeTab === 'ingredients' && !hasSearched && (
          <div className={styles.ingredientsPanel}>
            <div className={styles.panelHeader}>
              <h3>Выберите продукты</h3>
              <p>Добавьте ингредиенты, которые у вас есть</p>
            </div>

            <div className={styles.ingredientInput}>
              <input
                type="text"
                value={currentIngredient}
                onChange={(e) => setCurrentIngredient(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleAddIngredient()}
                placeholder="Например: помидоры, сыр, курица"
                className={styles.input}
              />
              <button onClick={handleAddIngredient} className={styles.addButton}>
                + Добавить
              </button>
            </div>

            <div className={styles.popularSection}>
              <h4>🔥 Популярные</h4>
              <div className={styles.popularGrid}>
                {popularIngredients.map(ing => (
                  <button
                    key={ing}
                    className={`${styles.popularChip} ${ingredients.includes(ing) ? styles.active : ''}`}
                    onClick={() => handleTogglePopular(ing)}
                  >
                    {ing}
                  </button>
                ))}
              </div>
            </div>

            {ingredients.length > 0 && (
              <div className={styles.selectedSection}>
                <h4>✅ Выбрано ({ingredients.length})</h4>
                <div className={styles.selectedGrid}>
                  {ingredients.map(ing => (
                    <div key={ing} className={styles.selectedChip}>
                      <span>{ing}</span>
                      <button 
                        onClick={() => handleRemoveIngredient(ing)}
                        className={styles.removeButton}
                      >
                        🗑️
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            )}

            <button 
              className={styles.findButton}
              onClick={handleFindByIngredients}
              disabled={ingredients.length === 0}
            >
              🔍 Найти рецепты
            </button>
          </div>
        )}

        {/* ВКЛАДКА: КАМЕРА / ЗАГРУЗКА ФОТО */}
        {activeTab === 'camera' && !hasSearched && (
          <div className={styles.cameraPanel}>
            <div className={styles.panelHeader}>
              <h3>Сфотографируйте или загрузите фото</h3>
              <p>ИИ распознает ингредиенты на фото</p>
            </div>

            {!stream && capturedImages.length === 0 && (
              <div className={styles.modeSelector}>
                <button 
                  className={`${styles.modeButton} ${!isMultiMode ? styles.active : ''}`}
                  onClick={() => setIsMultiMode(false)}
                >
                  📸 Одно фото
                </button>
                <button 
                  className={`${styles.modeButton} ${isMultiMode ? styles.active : ''}`}
                  onClick={() => setIsMultiMode(true)}
                >
                  🖼️ Несколько фото
                </button>
              </div>
            )}

            {!stream && capturedImages.length === 0 && (
              <div className={styles.cameraOptions}>
                <div className={styles.optionCard} onClick={startCamera}>
                  <div className={styles.optionIcon}>📸</div>
                  <h4>Сделать фото</h4>
                  <p>Использовать камеру устройства</p>
                </div>
        
                <div className={styles.optionDivider}>
                  <span>или</span>
                </div>
        
                <div className={styles.optionCard} onClick={() => fileInputRef.current?.click()}>
                  <div className={styles.optionIcon}>🖼️</div>
                  <h4>Загрузить с устройства</h4>
                  <p>Выбрать готовые фото</p>
                </div>
        
                <input
                  ref={fileInputRef}
                  type="file"
                  accept="image/*"
                  multiple={isMultiMode}
                  onChange={handleFileUpload}
                  style={{ display: 'none' }}
                />
              </div>
            )}

            {stream && (
              <div className={styles.cameraPreview}>
                <video 
                  ref={videoRef} 
                  autoPlay 
                  playsInline 
                  muted
                  className={styles.cameraVideo}
                />
                <div className={styles.cameraOverlay}>
                  <button className={styles.captureButton} onClick={capturePhoto}>
                    📸 {isMultiMode ? 'Сделать фото' : 'Сфотографировать'}
                  </button>
                  <button className={styles.backToOptionsButton} onClick={() => {
                    stopCamera();
                  }}>
                    ← {capturedImages.length > 0 ? 'Готово' : 'Назад'}
                  </button>
                </div>
        
                {isMultiMode && (
                  <div className={styles.multiModeHint}>
                    📸 Можно сделать несколько фото
                  </div>
                )}
              </div>
            )}

            {capturedImages.length > 0 && (
              <div className={styles.galleryContainer}>
                <div className={styles.galleryHeader}>
                  <h4>Выбрано фото: {capturedImages.length}</h4>
                  <div className={styles.galleryActions}>
                    {isMultiMode && !stream && (
                      <button className={styles.addMoreButton} onClick={startCamera}>
                        📸 Добавить с камеры
                      </button>
                    )}
                    {isMultiMode && (
                      <button className={styles.addMoreButton} onClick={() => fileInputRef.current?.click()}>
                        🖼️ Добавить файлы
                      </button>
                    )}
                    <button className={styles.clearAllButton} onClick={clearAllImages}>
                      🗑️ Очистить
                    </button>
                  </div>
                </div>
        
                <div className={styles.galleryGrid}>
                  {capturedImages.map((img, index) => (
                    <div key={index} className={styles.galleryItem}>
                      <img src={img} alt={`Фото ${index + 1}`} />
                      <button 
                        className={styles.removeImageButton}
                        onClick={() => removeImage(index)}
                      >
                        ✕
                      </button>
                      <span className={styles.imageNumber}>{index + 1}</span>
                    </div>
                  ))}
                </div>
        
                <button 
                  className={styles.analyzeButton}
                  onClick={handleFindByImages}
                  disabled={capturedImages.length === 0}
                >
                  🤖 Распознать ({capturedImages.length} фото)
                </button>
              </div>
            )}

            <canvas ref={canvasRef} style={{ display: 'none' }} />
          </div>
        )}

        {/* ЗАГРУЗКА */}
        {isLoading && (
          <div className={styles.loadingContainer}>
            <LoadingSpinner text={activeTab === 'camera' ? 'ИИ распознаёт продукты...' : 'Ищем рецепты...'} />
          </div>
        )}

        {/* РЕЗУЛЬТАТЫ */}
        {hasSearched && !isLoading && (
          <div className={styles.resultsPanel}>
            <div className={styles.resultsHeader}>
              <h3>🍽️ Найденные рецепты</h3>
              <button type="button" className={styles.backButton} onClick={resetSearchResults}>
                ← Назад к поиску
              </button>
            </div>

            {searchError ? (
              <div className={styles.emptyState} role="alert">
                <div className={styles.emptyIcon}>⚠️</div>
                <h3>Ошибка загрузки</h3>
                <p>{searchError}</p>
              </div>
            ) : recipes.length === 0 ? (
              <div className={styles.emptyState}>
                <div className={styles.emptyIcon}>😕</div>
                <h3>Ничего не найдено</h3>
                <p>Попробуйте другие ингредиенты или фото</p>
              </div>
            ) : (
              <MatchRecipeList items={recipes} />
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default RecipeFinderPage;