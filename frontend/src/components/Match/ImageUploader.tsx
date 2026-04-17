import React, { useState, useRef } from 'react';
import { uploadService } from '../../services/uploadService';
import styles from './Match.module.css';

interface ImageUploaderProps {
  onUploadSuccess: (matches: any[]) => void;
}

export const ImageUploader: React.FC<ImageUploaderProps> = ({ onUploadSuccess }) => {
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [preview, setPreview] = useState<string | null>(null);
  const [aiStatus, setAiStatus] = useState('');
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFile = async (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => setPreview(e.target?.result as string);
    reader.readAsDataURL(file);

    setIsUploading(true);
    setAiStatus('🤖 ИИ анализирует фото...');
    
    try {
      setTimeout(() => setAiStatus('🔍 Распознаём ингредиенты...'), 1000);
      setTimeout(() => setAiStatus('📚 Подбираем рецепты...'), 2000);
      
      const result = await uploadService.uploadAndMatch(file);
      // Исправление: проверяем, есть ли matches
      onUploadSuccess(result.matches || []);
    } catch (error) {
      console.error('Upload failed:', error);
      setAiStatus('');
      // Показываем ошибку пользователю
      alert('Не удалось загрузить фото. Попробуйте ещё раз.');
    } finally {
      setIsUploading(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    const file = e.dataTransfer.files[0];
    if (file && file.type.startsWith('image/')) {
      handleFile(file);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      handleFile(file);
    }
  };

  if (isUploading) {
    return (
      <div className={styles.aiProcessing}>
        <div className={styles.aiIcon}>🤖</div>
        <p className={styles.aiText}>{aiStatus}</p>
        <p className={styles.aiSubtext}>Это займёт несколько секунд</p>
      </div>
    );
  }

  return (
    <div>
      <div
        className={`${styles.uploadZone} ${isDragging ? styles.dragging : ''}`}
        onDrop={handleDrop}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onClick={handleClick}
      >
        <div className={styles.uploadIcon}>📸🍲</div>
        <p className={styles.uploadText}>
          {isDragging ? 'Отпустите файл здесь' : 'Загрузите фото блюда'}
        </p>
        <p className={styles.uploadHint}>
          Перетащите файл или нажмите для выбора
        </p>
        <p className={styles.uploadHint} style={{ marginTop: '8px' }}>
          Поддерживаются JPG, PNG до 10MB
        </p>
        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          onChange={handleFileInput}
          style={{ display: 'none' }}
        />
      </div>
      
      {preview && !isUploading && (
        <div style={{ textAlign: 'center', marginTop: '24px' }}>
          <img src={preview} alt="Preview" className={styles.previewImage} />
        </div>
      )}
    </div>
  );
};