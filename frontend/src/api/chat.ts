import axios from 'axios';

const API_URL = '/api';

export async function getChats() {
  const res = await axios.get(`${API_URL}/conversations`);
  return res.data as { id: number; title: string }[];
}

export async function createChat(title: string) {
  const res = await axios.post(`${API_URL}/conversations/create`, { title });
  return res.data as { id: number; title: string };
}

export async function getMessages(conversationId: number) {
  const res = await axios.get(`${API_URL}/messages?conversation_id=${conversationId}`);
  return res.data as { id: number; content: string; role: string }[];
}

export async function sendMessage(conversationId: number, message: string) {
  const res = await axios.post(`${API_URL}/chat`, { conversation_id: conversationId, message });
  return res.data as { message: string; timestamp: string };
}

export async function deleteChat(id: number) {
  await axios.post('/api/conversations/delete', { id });
}

export async function renameChat(id: number, title: string) {
  await axios.post('/api/conversations/rename', { id, title });
} 