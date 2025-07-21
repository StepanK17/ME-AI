import React, { useEffect, useState, useRef } from 'react';
import { Box, Button, TextField, Typography, Paper, List, CircularProgress, InputAdornment, IconButton } from '@mui/material';
import { useParams } from 'react-router-dom';
import { getMessages, sendMessage } from '../api/chat';
import ChatBubble from '../components/ChatBubble';
import SendIcon from '@mui/icons-material/Send';

const ChatPage: React.FC = () => {
  const { id } = useParams();
  const conversationId = Number(id);
  const [message, setMessage] = useState('');
  const [messages, setMessages] = useState<{ id: number; content: string; role: string }[]>([]);
  const [loading, setLoading] = useState(true);
  const [sending, setSending] = useState(false);
  const [error, setError] = useState('');
  const listRef = useRef<HTMLUListElement>(null);

  useEffect(() => {
    if (!conversationId) return;
    setLoading(true);
    getMessages(conversationId)
      .then(data => setMessages(Array.isArray(data) ? data : []))
      .catch(() => setError('Ошибка загрузки сообщений'))
      .finally(() => setLoading(false));
  }, [conversationId]);

  useEffect(() => {
    if (listRef.current) {
      listRef.current.scrollTop = listRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSend = async () => {
    if (!message.trim() || !conversationId) return;
    setSending(true);
    setError('');
    try {
      setMessages(prev => [...prev, { id: Date.now(), content: message, role: 'user' }]);
      setMessage('');
      const res = await sendMessage(conversationId, message);
      setMessages(prev => [...prev, { id: Date.now() + 1, content: res.message, role: 'assistant' }]);
    } catch {
      setError('Ошибка отправки сообщения');
    } finally {
      setSending(false);
    }
  };

  return (
    <Box component={Paper} elevation={3} sx={{ p: { xs: 1, sm: 3 }, maxWidth: 700, mx: 'auto', minHeight: 400, display: 'flex', flexDirection: 'column', borderRadius: 4 }}>
      <Typography variant="h6" mb={2} align="center">Чат #{id}</Typography>
      {loading ? <CircularProgress sx={{ display: 'block', mx: 'auto', my: 4 }} /> : (
        <List ref={listRef} sx={{ flex: 1, overflowY: 'auto', mb: 2, maxHeight: { xs: 300, sm: 400 }, px: 0 }}>
          {messages.length === 0 && <Typography align="center" color="text.secondary" mt={4}>Нет сообщений</Typography>}
          {messages.map((msg, idx) => (
            <li key={msg.id + '-' + idx} style={{ listStyle: 'none' }}>
              <ChatBubble content={msg.content} role={msg.role as 'user' | 'assistant'} />
            </li>
          ))}
        </List>
      )}
      {error && <Typography color="error" align="center">{error}</Typography>}
      <Box component="form" onSubmit={e => { e.preventDefault(); handleSend(); }} sx={{ display: 'flex', gap: 1, mt: 1 }}>
        <TextField
          value={message}
          onChange={e => setMessage(e.target.value)}
          placeholder="Введите сообщение..."
          fullWidth
          size="small"
          disabled={sending}
          autoFocus
          sx={{ borderRadius: 3, bgcolor: 'background.paper' }}
          InputProps={{
            endAdornment: (
              <InputAdornment position="end">
                <IconButton type="submit" color="primary" disabled={sending || !message.trim()}>
                  <SendIcon />
                </IconButton>
              </InputAdornment>
            ),
          }}
        />
      </Box>
    </Box>
  );
};

export default ChatPage; 