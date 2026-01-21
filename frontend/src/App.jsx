import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './components/Login';
import Dashboard, { DashboardHome } from './components/Dashboard';
import ArticleList from './components/ArticleList';
import ArticleForm from './components/ArticleForm';
import PageList from './components/PageList';
import PageForm from './components/PageForm';
import MediaManager from './components/MediaManager';

function PrivateRoute({ children }) {
  const token = localStorage.getItem('token');
  return token ? children : <Navigate to="/login" />;
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/"
          element={
            <PrivateRoute>
              <Dashboard />
            </PrivateRoute>
          }
        >
          <Route index element={<DashboardHome />} />
          <Route path="articles" element={<ArticleList />} />
          <Route path="articles/new" element={<ArticleForm />} />
          <Route path="articles/:id/edit" element={<ArticleForm />} />
          <Route path="pages" element={<PageList />} />
          <Route path="pages/new" element={<PageForm />} />
          <Route path="pages/:id/edit" element={<PageForm />} />
          <Route path="media" element={<MediaManager />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
