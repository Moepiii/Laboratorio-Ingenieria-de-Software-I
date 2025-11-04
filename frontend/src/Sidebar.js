import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import './AdminDashboard.css'; // Asumiendo que AdminDashboard.css está un nivel arriba
import { useAuth } from './context/AuthContext';

const Sidebar = () => {
  const { logout, userRole } = useAuth();
  const location = useLocation();
  const currentPath = location.pathname;

  // Verifica si estamos en una sub-página de configuración
  const configMatch = currentPath.match(/\/admin\/configuraciones\/proyecto\/(\d+)/);
  const isConfigProyectoPage = !!configMatch;
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

      {/* Botón Configuraciones */}
      <Link
        to="/admin/configuraciones"
        className={`nav-button ${currentPath.startsWith('/admin/configuraciones') ? 'active' : ''}`}
      >
        Configuraciones
      </Link>

      {/* ⭐️ SUB-MENÚ CONDICIONAL (CON ESPACIO) ⭐️ */}
      {isConfigProyectoPage && (
        <div style={{ paddingLeft: '1.5rem', borderLeft: '3px solid #4f46e5', margin: '0.5rem 0' }}>

          {/* Sub-Botón Labores */}
          <Link
            to={`/admin/configuraciones/proyecto/${proyectoId}/labores`}
            className={`nav-button ${currentPath.includes('/labores') ? 'active' : ''}`}
            style={{
              fontSize: '0.9rem',
              padding: '0.5rem 1rem',
              display: 'block', // Asegura que sea un bloque
              marginBottom: '0.5rem' // ⭐️ ESTE ES EL ESPACIO QUE PEDISTE ⭐️
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
              display: 'block' // Asegura que sea un bloque
            }}
          >
            ↳ Equipos e Implementos
          </Link>
        </div>
      )}

      {/* Botón Logs */}
      {userRole === 'admin' && (
        <Link
          to="/admin/logs"
          className={`nav-button ${currentPath.startsWith('/admin/logs') ? 'active' : ''}`}
        >
          Logger de eventos
        </Link>
      )}

      {/* Botón de Logout */}
      <button className="nav-button logout-button" onClick={logout}>
        Cerrar Sesión
      </button>
    </nav>
  );
};

export default Sidebar;