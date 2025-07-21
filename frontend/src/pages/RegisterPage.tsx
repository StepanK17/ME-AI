import React, { useState } from 'react';
import { Box, Button, TextField, Typography, Link, Paper } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { register } from '../api/auth';
import { useAuth } from '../context/AuthContext';

const RegisterPage: React.FC = () => {
  const [email, setEmail] = useState('');
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();
  const { setToken } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    if (!email || !password || !name) {
      setError('Заполните все поля');
      return;
    }
    try {
      const token = await register(email, name, password);
      setToken(token);
      navigate('/chats');
    } catch (err: any) {
      setError(err?.response?.data || 'Ошибка регистрации');
    }
  };

  return (
    <Box component={Paper} elevation={3} sx={{ p: 4, maxWidth: 400, mx: 'auto' }}>
      <Typography variant="h5" mb={2} align="center">Регистрация</Typography>
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
          label="Имя"
          value={name}
          onChange={e => setName(e.target.value)}
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
          Зарегистрироваться
        </Button>
      </form>
      <Box mt={2} textAlign="center">
        <Link href="/login" underline="hover">Уже есть аккаунт? Войти</Link>
      </Box>
    </Box>
  );
};

export default RegisterPage; 