import { createTheme } from '@mui/material/styles';

export const lightTheme = createTheme({
  palette: {
    mode: 'light',
    primary: { main: '#10a37f' },
    background: { default: '#f5f6fa', paper: '#fff' },
  },
  shape: { borderRadius: 14 },
  typography: { fontFamily: 'Inter, Arial, sans-serif' },
});

export const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: { main: '#10a37f' },
    background: { default: '#181a20', paper: '#23272f' },
  },
  shape: { borderRadius: 14 },
  typography: { fontFamily: 'Inter, Arial, sans-serif' },
}); 