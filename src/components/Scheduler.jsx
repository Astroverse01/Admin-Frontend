import { useState } from 'react';
import { schedulerAPI } from '../services/api';
import { Calendar, Download, CheckCircle, XCircle, FileText } from 'lucide-react';

const Scheduler = () => {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [generatedFiles, setGeneratedFiles] = useState([]);

  const handleGenerateReports = async () => {
    if (!startDate || !endDate) {
      alert('Please select both start and end dates');
      return;
    }

    // Validate date range
    const start = new Date(startDate);
    const end = new Date(endDate);
    if (start > end) {
      alert('Start date must be before or equal to end date');
      return;
    }

    setLoading(true);
    setResult(null);
    setGeneratedFiles([]);
    
    try {
      const response = await schedulerAPI.generateReports(startDate, endDate);
      setResult({
        success: true,
        message: response.message || 'Reports generated successfully',
      });
      setGeneratedFiles(response.files || []);
    } catch (error) {
      setResult({
        success: false,
        message: error.response?.data?.message || 'Failed to generate reports. Please try again.',
      });
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = async (fileName) => {
    try {
      const response = await schedulerAPI.downloadReport(fileName);
      
      // Create a blob from the response data
      const blob = new Blob([response.data], { type: 'text/csv' });
      
      // Create a temporary URL for the blob
      const url = window.URL.createObjectURL(blob);
      
      // Create a temporary anchor element and trigger download
      const link = document.createElement('a');
      link.href = url;
      link.download = fileName;
      document.body.appendChild(link);
      link.click();
      
      // Clean up
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Error downloading file:', error);
      alert('Failed to download file. Please try again.');
    }
  };


  // Get today's date in YYYY-MM-DD format
  const getTodayDate = () => {
    const today = new Date();
    return today.toISOString().split('T')[0];
  };

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-800">Report Generator</h1>
        <p className="text-gray-600 mt-2">Generate and download CSV reports for a specific date range</p>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-8">
        <div className="max-w-3xl mx-auto">
          <div className="flex items-center justify-center mb-6">
            <div className="bg-primary-100 p-4 rounded-full">
              <Calendar className="w-12 h-12 text-primary-600" />
            </div>
          </div>

          <div className="text-center mb-8">
            <h2 className="text-2xl font-semibold text-gray-800 mb-2">Date Range Report Generator</h2>
            <p className="text-gray-600">
              Select a date range to generate CSV reports from all collections.
              Data will be collected from 12:05 AM of start date to 11:55 PM of end date.
            </p>
          </div>

          {/* Date Range Selection */}
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-6 mb-6">
            <h3 className="font-semibold text-blue-800 mb-4">Select Date Range</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Start Date
                </label>
                <input
                  type="date"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                  max={getTodayDate()}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                />
                <p className="text-xs text-gray-500 mt-1">Data from 12:05 AM onwards</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  End Date
                </label>
                <input
                  type="date"
                  value={endDate}
                  onChange={(e) => setEndDate(e.target.value)}
                  max={getTodayDate()}
                  min={startDate}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none"
                />
                <p className="text-xs text-gray-500 mt-1">Data until 11:55 PM</p>
              </div>
            </div>
          </div>

          <div className="flex justify-center mb-6">
            <button
              onClick={handleGenerateReports}
              disabled={loading || !startDate || !endDate}
              className="px-8 py-3 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition flex items-center space-x-2 disabled:opacity-50 disabled:cursor-not-allowed font-semibold"
            >
              {loading ? (
                <>
                  <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                  <span>Generating Reports...</span>
                </>
              ) : (
                <>
                  <FileText className="w-5 h-5" />
                  <span>Generate Reports</span>
                </>
              )}
            </button>
          </div>

          {/* Result Message */}
          {result && (
            <div
              className={`mb-6 p-4 rounded-lg flex items-start space-x-3 ${
                result.success
                  ? 'bg-green-50 border border-green-200'
                  : 'bg-red-50 border border-red-200'
              }`}
            >
              {result.success ? (
                <CheckCircle className="w-5 h-5 text-green-600 mt-0.5" />
              ) : (
                <XCircle className="w-5 h-5 text-red-600 mt-0.5" />
              )}
              <div>
                <p
                  className={`font-semibold ${
                    result.success ? 'text-green-800' : 'text-red-800'
                  }`}
                >
                  {result.success ? 'Success' : 'Error'}
                </p>
                <p
                  className={`text-sm ${
                    result.success ? 'text-green-700' : 'text-red-700'
                  }`}
                >
                  {result.message}
                </p>
              </div>
            </div>
          )}

          {/* Generated Files List */}
          {generatedFiles.length > 0 && (
            <div className="bg-gray-50 rounded-lg p-6">
              <h3 className="font-semibold text-gray-800 mb-4">
                Generated Reports ({generatedFiles.length})
              </h3>
              <div className="space-y-2">
                {generatedFiles.map((file) => (
                  <div
                    key={file.fileName}
                    className="flex items-center justify-between bg-white p-4 rounded-lg border border-gray-200 hover:border-primary-300 transition"
                  >
                    <div className="flex items-center space-x-3">
                      <FileText className="w-5 h-5 text-gray-500" />
                      <div>
                        <p className="font-medium text-gray-800">{file.collection}</p>
                        <p className="text-xs text-gray-500">{file.fileName}</p>
                      </div>
                    </div>
                    <button
                      onClick={() => handleDownload(file.fileName)}
                      className="flex items-center space-x-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition"
                    >
                      <Download className="w-4 h-4" />
                      <span>Download</span>
                    </button>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Collections Info */}
          <div className="mt-8 bg-gray-50 rounded-lg p-4">
            <h3 className="font-semibold text-gray-800 mb-2">Collections Included:</h3>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-2 text-sm text-gray-600">
              {[
                'appointments',
                'astrologer',
                'chat',
                'conversionHistory',
                'feedback',
                'ivrCall',
                'orderScore',
                'rewards',
                'serviceReports',
                'user',
                'userPayment',
                'userProblem',
                'videoCall',
              ].map((collection) => (
                <div key={collection} className="flex items-center space-x-2">
                  <div className="w-2 h-2 bg-primary-500 rounded-full"></div>
                  <span>{collection}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Scheduler;
