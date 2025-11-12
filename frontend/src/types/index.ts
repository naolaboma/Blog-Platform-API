export interface User {
  id: string;
  username: string;
  email: string;
  role: 'admin' | 'user';
  emailVerified: boolean;
  profilePicture?: Photo;
  bio?: string;
  createdAt: string;
  updatedAt: string;
  oauthProvider?: string;
  oauthId?: string;
}

export interface Photo {
  filename: string;
  filePath: string;
  publicId: string;
  uploadedAt: string;
}

export interface Blog {
  id: string;
  title: string;
  content: string;
  authorId: string;
  authorUsername: string;
  tags: string[];
  viewCount: number;
  likeCount: number;
  commentCount: number;
  likes: string[];
  dislikes: string[];
  comments: Comment[];
  createdAt: string;
  updatedAt: string;
}

export interface Comment {
  id: string;
  authorId: string;
  authorUsername: string;
  content: string;
  createdAt: string;
  updatedAt: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
}

export interface UpdateProfileRequest {
  username?: string;
  email?: string;
  bio?: string;
}

export interface CreateBlogRequest {
  title: string;
  content: string;
  tags: string[];
}

export interface UpdateBlogRequest {
  title?: string;
  content?: string;
  tags?: string[];
}

export interface CreateCommentRequest {
  content: string;
}

export interface PaginationResponse<T> {
  data: T[];
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}

export interface AIGenerateRequest {
  topic: string;
}

export interface AIEnhanceRequest {
  content: string;
}

export interface AISuggestIdeasRequest {
  keywords: string[];
}

export interface ErrorResponse {
  error: string;
}
