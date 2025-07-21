import React, { useState } from 'react';
import { Box, List, ListItem, ListItemButton, ListItemAvatar, Avatar, ListItemText, Typography, Button, Divider, CircularProgress, IconButton, TextField } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import { renameChat } from '../api/chat';

interface ChatSidebarProps {
  chats: { id: number; title: string }[];
  loading: boolean;
  selectedId?: number;
  onSelect: (id: number) => void;
  onNewChat: () => void;
  onDelete: (id: number) => void;
  onRefresh?: () => void;
}

const ChatSidebar: React.FC<ChatSidebarProps> = ({ chats, loading, selectedId, onSelect, onNewChat, onDelete, onRefresh }) => {
  const safeChats = Array.isArray(chats) ? chats : [];
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editValue, setEditValue] = useState('');
  const [renaming, setRenaming] = useState(false);

  const handleEdit = (chat: { id: number; title: string }) => {
    setEditingId(chat.id);
    setEditValue(chat.title);
  };

  const handleRename = async (chat: { id: number; title: string }) => {
    if (!editValue.trim() || editValue === chat.title) {
      setEditingId(null);
      return;
    }
    setRenaming(true);
    try {
      await renameChat(chat.id, editValue.trim());
      setEditingId(null);
      if (typeof onRefresh === 'function') onRefresh();
    } finally {
      setRenaming(false);
    }
  };

  return (
    <Box sx={{ width: { xs: '100%', sm: 300 }, bgcolor: 'background.paper', borderRadius: 3, p: 2, boxShadow: 2, minHeight: 400, maxHeight: 600, display: 'flex', flexDirection: 'column' }}>
      <Button
        variant="contained"
        color="primary"
        startIcon={<AddIcon />}
        onClick={onNewChat}
        sx={{ mb: 2, borderRadius: 2, fontWeight: 600 }}
        fullWidth
      >
        Новый чат
      </Button>
      <Divider sx={{ mb: 1 }} />
      {loading ? <CircularProgress sx={{ mx: 'auto', my: 4 }} /> : (
        <List sx={{ flex: 1, overflowY: 'auto', pr: 1 }}>
          {safeChats.length === 0 && <Typography align="center" color="text.secondary" mt={4}>Нет чатов</Typography>}
          {safeChats.map(chat => (
            <ListItem
              key={chat.id}
              disablePadding
              sx={{ borderRadius: 2, mb: 0.5, bgcolor: selectedId === chat.id ? 'primary.light' : 'transparent', transition: 'background 0.2s' }}
            >
              <ListItemButton onClick={() => onSelect(chat.id)} sx={{ borderRadius: 2 }} selected={selectedId === chat.id}>
                <ListItemAvatar>
                  <Avatar sx={{ bgcolor: 'primary.main' }}>{chat.title[0]?.toUpperCase() || 'C'}</Avatar>
                </ListItemAvatar>
                {editingId === chat.id ? (
                  <TextField
                    value={editValue}
                    onChange={e => setEditValue(e.target.value)}
                    onBlur={() => handleRename(chat)}
                    onKeyDown={e => {
                      if (e.key === 'Enter') handleRename(chat);
                      if (e.key === 'Escape') setEditingId(null);
                    }}
                    size="small"
                    autoFocus
                    disabled={renaming}
                    sx={{ minWidth: 80, maxWidth: 120, mr: 1 }}
                  />
                ) : (
                  <ListItemText primary={chat.title} />
                )}
                <IconButton size="small" onClick={e => { e.stopPropagation(); handleEdit(chat); }} sx={{ ml: 1 }}>
                  <EditIcon fontSize="small" />
                </IconButton>
                <Button size="small" color="error" onClick={e => { e.stopPropagation(); onDelete(chat.id); }} sx={{ ml: 1, minWidth: 0, px: 1 }}>
                  ×
                </Button>
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      )}
    </Box>
  );
};

export default ChatSidebar; 