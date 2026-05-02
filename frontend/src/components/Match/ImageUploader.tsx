import React, { useState, useRef } from 'react';
import { uploadService } from '../../services/uploadService';
import type { MatchResponse } from '../../types';
import styles from './Match.module.css';

interface ImageUploaderProps {
  onUploadSuccess: (items: MatchResponse[]) => void;
  /** При ошибке не вызывать onUploadSuccess; показать сообщение без «пустого» успеха. */
  onUploadError?: (message: string) => void;
}

export const ImageUploader: React.FC<ImageUploaderProps> = ({
  onUploadSuccess,
  onUploadError,
}) => {
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
      
      const items = await uploadService.uploadAndMatch(file);
      onUploadSuccess(items);
    } catch (error) {
      console.error('Upload failed:', error);
      setAiStatus('');
      const msg = 'Не удалось загрузить фото. Попробуйте ещё раз.';
      if (onUploadError) {
        onUploadError(msg);
      } else {
        alert(msg);
      }
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
          <img src={preview} alt="Предпросмотр" className={styles.previewImage} />
        </div>
      )}
    </div>
  );
};