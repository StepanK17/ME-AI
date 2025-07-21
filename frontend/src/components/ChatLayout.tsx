import React, { useEffect, useState } from 'react';
import { Box, CircularProgress } from '@mui/material';
import ChatSidebar from './ChatSidebar';
import { getChats, createChat, deleteChat } from '../api/chat';
import { useNavigate, useLocation, useParams } from 'react-router-dom';

interface ChatLayoutProps {
  children: React.ReactNode;
}

const ChatLayout: React.FC<ChatLayoutProps> = ({ children }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const params = useParams();
  const [chats, setChats] = useState<{ id: number; title: string }[]>([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState('');
  const [selectedId, setSelectedId] = useState<number | null>(null);

  useEffect(() => {
    setLoading(true);
    getChats()
      .then(data => setChats(Array.isArray(data) ? data : []))
      .catch(() => setChats([]))
      .finally(() => setLoading(false));
    // eslint-disable-next-line
  }, []);

  useEffect(() => {
    // Выделяем чат, если в URL есть id
    if (params.id) {
      setSelectedId(Number(params.id));
    }
  }, [params.id]);

  const handleNewChat = async () => {
    setCreating(true);
    setError('');
    try {
      const chat = await createChat('Новый чат');
      setChats(prev => [chat, ...prev]);
      setSelectedId(chat.id);
      navigate(`/chat/${chat.id}`);
    } catch {
      setError('Ошибка создания чата');
    } finally {
      setCreating(false);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteChat(id);
      setChats(prev => prev.filter(chat => chat.id !== id));
      if (selectedId === id && chats.length > 1) {
        const next = chats.find(c => c.id !== id);
        if (next) {
          setSelectedId(next.id);
          navigate(`/chat/${next.id}`);
        }
      }
      if (chats.length === 1) {
        setSelectedId(null);
        navigate('/chats');
      }
    } catch {
      setError('Ошибка удаления чата');
    }
  };

  const handleSelect = (id: number) => {
    setSelectedId(id);
    navigate(`/chat/${id}`);
  };

  const refreshChats = () => {
    setLoading(true);
    getChats()
      .then(data => setChats(Array.isArray(data) ? data : []))
      .catch(() => setChats([]))
      .finally(() => setLoading(false));
  };

  return (
    <Box sx={{ display: 'flex', flexDirection: { xs: 'column', sm: 'row' }, gap: 3, minHeight: 500 }}>
      <ChatSidebar
        chats={chats}
        loading={loading}
        selectedId={selectedId ?? undefined}
        onSelect={handleSelect}
        onNewChat={handleNewChat}
        onDelete={handleDelete}
        onRefresh={refreshChats}
      />
      <Box sx={{ flex: 1, minWidth: 0 }}>
        {children}
      </Box>
    </Box>
  );
};

export default ChatLayout; 