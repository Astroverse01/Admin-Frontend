import { useState, useEffect } from 'react';
import { complaintsAPI } from '../services/api';
import { ChevronLeft, ChevronRight, Check, X, RefreshCw, Copy } from 'lucide-react';
import { format } from 'date-fns';

const UserServiceComplaints = () => {
  const [complaints, setComplaints] = useState([]);
  const [loading, setLoading] = useState(true);
  const [pagination, setPagination] = useState({ page: 1, limit: 10, total: 0, totalPages: 0 });
  const [filters, setFilters] = useState({ serviceType: '', status: '' });
  const [actionLoading, setActionLoading] = useState(null);
  const [refundAmounts, setRefundAmounts] = useState({});
  const [refreshing, setRefreshing] = useState(false);
  const [showDetailsModal, setShowDetailsModal] = useState(false);
  const [complaintDetails, setComplaintDetails] = useState(null);
  const [detailsLoading, setDetailsLoading] = useState(false);
  const [copiedUrl, setCopiedUrl] = useState(false);

  const fetchComplaints = async () => {
    setLoading(true);
    try {
      const response = await complaintsAPI.listUserServiceComplaints({
        page: pagination.page,
        limit: pagination.limit,
        serviceType: filters.serviceType || undefined,
        status: filters.status || undefined,
      });
      setComplaints(response.data || []);
      setPagination(response.pagination || pagination);
    } catch (error) {
      console.error('Error fetching complaints:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchComplaints();
  }, [pagination.page, pagination.limit, filters]);

  const handleRefundChange = (orderId, field, value) => {
    // Only allow non-negative integers
    const numValue = parseInt(value);
    if (value === '' || (numValue >= 0 && !isNaN(numValue))) {
      setRefundAmounts({
        ...refundAmounts,
        [orderId]: {
          ...refundAmounts[orderId],
          [field]: value === '' ? 0 : numValue,
        },
      });
    }
  };

  const handleAccept = async (orderId) => {
    const userRefund = refundAmounts[orderId]?.userRefundMoney || 0;
    const astroRefund = refundAmounts[orderId]?.astroRefundMoney || 0;

    // Validate that values are non-negative integers
    if (userRefund < 0 || astroRefund < 0) {
      alert('Refund amounts cannot be negative');
      return;
    }

    const reason = prompt('Please enter a reason for accepting this complaint:');
    if (!reason || reason.trim() === '') {
      alert('Reason is required');
      return;
    }

    setActionLoading(orderId);
    try {
      await complaintsAPI.acceptRejectComplaint(orderId, {
        action: 'accept',
        reason: reason.trim(),
        userRefundMoney: userRefund,
        astroRefundMoney: astroRefund,
      });
      alert('Complaint accepted successfully');
      fetchComplaints();
      // Clear refund amounts for this order
      const newRefunds = { ...refundAmounts };
      delete newRefunds[orderId];
      setRefundAmounts(newRefunds);
    } catch (error) {
      console.error('Error accepting complaint:', error);
      alert(error.response?.data?.message || 'Failed to accept complaint');
    } finally {
      setActionLoading(null);
    }
  };

  const handleReject = async (orderId) => {
    const reason = prompt('Please enter a reason for rejecting this complaint:');
    if (!reason || reason.trim() === '') {
      alert('Reason is required');
      return;
    }

    setActionLoading(orderId);
    try {
      await complaintsAPI.acceptRejectComplaint(orderId, {
        action: 'reject',
        reason: reason.trim(),
        userRefundMoney: 0,
        astroRefundMoney: 0,
      });
      alert('Complaint rejected successfully');
      fetchComplaints();
    } catch (error) {
      console.error('Error rejecting complaint:', error);
      alert(error.response?.data?.message || 'Failed to reject complaint');
    } finally {
      setActionLoading(null);
    }
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      await fetchComplaints();
    } finally {
      setRefreshing(false);
    }
  };

  // Map service type from GET /admin/user-service-complaints response format
  // to the format required by GET /admin/user-service-complaints/:serviceType/:orderId
  // chat -> chat, call -> ivrCall, video -> videoCall
  const mapServiceTypeToAPI = (serviceType) => {
    if (!serviceType) return 'chat'; // default fallback
    
    const normalized = serviceType.toLowerCase();
    const mapping = {
      'chat': 'chat',      // chat stays as chat
      'call': 'ivrCall',   // call maps to ivrCall for the API
      'video': 'videoCall', // video maps to videoCall for the API
    };
    return mapping[normalized] || normalized;
  };

  const handleOrderIdClick = async (orderId, serviceType) => {
    setDetailsLoading(true);
    setShowDetailsModal(true);
    try {
      const apiServiceType = mapServiceTypeToAPI(serviceType);
      const response = await complaintsAPI.getComplaintDetails(apiServiceType, orderId);
      setComplaintDetails(response.data || response);
    } catch (error) {
      console.error('Error fetching complaint details:', error);
      const errorMessage = error.response?.data?.message || 'Failed to fetch complaint details';
      
      // Check if error is about serviceType mismatch
      if (errorMessage.includes('serviceType mismatch') || errorMessage.includes('but serviceType mismatch')) {
        alert('no url found for this orderId');
      } else {
        alert(errorMessage);
      }
      setShowDetailsModal(false);
    } finally {
      setDetailsLoading(false);
    }
  };

  const handleCopyUrl = async (url) => {
    try {
      await navigator.clipboard.writeText(url);
      setCopiedUrl(true);
      setTimeout(() => setCopiedUrl(false), 2000);
    } catch (error) {
      console.error('Failed to copy URL:', error);
      // Fallback for older browsers
      const textArea = document.createElement('textarea');
      textArea.value = url;
      textArea.style.position = 'fixed';
      textArea.style.opacity = '0';
      document.body.appendChild(textArea);
      textArea.select();
      try {
        document.execCommand('copy');
        setCopiedUrl(true);
        setTimeout(() => setCopiedUrl(false), 2000);
      } catch (err) {
        alert('Failed to copy URL. Please copy it manually.');
      }
      document.body.removeChild(textArea);
    }
  };

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-800">User Service Complaints</h1>
        <p className="text-gray-600 mt-2">Manage service-related complaints from users</p>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-xl shadow-sm p-4 mb-6 border border-gray-200">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            setPagination({ ...pagination, page: 1 });
            fetchComplaints();
          }}
          className="flex flex-col md:flex-row gap-4"
        >
          <select
            value={filters.serviceType}
            onChange={(e) => setFilters({ ...filters, serviceType: e.target.value })}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
          >
            <option value="">All Service Types</option>
            <option value="chat">Chat</option>
            <option value="ivrCall">IVR Call</option>
            <option value="videoCall">Video Call</option>
          </select>
          <select
            value={filters.status}
            onChange={(e) => setFilters({ ...filters, status: e.target.value })}
            className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
          >
            <option value="">All Statuses</option>
            <option value="open">Open</option>
            <option value="closed">Closed</option>
            <option value="rejected">Rejected</option>
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
            Filter
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

      {/* Complaints Table */}
      <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          </div>
        ) : complaints.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500">No complaints found</p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 border-b border-gray-200">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Order ID
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Service Type
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      User Name
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Astrologer Name
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Created At
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      User Refund (₹)
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Astro Refund (₹)
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {complaints.map((complaint) => {
                    const isOpen = complaint.status === 'open';
                    const isLoading = actionLoading === complaint.orderId;
                    
                    return (
                      <tr key={complaint.orderId} className="hover:bg-gray-50">
                        <td className="px-4 py-4 text-sm text-gray-900">
                          <button
                            onClick={() => handleOrderIdClick(complaint.orderId, complaint.serviceType)}
                            className="max-w-[150px] truncate text-primary-600 hover:text-primary-800 hover:underline cursor-pointer"
                            title={`Click to view details: ${complaint.orderId}`}
                          >
                            {complaint.orderId}
                          </button>
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-900">
                          {complaint.serviceType}
                        </td>
                        <td className="px-4 py-4 text-sm text-gray-900">
                          <div className="max-w-[150px]" title={`${complaint.userName} (${complaint.userId})`}>
                            <span className="font-medium">{complaint.userName}</span>
                            <p className="text-xs text-gray-500 truncate">{complaint.userId}</p>
                          </div>
                        </td>
                        <td className="px-4 py-4 text-sm text-gray-900">
                          <div className="max-w-[150px]" title={`${complaint.astroName} (${complaint.astroId})`}>
                            <span className="font-medium">{complaint.astroName}</span>
                            <p className="text-xs text-gray-500 truncate">{complaint.astroId}</p>
                          </div>
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap">
                          <span
                            className={`px-2 py-1 text-xs font-semibold rounded-full ${
                              complaint.status === 'open'
                                ? 'bg-yellow-100 text-yellow-800'
                                : complaint.status === 'closed'
                                ? 'bg-green-100 text-green-800'
                                : 'bg-red-100 text-red-800'
                            }`}
                          >
                            {complaint.status}
                          </span>
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap text-sm text-gray-500">
                          {complaint.createdOn ? format(new Date(complaint.createdOn), 'MMM dd, yyyy') : 'N/A'}
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap">
                          <input
                            type="number"
                            min="0"
                            step="1"
                            value={refundAmounts[complaint.orderId]?.userRefundMoney || ''}
                            onChange={(e) => handleRefundChange(complaint.orderId, 'userRefundMoney', e.target.value)}
                            disabled={!isOpen || isLoading}
                            placeholder="0"
                            className="w-24 px-2 py-1 border border-gray-300 rounded focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none disabled:bg-gray-100 disabled:cursor-not-allowed"
                          />
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap">
                          <input
                            type="number"
                            min="0"
                            step="1"
                            value={refundAmounts[complaint.orderId]?.astroRefundMoney || ''}
                            onChange={(e) => handleRefundChange(complaint.orderId, 'astroRefundMoney', e.target.value)}
                            disabled={!isOpen || isLoading}
                            placeholder="0"
                            className="w-24 px-2 py-1 border border-gray-300 rounded focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none disabled:bg-gray-100 disabled:cursor-not-allowed"
                          />
                        </td>
                        <td className="px-4 py-4 whitespace-nowrap text-sm">
                          {isOpen ? (
                            <div className="flex space-x-2">
                              <button
                                onClick={() => handleAccept(complaint.orderId)}
                                disabled={isLoading}
                                className="flex items-center space-x-1 px-3 py-1 bg-green-100 text-green-700 rounded-lg hover:bg-green-200 transition disabled:opacity-50 disabled:cursor-not-allowed"
                                title="Accept"
                              >
                                <Check className="w-4 h-4" />
                                <span>Accept</span>
                              </button>
                              <button
                                onClick={() => handleReject(complaint.orderId)}
                                disabled={isLoading}
                                className="flex items-center space-x-1 px-3 py-1 bg-red-100 text-red-700 rounded-lg hover:bg-red-200 transition disabled:opacity-50 disabled:cursor-not-allowed"
                                title="Reject"
                              >
                                <X className="w-4 h-4" />
                                <span>Reject</span>
                              </button>
                            </div>
                          ) : (
                            <span className="text-gray-400 text-sm">Processed</span>
                          )}
                        </td>
                      </tr>
                    );
                  })}
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

      {/* Details Modal */}
      {showDetailsModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
            <div className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-2xl font-bold text-gray-800">Complaint Details</h2>
                <button
                  onClick={() => {
                    setShowDetailsModal(false);
                    setComplaintDetails(null);
                  }}
                  className="text-gray-500 hover:text-gray-700"
                >
                  <X className="w-6 h-6" />
                </button>
              </div>

              {detailsLoading ? (
                <div className="flex items-center justify-center h-64">
                  <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
                </div>
              ) : complaintDetails ? (
                <div className="space-y-4">
                  {/* Basic Information */}
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Order ID</label>
                      <p className="text-sm text-gray-900 break-all">{complaintDetails.ivrId || complaintDetails.chatId || complaintDetails.videoId || 'N/A'}</p>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">User ID</label>
                      <p className="text-sm text-gray-900 break-all">{complaintDetails.userId || 'N/A'}</p>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Astrologer ID</label>
                      <p className="text-sm text-gray-900 break-all">{complaintDetails.astroId || 'N/A'}</p>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
                      <p className="text-sm text-gray-900">{complaintDetails.lastStatus || complaintDetails.status || 'N/A'}</p>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Type</label>
                      <p className="text-sm text-gray-900">{complaintDetails.type || 'N/A'}</p>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Category Type</label>
                      <p className="text-sm text-gray-900">{complaintDetails.categoryType || 'N/A'}</p>
                    </div>
                    {complaintDetails.ratePerMintue && (
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Rate Per Minute</label>
                        <p className="text-sm text-gray-900">₹{complaintDetails.ratePerMintue}</p>
                      </div>
                    )}
                    {complaintDetails.paymentReceived !== undefined && (
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Payment Received</label>
                        <p className="text-sm text-gray-900">{complaintDetails.paymentReceived ? 'Yes' : 'No'}</p>
                      </div>
                    )}
                    {complaintDetails.createdOn && (
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Created On</label>
                        <p className="text-sm text-gray-900">{format(new Date(complaintDetails.createdOn), 'MMM dd, yyyy HH:mm:ss')}</p>
                      </div>
                    )}
                    {complaintDetails.updatedOn && (
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Updated On</label>
                        <p className="text-sm text-gray-900">{format(new Date(complaintDetails.updatedOn), 'MMM dd, yyyy HH:mm:ss')}</p>
                      </div>
                    )}
                  </div>

                  {/* URL/Audio Link */}
                  {complaintDetails.url && (
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Audio/Media URL</label>
                      <div className="flex items-center gap-2">
                        <input
                          type="text"
                          value={complaintDetails.url}
                          readOnly
                          className="flex-1 px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 text-sm text-gray-900 break-all focus:outline-none focus:ring-2 focus:ring-primary-500"
                        />
                        <button
                          onClick={() => handleCopyUrl(complaintDetails.url)}
                          className="flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition"
                          title="Copy URL"
                        >
                          <Copy className="w-4 h-4" />
                          {copiedUrl ? 'Copied!' : 'Copy'}
                        </button>
                      </div>
                    </div>
                  )}

                  {/* Conversation */}
                  {complaintDetails.conversation && complaintDetails.conversation.length > 0 && (
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Conversation</label>
                      <div className="border border-gray-200 rounded-lg p-4 max-h-64 overflow-y-auto">
                        {complaintDetails.conversation.map((msg, index) => (
                          <div key={index} className="mb-3 last:mb-0">
                            <div className="flex items-start gap-2">
                              <span className={`text-xs font-semibold ${msg.isAstrolger ? 'text-blue-600' : 'text-green-600'}`}>
                                {msg.isAstrolger ? 'Astrologer' : 'User'}
                              </span>
                              <span className="text-xs text-gray-500">
                                {msg.timestamp ? format(new Date(msg.timestamp), 'MMM dd, yyyy HH:mm:ss') : 'N/A'}
                              </span>
                            </div>
                            {msg.message && (
                              <div className="mt-1 text-sm text-gray-900 whitespace-pre-wrap">
                                {msg.message.text || JSON.stringify(msg.message)}
                              </div>
                            )}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Status History */}
                  {complaintDetails.status && Array.isArray(complaintDetails.status) && complaintDetails.status.length > 0 && (
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Status History</label>
                      <div className="border border-gray-200 rounded-lg p-4">
                        {complaintDetails.status.map((statusItem, index) => (
                          <div key={index} className="mb-2 last:mb-0">
                            <span className="text-sm font-medium text-gray-900">{statusItem.type}</span>
                            {statusItem.createdOn && (
                              <span className="text-xs text-gray-500 ml-2">
                                - {format(new Date(statusItem.createdOn), 'MMM dd, yyyy HH:mm:ss')}
                              </span>
                            )}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Raw Data (for debugging) */}
                  <details className="mt-4">
                    <summary className="cursor-pointer text-sm font-medium text-gray-700 hover:text-gray-900">
                      View Raw Data
                    </summary>
                    <pre className="mt-2 p-4 bg-gray-50 rounded-lg text-xs overflow-x-auto">
                      {JSON.stringify(complaintDetails, null, 2)}
                    </pre>
                  </details>
                </div>
              ) : (
                <div className="text-center py-8">
                  <p className="text-gray-500">No details available</p>
                </div>
              )}

              <div className="mt-6 flex justify-end">
                <button
                  onClick={() => {
                    setShowDetailsModal(false);
                    setComplaintDetails(null);
                  }}
                  className="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition"
                >
                  Close
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default UserServiceComplaints;
