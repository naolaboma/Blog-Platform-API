import api from './api';
import {
  Blog,
  CreateBlogRequest,
  UpdateBlogRequest,
  CreateCommentRequest,
  PaginationResponse,
} from '@/types';

export const blogAPI = {
  getAllBlogs: async (
    page: number = 1,
    limit: number = 10,
    sort: string = 'newest'
  ): Promise<PaginationResponse<Blog>> => {
    const response = await api.get(`/blogs?page=${page}&limit=${limit}&sort=${sort}`);
    return response.data;
  },

  getBlog: async (id: string): Promise<Blog> => {
    const response = await api.get(`/blogs/${id}`);
    return response.data;
  },

  createBlog: async (data: CreateBlogRequest): Promise<Blog> => {
    const response = await api.post('/blogs', data);
    return response.data;
  },

  updateBlog: async (id: string, data: UpdateBlogRequest): Promise<Blog> => {
    const response = await api.put(`/blogs/${id}`, data);
    return response.data;
  },

  deleteBlog: async (id: string): Promise<void> => {
    await api.delete(`/blogs/${id}`);
  },

  searchBlogsByTitle: async (
    title: string,
    page: number = 1,
    limit: number = 10
  ): Promise<PaginationResponse<Blog>> => {
    const response = await api.get(
      `/blogs/search/title?title=${encodeURIComponent(title)}&page=${page}&limit=${limit}`
    );
    return response.data;
  },

  searchBlogsByAuthor: async (
    author: string,
    page: number = 1,
    limit: number = 10
  ): Promise<PaginationResponse<Blog>> => {
    const response = await api.get(
      `/blogs/search/author?author=${encodeURIComponent(author)}&page=${page}&limit=${limit}`
    );
    return response.data;
  },

  filterBlogsByTags: async (
    tags: string[],
    page: number = 1,
    limit: number = 10
  ): Promise<PaginationResponse<Blog>> => {
    const response = await api.get(
      `/blogs/filter/tags?tags=${tags.join(',')}&page=${page}&limit=${limit}`
    );
    return response.data;
  },

  filterBlogsByDate: async (
    startDate: string,
    endDate: string,
    page: number = 1,
    limit: number = 10
  ): Promise<PaginationResponse<Blog>> => {
    const response = await api.get(
      `/blogs/filter/date?start_date=${startDate}&end_date=${endDate}&page=${page}&limit=${limit}`
    );
    return response.data;
  },

  getPopularBlogs: async (limit: number = 10): Promise<Blog[]> => {
    const response = await api.get(`/blogs/popular?limit=${limit}`);
    return response.data;
  },

  addComment: async (blogId: string, data: CreateCommentRequest): Promise<void> => {
    await api.post(`/blogs/${blogId}/comments`, data);
  },

  updateComment: async (
    blogId: string,
    commentId: string,
    data: { content: string }
  ): Promise<void> => {
    await api.put(`/blogs/${blogId}/comments/${commentId}`, data);
  },

  deleteComment: async (blogId: string, commentId: string): Promise<void> => {
    await api.delete(`/blogs/${blogId}/comments/${commentId}`);
  },

  likeBlog: async (blogId: string): Promise<void> => {
    await api.post(`/blogs/${blogId}/like`);
  },

  dislikeBlog: async (blogId: string): Promise<void> => {
    await api.post(`/blogs/${blogId}/dislike`);
  },
};
