import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { blogAPI } from '@/lib/blog';
import { useAuth } from '@/contexts/AuthContext';
import { Blog } from '@/types';
import { formatRelativeTime } from '@/lib/utils';
import {
  Calendar,
  Eye,
  Heart,
  MessageCircle,
  Edit,
  Trash2,
  User,
  ArrowLeft,
} from 'lucide-react';
import toast from 'react-hot-toast';

const BlogDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [blog, setBlog] = useState<Blog | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isDeleting, setIsDeleting] = useState(false);
  const { user } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (id) {
      fetchBlog(id);
    }
  }, [id]);

  const fetchBlog = async (blogId: string) => {
    try {
      const blogData = await blogAPI.getBlog(blogId);
      setBlog(blogData);
    } catch (error) {
      console.error('Error fetching blog:', error);
      toast.error('Failed to load blog');
    } finally {
      setIsLoading(false);
    }
  };

  const handleEdit = () => {
    if (blog) {
      navigate(`/blogs/${blog.id}/edit`);
    }
  };

  const handleDelete = async () => {
    if (!blog || !user) return;

    if (!window.confirm('Are you sure you want to delete this blog? This action cannot be undone.')) {
      return;
    }

    setIsDeleting(true);
    try {
      await blogAPI.deleteBlog(blog.id);
      toast.success('Blog deleted successfully');
      navigate('/blogs');
    } catch (error) {
      console.error('Error deleting blog:', error);
      toast.error('Failed to delete blog');
    } finally {
      setIsDeleting(false);
    }
  };

  const handleLike = async () => {
    if (!blog || !user) return;

    try {
      if ((blog.likes || []).includes(user.id)) {
        await blogAPI.dislikeBlog(blog.id);
      } else {
        await blogAPI.likeBlog(blog.id);
      }
      fetchBlog(blog.id);
    } catch (error) {
      console.error('Error liking blog:', error);
      toast.error('Failed to update like');
    }
  };

  const canEdit = user && (blog?.authorId === user.id || user.role === 'admin');
  const hasLiked = !!(user && (blog?.likes || []).includes(user.id));

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (!blog) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">Blog not found</h2>
          <p className="text-gray-600 mb-4">The blog you're looking for doesn't exist.</p>
          <button
            onClick={() => navigate('/blogs')}
            className="btn btn-primary"
          >
            Back to Blogs
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <button
          onClick={() => navigate(-1)}
          className="flex items-center space-x-2 text-gray-600 hover:text-gray-900 mb-6"
        >
          <ArrowLeft size={20} />
          <span>Back</span>
        </button>

        <article className="card p-8">
          <header className="mb-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-4">{blog.title}</h1>
            
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center space-x-4 text-sm text-gray-500">
                <div className="flex items-center space-x-1">
                  <User size={16} />
                  <span>{blog.authorUsername}</span>
                </div>
                <div className="flex items-center space-x-1">
                  <Calendar size={16} />
                  <span>{formatRelativeTime(blog.createdAt)}</span>
                </div>
                <div className="flex items-center space-x-1">
                  <Eye size={16} />
                  <span>{blog.viewCount} views</span>
                </div>
              </div>

              {canEdit && (
                <div className="flex space-x-2">
                  <button
                    onClick={handleEdit}
                    className="btn btn-secondary text-sm"
                  >
                    <Edit size={16} className="mr-1" />
                    Edit
                  </button>
                  <button
                    onClick={handleDelete}
                    disabled={isDeleting}
                    className="btn btn-danger text-sm disabled:opacity-50"
                  >
                    <Trash2 size={16} className="mr-1" />
                    {isDeleting ? 'Deleting...' : 'Delete'}
                  </button>
                </div>
              )}
            </div>

            {blog.tags && blog.tags.length > 0 && (
              <div className="flex flex-wrap gap-2 mb-4">
                {blog.tags.map((tag) => (
                  <span
                    key={tag}
                    className="px-3 py-1 bg-primary-100 text-primary-700 text-sm rounded-full"
                  >
                    #{tag}
                  </span>
                ))}
              </div>
            )}
          </header>

          <div className="prose prose-lg max-w-none mb-8">
            <div dangerouslySetInnerHTML={{ __html: blog.content.replace(/\n/g, '<br>') }} />
          </div>

          <footer className="border-t border-gray-200 pt-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-4">
                <button
                  onClick={handleLike}
                  className={`flex items-center space-x-2 px-4 py-2 rounded-md transition-colors ${
                    hasLiked
                      ? 'bg-red-100 text-red-700'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                >
                  <Heart size={16} className={hasLiked ? 'fill-current' : ''} />
                  <span>{blog.likeCount}</span>
                </button>

                <div className="flex items-center space-x-2 text-gray-600">
                  <MessageCircle size={16} />
                  <span>{blog.commentCount} comments</span>
                </div>
              </div>

              <div className="text-sm text-gray-500">
                Last updated: {formatRelativeTime(blog.updatedAt)}
              </div>
            </div>
          </footer>
        </article>

        {blog.comments && blog.comments.length > 0 && (
          <div className="mt-8 card p-6">
            <h3 className="text-xl font-semibold text-gray-900 mb-4">
              Comments ({blog.commentCount})
            </h3>
            <div className="space-y-4">
              {blog.comments.map((comment) => (
                <div key={comment.id} className="border-b border-gray-200 pb-4 last:border-b-0">
                  <div className="flex items-start space-x-3">
                    <div className="flex-shrink-0">
                      <div className="w-8 h-8 bg-primary-100 rounded-full flex items-center justify-center">
                        <User size={16} className="text-primary-600" />
                      </div>
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center space-x-2 mb-1">
                        <span className="font-medium text-gray-900">{comment.authorUsername}</span>
                        <span className="text-sm text-gray-500">
                          {formatRelativeTime(comment.createdAt)}
                        </span>
                      </div>
                      <p className="text-gray-700">{comment.content}</p>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default BlogDetail;
