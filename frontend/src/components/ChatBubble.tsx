import React from 'react';
import { Box, Typography, Avatar, useTheme } from '@mui/material';

interface ChatBubbleProps {
  content: string;
  role: 'user' | 'assistant';
}

const ChatBubble: React.FC<ChatBubbleProps> = ({ content, role }) => {
  const theme = useTheme();
  const isUser = role === 'user';
  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: isUser ? 'row-reverse' : 'row',
        alignItems: 'flex-end',
        mb: 1.5,
      }}
    >
      <Avatar
        sx={{
          bgcolor: isUser ? theme.palette.primary.main : theme.palette.grey[700],
          width: 36,
          height: 36,
          fontSize: 18,
          ml: isUser ? 2 : 0,
          mr: isUser ? 0 : 2,
        }}
      >
        {isUser ? 'Я' : <span role="img" aria-label="AI">🤖</span>}
      </Avatar>
      <Box
        sx={{
          bgcolor: isUser ? theme.palette.primary.light : theme.palette.background.paper,
          color: isUser ? theme.palette.primary.contrastText : theme.palette.text.primary,
          borderRadius: 3,
          px: 2,
          py: 1.2,
          maxWidth: '70%',
          boxShadow: 1,
          borderTopRightRadius: isUser ? 6 : 24,
          borderTopLeftRadius: isUser ? 24 : 6,
        }}
      >
        <Typography variant="body1" sx={{ whiteSpace: 'pre-line' }}>{content}</Typography>
      </Box>
    </Box>
  );
};

export default ChatBubble; 