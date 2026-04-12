import React from 'react';
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import './App.module.css';

// Страницы (пока заглушки, но уже с нормальным содержимым)
const HomePage = () => (
  <div style={{ textAlign: 'center', padding: '2rem' }}>
    <h1 style={{ fontSize: '2.5rem', color: '#333', marginBottom: '1rem' }}>
      🍳 Добро пожаловать в Receptor AI
    </h1>
    <p style={{ fontSize: '1.2rem', color: '#666', marginBottom: '2rem' }}>
      Придумаем рецепт из того, что есть в холодильнике
    </p>
    <div style={{ display: 'flex', gap: '1rem', justifyContent: 'center' }}>
      <Link to="/match" style={{
        padding: '1rem 2rem',
        background: '#2563eb',
        color: 'white',
        textDecoration: 'none',
        borderRadius: '8px',
        fontSize: '1.1rem'
      }}>
        🔍 Подобрать рецепт
      </Link>
      <Link to="/recipes" style={{
        padding: '1rem 2rem',
        background: '#10b981',
        color: 'white',
        textDecoration: 'none',
        borderRadius: '8px',
        fontSize: '1.1rem'
      }}>
        📖 Все рецепты
      </Link>
    </div>
  </div>
);

const RecipesPage = () => (
  <div>
    <h1 style={{ fontSize: '2rem', marginBottom: '1.5rem' }}>📖 Все рецепты</h1>
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '1.5rem' }}>
      {/* Здесь будет список рецептов */}
      <div style={{ border: '1px solid #e5e7eb', borderRadius: '12px', padding: '1rem' }}>
        <div style={{ height: '150px', background: '#f3f4f6', borderRadius: '8px', marginBottom: '1rem' }}></div>
        <h3>Паста Карбонара</h3>
        <p style={{ color: '#6b7280' }}>⏱️ 30 мин • Средне</p>
      </div>
      <div style={{ border: '1px solid #e5e7eb', borderRadius: '12px', padding: '1rem' }}>
        <div style={{ height: '150px', background: '#f3f4f6', borderRadius: '8px', marginBottom: '1rem' }}></div>
        <h3>Овощной суп</h3>
        <p style={{ color: '#6b7280' }}>⏱️ 45 мин • Легко</p>
      </div>
      <div style={{ border: '1px solid #e5e7eb', borderRadius: '12px', padding: '1rem' }}>
        <div style={{ height: '150px', background: '#f3f4f6', borderRadius: '8px', marginBottom: '1rem' }}></div>
        <h3>Курица с рисом</h3>
        <p style={{ color: '#6b7280' }}>⏱️ 40 мин • Средне</p>
      </div>
    </div>
  </div>
);

const FavoritesPage = () => (
  <div>
    <h1 style={{ fontSize: '2rem', marginBottom: '1.5rem' }}>❤️ Избранные рецепты</h1>
    <p style={{ color: '#6b7280' }}>Здесь будут ваши любимые рецепты</p>
  </div>
);

const MatchPage = () => (
  <div>
    <h1 style={{ fontSize: '2rem', marginBottom: '1.5rem' }}>🔍 Подбор рецепта</h1>
    
    <div style={{ marginBottom: '2rem' }}>
      <h2 style={{ fontSize: '1.3rem', marginBottom: '1rem' }}>Выберите ингредиенты:</h2>
      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem', marginBottom: '1.5rem' }}>
        {['Помидоры', 'Сыр', 'Курица', 'Макароны', 'Лук', 'Чеснок', 'Сливки', 'Яйца'].map(ing => (
          <label key={ing} style={{
            padding: '0.5rem 1rem',
            background: '#f3f4f6',
            borderRadius: '20px',
            cursor: 'pointer',
            display: 'flex',
            alignItems: 'center',
            gap: '0.5rem'
          }}>
            <input type="checkbox" /> {ing}
          </label>
        ))}
      </div>
      
      <button style={{
        padding: '1rem 2rem',
        background: '#2563eb',
        color: 'white',
        border: 'none',
        borderRadius: '8px',
        fontSize: '1.1rem',
        cursor: 'pointer'
      }}>
        🔍 Найти рецепты
      </button>
    </div>

    <div>
      <h2 style={{ fontSize: '1.3rem', marginBottom: '1rem' }}>Или загрузите фото продуктов:</h2>
      <div style={{
        border: '2px dashed #d1d5db',
        borderRadius: '12px',
        padding: '2rem',
        textAlign: 'center',
        cursor: 'pointer'
      }}>
        <p style={{ fontSize: '2rem', marginBottom: '0.5rem' }}>📸</p>
        <p>Нажмите для загрузки фото</p>
      </div>
    </div>
  </div>
);

const ProfilePage = () => (
  <div>
    <h1 style={{ fontSize: '2rem', marginBottom: '1.5rem' }}>👤 Профиль</h1>
    
    <div style={{ display: 'flex', gap: '2rem', alignItems: 'flex-start' }}>
      <div style={{
        width: '150px',
        height: '150px',
        background: '#e5e7eb',
        borderRadius: '50%',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        fontSize: '3rem'
      }}>
        👤
      </div>
      
      <div>
        <h2 style={{ marginBottom: '1rem' }}>Имя пользователя</h2>
        <p style={{ color: '#6b7280', marginBottom: '0.5rem' }}>email@example.com</p>
        <p style={{ color: '#6b7280' }}>
          <a href="#" style={{ color: '#2563eb' }}>📱 Telegram: @username</a>
        </p>
      </div>
    </div>
  </div>
);

const AuthPage = () => (
  <div style={{ maxWidth: '400px', margin: '0 auto' }}>
    <h1 style={{ fontSize: '2rem', marginBottom: '1.5rem', textAlign: 'center' }}>🔐 Вход</h1>
    
    <form style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
      <input
        type="email"
        placeholder="Email"
        style={{
          padding: '0.75rem',
          border: '1px solid #d1d5db',
          borderRadius: '8px',
          fontSize: '1rem'
        }}
      />
      <input
        type="password"
        placeholder="Пароль"
        style={{
          padding: '0.75rem',
          border: '1px solid #d1d5db',
          borderRadius: '8px',
          fontSize: '1rem'
        }}
      />
      <button style={{
        padding: '0.75rem',
        background: '#2563eb',
        color: 'white',
        border: 'none',
        borderRadius: '8px',
        fontSize: '1rem',
        cursor: 'pointer'
      }}>
        Войти
      </button>
    </form>
    
    <p style={{ textAlign: 'center', marginTop: '1rem', color: '#6b7280' }}>
      Нет аккаунта? <a href="#" style={{ color: '#2563eb' }}>Зарегистрироваться</a>
    </p>
  </div>
);

function App() {
  return (
    <BrowserRouter>
      <div style={{ minHeight: '100vh', background: '#f9fafb' }}>
        {/* Шапка с навигацией */}
        <nav style={{
          padding: '1rem 2rem',
          background: 'white',
          boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
          display: 'flex',
          gap: '2rem',
          flexWrap: 'wrap',
          alignItems: 'center',
          position: 'sticky',
          top: 0,
          zIndex: 10
        }}>
          <Link to="/" style={{ 
            fontSize: '1.5rem', 
            fontWeight: 'bold',
            color: '#2563eb',
            textDecoration: 'none',
            marginRight: '1rem'
          }}>
            🍳 Receptor
          </Link>
          
          <div style={{ display: 'flex', gap: '1.5rem', flex: 1 }}>
            <Link to="/" style={{ color: '#4b5563', textDecoration: 'none' }}>Главная</Link>
            <Link to="/recipes" style={{ color: '#4b5563', textDecoration: 'none' }}>Рецепты</Link>
            <Link to="/favorites" style={{ color: '#4b5563', textDecoration: 'none' }}>Избранное</Link>
            <Link to="/match" style={{ color: '#4b5563', textDecoration: 'none' }}>Подбор</Link>
          </div>
          
          <div style={{ display: 'flex', gap: '1rem' }}>
            <Link to="/profile" style={{ color: '#4b5563', textDecoration: 'none' }}>Профиль</Link>
            <Link to="/auth" style={{ 
              padding: '0.5rem 1rem',
              background: '#2563eb',
              color: 'white',
              textDecoration: 'none',
              borderRadius: '6px'
            }}>
              Войти
            </Link>
          </div>
        </nav>

        {/* Основной контент */}
        <main style={{ padding: '2rem', maxWidth: '1200px', margin: '0 auto' }}>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/recipes" element={<RecipesPage />} />
            <Route path="/favorites" element={<FavoritesPage />} />
            <Route path="/match" element={<MatchPage />} />
            <Route path="/profile" element={<ProfilePage />} />
            <Route path="/auth" element={<AuthPage />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}

export default App;