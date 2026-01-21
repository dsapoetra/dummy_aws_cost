import { useState, useEffect } from 'react';
import { Link, Outlet, useNavigate, useLocation } from 'react-router-dom';
import { auth } from '../api/client';

export default function Dashboard() {
  const [user, setUser] = useState(null);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('/login');
      return;
    }

    auth.me()
      .then((res) => setUser(res.data))
      .catch(() => {
        localStorage.removeItem('token');
        navigate('/login');
      });
  }, [navigate]);

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  const isActive = (path) => location.pathname === path || location.pathname.startsWith(path + '/');

  if (!user) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-gray-600">Loading...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4">
          <div className="flex justify-between h-16">
            <div className="flex items-center space-x-8">
              <span className="text-xl font-bold text-gray-800">CMS</span>
              <Link
                to="/"
                className={`px-3 py-2 rounded ${isActive('/') && location.pathname === '/' ? 'bg-blue-100 text-blue-700' : 'text-gray-600 hover:text-gray-800'}`}
              >
                Dashboard
              </Link>
              <Link
                to="/articles"
                className={`px-3 py-2 rounded ${isActive('/articles') ? 'bg-blue-100 text-blue-700' : 'text-gray-600 hover:text-gray-800'}`}
              >
                Articles
              </Link>
              <Link
                to="/pages"
                className={`px-3 py-2 rounded ${isActive('/pages') ? 'bg-blue-100 text-blue-700' : 'text-gray-600 hover:text-gray-800'}`}
              >
                Pages
              </Link>
              <Link
                to="/media"
                className={`px-3 py-2 rounded ${isActive('/media') ? 'bg-blue-100 text-blue-700' : 'text-gray-600 hover:text-gray-800'}`}
              >
                Media
              </Link>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-gray-600">Welcome, {user.username}</span>
              <button
                onClick={handleLogout}
                className="text-red-600 hover:text-red-800"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 px-4">
        <Outlet />
      </main>
    </div>
  );
}

export function DashboardHome() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-800 mb-6">Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Link to="/articles" className="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow">
          <h2 className="text-lg font-semibold text-gray-800">Articles</h2>
          <p className="text-gray-600 mt-2">Manage blog posts and articles</p>
        </Link>
        <Link to="/pages" className="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow">
          <h2 className="text-lg font-semibold text-gray-800">Pages</h2>
          <p className="text-gray-600 mt-2">Manage static pages</p>
        </Link>
        <Link to="/media" className="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow">
          <h2 className="text-lg font-semibold text-gray-800">Media</h2>
          <p className="text-gray-600 mt-2">Upload and manage media files</p>
        </Link>
      </div>
    </div>
  );
}
