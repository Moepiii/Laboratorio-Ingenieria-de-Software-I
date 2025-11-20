import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import './AdminDashboard.css';
import { useAuth } from './context/AuthContext';

const Sidebar = () => {
  const { logout, userRole } = useAuth();
  const location = useLocation();
  const currentPath = location.pathname;

  // Verifica si estamos en una sub-página de configuración de proyecto
  // Esto detecta URLs como: /admin/configuraciones/proyecto/1/labores
  const configMatch = currentPath.match(/\/admin\/configuraciones\/proyecto\/(\d+)/);
  const proyectoId = configMatch ? configMatch[1] : null;

  return (
    <nav className="sidebar">
      <h3>Menú</h3>

      {/* Botón Portafolio */}
      <Link
        to="/admin/proyectos"
        className={`nav-button ${currentPath.startsWith('/admin/proyectos') ? 'active' : ''}`}
      >
        Portafolio de Proyectos
      </Link>

      {/* Botón Usuarios */}
      <Link
        to="/admin/usuarios"
        className={`nav-button ${currentPath.startsWith('/admin/usuarios') ? 'active' : ''}`}
      >
        Perfiles de usuarios
      </Link>

      {/* Botón Configuraciones (Lista de Proyectos) */}
      <Link
        to="/admin/configuraciones"
        className={`nav-button ${currentPath.startsWith('/admin/configuraciones') ? 'active' : ''}`}
      >
        Configuraciones
      </Link>

      {/* ⭐️ SUB-MENÚ DEL PROYECTO SELECCIONADO ⭐️ */}
      {proyectoId && (
        <div style={{ marginLeft: '1rem', borderLeft: '2px solid #e5e7eb', padding: '0.5rem 0' }}>
          
          {/* 1. Labores */}
          <Link
            to={`/admin/configuraciones/proyecto/${proyectoId}/labores`}
            className={`nav-button ${currentPath.includes('/labores') ? 'active' : ''}`}
            style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block', marginBottom: '0.2rem' }}
          >
            ↳ Labores Agronómicas
          </Link>

          {/* 2. Equipos */}
          <Link
            to={`/admin/configuraciones/proyecto/${proyectoId}/equipos`}
            className={`nav-button ${currentPath.includes('/equipos') ? 'active' : ''}`}
            style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block', marginBottom: '0.2rem' }}
          >
            ↳ Equipos e Implementos
          </Link>

          {/* ⭐️ 3. UNIDADES DE MEDIDA (Ahora aquí dentro) ⭐️ */}
          <Link
            to={`/admin/configuraciones/proyecto/${proyectoId}/unidades`}
            className={`nav-button ${currentPath.includes('/unidades') ? 'active' : ''}`}
            style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block' }}
          >
            ↳ Unidades de Medida
          </Link>

        </div>
      )}

      {/* Botón Logs (Solo Admin) */}
      {userRole === 'admin' && (
        <Link
          to="/admin/logs"
          className={`nav-button ${currentPath.startsWith('/admin/logs') ? 'active' : ''}`}
        >
          Logger de eventos
        </Link>
      )}

      <div className="sidebar-spacer"></div>

      <button onClick={logout} className="logout-button">
        Cerrar Sesión
      </button>
    </nav>
  );
};

export default Sidebar;