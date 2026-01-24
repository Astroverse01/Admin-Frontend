import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'https://api-admin.astrosway.com';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Handle 401 errors (unauthorized)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: async (username, password) => {
    const response = await api.post('/admin/login', { username, password });
    return response.data;
  },
};

// Users API
export const usersAPI = {
  listUsers: async (params = {}) => {
    const response = await api.get('/admin/users', { params });
    return response.data;
  },
  deactivateUser: async (userId, status) => {
    const response = await api.patch(`/admin/users/${userId}/deactivate`, { status });
    return response.data;
  },
};

// Astros API
export const astrosAPI = {
  listAstros: async (params = {}) => {
    const response = await api.get('/admin/astros', { params });
    return response.data;
  },
  updateStatus: async (astroId, status) => {
    const response = await api.patch(`/admin/astros/${astroId}/status`, { status });
    return response.data;
  },
  toggleVisibility: async (astroId, visible) => {
    const response = await api.patch(`/admin/astros/${astroId}/visibility`, { visible });
    return response.data;
  },
};

// Complaints API
export const complaintsAPI = {
  listUserServiceComplaints: async (params = {}) => {
    const response = await api.get('/admin/user-service-complaints', { params });
    return response.data;
  },
  getComplaintDetails: async (serviceType, orderId) => {
    const response = await api.get(`/admin/user-service-complaints/${serviceType}/${orderId}`);
    return response.data;
  },
  acceptRejectComplaint: async (orderId, data) => {
    const response = await api.patch(`/admin/user-service-complaints/${orderId}`, data);
    return response.data;
  },
};

// User Problems API
export const userProblemsAPI = {
  listUserGeneralComplaints: async (params = {}) => {
    const response = await api.get('/admin/user-general-complaints', { params });
    return response.data;
  },
  closeComplaint: async (problemId, reason) => {
    const response = await api.patch(`/admin/user-general-complaints/${problemId}/close`, { reason });
    return response.data;
  },
};

// Astro Problems API
export const astroProblemsAPI = {
  listAstroGeneralComplaints: async (params = {}) => {
    const response = await api.get('/admin/astro-general-complaints', { params });
    return response.data;
  },
  closeComplaint: async (problemId, reason) => {
    const response = await api.patch(`/admin/astro-general-complaints/${problemId}/close`, { reason });
    return response.data;
  },
};

// Horoscopes API
export const horoscopesAPI = {
  listHoroscopes: async (params = {}) => {
    const response = await api.get('/admin/horoscopes', { params });
    return response.data;
  },
  bulkCreate: async (horoscopes) => {
    const response = await api.post('/admin/horoscopes/bulk', horoscopes);
    return response.data;
  },
  updateHoroscope: async (horoscopeId, data) => {
    const response = await api.patch(`/admin/horoscopes/${horoscopeId}`, data);
    return response.data;
  },
  deleteHoroscope: async (horoscopeId) => {
    const response = await api.delete(`/admin/horoscopes/${horoscopeId}`);
    return response.data;
  },
};

// Scheduler API
export const schedulerAPI = {
  triggerManually: async () => {
    const response = await api.post('/admin/scheduler/trigger');
    return response.data;
  },
  generateReports: async (startDate, endDate) => {
    const response = await api.post('/admin/scheduler/generate', { startDate, endDate });
    return response.data;
  },
  downloadReport: async (fileName) => {
    const response = await api.get(`/admin/scheduler/download/${fileName}`, {
      responseType: 'blob', // Important for file download
    });
    return response;
  },
};

// Dashboard API
export const dashboardAPI = {
  getMetrics: async () => {
    const response = await api.get('/admin/dashboard/metrics');
    return response.data;
  },
};

export default api;

