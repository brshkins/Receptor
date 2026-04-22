import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { Header } from './components/Layout/Header';
import { Toaster } from 'react-hot-toast';
import styles from './App.module.css';

// Обычные импорты вместо lazy
import HomePage from './pages/HomePage';
import RecipesPage from './pages/RecipesPage';
import RecipeDetailsPage from './pages/RecipeDetailsPage';
import FavoritesPage from './pages/FavoritesPage';
import MatchPage from './pages/MatchPage';
import ProfilePage from './pages/ProfilePage';
import AuthPage from './pages/AuthPage';

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <div className={styles.app}>
          <Header />
          <main className={styles.mainContent}>
            <Routes>
              <Route path="/" element={<HomePage />} />
              <Route path="/recipes" element={<RecipesPage />} />
              <Route path="/recipes/:id" element={<RecipeDetailsPage />} />
              <Route path="/favorites" element={<FavoritesPage />} />
              <Route path="/match" element={<MatchPage />} />
              <Route path="/profile" element={<ProfilePage />} />
              <Route path="/auth" element={<AuthPage />} />
            </Routes>
          </main>
          <Toaster 
            position="top-right"
            toastOptions={{
              style: {
                background: 'linear-gradient(135deg, #f6d365 0%, #fda085 100%)',
                color: 'white',
                fontWeight: 500,
                padding: '16px',
                borderRadius: '24px',
              },
            }}
          />
        </div>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;