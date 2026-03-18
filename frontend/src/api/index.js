import api from './axios';

export const authApi = {
  register: async (data) => {
    const response = await api.post('/auth/register', data);
    return response.data;
  },

  login: async (data) => {
    const response = await api.post('/auth/login', data);
    return response.data;
  },

  getCurrentUser: async () => {
    const response = await api.get('/auth/me');
    return response.data;
  },

  refreshToken: async (refreshToken) => {
    const response = await api.post('/auth/refresh', { refresh_token: refreshToken });
    return response.data;
  },
};

export const leaveApi = {
  createLeave: async (data) => {
    const response = await api.post('/leaves', data);
    return response.data;
  },

  getUserLeaves: async (page = 1, pageSize = 10) => {
    const response = await api.get(`/leaves?page=${page}&page_size=${pageSize}`);
    return response.data;
  },

  getLeaveById: async (id) => {
    const response = await api.get(`/leaves/${id}`);
    return response.data;
  },

  cancelLeave: async (id) => {
    const response = await api.delete(`/leaves/${id}`);
    return response.data;
  },

  getLeaveTypes: async () => {
    const response = await api.get('/leave-types');
    return response.data;
  },

  getUserBalance: async () => {
    const response = await api.get('/profile/balance');
    return response.data;
  },
};

export const adminApi = {
  getAllLeaves: async (status = '', page = 1, pageSize = 10) => {
    const response = await api.get(`/admin/leaves?status=${status}&page=${page}&page_size=${pageSize}`);
    return response.data;
  },

  approveLeave: async (id, reviewNotes = '') => {
    const response = await api.put(`/admin/leaves/${id}/approve`, { review_notes: reviewNotes });
    return response.data;
  },

  rejectLeave: async (id, reviewNotes = '') => {
    const response = await api.put(`/admin/leaves/${id}/reject`, { review_notes: reviewNotes });
    return response.data;
  },

  getAllUsers: async (page = 1, pageSize = 10) => {
    const response = await api.get(`/admin/users?page=${page}&page_size=${pageSize}`);
    return response.data;
  },

  createUser: async (data) => {
    const response = await api.post('/admin/users', data);
    return response.data;
  },

  getDashboardStats: async () => {
    const response = await api.get('/admin/dashboard');
    return response.data;
  },

  allocateBalance: async (data) => {
    const response = await api.post('/admin/balances', data);
    return response.data;
  },
};
