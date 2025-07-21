import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { CssBaseline, Container, Box } from '@mui/material';
import LoginPage from './pages/LoginPage.tsx';
import RegisterPage from './pages/RegisterPage.tsx';
import ChatListPage from './pages/ChatListPage.tsx';
import ChatPage from './pages/ChatPage.tsx';
import TopBar from './components/TopBar';
import ChatLayout from './components/ChatLayout';
import { useAuth } from './context/AuthContext';

const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { token } = useAuth();
  if (!token) return <Navigate to="/login" replace />;
  return <>{children}</>;
};

function App() {
  return (
    <Router>
      <CssBaseline />
      <Container maxWidth="md" sx={{ minHeight: '100vh', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
        <Box sx={{ flex: 1, display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route path="/chats" element={<PrivateRoute><TopBar /><ChatLayout><ChatListPage /></ChatLayout></PrivateRoute>} />
            <Route path="/chat/:id" element={<PrivateRoute><TopBar /><ChatLayout><ChatPage /></ChatLayout></PrivateRoute>} />
            <Route path="*" element={<Navigate to="/login" replace />} />
          </Routes>
        </Box>
      </Container>
    </Router>
  );
}

export default App;
