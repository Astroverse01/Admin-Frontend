import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { dashboardAPI, usersAPI, astrosAPI, complaintsAPI, userProblemsAPI, astroProblemsAPI } from '../services/api';
import { Users, UserCircle, MessageSquare, AlertCircle, TrendingUp, MessageCircle, Phone, Video } from 'lucide-react';

const Dashboard = () => {
  const [stats, setStats] = useState({
    users: { total: 0, active: 0 },
    astros: { total: 0, active: 0 },
    userServiceComplaints: 0,
    userGeneralComplaints: 0,
    astroGeneralComplaints: 0,
  });
  const [serviceMetrics, setServiceMetrics] = useState({
    date: '',
    chat: { failed: 0, request: 0, complete: 0, issue: 0, reject: 0, total: 0 },
    ivrCall: { failed: 0, request: 0, complete: 0, issue: 0, reject: 0, total: 0 },
    videoCall: { failed: 0, request: 0, complete: 0, issue: 0, reject: 0, total: 0 },
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        // Fetch service metrics (daily) and global overview stats in parallel
        const [metricsResponse, usersRes, astrosRes, complaintsRes, userProblemsRes, astroProblemsRes] = await Promise.all([
          dashboardAPI.getMetrics(),
          usersAPI.listUsers({ limit: 1 }),
          astrosAPI.listAstros({ limit: 1 }),
          complaintsAPI.listUserServiceComplaints({ limit: 1 }),
          userProblemsAPI.listUserGeneralComplaints({ limit: 1 }),
          astroProblemsAPI.listAstroGeneralComplaints({ limit: 1 }),
        ]);

        // Set service metrics (daily)
        if (metricsResponse.success && metricsResponse.data) {
          setServiceMetrics({
            date: metricsResponse.data.date || '',
            chat: metricsResponse.data.chat || { failed: 0, request: 0, complete: 0, issue: 0, reject: 0, total: 0 },
            ivrCall: metricsResponse.data.ivrCall || { failed: 0, request: 0, complete: 0, issue: 0, reject: 0, total: 0 },
            videoCall: metricsResponse.data.videoCall || { failed: 0, request: 0, complete: 0, issue: 0, reject: 0, total: 0 },
          });
        }

        // Set global overview stats
        setStats({
          users: {
            total: usersRes.pagination?.total || 0,
            active: usersRes.data?.filter((u) => u.status === 'active').length || 0,
          },
          astros: {
            total: astrosRes.pagination?.total || 0,
            active: astrosRes.data?.filter((a) => a.status === 'active').length || 0,
          },
          userServiceComplaints: complaintsRes.pagination?.total || 0,
          userGeneralComplaints: userProblemsRes.pagination?.total || 0,
          astroGeneralComplaints: astroProblemsRes.pagination?.total || 0,
        });
      } catch (error) {
        console.error('Error fetching dashboard data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  const statCards = [
    {
      title: 'Total Users',
      value: stats.users.total,
      subtitle: `${stats.users.active} active`,
      icon: Users,
      color: 'bg-blue-500',
      link: '/users',
    },
    {
      title: 'Total Astrologers',
      value: stats.astros.total,
      subtitle: `${stats.astros.active} active`,
      icon: UserCircle,
      color: 'bg-purple-500',
      link: '/astros',
    },
    {
      title: 'User Service Complaints',
      value: stats.userServiceComplaints,
      subtitle: 'Pending review',
      icon: MessageSquare,
      color: 'bg-orange-500',
      link: '/user-service-complaints',
    },
    {
      title: 'User General Complaints',
      value: stats.userGeneralComplaints,
      subtitle: 'Open issues',
      icon: AlertCircle,
      color: 'bg-red-500',
      link: '/user-general-complaints',
    },
    {
      title: 'Astro General Complaints',
      value: stats.astroGeneralComplaints,
      subtitle: 'Open issues',
      icon: AlertCircle,
      color: 'bg-yellow-500',
      link: '/astro-general-complaints',
    },
  ];

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  const serviceCards = [
    {
      title: 'Chat Services',
      data: serviceMetrics.chat,
      icon: MessageCircle,
      color: 'bg-blue-500',
    },
    {
      title: 'IVR Call Services',
      data: serviceMetrics.ivrCall,
      icon: Phone,
      color: 'bg-green-500',
    },
    {
      title: 'Video Call Services',
      data: serviceMetrics.videoCall,
      icon: Video,
      color: 'bg-purple-500',
    },
  ];

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-800">Dashboard</h1>
        <p className="text-gray-600 mt-2">
          Welcome to Astro Admin Dashboard
          {serviceMetrics.date && (
            <span className="ml-2 text-sm text-gray-500">({serviceMetrics.date})</span>
          )}
        </p>
      </div>

      {/* Service Metrics Section */}
      <div className="mb-8">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">Service Metrics</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {serviceCards.map((card, index) => {
            const Icon = card.icon;
            const { total, complete, request, failed, issue, reject } = card.data;
            return (
              <div
                key={index}
                className="bg-white rounded-xl shadow-sm hover:shadow-md transition-shadow p-6 border border-gray-200"
              >
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-semibold text-gray-800">{card.title}</h3>
                  <div className={`${card.color} p-3 rounded-lg`}>
                    <Icon className="w-6 h-6 text-white" />
                  </div>
                </div>
                <div className="space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-2xl font-bold text-gray-800">{total}</span>
                    <span className="text-sm text-gray-600">Total</span>
                  </div>
                  <div className="grid grid-cols-2 gap-2 pt-2 border-t border-gray-200">
                    <div>
                      <p className="text-xs text-gray-500">Complete</p>
                      <p className="text-lg font-semibold text-green-600">{complete}</p>
                    </div>
                    <div>
                      <p className="text-xs text-gray-500">Request</p>
                      <p className="text-lg font-semibold text-blue-600">{request}</p>
                    </div>
                    <div>
                      <p className="text-xs text-gray-500">Failed</p>
                      <p className="text-lg font-semibold text-red-600">{failed}</p>
                    </div>
                    <div>
                      <p className="text-xs text-gray-500">Issue</p>
                      <p className="text-lg font-semibold text-orange-600">{issue}</p>
                    </div>
                    <div className="col-span-2">
                      <p className="text-xs text-gray-500">Reject</p>
                      <p className="text-lg font-semibold text-gray-600">{reject}</p>
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Existing Stats Cards */}
      <div className="mb-8">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {statCards.map((card, index) => {
            const Icon = card.icon;
            return (
              <Link
                key={index}
                to={card.link}
                className="bg-white rounded-xl shadow-sm hover:shadow-md transition-shadow p-6 border border-gray-200"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600">{card.title}</p>
                    <p className="text-3xl font-bold text-gray-800 mt-2">{card.value}</p>
                    <p className="text-sm text-gray-500 mt-1">{card.subtitle}</p>
                  </div>
                  <div className={`${card.color} p-4 rounded-lg`}>
                    <Icon className="w-8 h-8 text-white" />
                  </div>
                </div>
              </Link>
            );
          })}
        </div>
      </div>

      <div className="mt-8 bg-white rounded-xl shadow-sm p-6 border border-gray-200">
        <div className="flex items-center space-x-2 mb-4">
          <TrendingUp className="w-5 h-5 text-primary-600" />
          <h2 className="text-xl font-semibold text-gray-800">Quick Actions</h2>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Link
            to="/users"
            className="p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <p className="font-medium text-gray-800">Manage Users</p>
            <p className="text-sm text-gray-600 mt-1">View and manage user accounts</p>
          </Link>
          <Link
            to="/astros"
            className="p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <p className="font-medium text-gray-800">Manage Astrologers</p>
            <p className="text-sm text-gray-600 mt-1">View and manage astrologer accounts</p>
          </Link>
          <Link
            to="/horoscopes"
            className="p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <p className="font-medium text-gray-800">Manage Horoscopes</p>
            <p className="text-sm text-gray-600 mt-1">Update daily horoscopes</p>
          </Link>
          <Link
            to="/scheduler"
            className="p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <p className="font-medium text-gray-800">Scheduler</p>
            <p className="text-sm text-gray-600 mt-1">Trigger daily reports</p>
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;

