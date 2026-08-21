import axios from 'axios';

const API_URL = 'http://localhost:8082';

// Настройка axios instance
const api = axios.create({
  baseURL: API_URL,
});

// Добавляем перехватчик, который берет токен из localStorage перед каждым запросом
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const authService = {
  signIn: async (data: any) => await axios.post(`${API_URL}/sign-in`, data),
  signUp: async (data: any) => await axios.post(`${API_URL}/sign-up`, data),
  
  // Workspaces
  createWorkspace: async (data: { name: string }) => {
    return await api.post('/Workspace', data);
  },
  getAllWorkspaces: async () => {
    return await api.get('/Workspaces');
  },
  getWorkspaceBlocks: async (id: string) => {
    return await api.get(`/Workspaces/${id}/blocks`);
  },
  updateWorkspace: async (id: string, data: { name: string }) => {
    return await api.patch(`/Workspaces/${id}`, data);
  },
  deleteWorkspace: async (id: string) => {
    return await api.delete(`/Workspaces/${id}`);
  },
  
  // Users
  deleteUser: async () => {
    return await api.delete('/Users/');
  },
  
  saveToken: (token: string) => localStorage.setItem('auth_token', token),
};
