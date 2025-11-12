import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { blogAPI } from '@/lib/blog';
import { Blog } from '@/types';
import BlogCard from '@/components/BlogCard';
import Pagination from '@/components/Pagination';
import LoadingSpinner from '@/components/LoadingSpinner';
import { BookOpen, TrendingUp, Calendar } from 'lucide-react';

const Home: React.FC = () => {
  const [blogs, setBlogs] = useState<Blog[]>([]);
  const [popularBlogs, setPopularBlogs] = useState<Blog[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [isLoading, setIsLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'recent' | 'popular'>('recent');
  const [error, setError] = useState<string | null>(null);

  const fetchBlogs = async (page: number) => {
    try {
      const response = await blogAPI.getAllBlogs(page, 6, 'newest');
      setBlogs(response.data || []);
      setTotalPages(response.totalPages || 1);
      setError(null);
    } catch (error) {
      console.error('Error fetching blogs:', error);
      setError('Failed to load blogs. Please try again later.');
      setBlogs([]);
    }
  };

  const fetchPopularBlogs = async () => {
    try {
      const popular = await blogAPI.getPopularBlogs(6);
      setPopularBlogs(popular || []);
    } catch (error) {
      console.error('Error fetching popular blogs:', error);
      setPopularBlogs([]);
    }
  };

  useEffect(() => {
    const loadData = async () => {
      setIsLoading(true);
      await Promise.all([
        fetchBlogs(currentPage),
        fetchPopularBlogs(),
      ]);
      setIsLoading(false);
    };

    loadData();
  }, [currentPage]);

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const displayBlogs = activeTab === 'recent' ? blogs : popularBlogs;

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Welcome to BlogPlatform
          </h1>
          <p className="text-xl text-gray-600 max-w-2xl mx-auto">
            Discover insightful articles, share your thoughts, and connect with a community of passionate writers and readers.
          </p>
          <div className="mt-8 flex justify-center space-x-4">
            <Link
              to="/blogs"
              className="btn btn-primary"
            >
              <BookOpen size={20} className="mr-2" />
              Explore Blogs
            </Link>
            <Link
              to="/create"
              className="btn btn-secondary"
            >
              Start Writing
            </Link>
          </div>
        </div>

        <div className="mb-8">
          <div className="border-b border-gray-200">
            <nav className="-mb-px flex space-x-8">
              <button
                onClick={() => setActiveTab('recent')}
                className={`py-2 px-1 border-b-2 font-medium text-sm ${
                  activeTab === 'recent'
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                <Calendar size={16} className="inline mr-2" />
                Recent Posts
              </button>
              <button
                onClick={() => setActiveTab('popular')}
                className={`py-2 px-1 border-b-2 font-medium text-sm ${
                  activeTab === 'popular'
                    ? 'border-primary-500 text-primary-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                <TrendingUp size={16} className="inline mr-2" />
                Popular Posts
              </button>
            </nav>
          </div>
        </div>

        {error ? (
          <div className="text-center py-12">
            <div className="text-red-500 mb-4">
              <svg className="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
            </div>
            <h3 className="mt-2 text-sm font-medium text-gray-900">Error Loading Blogs</h3>
            <p className="mt-1 text-sm text-gray-500">{error}</p>
            <div className="mt-6">
              <button
                onClick={() => {
                  setError(null);
                  fetchBlogs(currentPage);
                }}
                className="btn btn-primary"
              >
                Try Again
              </button>
            </div>
          </div>
        ) : isLoading ? (
          <div className="flex justify-center items-center py-12">
            <LoadingSpinner size="lg" />
          </div>
        ) : (
          <>
            {displayBlogs.length === 0 ? (
              <div className="text-center py-12">
                <BookOpen className="mx-auto h-12 w-12 text-gray-400" />
                <h3 className="mt-2 text-sm font-medium text-gray-900">No blogs found</h3>
                <p className="mt-1 text-sm text-gray-500">
                  {activeTab === 'recent' 
                    ? 'No recent blogs available. Be the first to write one!' 
                    : 'No popular blogs available yet.'
                  }
                </p>
                <div className="mt-6">
                  <Link
                    to="/create"
                    className="btn btn-primary"
                  >
                    Write your first blog
                  </Link>
                </div>
              </div>
            ) : (
              <>
                <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                  {displayBlogs.map((blog) => (
                    <BlogCard key={blog.id} blog={blog} />
                  ))}
                </div>

                {activeTab === 'recent' && totalPages > 1 && (
                  <Pagination
                    currentPage={currentPage}
                    totalPages={totalPages}
                    onPageChange={handlePageChange}
                  />
                )}
              </>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default Home;
