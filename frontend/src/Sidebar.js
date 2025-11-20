import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import './AdminDashboard.css';
import { useAuth } from './context/AuthContext';

const Sidebar = () => {
  const { logout, userRole } = useAuth();
  const location = useLocation();
  const currentPath = location.pathname;

  // Verifica si estamos en una sub-página de configuración de proyecto
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

      {/* Botón Configuraciones (Proyectos) */}
      <Link
        to="/admin/configuraciones"
        className={`nav-button ${currentPath.startsWith('/admin/configuraciones') ? 'active' : ''}`}
      >
        Configuraciones
      </Link>

      {/* Sub-menú condicional (solo si estamos editando un proyecto específico) */}
      {proyectoId && (
        <div style={{ marginLeft: '1rem', borderLeft: '2px solid #e5e7eb', padding: '0.5rem 0' }}>

          {/* Sub-Botón Labores */}
          <Link
            to={`/admin/configuraciones/proyecto/${proyectoId}/labores`}
            className={`nav-button ${currentPath.includes('/labores') ? 'active' : ''}`}
            style={{
              fontSize: '0.9rem',
              padding: '0.5rem 1rem',
              display: 'block',
              marginBottom: '0.5rem'
            }}
          >
            ↳ Labores Agronómicas
          </Link>

          {/* Sub-Botón Equipos */}
          <Link
            to={`/admin/configuraciones/proyecto/${proyectoId}/equipos`}
            className={`nav-button ${currentPath.includes('/equipos') ? 'active' : ''}`}
            style={{
              fontSize: '0.9rem',
              padding: '0.5rem 1rem',
              display: 'block'
            }}
          >
            ↳ Equipos e Implementos
          </Link>
        </div>
      )}

      {/* ⭐️ NUEVO BOTÓN: UNIDADES DE MEDIDA ⭐️ */}
      {/* Se muestra para admins y gerentes (igual que los handlers) */}
      <Link
        to="/admin/unidades"
        className={`nav-button ${currentPath.startsWith('/admin/unidades') ? 'active' : ''}`}
      >
        Unidades de Medida
      </Link>

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

      {/* Botón Salir */}
      <button onClick={logout} className="logout-button">
        Cerrar Sesión
      </button>
    </nav>
  );
};

export default Sidebar;