import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { authService } from '../services/authService';
import { getApiErrorMessage } from '../utils/apiError';

interface User {
  id: number | string;
  name: string;
  email: string;
}

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const initAuth = async () => {
      const token = localStorage.getItem('token');
      if (token) {
        try {
          const userData = await authService.getMe();
          setUser(userData);
        } catch (error) {
          console.error('Failed to get user:', error);
          localStorage.removeItem('token');
          setUser(null);
        }
      }
      setIsLoading(false);
    };

    initAuth();
  }, []);

  const login = async (email: string, password: string) => {
    await authService.login({ email, password });
    try {
      const userData = await authService.getMe();
      setUser(userData);
    } catch (error) {
      authService.logout();
      setUser(null);
      const msg = getApiErrorMessage(
        error,
        error instanceof Error ? error.message : 'Не удалось загрузить профиль'
      );
      throw new Error(msg);
    }
  };

  const register = async (name: string, email: string, password: string) => {
    await authService.register({ name, email, password });
    try {
      const userData = await authService.getMe();
      setUser(userData);
    } catch (error) {
      authService.logout();
      setUser(null);
      const msg = getApiErrorMessage(
        error,
        error instanceof Error ? error.message : 'Не удалось загрузить профиль'
      );
      throw new Error(msg);
    }
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
    window.location.href = '/';
  };

  return (
    <AuthContext.Provider value={{
      user,
      isLoading,
      login,
      register,
      logout,
      isAuthenticated: !!user,
    }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error('useAuth must be used within AuthProvider');
  return context;
};
