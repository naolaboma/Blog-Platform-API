import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { blogAPI } from '@/lib/blog';
import { Blog } from '@/types';
import BlogCard from '@/components/BlogCard';
import Pagination from '@/components/Pagination';
import LoadingSpinner from '@/components/LoadingSpinner';
import { Search, Filter, X } from 'lucide-react';
import toast from 'react-hot-toast';

const Blogs: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [blogs, setBlogs] = useState<Blog[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [isLoading, setIsLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [authorFilter, setAuthorFilter] = useState('');
  const [tagsFilter, setTagsFilter] = useState('');
  const [sortBy, setSortBy] = useState('newest');

  const query = searchParams.get('search') || '';
  const author = searchParams.get('author') || '';
  const tags = searchParams.get('tags') || '';

  useEffect(() => {
    setSearchTerm(query);
    setAuthorFilter(author);
    setTagsFilter(tags);
  }, [query, author, tags]);

  useEffect(() => {
    fetchBlogs();
  }, [currentPage, searchParams]);

  const fetchBlogs = async () => {
    setIsLoading(true);
    try {
      let response;
      
      if (query) {
        response = await blogAPI.searchBlogsByTitle(query, currentPage, 6);
      } else if (author) {
        response = await blogAPI.searchBlogsByAuthor(author, currentPage, 6);
      } else if (tags) {
        const tagArray = tags.split(',').map(tag => tag.trim());
        response = await blogAPI.filterBlogsByTags(tagArray, currentPage, 6);
      } else {
        response = await blogAPI.getAllBlogs(currentPage, 6, sortBy);
      }
      
      setBlogs(response.data);
      setTotalPages(response.totalPages);
    } catch (error) {
      console.error('Error fetching blogs:', error);
      toast.error('Failed to load blogs');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    const params = new URLSearchParams();
    if (searchTerm.trim()) {
      params.set('search', searchTerm.trim());
    }
    params.set('page', '1');
    setSearchParams(params);
    setCurrentPage(1);
  };

  const handleFilter = () => {
    const params = new URLSearchParams();
    if (authorFilter.trim()) {
      params.set('author', authorFilter.trim());
    }
    if (tagsFilter.trim()) {
      params.set('tags', tagsFilter.trim());
    }
    params.set('page', '1');
    setSearchParams(params);
    setCurrentPage(1);
  };

  const clearFilters = () => {
    setSearchTerm('');
    setAuthorFilter('');
    setTagsFilter('');
    setSortBy('newest');
    setSearchParams({});
    setCurrentPage(1);
  };

  const handlePageChange = (page: number) => {
    const params = new URLSearchParams(searchParams);
    params.set('page', page.toString());
    setSearchParams(params);
    setCurrentPage(page);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleSortChange = (newSortBy: string) => {
    setSortBy(newSortBy);
    setCurrentPage(1);
  };

  const hasActiveFilters = query || author || tags;

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">Explore Blogs</h1>

        <div className="card p-6 mb-8">
          <form onSubmit={handleSearch} className="mb-4">
            <div className="flex space-x-4">
              <div className="flex-1 relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Search className="h-5 w-5 text-gray-400" />
                </div>
                <input
                  type="text"
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="input pl-10"
                  placeholder="Search blogs by title..."
                />
              </div>
              <button type="submit" className="btn btn-primary">
                Search
              </button>
            </div>
          </form>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
            <div>
              <input
                type="text"
                value={authorFilter}
                onChange={(e) => setAuthorFilter(e.target.value)}
                className="input"
                placeholder="Filter by author..."
              />
            </div>
            <div>
              <input
                type="text"
                value={tagsFilter}
                onChange={(e) => setTagsFilter(e.target.value)}
                className="input"
                placeholder="Filter by tags (comma separated)..."
              />
            </div>
            <div className="flex space-x-2">
              <select
                value={sortBy}
                onChange={(e) => handleSortChange(e.target.value)}
                className="input"
              >
                <option value="newest">Newest First</option>
                <option value="oldest">Oldest First</option>
                <option value="popular">Most Popular</option>
              </select>
              <button
                onClick={handleFilter}
                className="btn btn-secondary"
              >
                <Filter size={16} className="mr-1" />
                Filter
              </button>
            </div>
          </div>

          {hasActiveFilters && (
            <div className="flex items-center justify-between">
              <div className="flex flex-wrap gap-2">
                {query && (
                  <span className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-primary-100 text-primary-700">
                    Search: {query}
                    <button
                      onClick={() => {
                        setSearchTerm('');
                        const params = new URLSearchParams(searchParams);
                        params.delete('search');
                        params.set('page', '1');
                        setSearchParams(params);
                      }}
                      className="ml-2 hover:text-primary-900"
                    >
                      <X size={14} />
                    </button>
                  </span>
                )}
                {author && (
                  <span className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-primary-100 text-primary-700">
                    Author: {author}
                    <button
                      onClick={() => {
                        setAuthorFilter('');
                        const params = new URLSearchParams(searchParams);
                        params.delete('author');
                        params.set('page', '1');
                        setSearchParams(params);
                      }}
                      className="ml-2 hover:text-primary-900"
                    >
                      <X size={14} />
                    </button>
                  </span>
                )}
                {tags && (
                  <span className="inline-flex items-center px-3 py-1 rounded-full text-sm bg-primary-100 text-primary-700">
                    Tags: {tags}
                    <button
                      onClick={() => {
                        setTagsFilter('');
                        const params = new URLSearchParams(searchParams);
                        params.delete('tags');
                        params.set('page', '1');
                        setSearchParams(params);
                      }}
                      className="ml-2 hover:text-primary-900"
                    >
                      <X size={14} />
                    </button>
                  </span>
                )}
              </div>
              <button
                onClick={clearFilters}
                className="text-sm text-primary-600 hover:text-primary-700"
              >
                Clear all filters
              </button>
            </div>
          )}
        </div>

        {isLoading ? (
          <div className="flex justify-center items-center py-12">
            <LoadingSpinner size="lg" />
          </div>
        ) : (
          <>
            {blogs.length === 0 ? (
              <div className="text-center py-12">
                <div className="text-gray-400 mb-4">
                  <Search className="mx-auto h-12 w-12" />
                </div>
                <h3 className="text-lg font-medium text-gray-900 mb-2">No blogs found</h3>
                <p className="text-gray-500 mb-4">
                  {hasActiveFilters
                    ? 'Try adjusting your search or filters'
                    : 'No blogs available yet. Be the first to write one!'}
                </p>
                {hasActiveFilters && (
                  <button onClick={clearFilters} className="btn btn-primary">
                    Clear Filters
                  </button>
                )}
              </div>
            ) : (
              <>
                <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                  {blogs.map((blog) => (
                    <BlogCard key={blog.id} blog={blog} />
                  ))}
                </div>

                {totalPages > 1 && (
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

export default Blogs;
