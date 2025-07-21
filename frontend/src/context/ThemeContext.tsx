import React, { createContext, useContext, useState } from 'react';
import { ThemeProvider } from '@mui/material/styles';
import { lightTheme, darkTheme } from '../theme';

const ThemeContext = createContext<{ dark: boolean; toggle: () => void }>({ dark: false, toggle: () => {} });

export const useThemeMode = () => useContext(ThemeContext);

export const CustomThemeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [dark, setDark] = useState(() => localStorage.getItem('theme') === 'dark');
  const toggle = () => {
    setDark(d => {
      localStorage.setItem('theme', !d ? 'dark' : 'light');
      return !d;
    });
  };
  return (
    <ThemeContext.Provider value={{ dark, toggle }}>
      <ThemeProvider theme={dark ? darkTheme : lightTheme}>{children}</ThemeProvider>
    </ThemeContext.Provider>
  );
}; 