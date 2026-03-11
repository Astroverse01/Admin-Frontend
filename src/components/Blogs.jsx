import { useEffect, useState } from 'react';
import { blogsAPI } from '../services/api';
import { Search, Plus, Edit, Trash2, ChevronLeft, ChevronRight, X, Image as ImageIcon, RefreshCw } from 'lucide-react';

const defaultFormState = {
  title: '',
  slug: '',
  contentType: '',
  readingTimeMin: '',
  contentBody: '',
  seoTitle: '',
  seoDescription: '',
  status: 'draft',
};

const Blogs = () => {
  const [blogs, setBlogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [pagination, setPagination] = useState({ page: 1, limit: 10, total: 0, totalPages: 0 });
  const [filters, setFilters] = useState({ q: '', status: '' });
  const [showModal, setShowModal] = useState(false);
  const [editingBlog, setEditingBlog] = useState(null);
  const [formData, setFormData] = useState(defaultFormState);
  const [coverFile, setCoverFile] = useState(null);
  const [actionLoading, setActionLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);

  const fetchBlogs = async () => {
    setLoading(true);
    try {
      const response = await blogsAPI.listBlogs({
        page: pagination.page,
        limit: pagination.limit,
        q: filters.q || '',
        status: filters.status || '',
      });

      setBlogs(response.data || []);
      setPagination(response.pagination || pagination);
    } catch (error) {
      console.error('Error fetching blogs:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchBlogs();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pagination.page, pagination.limit, filters.q, filters.status]);

  const handleOpenCreate = () => {
    setEditingBlog(null);
    setFormData(defaultFormState);
    setCoverFile(null);
    setShowModal(true);
  };

  const handleOpenEdit = (blog) => {
    setEditingBlog(blog);
    setFormData({
      title: blog.title || '',
      slug: blog.slug || '',
      contentType: blog.contentType || '',
      readingTimeMin: blog.readingTimeMin != null ? String(blog.readingTimeMin) : '',
      // Prefer contentBody, fallback to content for safety
      contentBody: blog.contentBody || blog.content || '',
      seoTitle: blog.seoTitle || '',
      seoDescription: blog.seoDescription || '',
      status: blog.status || 'draft',
    });
    setCoverFile(null);
    setShowModal(true);
  };

  const handleCloseModal = () => {
    if (actionLoading) return;
    setShowModal(false);
    setEditingBlog(null);
    setFormData(defaultFormState);
    setCoverFile(null);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Basic front-end validation for required fields
    if (
      !formData.title.trim() ||
      !formData.slug.trim() ||
      !formData.seoTitle.trim() ||
      !formData.seoDescription.trim() ||
      !formData.contentBody.trim()
    ) {
      alert('Please fill all mandatory fields.');
      return;
    }

    if (formData.contentType && !formData.readingTimeMin) {
      alert('Please provide reading time when content type is selected.');
      return;
    }

    if (!editingBlog && !coverFile) {
      alert('Cover image is mandatory for creating a blog.');
      return;
    }

    setActionLoading(true);
    try {
      const payload = {
        title: formData.title.trim(),
        slug: formData.slug.trim(),
        contentType: formData.contentType || undefined,
        readingTimeMin: formData.contentType
          ? Number.isNaN(Number(formData.readingTimeMin))
            ? formData.readingTimeMin
            : Number(formData.readingTimeMin)
          : undefined,
        // Align with backend "contentBody" field
        contentBody: formData.contentBody.trim(),
        seoTitle: formData.seoTitle.trim(),
        seoDescription: formData.seoDescription.trim(),
        status: formData.status,
      };

      if (!editingBlog) {
        await blogsAPI.createBlog({
          ...payload,
          coverImage: coverFile,
        });
        alert('Blog created successfully');
      } else {
        await blogsAPI.updateBlog(editingBlog.blogId || editingBlog.id, payload);
        alert('Blog updated successfully');
      }

      handleCloseModal();
      fetchBlogs();
    } catch (error) {
      console.error('Error saving blog:', error);
      alert('Failed to save blog');
    } finally {
      setActionLoading(false);
    }
  };

  const handleDelete = async (blog) => {
    if (!window.confirm('Are you sure you want to delete this blog? This will also delete its cover image.')) {
      return;
    }

    try {
      await blogsAPI.deleteBlog(blog.blogId || blog.id);
      alert('Blog deleted successfully');
      fetchBlogs();
    } catch (error) {
      console.error('Error deleting blog:', error);
      alert('Failed to delete blog');
    }
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await fetchBlogs();
    } finally {
      setRefreshing(false);
    }
  };

  const isReadingTimeDisabled = !formData.contentType;

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-800">Blogs Management</h1>
          <p className="text-gray-600 mt-2">Create, edit, and manage blog posts</p>
        </div>
        <div className="flex space-x-2">
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition flex items-center space-x-2 disabled:opacity-50 disabled:cursor-not-allowed"
            title="Refresh data from database"
          >
            <RefreshCw className={`w-4 h-4 ${refreshing ? 'animate-spin' : ''}`} />
            <span>Refresh</span>
          </button>
          <button
            onClick={handleOpenCreate}
            className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition flex items-center space-x-2"
          >
            <Plus className="w-4 h-4" />
            <span>Create Blog</span>
          </button>
        </div>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl shadow-sm p-4 mb-6 border border-gray-200">
        <div className="flex flex-col md:flex-row gap-4 items-center">
          <div className="flex-1 w-full relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Search by title or slug..."
              value={filters.q}
              onChange={(e) => setFilters({ ...filters, q: e.target.value })}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            />
          </div>
          <select
            value={filters.status}
            onChange={(e) => setFilters({ ...filters, status: e.target.value })}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
          >
            <option value="">All Statuses</option>
            <option value="draft">Draft</option>
            <option value="published">Published</option>
          </select>
          <select
            value={pagination.limit}
            onChange={(e) => setPagination({ ...pagination, limit: parseInt(e.target.value, 10), page: 1 })}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
          >
            <option value="10">10 per page</option>
            <option value="25">25 per page</option>
            <option value="50">50 per page</option>
          </select>
        </div>
      </div>

      {/* Blogs list */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600" />
          </div>
        ) : blogs.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500">No blogs found</p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 border-b border-gray-200">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Title
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Slug
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Content Type
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Reading Time
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {blogs.map((blog) => (
                    <tr key={blog.blogId || blog.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {blog.title}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                        {blog.slug}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                        {blog.contentType || '-'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                        {blog.readingTimeMin != null ? `${blog.readingTimeMin} min` : '-'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span
                          className={`px-2 py-1 text-xs font-semibold rounded-full ${
                            blog.status === 'published'
                              ? 'bg-green-100 text-green-800'
                              : 'bg-yellow-100 text-yellow-800'
                          }`}
                        >
                          {blog.status || 'draft'}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                        <button
                          onClick={() => handleOpenEdit(blog)}
                          className="px-3 py-1 bg-blue-100 text-blue-700 rounded-lg hover:bg-blue-200 transition flex items-center space-x-1"
                        >
                          <Edit className="w-4 h-4" />
                          <span>Edit</span>
                        </button>
                        <button
                          onClick={() => handleDelete(blog)}
                          className="px-3 py-1 bg-red-100 text-red-700 rounded-lg hover:bg-red-200 transition flex items-center space-x-1"
                        >
                          <Trash2 className="w-4 h-4" />
                          <span>Delete</span>
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {pagination.totalPages > 1 && (
              <div className="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
                <div className="text-sm text-gray-700">
                  Showing {((pagination.page - 1) * pagination.limit) + 1} to{' '}
                  {Math.min(pagination.page * pagination.limit, pagination.total)} of {pagination.total} results
                </div>
                <div className="flex items-center space-x-2">
                  <button
                    onClick={() => setPagination({ ...pagination, page: pagination.page - 1 })}
                    disabled={pagination.page === 1}
                    className="p-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <ChevronLeft className="w-5 h-5" />
                  </button>
                  <span className="text-sm text-gray-700">
                    Page {pagination.page} of {pagination.totalPages}
                  </span>
                  <button
                    onClick={() => setPagination({ ...pagination, page: pagination.page + 1 })}
                    disabled={pagination.page >= pagination.totalPages}
                    className="p-2 border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <ChevronRight className="w-5 h-5" />
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </div>

      {/* Create / Edit Blog Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">
            <div className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-2xl font-bold text-gray-800">
                  {editingBlog ? 'Edit Blog Post' : 'Create New Blog Post'}
                </h2>
                <button
                  onClick={handleCloseModal}
                  className="text-gray-500 hover:text-gray-700"
                >
                  <X className="w-6 h-6" />
                </button>
              </div>

              <form onSubmit={handleSubmit} className="space-y-6">
                {/* Title & Slug */}
                <div className="grid md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Title <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="text"
                      value={formData.title}
                      onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                      placeholder="Enter a compelling blog title..."
                      required
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Slug <span className="text-red-500">*</span>
                    </label>
                    <input
                      type="text"
                      value={formData.slug}
                      onChange={(e) => setFormData({ ...formData, slug: e.target.value })}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                      placeholder="url-friendly-slug"
                      required
                    />
                    <p className="mt-1 text-xs text-gray-500">URL-friendly version of your title</p>
                  </div>
                </div>

                {/* Content settings */}
                <div className="border border-gray-200 rounded-lg p-4">
                  <h3 className="text-sm font-semibold text-gray-800 mb-3">Content</h3>
                  <div className="grid md:grid-cols-3 gap-4 items-end">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Content Type</label>
                      <select
                        value={formData.contentType}
                        onChange={(e) =>
                          setFormData({
                            ...formData,
                            contentType: e.target.value,
                            // Clear reading time if content type is cleared
                            readingTimeMin: e.target.value ? formData.readingTimeMin : '',
                          })
                        }
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                      >
                        <option value="">Select type</option>
                        <option value="mdx">MDX</option>
                        <option value="markdown">Markdown</option>
                        <option value="html">HTML</option>
                      </select>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Reading Time (minutes) <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="number"
                        min="1"
                        value={formData.readingTimeMin}
                        onChange={(e) => setFormData({ ...formData, readingTimeMin: e.target.value })}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none disabled:bg-gray-100"
                        disabled={isReadingTimeDisabled}
                        required={!!formData.contentType}
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Status
                      </label>
                      <select
                        value={formData.status}
                        onChange={(e) => setFormData({ ...formData, status: e.target.value })}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                      >
                        <option value="draft">Draft</option>
                        <option value="published">Published</option>
                      </select>
                    </div>
                  </div>

                  <div className="mt-4">
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Content Body <span className="text-red-500">*</span>
                    </label>
                    <textarea
                      value={formData.contentBody}
                      onChange={(e) => setFormData({ ...formData, contentBody: e.target.value })}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                      rows="6"
                      placeholder="Write your MDX content here..."
                      required
                    />
                  </div>
                </div>

                {/* Cover image */}
                {!editingBlog && (
                  <div className="border border-dashed border-gray-300 rounded-lg p-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Cover Image <span className="text-red-500">*</span>
                    </label>
                    <div className="flex items-center space-x-4">
                      <div className="flex-1">
                        <label className="flex flex-col items-center justify-center px-4 py-6 border-2 border-dashed border-gray-300 rounded-lg cursor-pointer hover:border-primary-500 transition">
                          <ImageIcon className="w-8 h-8 text-gray-400 mb-2" />
                          <span className="text-sm text-gray-600">
                            Click to upload cover image
                          </span>
                          <input
                            type="file"
                            accept="image/*"
                            className="hidden"
                            onChange={(e) => {
                              const file = e.target.files?.[0];
                              if (file) {
                                setCoverFile(file);
                              }
                            }}
                          />
                        </label>
                        {coverFile && (
                          <p className="mt-2 text-xs text-gray-600">
                            Selected: <span className="font-medium">{coverFile.name}</span>
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                )}

                {/* SEO settings */}
                <div className="border border-gray-200 rounded-lg p-4">
                  <h3 className="text-sm font-semibold text-gray-800 mb-3">SEO Settings</h3>
                  <div className="grid md:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        SEO Title <span className="text-red-500">*</span>
                      </label>
                      <input
                        type="text"
                        value={formData.seoTitle}
                        onChange={(e) => setFormData({ ...formData, seoTitle: e.target.value })}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                        placeholder="Custom meta title for SEO..."
                        required
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        SEO Description <span className="text-red-500">*</span>
                      </label>
                      <textarea
                        value={formData.seoDescription}
                        onChange={(e) => setFormData({ ...formData, seoDescription: e.target.value })}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                        rows="3"
                        placeholder="Custom meta description for SEO..."
                        required
                      />
                    </div>
                  </div>
                </div>

                <div className="flex justify-end space-x-4 pt-2">
                  <button
                    type="button"
                    onClick={handleCloseModal}
                    className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition"
                    disabled={actionLoading}
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={actionLoading}
                    className="px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition disabled:opacity-50"
                  >
                    {actionLoading ? 'Saving...' : editingBlog ? 'Update Blog' : 'Create Blog'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Blogs;

