import axios from "axios";
import type {
  APIResponse,
  PaginatedResponse,
  LoginRequest,
  RegisterRequest,
  LoginResponse,
  TokenPair,
  UserResponse,
  CreateUserRequest,
  UpdateUserRequest,
  HealthResponse,
} from "../types/api";

const API_BASE_URL = "http://localhost:8080";

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("accessToken");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor to handle token refresh
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      const refreshToken = localStorage.getItem("refreshToken");
      if (refreshToken) {
        try {
          const response = await axios.post(
            `${API_BASE_URL}/api/v1/auth/refresh`,
            {
              refresh_token: refreshToken,
            }
          );

          const { access_token, refresh_token } = response.data.data;
          localStorage.setItem("accessToken", access_token);
          localStorage.setItem("refreshToken", refresh_token);

          // Retry original request with new token
          originalRequest.headers.Authorization = `Bearer ${access_token}`;
          return api(originalRequest);
        } catch {
          // Refresh failed, redirect to login
          localStorage.removeItem("accessToken");
          localStorage.removeItem("refreshToken");
          localStorage.removeItem("user");
          window.location.href = "/login";
        }
      }
    }

    return Promise.reject(error);
  }
);

// Health Check
export const healthCheck = async (): Promise<HealthResponse> => {
  const response = await api.get<HealthResponse>("/health");
  return response.data;
};

// API Info
export const getApiInfo = async (): Promise<APIResponse> => {
  const response = await api.get<APIResponse>("/api/v1");
  return response.data;
};

// Authentication endpoints
export const authApi = {
  register: async (
    data: RegisterRequest
  ): Promise<APIResponse<LoginResponse>> => {
    const response = await api.post<APIResponse<LoginResponse>>(
      "/api/v1/auth/register",
      data
    );
    return response.data;
  },

  login: async (data: LoginRequest): Promise<APIResponse<LoginResponse>> => {
    const response = await api.post<APIResponse<LoginResponse>>(
      "/api/v1/auth/login",
      data
    );
    return response.data;
  },

  refreshToken: async (
    refreshToken: string
  ): Promise<APIResponse<TokenPair>> => {
    const response = await api.post<APIResponse<TokenPair>>(
      "/api/v1/auth/refresh",
      {
        refresh_token: refreshToken,
      }
    );
    return response.data;
  },

  logout: async (refreshToken: string): Promise<APIResponse<null>> => {
    const response = await api.post<APIResponse<null>>("/api/v1/auth/logout", {
      refresh_token: refreshToken,
    });
    return response.data;
  },

  me: async (): Promise<APIResponse<UserResponse>> => {
    const response = await api.get<APIResponse<UserResponse>>(
      "/api/v1/auth/me"
    );
    return response.data;
  },
};

// Users endpoints
export const usersApi = {
  getUsers: async (
    page = 1,
    limit = 10
  ): Promise<PaginatedResponse<UserResponse[]>> => {
    const response = await api.get<PaginatedResponse<UserResponse[]>>(
      `/api/v1/users?page=${page}&limit=${limit}`
    );
    return response.data;
  },

  getUser: async (id: number): Promise<APIResponse<UserResponse>> => {
    const response = await api.get<APIResponse<UserResponse>>(
      `/api/v1/users/${id}`
    );
    return response.data;
  },

  createUser: async (
    data: CreateUserRequest
  ): Promise<APIResponse<UserResponse>> => {
    const response = await api.post<APIResponse<UserResponse>>(
      "/api/v1/users",
      data
    );
    return response.data;
  },

  updateUser: async (
    id: number,
    data: UpdateUserRequest
  ): Promise<APIResponse<UserResponse>> => {
    const response = await api.put<APIResponse<UserResponse>>(
      `/api/v1/users/${id}`,
      data
    );
    return response.data;
  },

  deleteUser: async (id: number): Promise<APIResponse<null>> => {
    const response = await api.delete<APIResponse<null>>(`/api/v1/users/${id}`);
    return response.data;
  },
};

export default api;
