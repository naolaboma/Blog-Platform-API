import React from 'react';
import { Link } from 'react-router-dom';
import { Calendar, Eye, Heart, MessageCircle, User } from 'lucide-react';
import { Blog } from '@/types';
import { formatRelativeTime, truncateText } from '@/lib/utils';

interface BlogCardProps {
  blog: Blog;
}

const BlogCard: React.FC<BlogCardProps> = ({ blog }) => {
  // Defensive checks for missing data
  if (!blog) {
    return null;
  }

  return (
    <div className="card p-6 hover:shadow-lg transition-shadow duration-200">
      <div className="flex items-start justify-between mb-4">
        <div className="flex-1">
          <Link
            to={`/blogs/${blog.id}`}
            className="text-xl font-semibold text-gray-900 hover:text-primary-600 transition-colors"
          >
            {truncateText(blog.title || 'Untitled', 80)}
          </Link>
          <div className="flex items-center space-x-4 mt-2 text-sm text-gray-500">
            <div className="flex items-center space-x-1">
              <User size={14} />
              <Link
                to={`/blogs?author=${blog.authorUsername || 'unknown'}`}
                className="hover:text-primary-600"
              >
                {blog.authorUsername || 'Unknown Author'}
              </Link>
            </div>
            <div className="flex items-center space-x-1">
              <Calendar size={14} />
              <span>{formatRelativeTime(blog.createdAt || new Date().toISOString())}</span>
            </div>
          </div>
        </div>
      </div>

      <p className="text-gray-600 mb-4 line-clamp-3">
        {truncateText(blog.content || 'No content available', 200)}
      </p>

      {blog.tags && blog.tags.length > 0 && (
        <div className="flex flex-wrap gap-2 mb-4">
          {blog.tags.slice(0, 5).map((tag) => (
            <Link
              key={tag}
              to={`/blogs?tags=${tag}`}
              className="px-2 py-1 bg-primary-100 text-primary-700 text-xs rounded-full hover:bg-primary-200 transition-colors"
            >
              #{tag}
            </Link>
          ))}
          {blog.tags && blog.tags.length > 5 && (
            <span className="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-full">
              +{blog.tags.length - 5}
            </span>
          )}
        </div>
      )}

      <div className="flex items-center justify-between text-sm text-gray-500">
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-1">
            <Eye size={14} />
            <span>{blog.viewCount || 0}</span>
          </div>
          <div className="flex items-center space-x-1">
            <Heart size={14} />
            <span>{blog.likeCount || 0}</span>
          </div>
          <div className="flex items-center space-x-1">
            <MessageCircle size={14} />
            <span>{blog.commentCount || 0}</span>
          </div>
        </div>
        <Link
          to={`/blogs/${blog.id}`}
          className="text-primary-600 hover:text-primary-700 font-medium"
        >
          Read more →
        </Link>
      </div>
    </div>
  );
};

export default BlogCard;
