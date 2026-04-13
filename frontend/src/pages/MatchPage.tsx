// src/pages/MatchPage.tsx
import React, { useState } from 'react';
import styles from './Pages.module.css';

const MatchPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'ingredients' | 'photo'>('ingredients');
  const [ingredients, setIngredients] = useState<string[]>([]);
  const [currentIngredient, setCurrentIngredient] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [previewImage, setPreviewImage] = useState<string | null>(null);
  const [dragActive, setDragActive] = useState(false);

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
      // Здесь будет реальный запрос к API
      console.log('Поиск рецептов с ингредиентами:', ingredients);
    }, 2000);
  };

  const handleImageUpload = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      setPreviewImage(e.target?.result as string);
    };
    reader.readAsDataURL(file);
    
    // Имитация анализа фото ИИ
    setIsLoading(true);
    setTimeout(() => {
      setIsLoading(false);
      // Здесь будет реальный запрос к ML сервису
      console.log('Анализ фото:', file.name);
    }, 2000);
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
    
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleImageUpload(e.dataTransfer.files[0]);
    }
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      handleImageUpload(e.target.files[0]);
    }
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
            <div 
              className={`${styles.uploadZone} ${dragActive ? styles.dragActive : ''} ${previewImage ? styles.hasPreview : ''}`}
              onDragEnter={handleDrag}
              onDragLeave={handleDrag}
              onDragOver={handleDrag}
              onDrop={handleDrop}
              onClick={() => !previewImage && document.getElementById('fileInput')?.click()}
            >
              {isLoading ? (
                <div className={styles.aiProcessing}>
                  <div className={styles.aiIcon}>🤖</div>
                  <p className={styles.aiText}>ИИ анализирует фото...</p>
                  <p className={styles.aiSubtext}>Это займёт несколько секунд</p>
                  <div className={styles.loadingBar}>
                    <div className={styles.loadingProgress}></div>
                  </div>
                </div>
              ) : previewImage ? (
                <div className={styles.previewContainer}>
                  <img src={previewImage} alt="Preview" className={styles.previewImage} />
                  <button 
                    className={styles.changePhotoButton}
                    onClick={() => {
                      setPreviewImage(null);
                      document.getElementById('fileInput')?.click();
                    }}
                  >
                    📸 Изменить фото
                  </button>
                </div>
              ) : (
                <>
                  <div className={styles.uploadIcon}>📸🍲</div>
                  <p className={styles.uploadText}>
                    {dragActive ? 'Отпустите файл здесь' : 'Перетащите фото сюда'}
                  </p>
                  <p className={styles.uploadHint}>
                    или нажмите для выбора
                  </p>
                  <p className={styles.uploadFormats}>
                    Поддерживаются JPG, PNG до 10MB
                  </p>
                </>
              )}
              <input
                id="fileInput"
                type="file"
                accept="image/*"
                onChange={handleFileInput}
                style={{ display: 'none' }}
              />
            </div>

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