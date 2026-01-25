import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import Layout from './components/Layout';
import Login from './components/Login';
import Dashboard from './components/Dashboard';
import Users from './components/Users';
import Astros from './components/Astros';
import UserServiceComplaints from './components/UserServiceComplaints';
import UserGeneralComplaints from './components/UserGeneralComplaints';
import AstroGeneralComplaints from './components/AstroGeneralComplaints';
import Horoscopes from './components/Horoscopes';
import Feedbacks from './components/Feedbacks';
import Scheduler from './components/Scheduler';

const PrivateRoute = ({ children }) => {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? children : <Navigate to="/login" />;
};

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route
            path="/dashboard"
            element={
              <PrivateRoute>
                <Layout>
                  <Dashboard />
                </Layout>
              </PrivateRoute>
            }
          />
          <Route
            path="/users"
            element={
              <PrivateRoute>
                <Layout>
                  <Users />
                </Layout>
              </PrivateRoute>
            }
          />
          <Route
            path="/astros"
            element={
              <PrivateRoute>
                <Layout>
                  <Astros />
                </Layout>
              </PrivateRoute>
            }
          />
          <Route
            path="/user-service-complaints"
            element={
              <PrivateRoute>
                <Layout>
                  <UserServiceComplaints />
                </Layout>
              </PrivateRoute>
            }
          />
          <Route
            path="/user-general-complaints"
            element={
              <PrivateRoute>
                <Layout>
                  <UserGeneralComplaints />
                </Layout>
              </PrivateRoute>
            }
          />
          <Route
            path="/astro-general-complaints"
            element={
              <PrivateRoute>
                <Layout>
                  <AstroGeneralComplaints />
                </Layout>
              </PrivateRoute>
            }
          />
          <Route
            path="/horoscopes"
            element={
              <PrivateRoute>
                <Layout>
                  <Horoscopes />
                </Layout>
              </PrivateRoute>
            }
          />
          <Route
            path="/feedbacks"
            element={
              <PrivateRoute>
                <Layout>
                  <Feedbacks />
                </Layout>
              </PrivateRoute>
            }
          />
          <Route
            path="/scheduler"
            element={
              <PrivateRoute>
                <Layout>
                  <Scheduler />
                </Layout>
              </PrivateRoute>
            }
          />
          <Route path="/" element={<Navigate to="/dashboard" />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;

