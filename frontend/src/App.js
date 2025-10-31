import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext';


import AdminDashboard from './pages/AdminDashboard';
import UserDashboard from './pages/UserDashboard';
import LoginPage from './pages/LoginPage';

/**
 * Componente de Ruta Protegida
 * Si el usuario no está logueado, lo redirige a /login
 */
const ProtectedRoute = ({ children }) => {
  const { isLoggedIn } = useAuth();
  return isLoggedIn ? children : <Navigate to="/login" replace />;
};

/**
 * Componente de Ruta de Administrador
 * Verifica si está logueado Y si tiene el rol correcto.
 */
const AdminRoute = ({ children }) => {
  const { isLoggedIn, userRole } = useAuth();

  if (!isLoggedIn) {
    return <Navigate to="/login" replace />;
  }
  if (userRole === 'admin' || userRole === 'gerente') {
    return children;
  }
  // Si es un usuario 'user' intentando entrar a /admin
  return <Navigate to="/dashboard" replace />;
};

/**
 * Componente de Ruta Pública (para el Login)
 * Si el usuario YA está logueado, lo redirige a su dashboard
 */
const PublicRoute = ({ children }) => {
  const { isLoggedIn, userRole } = useAuth();

  if (isLoggedIn) {
    const redirectTo = (userRole === 'admin' || userRole === 'gerente') ? '/admin' : '/dashboard';
    return <Navigate to={redirectTo} replace />;
  }
  return children;
};

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Ruta de Login (Pública) */}
        <Route
          path="/login"
          element={
            <PublicRoute>
              <LoginPage />
            </PublicRoute>
          }
        />

        {/* Rutas de Administrador (Protegidas por Rol) */}
        <Route
          path="/admin/*" // El '/*' permite rutas anidadas (ej. /admin/proyectos)
          element={
            <AdminRoute>
              <AdminDashboard />
            </AdminRoute>
          }
        />

        {/* Ruta de Usuario (Protegida) */}
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <UserDashboard />
            </ProtectedRoute>
          }
        />

        {/* Ruta Raíz (/) */}
        <Route
          path="/"
          element={<Navigate to="/login" replace />} // Redirige a /login por defecto
        />

        {/* Ruta "catch-all" para 404 */}
        <Route
          path="*"
          element={<Navigate to="/" replace />} // Redirige a la raíz
        />

      </Routes>
    </BrowserRouter>
  );
}

export default App;