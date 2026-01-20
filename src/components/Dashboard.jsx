import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { usersAPI, astrosAPI, complaintsAPI, userProblemsAPI, astroProblemsAPI } from '../services/api';
import { Users, UserCircle, MessageSquare, AlertCircle, TrendingUp } from 'lucide-react';

const Dashboard = () => {
  const [stats, setStats] = useState({
    users: { total: 0, active: 0 },
    astros: { total: 0, active: 0 },
    userServiceComplaints: 0,
    userGeneralComplaints: 0,
    astroGeneralComplaints: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [usersRes, astrosRes, complaintsRes, userProblemsRes, astroProblemsRes] = await Promise.all([
          usersAPI.listUsers({ limit: 1 }),
          astrosAPI.listAstros({ limit: 1 }),
          complaintsAPI.listUserServiceComplaints({ limit: 1 }),
          userProblemsAPI.listUserGeneralComplaints({ limit: 1 }),
          astroProblemsAPI.listAstroGeneralComplaints({ limit: 1 }),
        ]);

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
        console.error('Error fetching stats:', error);
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

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-800">Dashboard</h1>
        <p className="text-gray-600 mt-2">Welcome to Astro Admin Dashboard</p>
      </div>

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

