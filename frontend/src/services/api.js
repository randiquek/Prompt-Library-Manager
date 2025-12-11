import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api'

// axios instance with base config

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// add token to requests if exists

api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Auth api

export const login = async (username, password) => {
    const response = await api.post('/login', { username, password });
    return response.data;
};

// Prompts api

export const getAllPrompts = async () => {
    const response = await api.get('/prompts');
    return response.data;
};

export const createPrompt = async (promptData) => {
    const response = await api.post('/prompts', promptData);
    return response.data;
};

export const updatePrompt = async (promptData) => {
    const response = await api.put(`/prompts/${id}`, promptData);
    return response.data;
};

export const deletePrompt = async (id) => {
    const response = await api.delete(`/prompts/${id}`);
    return response.data;
};

// Audit logs api

export const getAuditLogs = async () => {
    const response = await api.get('/audit-logs');
    return response.data;
};


export default api;