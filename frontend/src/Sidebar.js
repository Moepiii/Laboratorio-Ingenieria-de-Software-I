import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import './AdminDashboard.css';
import { useAuth } from './context/AuthContext';

const Sidebar = () => {
  const { logout, userRole } = useAuth();
  const location = useLocation();
  const currentPath = location.pathname;

  // Estado para el menú desplegable
  const [showPlanesAccion, setShowPlanesAccion] = useState(false);

  const configMatch = currentPath.match(/\/admin\/configuraciones\/proyecto\/(\d+)/);
  const proyectoId = configMatch ? configMatch[1] : null;

  const isActivePlanes = currentPath.includes('/admin/planes-accion');

  return (
    <nav className="sidebar">
      <h3>Menú</h3>

      {/* --- Enlaces Principales --- */}
      <Link
        to="/admin/proyectos"
        className={`nav-button ${currentPath.startsWith('/admin/proyectos') ? 'active' : ''}`}
      >
        Portafolio de Proyectos
      </Link>

      <Link
        to="/admin/usuarios"
        className={`nav-button ${currentPath.startsWith('/admin/usuarios') ? 'active' : ''}`}
      >
        Perfiles de usuarios
      </Link>

      <Link
        to="/admin/configuraciones"
        className={`nav-button ${currentPath.startsWith('/admin/configuraciones') ? 'active' : ''}`}
      >
        Configuraciones
      </Link>

      {/* Sub-menú de Configuración de Proyecto */}
      {proyectoId && (
        <div style={{ marginLeft: '1rem', borderLeft: '2px solid #e5e7eb', padding: '0.5rem 0' }}>
          <Link
            to={`/admin/configuraciones/proyecto/${proyectoId}/labores`}
            className={`nav-button ${currentPath.includes('/labores') ? 'active' : ''}`}
            style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block', marginBottom: '0.2rem' }}
          >
            ↳ Labores Agronómicas
          </Link>
          <Link
            to={`/admin/configuraciones/proyecto/${proyectoId}/equipos`}
            className={`nav-button ${currentPath.includes('/equipos') ? 'active' : ''}`}
            style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block', marginBottom: '0.2rem' }}
          >
            ↳ Equipos e Implementos
          </Link>
          <Link
            to={`/admin/configuraciones/proyecto/${proyectoId}/unidades`}
            className={`nav-button ${currentPath.includes('/unidades') ? 'active' : ''}`}
            style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block' }}
          >
            ↳ Unidades de Medida
          </Link>
        </div>
      )}

      {/* ⭐️⭐️⭐️ SECCIÓN DESPLEGABLE: PLANES DE ACCIÓN ⭐️⭐️⭐️ */}
      
      {/* Botón Principal */}
      <button
        className={`nav-button ${isActivePlanes ? 'active' : ''}`}
        onClick={() => setShowPlanesAccion(!showPlanesAccion)}
        style={{ 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center',
            cursor: 'pointer', 
            width: '100%', 
            textAlign: 'left',
            border: 'none', 
            background: isActivePlanes ? '#2563eb' : 'transparent',
            color: 'inherit', // Hereda el color del texto (blanco)
            padding: '0.75rem 1rem', // Padding estándar de tus botones
            fontSize: '1rem',
            fontWeight: '500'
        }}
      >
        Planes de Acción
        <span style={{ fontSize: '0.8rem' }}>{showPlanesAccion ? '▲' : '▼'}</span>
      </button>

      {/* Contenedor Desplegable (CORREGIDO) */}
      {showPlanesAccion && (
        // ⭐️ CAMBIO: Eliminado backgroundColor: '#f9fafb' para que sea transparente
        <div style={{ marginLeft: '1rem', borderLeft: '2px solid rgba(255,255,255,0.3)', padding: '0.5rem 0' }}>
            <Link
                to="/admin/planes-accion/general" 
                className={`nav-button ${currentPath === '/admin/planes-accion/general' ? 'active' : ''}`}
                style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block' }}
            >
                ↳ Plan de Acción
            </Link>
            <Link
                to="/admin/planes-accion/recursos"
                className={`nav-button ${currentPath === '/admin/planes-accion/recursos' ? 'active' : ''}`}
                style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block' }}
            >
                ↳ Recurso Humano
            </Link>
            <Link
                to="/admin/planes-accion/materiales"
                className={`nav-button ${currentPath === '/admin/planes-accion/materiales' ? 'active' : ''}`}
                style={{ fontSize: '0.9rem', padding: '0.5rem 1rem', display: 'block' }}
            >
                ↳ Materiales
            </Link>
        </div>
      )}
      {/* ⭐️⭐️⭐️ FIN SECCIÓN ⭐️⭐️⭐️ */}


      {userRole === 'admin' && (
        <Link
          to="/admin/logs"
          className={`nav-button ${currentPath.startsWith('/admin/logs') ? 'active' : ''}`}
          style={{ marginTop: '0.5rem' }}
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