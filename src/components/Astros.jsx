import { useState, useEffect } from 'react';
import { astrosAPI } from '../services/api';
import { Search, ChevronLeft, ChevronRight, Eye, EyeOff, UserCheck, UserX, RefreshCw } from 'lucide-react';

const Astros = () => {
  const [astros, setAstros] = useState([]);
  const [loading, setLoading] = useState(true);
  const [pagination, setPagination] = useState({ page: 1, limit: 10, total: 0, totalPages: 0 });
  const [filters, setFilters] = useState({
    name: '',
    sort: 'asc',
    status: '',
    finalStatus: '',
    visible: '',
  });
  const [actionLoading, setActionLoading] = useState(null);
  const [refreshing, setRefreshing] = useState(false);
  const [astroAmountDisbursed, setAstroAmountDisbursed] = useState({});

  const fetchAstros = async () => {
    setLoading(true);
    try {
      const response = await astrosAPI.listAstros({
        page: pagination.page,
        limit: pagination.limit,
        name: filters.name,
        sort: filters.sort,
        ...(filters.status ? { status: filters.status } : {}),
        ...(filters.finalStatus ? { finalStatus: filters.finalStatus } : {}),
        ...(filters.visible ? { visible: filters.visible } : {}),
      });
      let data = response.data || [];

      // Frontend filtering to ensure dropdowns always work,
      // even if backend ignores some query params.
      if (filters.status) {
        if (filters.status === 'online') {
          data = data.filter((a) => a.status === 'online');
        } else if (filters.status === 'offline') {
          // Treat anything not strictly "online" as offline
          data = data.filter((a) => a.status !== 'online');
        }
      }
      if (filters.finalStatus) {
        data = data.filter((a) => a.finalStatus === filters.finalStatus);
      }
      if (filters.visible) {
        data = data.filter((a) => a.visible === filters.visible);
      }

      setAstros(data);
      setPagination(response.pagination || pagination);
    } catch (error) {
      console.error('Error fetching astros:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAstros();
  }, [pagination.page, pagination.limit, filters]);

  const handleStatusUpdate = async (astroId, currentFinalStatus, currentStatus) => {
    const newFinalStatus = currentFinalStatus === 'active' ? 'inactive' : 'active';
    setActionLoading(`status-${astroId}`);
    try {
      await astrosAPI.updateStatus(astroId, currentStatus, newFinalStatus);
      fetchAstros();
    } catch (error) {
      console.error('Error updating astro status:', error);
      alert('Failed to update astrologer status');
    } finally {
      setActionLoading(null);
    }
  };

  const handleOnlineToggle = async (astroId, newStatus, currentFinalStatus) => {
    setActionLoading(`online-${astroId}`);
    try {
      await astrosAPI.updateStatus(astroId, newStatus, currentFinalStatus);
      fetchAstros();
    } catch (error) {
      console.error('Error updating online status:', error);
      alert('Failed to update online status');
    } finally {
      setActionLoading(null);
    }
  };

  const handleVisibilityToggle = async (astroId, currentVisible) => {
    setActionLoading(`visibility-${astroId}`);
    try {
      const newVisible = currentVisible === 'visible' ? 'hidden' : 'visible';
      await astrosAPI.toggleVisibility(astroId, newVisible);
      fetchAstros();
    } catch (error) {
      console.error('Error updating visibility:', error);
      alert('Failed to update visibility');
    } finally {
      setActionLoading(null);
    }
  };

  const handleDisbursedSubmit = async (astroId, currentStatus, currentFinalStatus) => {
    const amount = astroAmountDisbursed[astroId];
    if (amount === undefined || amount === null || amount === '') {
      return;
    }

    setActionLoading(`disbursed-${astroId}`);
    try {
      await astrosAPI.updateStatus(astroId, currentStatus, currentFinalStatus, amount);
      await fetchAstros();
      // Clear the amount input for this astro so the placeholder is empty
      setAstroAmountDisbursed((prev) => {
        const next = { ...prev };
        delete next[astroId];
        return next;
      });
    } catch (error) {
      console.error('Error updating astro amount disbursed:', error);
      alert('Failed to update astro amount disbursed');
    } finally {
      setActionLoading(null);
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    setPagination({ ...pagination, page: 1 });
    fetchAstros();
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await fetchAstros();
    } finally {
      setRefreshing(false);
    }
  };

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-800">Astrologers Management</h1>
        <p className="text-gray-600 mt-2">Manage astrologer accounts, status, and visibility</p>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl shadow-sm p-4 mb-6 border border-gray-200">
        <form onSubmit={handleSearch} className="flex flex-col md:flex-row gap-4">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Search by name..."
              value={filters.name}
              onChange={(e) => setFilters({ ...filters, name: e.target.value })}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            />
          </div>
          <select
            value={filters.sort}
            onChange={(e) => setFilters({ ...filters, sort: e.target.value })}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
          >
            <option value="asc">Sort: A-Z</option>
            <option value="desc">Sort: Z-A</option>
          </select>
          <select
            value={filters.status}
            onChange={(e) => setFilters({ ...filters, status: e.target.value })}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            title="Filter by online/offline status"
          >
            <option value="">Status: All</option>
            <option value="online">Online</option>
            <option value="offline">Offline</option>
          </select>
          <select
            value={filters.finalStatus}
            onChange={(e) => setFilters({ ...filters, finalStatus: e.target.value })}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            title="Filter by active/inactive"
          >
            <option value="">Final Status: All</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
          <select
            value={filters.visible}
            onChange={(e) => setFilters({ ...filters, visible: e.target.value })}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
            title="Filter by visibility"
          >
            <option value="">Visibility: All</option>
            <option value="visible">Visible</option>
            <option value="hidden">Hidden</option>
          </select>
          <select
            value={pagination.limit}
            onChange={(e) => setPagination({ ...pagination, limit: parseInt(e.target.value), page: 1 })}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
          >
            <option value="10">10 per page</option>
            <option value="25">25 per page</option>
            <option value="50">50 per page</option>
            <option value="100">100 per page</option>
          </select>
          <button
            type="submit"
            className="px-6 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition"
          >
            Search
          </button>
          <button
            type="button"
            onClick={handleRefresh}
            disabled={refreshing}
            className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
            title="Refresh data from database"
          >
            <RefreshCw className={`w-4 h-4 ${refreshing ? 'animate-spin' : ''}`} />
            Refresh
          </button>
        </form>
      </div>

      {/* Astros Table */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          </div>
        ) : astros.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500">No astrologers found</p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 border-b border-gray-200">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Astro ID
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Name
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Final Status
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Visibility
                    </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Amount Disbursed
                </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {astros.map((astro) => (
                    <tr key={astro.astroId} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        {astro.astroId}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {astro.name}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span
                          className={`px-2 py-1 text-xs font-semibold rounded-full ${
                            astro.status === 'online'
                              ? 'bg-green-100 text-green-800'
                              : 'bg-gray-100 text-gray-800'
                          }`}
                        >
                          {astro.status || 'offline'}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span
                          className={`px-2 py-1 text-xs font-semibold rounded-full ${
                            astro.finalStatus === 'active'
                              ? 'bg-green-100 text-green-800'
                              : 'bg-red-100 text-red-800'
                          }`}
                        >
                          {astro.finalStatus}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span
                          className={`px-2 py-1 text-xs font-semibold rounded-full ${
                            astro.visible === 'visible'
                              ? 'bg-blue-100 text-blue-800'
                              : 'bg-gray-100 text-gray-800'
                          }`}
                        >
                          {astro.visible}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <div className="flex items-center space-x-2">
                          <input
                            type="number"
                            className="w-24 px-2 py-1 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                            placeholder="Amount"
                            value={astroAmountDisbursed[astro.astroId] ?? ''}
                            onChange={(e) =>
                              setAstroAmountDisbursed((prev) => ({
                                ...prev,
                                [astro.astroId]: e.target.value,
                              }))
                            }
                          />
                          <button
                            type="button"
                            onClick={() => handleDisbursedSubmit(astro.astroId, astro.status, astro.finalStatus)}
                            disabled={
                              actionLoading === `disbursed-${astro.astroId}` ||
                              !astroAmountDisbursed[astro.astroId]
                            }
                            className="px-3 py-1 text-xs bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
                          >
                            Go
                          </button>
                        </div>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                        {/* Final status toggle (active / inactive) */}
                        <button
                          onClick={() => handleStatusUpdate(astro.astroId, astro.finalStatus, astro.status)}
                          disabled={actionLoading === `status-${astro.astroId}`}
                          className={`inline-flex items-center space-x-1 px-3 py-1 rounded-lg transition ${
                            astro.finalStatus === 'active'
                              ? 'bg-red-100 text-red-700 hover:bg-red-200'
                              : 'bg-green-100 text-green-700 hover:bg-green-200'
                          } disabled:opacity-50`}
                        >
                          {astro.finalStatus === 'active' ? (
                            <>
                              <UserX className="w-4 h-4" />
                              <span>Deactivate</span>
                            </>
                          ) : (
                            <>
                              <UserCheck className="w-4 h-4" />
                              <span>Activate</span>
                            </>
                          )}
                        </button>

                        {/* Online/offline toggle (status: online/offline) */}
                        <button
                          onClick={() =>
                            handleOnlineToggle(
                              astro.astroId,
                              astro.status === 'online' ? 'offline' : 'online',
                              astro.finalStatus
                            )
                          }
                          disabled={actionLoading === `online-${astro.astroId}`}
                          className={`inline-flex items-center space-x-1 px-3 py-1 rounded-lg border transition ${
                            astro.status === 'online'
                              ? 'bg-green-50 text-green-700 border-green-200 hover:bg-green-100'
                              : 'bg-gray-50 text-gray-700 border-gray-200 hover:bg-gray-100'
                          } disabled:opacity-50`}
                        >
                          <span className="text-xs font-semibold">
                            {astro.status === 'online' ? 'On' : 'Off'}
                          </span>
                        </button>

                        {/* Visibility toggle */}
                        <button
                          onClick={() => handleVisibilityToggle(astro.astroId, astro.visible)}
                          disabled={actionLoading === `visibility-${astro.astroId}`}
                          className={`inline-flex items-center space-x-1 px-3 py-1 rounded-lg transition ${
                            astro.visible === 'visible'
                              ? 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                              : 'bg-blue-100 text-blue-700 hover:bg-blue-200'
                          } disabled:opacity-50`}
                        >
                          {astro.visible === 'visible' ? (
                            <>
                              <EyeOff className="w-4 h-4" />
                              <span>Hide</span>
                            </>
                          ) : (
                            <>
                              <Eye className="w-4 h-4" />
                              <span>Show</span>
                            </>
                          )}
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            <div className="px-6 py-4 border-t border-gray-200 flex flex-col sm:flex-row items-center justify-between gap-4">
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
          </>
        )}
      </div>
    </div>
  );
};

export default Astros;

