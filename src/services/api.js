import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'https://api-admin.astrosway.com';
//'http://localhost:8082/'
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
  updateStatus: async (astroId, status, finalStatus, astroAmountDisbursed) => {
    const payload = {};

    // Only send fields that are explicitly provided from the UI.
    if (status !== undefined && status !== null && status !== '') {
      payload.status = status;
    }
    if (finalStatus !== undefined && finalStatus !== null && finalStatus !== '') {
      payload.finalStatus = finalStatus;
    }
    if (astroAmountDisbursed !== undefined && astroAmountDisbursed !== null && astroAmountDisbursed !== '') {
      // Send as a number when possible; backend treats this as optional.
      const numericAmount = Number(astroAmountDisbursed);
      payload.astroAmountDisbursed = Number.isNaN(numericAmount) ? astroAmountDisbursed : numericAmount;
    }

    const response = await api.patch(`/admin/astros/${astroId}/status`, payload);
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
    const response = await api.patch(`/admin/user-general-complaints/${problemId}/close`, {
      reason,
      response: reason,
    });
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
    const response = await api.patch(`/admin/astro-general-complaints/${problemId}/close`, {
      reason,
      response: reason,
    });
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

// Blogs API
export const blogsAPI = {
  listBlogs: async (params = {}) => {
    const response = await api.get('/admin/blogs', { params });
    return response.data;
  },
  getBlog: async (blogId) => {
    const response = await api.get(`/admin/blogs/${blogId}`);
    return response.data;
  },
  createBlog: async (data) => {
    const formData = new FormData();

    if (data.title) formData.append('title', data.title);
    if (data.slug) formData.append('slug', data.slug);
    if (data.status) formData.append('status', data.status);
    if (data.seoTitle) formData.append('seoTitle', data.seoTitle);
    if (data.seoDescription) formData.append('seoDescription', data.seoDescription);
    if (data.contentType) formData.append('contentType', data.contentType);

    // Backend expects "contentBody" field (per curl/Postman example)
    if (data.contentBody) {
      formData.append('contentBody', data.contentBody);
    } else if (data.content) {
      // Fallback if existing UI still uses "content"
      formData.append('contentBody', data.content);
    }

    // Optional metadata fields used by backend
    if (data.excerpt) formData.append('excerpt', data.excerpt);

    if (data.isFeatured !== undefined && data.isFeatured !== null && data.isFeatured !== '') {
      formData.append('isFeatured', String(data.isFeatured));
    }

    if (data.categories) {
      // Backend curl example uses plain string "Astrology"
      formData.append('categories', Array.isArray(data.categories) ? data.categories.join(',') : data.categories);
    }

    if (data.tags) {
      formData.append('tags', Array.isArray(data.tags) ? data.tags.join(',') : data.tags);
    }

    if (data.readingTimeMin !== undefined && data.readingTimeMin !== null && data.readingTimeMin !== '') {
      formData.append('readingTimeMin', data.readingTimeMin);
    }
    if (data.coverImage) {
      formData.append('coverImage', data.coverImage);
    }

    const response = await api.post('/admin/blogs', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
  updateBlog: async (blogId, data) => {
    const payload = {};

    if (data.title !== undefined) payload.title = data.title;
    if (data.slug !== undefined) payload.slug = data.slug;
    if (data.status !== undefined) payload.status = data.status;
    if (data.contentType !== undefined) payload.contentType = data.contentType;
    // Backend expects "contentBody" (align with createBlog and curl example)
    if (data.contentBody !== undefined) {
      payload.contentBody = data.contentBody;
    } else if (data.content !== undefined) {
      // Backwards compatibility if callers still use "content"
      payload.contentBody = data.content;
    }
    if (data.seoTitle !== undefined) payload.seoTitle = data.seoTitle;
    if (data.seoDescription !== undefined) payload.seoDescription = data.seoDescription;
    if (data.readingTimeMin !== undefined && data.readingTimeMin !== null && data.readingTimeMin !== '') {
      payload.readingTimeMin = data.readingTimeMin;
    }

    const response = await api.patch(`/admin/blogs/${blogId}`, payload);
    return response.data;
  },
  deleteBlog: async (blogId) => {
    const response = await api.delete(`/admin/blogs/${blogId}`);
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

// Dashboard API — pass date (YYYY-MM-DD) for daily metrics; backend defaults to today if omitted
export const dashboardAPI = {
  getDailyMetrics: async (date) => {
    const response = await api.get('/admin/dashboard/metrics', {
      params: date ? { date } : {},
    });
    return response.data;
  },
};

// Feedbacks API
export const feedbacksAPI = {
  bulkCreate: async (feedbacks) => {
    const response = await api.post('/admin/feedbacks/bulk', feedbacks);
    return response.data;
  },
};

export default api;

