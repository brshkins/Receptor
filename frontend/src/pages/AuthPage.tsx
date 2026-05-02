// src/pages/AuthPage.tsx
import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { getApiErrorMessage } from '../utils/apiError';
import styles from './Pages.module.css';

const AuthPage: React.FC = () => {
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: ''
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  
  const { login, register } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  
  const from = (location.state as any)?.from?.pathname || '/';

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    // Очищаем ошибку при вводе
    if (errors[name]) {
      setErrors(prev => ({ ...prev, [name]: '' }));
    }
  };

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};
    
    if (!isLogin && !formData.name.trim()) {
      newErrors.name = 'Введите имя';
    }
    
    if (!formData.email.trim()) {
      newErrors.email = 'Укажите эл. почту';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = 'Некорректный адрес почты';
    }
    
    if (!formData.password) {
      newErrors.password = 'Введите пароль';
    } else if (formData.password.length < 6) {
      newErrors.password = 'Пароль должен быть не менее 6 символов';
    }
    
    if (!isLogin && formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Пароли не совпадают';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) return;
    
    setIsLoading(true);
    
    try {
      if (isLogin) {
        await login(formData.email, formData.password);
      } else {
        await register(formData.name, formData.email, formData.password);
      }
      navigate(from, { replace: true });
    } catch (error: unknown) {
      setErrors({ submit: getApiErrorMessage(error, 'Произошла ошибка') });
    } finally {
      setIsLoading(false);
    }
  };

  const switchMode = () => {
    setIsLogin(!isLogin);
    setErrors({});
    setFormData({
      name: '',
      email: '',
      password: '',
      confirmPassword: ''
    });
  };

  return (
    <div className={styles.authContainer}>
      <div className={styles.authCard}>
        {/* Декоративные элементы */}
        <div className={styles.authDecoration}>
          <span className={styles.authEmoji}>🍳</span>
          <span className={styles.authEmoji}>🥘</span>
          <span className={styles.authEmoji}>🍲</span>
        </div>
        
        {/* Заголовок */}
        <div className={styles.authHeader}>
          <h1 className={styles.authTitle}>
            <span className={styles.authTitleEmoji}>
              {isLogin ? '🔐' : '📝'}
            </span>
            <span className={styles.authTitleText}>
              {isLogin ? 'Вход в аккаунт' : 'Регистрация'}
            </span>
          </h1>
          <p className={styles.authSubtitle}>
            {isLogin 
              ? 'Войдите, чтобы сохранять рецепты и получать персональные рекомендации' 
              : 'Создайте аккаунт и начните кулинарное путешествие'}
          </p>
        </div>
        
        {/* Ошибка отправки */}
        {errors.submit && (
          <div className={styles.authError}>
            <span className={styles.errorIcon}>⚠️</span>
            <span>{errors.submit}</span>
          </div>
        )}
        
        {/* Форма */}
        <form onSubmit={handleSubmit} className={styles.authForm}>
          {!isLogin && (
            <div className={styles.inputGroup}>
              <label htmlFor="name" className={styles.inputLabel}>
                <span className={styles.labelIcon}>👤</span>
                <span>Имя</span>
              </label>
              <div className={styles.inputWrapper}>
                <input
                  id="name"
                  name="name"
                  type="text"
                  value={formData.name}
                  onChange={handleChange}
                  placeholder="Введите ваше имя"
                  className={`${styles.authInput} ${errors.name ? styles.inputError : ''}`}
                  disabled={isLoading}
                />
                {errors.name && (
                  <span className={styles.inputErrorMessage}>{errors.name}</span>
                )}
              </div>
            </div>
          )}
          
          <div className={styles.inputGroup}>
            <label htmlFor="email" className={styles.inputLabel}>
              <span className={styles.labelIcon}>📧</span>
              <span>Эл. почта</span>
            </label>
            <div className={styles.inputWrapper}>
              <input
                id="email"
                name="email"
                type="email"
                value={formData.email}
                onChange={handleChange}
                  placeholder="pochta@primer.ru"
                className={`${styles.authInput} ${errors.email ? styles.inputError : ''}`}
                disabled={isLoading}
              />
              {errors.email && (
                <span className={styles.inputErrorMessage}>{errors.email}</span>
              )}
            </div>
          </div>
          
          <div className={styles.inputGroup}>
            <label htmlFor="password" className={styles.inputLabel}>
              <span className={styles.labelIcon}>🔒</span>
              <span>Пароль</span>
            </label>
            <div className={styles.inputWrapper}>
              <div className={styles.passwordWrapper}>
                <input
                  id="password"
                  name="password"
                  type={showPassword ? 'text' : 'password'}
                  value={formData.password}
                  onChange={handleChange}
                  placeholder={isLogin ? 'Введите пароль' : 'Придумайте пароль'}
                  className={`${styles.authInput} ${errors.password ? styles.inputError : ''}`}
                  disabled={isLoading}
                />
                <button
                  type="button"
                  className={styles.passwordToggle}
                  onClick={() => setShowPassword(!showPassword)}
                  tabIndex={-1}
                >
                  {showPassword ? '👁️' : '👁️‍🗨️'}
                </button>
              </div>
              {errors.password && (
                <span className={styles.inputErrorMessage}>{errors.password}</span>
              )}
            </div>
          </div>
          
          {!isLogin && (
            <div className={styles.inputGroup}>
              <label htmlFor="confirmPassword" className={styles.inputLabel}>
                <span className={styles.labelIcon}>✓</span>
                <span>Подтверждение пароля</span>
              </label>
              <div className={styles.inputWrapper}>
                <div className={styles.passwordWrapper}>
                  <input
                    id="confirmPassword"
                    name="confirmPassword"
                    type={showPassword ? 'text' : 'password'}
                    value={formData.confirmPassword}
                    onChange={handleChange}
                    placeholder="Повторите пароль"
                    className={`${styles.authInput} ${errors.confirmPassword ? styles.inputError : ''}`}
                    disabled={isLoading}
                  />
                </div>
                {errors.confirmPassword && (
                  <span className={styles.inputErrorMessage}>{errors.confirmPassword}</span>
                )}
              </div>
            </div>
          )}
          
          <button
            type="submit"
            className={styles.authSubmitButton}
            disabled={isLoading}
          >
            {isLoading ? (
              <>
                <span className={styles.spinnerSmall}></span>
                <span>{isLogin ? 'Входим...' : 'Регистрируемся...'}</span>
              </>
            ) : (
              <>
                <span>{isLogin ? '🔓' : '✨'}</span>
                <span>{isLogin ? 'Войти' : 'Зарегистрироваться'}</span>
              </>
            )}
          </button>
        </form>
        
        {/* Переключение режима */}
        <div className={styles.authSwitch}>
          <span>
            {isLogin ? 'Нет аккаунта?' : 'Уже есть аккаунт?'}
          </span>
          <button
            type="button"
            className={styles.authSwitchButton}
            onClick={switchMode}
            disabled={isLoading}
          >
            {isLogin ? 'Зарегистрироваться' : 'Войти'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default AuthPage;