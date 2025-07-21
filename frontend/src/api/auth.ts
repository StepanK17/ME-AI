import axios from 'axios';

const API_URL = '/api';

export async function login(email: string, password: string) {
  const res = await axios.post(`${API_URL}/auth/login`, { email, password });
  return res.data.token as string;
}

export async function register(email: string, name: string, password: string) {
  const res = await axios.post(`${API_URL}/auth/register`, { email, name, password });
  return res.data.token as string;
} 