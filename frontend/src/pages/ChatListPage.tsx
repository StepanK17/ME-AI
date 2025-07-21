import React from 'react';
import { Box, Paper } from '@mui/material';

const ChatListPage: React.FC = () => {
  return (
    <Paper sx={{ flex: 1, minHeight: 400, p: 3, borderRadius: 4, bgcolor: 'background.default', boxShadow: 0 }}>
      <Box sx={{ color: 'text.secondary', textAlign: 'center', mt: 10 }}>
        Выберите чат слева или создайте новый
      </Box>
    </Paper>
  );
};

export default ChatListPage; 