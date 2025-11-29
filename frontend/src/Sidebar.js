import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import './AdminDashboard.css';
import { useAuth } from './context/AuthContext';

const Sidebar = () => {
  const { logout, userRole } = useAuth();
  const location = useLocation();
  const currentPath = location.pathname;

  // 1. Detectar si estamos en CONFIGURACIONES y en qué proyecto
  const configMatch = currentPath.match(/\/admin\/configuraciones\/proyecto\/(\d+)/);
  const configProyectoId = configMatch ? configMatch[1] : null;

  // 2. Detectar si estamos en PLANES DE ACCIÓN y en qué proyecto
  const planesMatch = currentPath.match(/\/admin\/planes-accion\/proyecto\/(\d+)/);
  const planesProyectoId = planesMatch ? planesMatch[1] : null;

  return (
    <nav className="sidebar">
      <h3>Menú</h3>

      {/* --- Portafolio --- */}
      <Link
        to="/admin/proyectos"
        className={`nav-button ${currentPath === '/admin/proyectos' ? 'active' : ''}`}
      >
        Portafolio de Proyectos
      </Link>

      {/* --- Usuarios --- */}
      <Link
        to="/admin/usuarios"
        className={`nav-button ${currentPath.startsWith('/admin/usuarios') ? 'active' : ''}`}
      >
        Perfiles de usuarios
      </Link>

      {/* ---------------- SECCIÓN CONFIGURACIONES ---------------- */}
      <Link
        to="/admin/configuraciones"
        className={`nav-button ${currentPath.startsWith('/admin/configuraciones') ? 'active' : ''}`}
      >
        Configuraciones
      </Link>

      {/* Sub-menú Configuraciones (Solo si hay ID) */}
      {configProyectoId && (
        <div style={{ marginLeft: '1rem', borderLeft: '2px solid rgba(255,255,255,0.3)', padding: '0.5rem 0' }}>
          <Link
            to={`/admin/configuraciones/proyecto/${configProyectoId}/labores`}
            className={`nav-button ${currentPath.includes('/labores') ? 'active' : ''}`}
            style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block' }}
          >
            ↳ Labores Agronómicas
          </Link>
          <Link
            to={`/admin/configuraciones/proyecto/${configProyectoId}/equipos`}
            className={`nav-button ${currentPath.includes('/equipos') ? 'active' : ''}`}
            style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block' }}
          >
            ↳ Equipos e Implementos
          </Link>
          <Link
            to={`/admin/configuraciones/proyecto/${configProyectoId}/unidades`}
            className={`nav-button ${currentPath.includes('/unidades') ? 'active' : ''}`}
            style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block' }}
          >
            ↳ Unidades de Medida
          </Link>
        </div>
      )}

      {/* ---------------- SECCIÓN PLANES DE ACCIÓN ---------------- */}
      <Link
        to="/admin/planes-accion"
        className={`nav-button ${currentPath.startsWith('/admin/planes-accion') ? 'active' : ''}`}
      >
        Planes de Acción
      </Link>

      {/* Sub-menú Planes de Acción (Solo si hay ID detectado en la URL) */}
      {planesProyectoId && (
        <div style={{ marginLeft: '1rem', borderLeft: '2px solid rgba(255,255,255,0.3)', padding: '0.5rem 0' }}>
            <Link
                to={`/admin/planes-accion/proyecto/${planesProyectoId}/general`}
                className={`nav-button ${currentPath.includes('/general') ? 'active' : ''}`}
                style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block' }}
            >
                ↳ Plan de Acción
            </Link>
            <Link
                to={`/admin/planes-accion/proyecto/${planesProyectoId}/recursos`}
                className={`nav-button ${currentPath.includes('/recursos') ? 'active' : ''}`}
                style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block' }}
            >
                ↳ Recurso Humano
            </Link>
            <Link
                to={`/admin/planes-accion/proyecto/${planesProyectoId}/materiales`}
                className={`nav-button ${currentPath.includes('/materiales') ? 'active' : ''}`}
                style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block' }}
            >
                ↳ Materiales
            </Link>
        </div>
      )}

      {/* ⭐️ RESTAURADO: Logger de eventos (Solo Admin) ⭐️ */}
      {userRole === 'admin' && (
        <Link
          to="/admin/logs"
          className={`nav-button ${currentPath.startsWith('/admin/logs') ? 'active' : ''}`}
          style={{ marginTop: '1rem' }}
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