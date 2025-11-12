import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { useAuth } from '@/contexts/AuthContext';
import { blogAPI } from '@/lib/blog';
import { aiAPI } from '@/lib/ai';
import { CreateBlogRequest } from '@/types';
import { Save, Eye, Sparkles, X } from 'lucide-react';
import toast from 'react-hot-toast';

interface CreateBlogFormData {
  title: string;
  content: string;
  tags: string;
}

const CreateBlog: React.FC = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [isGenerating, setIsGenerating] = useState(false);
  const [isEnhancing, setIsEnhancing] = useState(false);
  const [showPreview, setShowPreview] = useState(false);
  const [aiSuggestions, setAiSuggestions] = useState<string[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const { user } = useAuth();
  const navigate = useNavigate();

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    watch,
    reset,
  } = useForm<CreateBlogFormData>();

  const watchedContent = watch('content');

  useEffect(() => {
    if (!user) {
      navigate('/login');
    }
  }, [user, navigate]);

  const onSubmit = async (data: CreateBlogFormData) => {
    if (!user) return;

    setIsLoading(true);
    try {
      const tagsArray = data.tags
        .split(',')
        .map(tag => tag.trim())
        .filter(tag => tag.length > 0);

      const blogData: CreateBlogRequest = {
        title: data.title,
        content: data.content,
        tags: tagsArray,
      };

      const newBlog = await blogAPI.createBlog(blogData);
      toast.success('Blog published successfully!');
      navigate(`/blogs/${newBlog.id}`);
    } catch (error) {
      console.error('Error creating blog:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const generateContent = async () => {
    const title = watch('title');
    if (!title) {
      toast.error('Please enter a title first');
      return;
    }

    setIsGenerating(true);
    try {
      const response = await aiAPI.generateBlog({ topic: title });
      setValue('content', response.content);
      toast.success('Content generated successfully!');
    } catch (error) {
      console.error('Error generating content:', error);
    } finally {
      setIsGenerating(false);
    }
  };

  const enhanceContent = async () => {
    const content = watch('content');
    if (!content) {
      toast.error('Please enter some content first');
      return;
    }

    setIsEnhancing(true);
    try {
      const response = await aiAPI.enhanceBlog({ content });
      setValue('content', response.content);
      toast.success('Content enhanced successfully!');
    } catch (error) {
      console.error('Error enhancing content:', error);
    } finally {
      setIsEnhancing(false);
    }
  };

  const suggestIdeas = async () => {
    const tags = watch('tags');
    const keywords = tags
      ? tags.split(',').map(tag => tag.trim()).filter(tag => tag.length > 0)
      : [];

    if (keywords.length === 0) {
      toast.error('Please enter some keywords or tags first');
      return;
    }

    setIsGenerating(true);
    try {
      const response = await aiAPI.suggestIdeas({ keywords });
      setAiSuggestions(response.ideas);
      setShowSuggestions(true);
    } catch (error) {
      console.error('Error getting suggestions:', error);
    } finally {
      setIsGenerating(false);
    }
  };

  const useSuggestion = (suggestion: string) => {
    setValue('title', suggestion);
    setShowSuggestions(false);
    setAiSuggestions([]);
  };

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Create New Blog</h1>
          <div className="flex space-x-2">
            <button
              type="button"
              onClick={() => setShowPreview(!showPreview)}
              className="btn btn-secondary"
            >
              <Eye size={16} className="mr-2" />
              {showPreview ? 'Edit' : 'Preview'}
            </button>
          </div>
        </div>

        <div className="card p-6">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
                Blog Title
              </label>
              <input
                {...register('title', {
                  required: 'Title is required',
                  minLength: {
                    value: 5,
                    message: 'Title must be at least 5 characters',
                  },
                  maxLength: {
                    value: 255,
                    message: 'Title must be less than 255 characters',
                  },
                })}
                type="text"
                className="input"
                placeholder="Enter an engaging title..."
                disabled={showPreview}
              />
              {errors.title && (
                <p className="mt-1 text-sm text-red-600">{errors.title.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="content" className="block text-sm font-medium text-gray-700 mb-2">
                Content
              </label>
              <div className="relative">
                <textarea
                  {...register('content', {
                    required: 'Content is required',
                    minLength: {
                      value: 20,
                      message: 'Content must be at least 20 characters',
                    },
                  })}
                  rows={12}
                  className={`input resize-none ${showPreview ? 'hidden' : ''}`}
                  placeholder="Write your blog content here..."
                  disabled={showPreview}
                />
                {showPreview && (
                  <div className="min-h-[300px] p-4 border border-gray-300 rounded-lg bg-white prose prose-sm max-w-none">
                    <div dangerouslySetInnerHTML={{ __html: watchedContent.replace(/\n/g, '<br>') }} />
                  </div>
                )}
              </div>
              {errors.content && (
                <p className="mt-1 text-sm text-red-600">{errors.content.message}</p>
              )}
              
              {!showPreview && (
                <div className="mt-2 flex space-x-2">
                  <button
                    type="button"
                    onClick={generateContent}
                    disabled={isGenerating}
                    className="btn btn-secondary text-sm disabled:opacity-50"
                  >
                    <Sparkles size={14} className="mr-1" />
                    {isGenerating ? 'Generating...' : 'Generate with AI'}
                  </button>
                  <button
                    type="button"
                    onClick={enhanceContent}
                    disabled={isEnhancing}
                    className="btn btn-secondary text-sm disabled:opacity-50"
                  >
                    <Sparkles size={14} className="mr-1" />
                    {isEnhancing ? 'Enhancing...' : 'Enhance with AI'}
                  </button>
                </div>
              )}
            </div>

            <div>
              <label htmlFor="tags" className="block text-sm font-medium text-gray-700 mb-2">
                Tags
              </label>
              <input
                {...register('tags')}
                type="text"
                className="input"
                placeholder="Enter tags separated by commas (e.g., technology, programming, web)"
                disabled={showPreview}
              />
              <p className="mt-1 text-sm text-gray-500">
                Separate multiple tags with commas
              </p>
              <button
                type="button"
                onClick={suggestIdeas}
                disabled={isGenerating}
                className="mt-2 btn btn-secondary text-sm disabled:opacity-50"
              >
                <Sparkles size={14} className="mr-1" />
                Get AI Suggestions
              </button>
            </div>

            <div className="flex justify-end space-x-4">
              <Link
                to="/blogs"
                className="btn btn-secondary"
              >
                Cancel
              </Link>
              <button
                type="submit"
                disabled={isLoading}
                className="btn btn-primary disabled:opacity-50"
              >
                <Save size={16} className="mr-2" />
                {isLoading ? 'Publishing...' : 'Publish Blog'}
              </button>
            </div>
          </form>
        </div>

        {showSuggestions && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-semibold">AI Title Suggestions</h3>
                <button
                  onClick={() => {
                    setShowSuggestions(false);
                    setAiSuggestions([]);
                  }}
                  className="text-gray-400 hover:text-gray-600"
                >
                  <X size={20} />
                </button>
              </div>
              <div className="space-y-2">
                {aiSuggestions.map((suggestion, index) => (
                  <button
                    key={index}
                    onClick={() => useSuggestion(suggestion)}
                    className="w-full text-left p-3 border border-gray-200 rounded-md hover:bg-gray-50 transition-colors"
                  >
                    {suggestion}
                  </button>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default CreateBlog;
