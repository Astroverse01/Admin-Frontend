import { useState, useEffect } from 'react';
import { horoscopesAPI } from '../services/api';
import { ChevronLeft, ChevronRight, Plus, Edit, Trash2, X, RefreshCw } from 'lucide-react';

const Horoscopes = () => {
  const [horoscopes, setHoroscopes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [pagination, setPagination] = useState({ page: 1, limit: 10, total: 0, totalPages: 0 });
  const [showModal, setShowModal] = useState(false);
  const [showBulkModal, setShowBulkModal] = useState(false);
  const [editingHoroscope, setEditingHoroscope] = useState(null);
  const [formData, setFormData] = useState({
    signName: '',
    description: '',
    date: '',
    isAsctive: 1,
  });
  const [bulkData, setBulkData] = useState('');
  const [actionLoading, setActionLoading] = useState(false);
  const [refreshing, setRefreshing] = useState(false);

  const fetchHoroscopes = async () => {
    setLoading(true);
    try {
      const response = await horoscopesAPI.listHoroscopes({
        page: pagination.page,
        limit: pagination.limit,
      });
      setHoroscopes(response.data || []);
      setPagination(response.pagination || pagination);
    } catch (error) {
      console.error('Error fetching horoscopes:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchHoroscopes();
  }, [pagination.page]);

  const handleEdit = (horoscope) => {
    setEditingHoroscope(horoscope);
    setFormData({
      signName: horoscope.signName,
      description: horoscope.description,
      date: horoscope.date,
      isAsctive: horoscope.isAsctive || 1,
    });
    setShowModal(true);
  };

  const handleCreate = () => {
    setEditingHoroscope(null);
    setFormData({
      signName: '',
      description: '',
      date: '',
      isAsctive: 1,
    });
    setShowModal(true);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setActionLoading(true);
    try {
      if (editingHoroscope) {
        await horoscopesAPI.updateHoroscope(editingHoroscope.horoscopeId, formData);
        alert('Horoscope updated successfully');
      } else {
        // Use bulk create for single horoscope
        await horoscopesAPI.bulkCreate([formData]);
        alert('Horoscope created successfully');
      }
      setShowModal(false);
      fetchHoroscopes();
    } catch (error) {
      console.error('Error saving horoscope:', error);
      alert('Failed to save horoscope');
    } finally {
      setActionLoading(false);
    }
  };

  const handleDelete = async (horoscopeId) => {
    if (!confirm('Are you sure you want to delete this horoscope?')) return;

    try {
      await horoscopesAPI.deleteHoroscope(horoscopeId);
      alert('Horoscope deleted successfully');
      fetchHoroscopes();
    } catch (error) {
      console.error('Error deleting horoscope:', error);
      alert('Failed to delete horoscope');
    }
  };

  const handleBulkCreate = async () => {
    try {
      const lines = bulkData.split('\n').filter((line) => line.trim());
      const horoscopes = lines.map((line) => {
        const parts = line.split('|').map((p) => p.trim());
        if (parts.length < 4) {
          throw new Error('Invalid format. Expected: signName|description|date|isActive');
        }
        return {
          signName: parts[0],
          description: parts[1],
          date: parts[2],
          isAsctive: parseInt(parts[3]) || 1,
        };
      });

      await horoscopesAPI.bulkCreate(horoscopes);
      alert('Horoscopes created successfully');
      setShowBulkModal(false);
      setBulkData('');
      fetchHoroscopes();
    } catch (error) {
      console.error('Error bulk creating horoscopes:', error);
      alert('Failed to create horoscopes. Please check the format.');
    }
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await fetchHoroscopes();
    } finally {
      setRefreshing(false);
    }
  };

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-800">Horoscopes Management</h1>
          <p className="text-gray-600 mt-2">Manage daily horoscopes for all zodiac signs</p>
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
            onClick={() => setShowBulkModal(true)}
            className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition flex items-center space-x-2"
          >
            <Plus className="w-4 h-4" />
            <span>Bulk Create</span>
          </button>
          <button
            onClick={handleCreate}
            className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition flex items-center space-x-2"
          >
            <Plus className="w-4 h-4" />
            <span>Add Horoscope</span>
          </button>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          </div>
        ) : horoscopes.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500">No horoscopes found</p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 border-b border-gray-200">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Sign Name
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Description
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Date
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
                  {horoscopes.map((horoscope) => (
                    <tr key={horoscope.horoscopeId} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {horoscope.signName}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-900 max-w-md truncate">
                        {horoscope.description}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {horoscope.date}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span
                          className={`px-2 py-1 text-xs font-semibold rounded-full ${
                            horoscope.isAsctive === 1
                              ? 'bg-green-100 text-green-800'
                              : 'bg-gray-100 text-gray-800'
                          }`}
                        >
                          {horoscope.isAsctive === 1 ? 'Active' : 'Inactive'}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                        <button
                          onClick={() => handleEdit(horoscope)}
                          className="px-3 py-1 bg-blue-100 text-blue-700 rounded-lg hover:bg-blue-200 transition flex items-center space-x-1"
                        >
                          <Edit className="w-4 h-4" />
                          <span>Edit</span>
                        </button>
                        <button
                          onClick={() => handleDelete(horoscope.horoscopeId)}
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
                  {Math.min(pagination.page * pagination.limit, pagination.total)} of{' '}
                  {pagination.total} results
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

      {/* Create/Edit Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full">
            <div className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-2xl font-bold text-gray-800">
                  {editingHoroscope ? 'Edit Horoscope' : 'Add Horoscope'}
                </h2>
                <button
                  onClick={() => {
                    setShowModal(false);
                    setEditingHoroscope(null);
                  }}
                  className="text-gray-500 hover:text-gray-700"
                >
                  <X className="w-6 h-6" />
                </button>
              </div>

              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Sign Name *</label>
                  <input
                    type="text"
                    value={formData.signName}
                    onChange={(e) => setFormData({ ...formData, signName: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Description *</label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                    rows="4"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Date *</label>
                  <input
                    type="date"
                    value={formData.date}
                    onChange={(e) => setFormData({ ...formData, date: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Status *</label>
                  <select
                    value={formData.isAsctive}
                    onChange={(e) => setFormData({ ...formData, isAsctive: parseInt(e.target.value) })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                  >
                    <option value={1}>Active</option>
                    <option value={0}>Inactive</option>
                  </select>
                </div>
                <div className="flex space-x-4">
                  <button
                    type="submit"
                    disabled={actionLoading}
                    className="flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition disabled:opacity-50"
                  >
                    {actionLoading ? 'Saving...' : editingHoroscope ? 'Update' : 'Create'}
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setShowModal(false);
                      setEditingHoroscope(null);
                    }}
                    className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition"
                  >
                    Cancel
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}

      {/* Bulk Create Modal */}
      {showBulkModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full">
            <div className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-2xl font-bold text-gray-800">Bulk Create Horoscopes</h2>
                <button
                  onClick={() => {
                    setShowBulkModal(false);
                    setBulkData('');
                  }}
                  className="text-gray-500 hover:text-gray-700"
                >
                  <X className="w-6 h-6" />
                </button>
              </div>

              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Format: signName|description|date|isActive (one per line)
                </label>
                <textarea
                  value={bulkData}
                  onChange={(e) => setBulkData(e.target.value)}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none font-mono text-sm"
                  rows="10"
                  placeholder="Aries|Today is a great day for new beginnings...|2024-01-15|1&#10;Taurus|Focus on your goals today...|2024-01-15|1"
                />
              </div>

              <div className="flex space-x-4">
                <button
                  onClick={handleBulkCreate}
                  className="flex-1 px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition"
                >
                  Create All
                </button>
                <button
                  onClick={() => {
                    setShowBulkModal(false);
                    setBulkData('');
                  }}
                  className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Horoscopes;

