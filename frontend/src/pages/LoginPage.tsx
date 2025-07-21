import React, { useState } from 'react';
import { Box, Button, TextField, Typography, Link, Paper } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { login } from '../api/auth';
import { useAuth } from '../context/AuthContext';

const LoginPage: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const { setToken } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    if (!email || !password) {
      setError('Введите email и пароль');
      return;
    }
    try {
      const token = await login(email, password);
      setToken(token);
      navigate('/chats');
    } catch (err: any) {
      setError(err?.response?.data || 'Ошибка входа');
    }
  };

  return (
    <Box component={Paper} elevation={3} sx={{ p: 4, maxWidth: 400, mx: 'auto' }}>
      <Typography variant="h5" mb={2} align="center">Вход</Typography>
      <form onSubmit={handleSubmit}>
        <TextField
          label="Email"
          type="email"
          value={email}
          onChange={e => setEmail(e.target.value)}
          fullWidth
          margin="normal"
          required
        />
        <TextField
          label="Пароль"
          type="password"
          value={password}
          onChange={e => setPassword(e.target.value)}
          fullWidth
          margin="normal"
          required
        />
        {error && <Typography color="error" variant="body2">{error}</Typography>}
        <Button type="submit" variant="contained" color="primary" fullWidth sx={{ mt: 2 }}>
          Войти
        </Button>
      </form>
      <Box mt={2} textAlign="center">
        <Link href="/register" underline="hover">Нет аккаунта? Зарегистрироваться</Link>
      </Box>
    </Box>
  );
};

export default LoginPage; 