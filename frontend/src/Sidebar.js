import React from 'react';
// 1. Importa Link y useLocation
import { Link, useLocation } from 'react-router-dom';
import './AdminDashboard.css'; // Reutilizamos los mismos estilos
import { useAuth } from './context/AuthContext'; // Importa useAuth

const Sidebar = () => {

  // 2. Obtiene 'logout' y 'userRole' del contexto
  const { logout, userRole } = useAuth();

  // 3. Obtiene la ubicación actual (la URL)
  const location = useLocation();
  const currentPath = location.pathname; // Ej: "/admin/proyectos"

  return (
    <nav className="sidebar">
      <h3>Menú</h3>

      {/* 4. Cambiamos <button> por <Link> */}
      <Link
        to="/admin/proyectos" // La ruta a la que navega
        // 5. La clase 'active' se asigna si la URL actual incluye esta ruta
        className={`nav-button ${currentPath.includes('/admin/proyectos') ? 'active' : ''}`}
      >
        Portafolio de Proyectos
      </Link>

      <Link
        to="/admin/usuarios"
        className={`nav-button ${currentPath.includes('/admin/usuarios') ? 'active' : ''}`}
      >
        Perfiles de usuarios
      </Link>

      {/* 6. Renderizado condicional del link de Logs (basado en lógica original) */}
      {userRole === 'admin' && (
        <Link
          to="/admin/logs"
          className={`nav-button ${currentPath.includes('/admin/logs') ? 'active' : ''}`}
        >
          Logger de eventos
        </Link>
      )}

      {/* El botón de Logout usa 'logout' del contexto */}
      <button className="nav-button logout-button" onClick={logout}>
        Cerrar Sesión
      </button>
    </nav>
  );
};

export default Sidebar;