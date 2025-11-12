import api from './api';
import {
  AIGenerateRequest,
  AIEnhanceRequest,
  AISuggestIdeasRequest,
} from '@/types';

export const aiAPI = {
  generateBlog: async (data: AIGenerateRequest): Promise<{ content: string }> => {
    const response = await api.post('/ai/generate-blog', data);
    return response.data;
  },

  enhanceBlog: async (data: AIEnhanceRequest): Promise<{ content: string }> => {
    const response = await api.post('/ai/enhance-blog', data);
    return response.data;
  },

  suggestIdeas: async (data: AISuggestIdeasRequest): Promise<{ ideas: string[] }> => {
    const response = await api.post('/ai/suggest-ideas', data);
    return response.data;
  },
};
