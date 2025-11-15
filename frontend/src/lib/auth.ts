import api from "./api";
import {
  User,
  LoginRequest,
  RegisterRequest,
  LoginResponse,
  UpdateProfileRequest,
} from "@/types";

export const authAPI = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post("/auth/login", data);
    return response.data;
  },

  register: async (data: RegisterRequest): Promise<LoginResponse> => {
    const response = await api.post("/auth/register", data);
    return response.data;
  },

  logout: async (): Promise<void> => {
    await api.post("/auth/logout");
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
  },

  refreshToken: async (refreshToken: string): Promise<LoginResponse> => {
    const response = await api.post("/auth/refresh", {
      refresh_token: refreshToken,
    });
    return response.data;
  },

  sendVerificationEmail: async (email: string): Promise<void> => {
    await api.post("/auth/send-verification", { email });
  },

  verifyEmail: async (
    token: string
  ): Promise<{ message: string; user: User }> => {
    const response = await api.get(`/auth/verify-email?token=${token}`);
    return response.data;
  },

  sendPasswordResetEmail: async (email: string): Promise<void> => {
    await api.post("/auth/forgot-password", { email });
  },

  resetPassword: async (token: string, newPassword: string): Promise<void> => {
    await api.post("/auth/reset-password", {
      token,
      new_password: newPassword,
    });
  },

  getGoogleOAuthURL: (): string => {
    return `${api.defaults.baseURL}/auth/google/login`;
  },

  getGitHubOAuthURL: (): string => {
    return `${api.defaults.baseURL}/auth/github/login`;
  },

  getProfile: async (): Promise<User> => {
    const response = await api.get("/users/profile");
    return response.data;
  },
};

export const userAPI = {
  getProfile: async (): Promise<User> => {
    const response = await api.get("/users/profile");
    return response.data;
  },

  updateProfile: async (data: UpdateProfileRequest): Promise<User> => {
    const response = await api.put("/users/profile", data);
    return response.data;
  },

  uploadProfilePicture: async (file: File): Promise<User> => {
    const formData = new FormData();
    formData.append("profile_picture", file);

    const response = await api.post("/users/profile/picture", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return response.data;
  },
};

export const adminAPI = {
  promoteUser: async (
    userId: string,
    role: "admin" | "user"
  ): Promise<void> => {
    await api.put(`/admin/users/${userId}/promote`, { role });
  },

  demoteUser: async (userId: string, role: "admin" | "user"): Promise<void> => {
    await api.put(`/admin/users/${userId}/demote`, { role });
  },

  getAllUsers: async (): Promise<{ users: User[] }> => {
    const response = await api.get("/admin/users");
    return response.data;
  },
};
