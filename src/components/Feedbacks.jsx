import { useState } from 'react';
import { feedbacksAPI } from '../services/api';
import { Upload, CheckCircle, XCircle } from 'lucide-react';

const Feedbacks = () => {
  const [feedbackData, setFeedbackData] = useState('');
  const [uploading, setUploading] = useState(false);
  const [uploadMessage, setUploadMessage] = useState({ type: '', text: '' });

  const handleFileUpload = (event) => {
    const file = event.target.files[0];
    if (file && file.type === 'application/json') {
      const reader = new FileReader();
      reader.onload = (e) => {
        try {
          const json = JSON.parse(e.target.result);
          setFeedbackData(JSON.stringify(json, null, 2));
          setUploadMessage({ type: '', text: '' });
        } catch (error) {
          setUploadMessage({ type: 'error', text: 'Invalid JSON file. Please check the format.' });
        }
      };
      reader.readAsText(file);
    } else {
      setUploadMessage({ type: 'error', text: 'Please upload a valid JSON file.' });
    }
  };

  const handleBulkUpload = async (e) => {
    e?.preventDefault?.();
    
    if (!feedbackData.trim()) {
      setUploadMessage({ type: 'error', text: 'Please provide feedback data.' });
      return;
    }

    try {
      setUploading(true);
      setUploadMessage({ type: '', text: '' });
      
      const parsedData = JSON.parse(feedbackData);
      
      // Ensure it's an array
      const feedbacksArray = Array.isArray(parsedData) ? parsedData : [parsedData];
      
      console.log('Uploading feedbacks:', feedbacksArray);
      const response = await feedbacksAPI.bulkCreate(feedbacksArray);
      console.log('Upload response:', response);
      
      setUploadMessage({ 
        type: 'success', 
        text: response?.message || 'Successfully uploaded feedback for astrologers.' 
      });
      setFeedbackData('');
      
      // Clear success message after 5 seconds
      setTimeout(() => {
        setUploadMessage({ type: '', text: '' });
      }, 5000);
    } catch (error) {
      console.error('Error uploading feedbacks:', error);
      console.error('Error details:', {
        message: error.message,
        response: error.response?.data,
        status: error.response?.status,
      });
      setUploadMessage({
        type: 'error',
        text: error.response?.data?.message || error.message || 'Failed to upload feedbacks. Please check the data format.',
      });
    } finally {
      setUploading(false);
    }
  };

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-800">Bulk Upload Feedbacks</h1>
        <p className="text-gray-600 mt-2">Upload feedback data for astrologers in bulk</p>
      </div>

      <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-200 max-w-4xl">
        <div className="flex items-center space-x-2 mb-6">
          <Upload className="w-5 h-5 text-primary-600" />
          <h2 className="text-xl font-semibold text-gray-800">Upload Feedback Data</h2>
        </div>
        
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Upload JSON File
            </label>
            <input
              type="file"
              accept=".json,application/json"
              onChange={handleFileUpload}
              className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-semibold file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100 cursor-pointer"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Or Paste JSON Data
            </label>
            <textarea
              value={feedbackData}
              onChange={(e) => {
                setFeedbackData(e.target.value);
                setUploadMessage({ type: '', text: '' });
              }}
              placeholder='[{"astroId": "...", "comment": "...", "name": "...", "rating": 4, "profilePic": "", "feedbackId": "..."}]'
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-sm font-mono"
              rows="12"
            />
          </div>

          {uploadMessage.text && (
            <div
              className={`flex items-center space-x-2 p-3 rounded-lg ${
                uploadMessage.type === 'success'
                  ? 'bg-green-50 text-green-800 border border-green-200'
                  : 'bg-red-50 text-red-800 border border-red-200'
              }`}
            >
              {uploadMessage.type === 'success' ? (
                <CheckCircle className="w-5 h-5" />
              ) : (
                <XCircle className="w-5 h-5" />
              )}
              <span className="text-sm">{uploadMessage.text}</span>
            </div>
          )}

          <button
            type="button"
            onClick={handleBulkUpload}
            disabled={uploading || !feedbackData.trim()}
            className="w-full bg-primary-600 text-white py-2 px-4 rounded-lg hover:bg-primary-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors flex items-center justify-center space-x-2"
          >
            {uploading ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                <span>Uploading...</span>
              </>
            ) : (
              <>
                <Upload className="w-4 h-4" />
                <span>Upload Feedbacks</span>
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

export default Feedbacks;

