import React from 'react';
import { AppBar, Toolbar, Typography, Button, Box, IconButton, Tooltip } from '@mui/material';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { useThemeMode } from '../context/ThemeContext';
import Brightness4Icon from '@mui/icons-material/Brightness4';
import Brightness7Icon from '@mui/icons-material/Brightness7';

const TopBar: React.FC = () => {
  const { logout } = useAuth();
  const navigate = useNavigate();
  const { dark, toggle } = useThemeMode();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <AppBar position="static" color="default" elevation={2} sx={{ mb: 3, borderRadius: 3, mt: 2 }}>
      <Toolbar>
        <Box sx={{ flexGrow: 1, display: 'flex', alignItems: 'center', gap: 2 }}>
          <Typography variant="h6" color="inherit" component="div" sx={{ fontWeight: 700, letterSpacing: 1 }}>
            <span style={{ color: '#10a37f' }}>me</span>Chat
          </Typography>
        </Box>
        <Tooltip title={dark ? 'Светлая тема' : 'Темная тема'}>
          <IconButton color="inherit" onClick={toggle} sx={{ mr: 1 }}>
            {dark ? <Brightness7Icon /> : <Brightness4Icon />}
          </IconButton>
        </Tooltip>
        <Button color="inherit" onClick={handleLogout} sx={{ fontWeight: 600 }}>
          Выйти
        </Button>
      </Toolbar>
    </AppBar>
  );
};

export default TopBar; 