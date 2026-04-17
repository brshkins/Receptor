// src/pages/MatchPage.tsx
import React, { useState, useRef } from 'react';
import styles from './Pages.module.css';

const MatchPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'ingredients' | 'photo'>('ingredients');
  const [ingredients, setIngredients] = useState<string[]>([]);
  const [currentIngredient, setCurrentIngredient] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [previewImages, setPreviewImages] = useState<string[]>([]);
  const [dragActive, setDragActive] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Популярные ингредиенты для быстрого выбора
  const popularIngredients = [
    'Помидоры', 'Сыр', 'Курица', 'Макароны', 
    'Лук', 'Чеснок', 'Сливки', 'Яйца',
    'Картофель', 'Морковь', 'Рис', 'Грибы'
  ];

  const handleAddIngredient = () => {
    if (currentIngredient.trim() && !ingredients.includes(currentIngredient.trim())) {
      setIngredients([...ingredients, currentIngredient.trim()]);
      setCurrentIngredient('');
    }
  };

  const handleRemoveIngredient = (ingredient: string) => {
    setIngredients(ingredients.filter(i => i !== ingredient));
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddIngredient();
    }
  };

  const handleTogglePopular = (ingredient: string) => {
    if (ingredients.includes(ingredient)) {
      handleRemoveIngredient(ingredient);
    } else {
      setIngredients([...ingredients, ingredient]);
    }
  };

  const handleFindRecipes = async () => {
    if (ingredients.length === 0) return;
    
    setIsLoading(true);
    // Имитация запроса к ИИ
    setTimeout(() => {
      setIsLoading(false);
      console.log('Поиск рецептов с ингредиентами:', ingredients);
    }, 2000);
  };

  const handleImageUpload = (files: FileList) => {
    const newImages: string[] = [];
    
    Array.from(files).forEach(file => {
      if (file.type.startsWith('image/')) {
        const reader = new FileReader();
        reader.onload = (e) => {
          const result = e.target?.result as string;
          newImages.push(result);
          
          // Когда все файлы загружены
          if (newImages.length === files.length) {
            setPreviewImages(prev => [...prev, ...newImages]);
          }
        };
        reader.readAsDataURL(file);
      }
    });
    
    // Имитация анализа фото ИИ
    if (files.length > 0) {
      setIsLoading(true);
      setTimeout(() => {
        setIsLoading(false);
        console.log('Анализ фото:', files.length, 'изображений');
      }, 2000);
    }
  };

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      handleImageUpload(e.dataTransfer.files);
    }
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      handleImageUpload(e.target.files);
    }
  };

  const handleRemoveImage = (index: number) => {
    setPreviewImages(prev => prev.filter((_, i) => i !== index));
  };

  const handleClearAllImages = () => {
    setPreviewImages([]);
  };

  const handleAddMorePhotos = () => {
    fileInputRef.current?.click();
  };

  const handleAnalyzePhotos = () => {
    if (previewImages.length === 0) return;
    
    setIsLoading(true);
    setTimeout(() => {
      setIsLoading(false);
      console.log('Анализ', previewImages.length, 'фото');
    }, 2000);
  };

  return (
    <div className={styles.pageContainer}>
      <div className={styles.pageHeader}>
        <h1 className={styles.pageTitle}>
          <span className={styles.pageTitleEmoji}>🤖</span>
          <span className={styles.pageTitleText}>ИИ-подбор рецептов</span>
        </h1>
        <p className={styles.pageDescription}>
          Искусственный интеллект поможет найти идеальный рецепт из ваших продуктов
        </p>
      </div>

      {/* Вкладки */}
      <div className={styles.matchTabs}>
        <button 
          className={`${styles.matchTab} ${activeTab === 'ingredients' ? styles.active : ''}`}
          onClick={() => setActiveTab('ingredients')}
        >
          <span className={styles.tabIcon}>📝</span>
          <span className={styles.tabText}>Список продуктов</span>
          <span className={styles.tabBadge}>вручную</span>
        </button>
        <button 
          className={`${styles.matchTab} ${activeTab === 'photo' ? styles.active : ''}`}
          onClick={() => setActiveTab('photo')}
        >
          <span className={styles.tabIcon}>📸</span>
          <span className={styles.tabText}>Загрузить фото</span>
          <span className={styles.tabBadge}>AI распознавание</span>
        </button>
      </div>

      {/* Контент вкладок */}
      <div className={styles.matchContent}>
        {activeTab === 'ingredients' ? (
          <div className={styles.ingredientsPanel}>
            <div className={styles.panelHeader}>
              <h3>Выберите или добавьте продукты</h3>
              <p>Добавьте ингредиенты, которые у вас есть, и ИИ подберёт подходящие рецепты</p>
            </div>

            {/* Поле ввода */}
            <div className={styles.ingredientInput}>
              <input
                type="text"
                value={currentIngredient}
                onChange={(e) => setCurrentIngredient(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="Например: помидоры, сыр, курица..."
                className={styles.input}
              />
              <button 
                onClick={handleAddIngredient}
                className={styles.addButton}
                disabled={!currentIngredient.trim()}
              >
                <span>+</span> Добавить
              </button>
            </div>

            {/* Популярные ингредиенты */}
            <div className={styles.popularSection}>
              <h4>🔥 Популярные ингредиенты</h4>
              <div className={styles.popularGrid}>
                {popularIngredients.map(ing => (
                  <button
                    key={ing}
                    className={`${styles.popularChip} ${ingredients.includes(ing) ? styles.active : ''}`}
                    onClick={() => handleTogglePopular(ing)}
                  >
                    <span className={styles.chipIcon}>{ingredients.includes(ing) ? '✓' : '+'}</span>
                    {ing}
                  </button>
                ))}
              </div>
            </div>

            {/* Выбранные ингредиенты */}
            {ingredients.length > 0 && (
              <div className={styles.selectedSection}>
                <h4>✅ Выбранные продукты ({ingredients.length})</h4>
                <div className={styles.selectedGrid}>
                  {ingredients.map(ing => (
                    <div key={ing} className={styles.selectedChip}>
                      {ing}
                      <button 
                        onClick={() => handleRemoveIngredient(ing)}
                        className={styles.removeButton}
                      >
                        ✕
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Кнопка поиска */}
            <button 
              className={styles.findButton}
              onClick={handleFindRecipes}
              disabled={ingredients.length === 0 || isLoading}
            >
              {isLoading ? (
                <>
                  <span className={styles.spinnerSmall}></span>
                  ИИ подбирает рецепты...
                </>
              ) : (
                <>
                  <span>🔍</span> Найти рецепты
                </>
              )}
            </button>
          </div>
        ) : (
          <div className={styles.photoPanel}>
            <div className={styles.panelHeader}>
              <h3>Загрузите фото продуктов</h3>
              <p>ИИ проанализирует фото и определит доступные ингредиенты</p>
            </div>

            {/* Зона загрузки */}
            {previewImages.length === 0 ? (
              <div 
                className={`${styles.uploadZone} ${dragActive ? styles.dragActive : ''}`}
                onDragEnter={handleDrag}
                onDragLeave={handleDrag}
                onDragOver={handleDrag}
                onDrop={handleDrop}
                onClick={() => fileInputRef.current?.click()}
              >
                <div className={styles.uploadIcon}>📸🍲</div>
                <p className={styles.uploadText}>
                  {dragActive ? 'Отпустите файлы здесь' : 'Перетащите фото сюда'}
                </p>
                <p className={styles.uploadHint}>
                  или нажмите для выбора
                </p>
                <p className={styles.uploadFormats}>
                  Поддерживаются JPG, PNG до 10MB
                </p>
                <p className={styles.uploadMultiple}>
                  Можно загрузить несколько фото
                </p>
              </div>
            ) : (
              <div className={styles.previewSection}>
                {/* Заголовок с информацией о фото */}
                <div className={styles.previewHeader}>
                  <h4>
                    <span className={styles.previewIcon}>🖼️</span>
                    Загружено фото: {previewImages.length}
                  </h4>
                  <div className={styles.previewActions}>
                    <button 
                      className={styles.addMoreButton}
                      onClick={handleAddMorePhotos}
                    >
                      <span>+</span> Добавить
                    </button>
                    <button 
                      className={styles.clearAllButton}
                      onClick={handleClearAllImages}
                    >
                      <span>🗑️</span> Очистить
                    </button>
                  </div>
                </div>

                {/* Сетка с превью */}
                <div className={styles.previewGrid}>
                  {previewImages.map((image, index) => (
                    <div key={index} className={styles.previewItem}>
                      <img 
                        src={image} 
                        alt={`Продукт ${index + 1}`} 
                        className={styles.previewImage} 
                      />
                      <button 
                        className={styles.removeImageButton}
                        onClick={() => handleRemoveImage(index)}
                        title="Удалить фото"
                      >
                        ✕
                      </button>
                      <span className={styles.imageNumber}>{index + 1}</span>
                    </div>
                  ))}
                  
                  {/* Кнопка добавления внутри сетки */}
                  <div 
                    className={styles.addMoreItem}
                    onClick={handleAddMorePhotos}
                  >
                    <span className={styles.addMoreIcon}>+</span>
                    <span className={styles.addMoreText}>Добавить фото</span>
                  </div>
                </div>

                {/* Кнопка анализа */}
                <button 
                  className={styles.analyzeButton}
                  onClick={handleAnalyzePhotos}
                  disabled={isLoading}
                >
                  {isLoading ? (
                    <>
                      <span className={styles.spinnerSmall}></span>
                      ИИ анализирует фото...
                    </>
                  ) : (
                    <>
                      <span>🤖</span> Распознать продукты
                    </>
                  )}
                </button>
              </div>
            )}

            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              multiple
              onChange={handleFileInput}
              style={{ display: 'none' }}
            />

            {/* Подсказки */}
            <div className={styles.photoTips}>
              <div className={styles.tipItem}>
                <span className={styles.tipIcon}>💡</span>
                <span>Снимайте при хорошем освещении</span>
              </div>
              <div className={styles.tipItem}>
                <span className={styles.tipIcon}>📐</span>
                <span>Продукты должны быть видны чётко</span>
              </div>
              <div className={styles.tipItem}>
                <span className={styles.tipIcon}>🎯</span>
                <span>Разместите продукты на контрастном фоне</span>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default MatchPage;