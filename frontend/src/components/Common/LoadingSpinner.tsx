import React from 'react';
import styles from './Common.module.css';

interface LoadingSpinnerProps {
  text?: string;
}

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({ 
  text = 'Загрузка...' 
}) => {
  return (
    <div className={styles.spinnerContainer}>
      <div className={styles.spinner} />
      <p className={styles.spinnerText}>{text}</p>
    </div>
  );
};