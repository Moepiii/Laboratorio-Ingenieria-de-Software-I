import React from 'react';
import './AdminDashboard.css'; // Reutilizamos los mismos estilos

// Recibe 'onNavigate' para cambiar la vista en el padre (AdminDashboard)
// Recibe 'currentView' para saber qué botón resaltar como "activo"
const Sidebar = ({ onNavigate, handleLogout, currentView }) => {
  return (
    <nav className="sidebar">
      <h3>Menú</h3>
      
      <button 
        className={`nav-button ${currentView === 'proyectos' ? 'active' : ''}`}
        onClick={() => onNavigate('proyectos')}
      >
        Portafolio de Proyectos
      </button>
      
      <button 
        className={`nav-button ${currentView === 'usuarios' ? 'active' : ''}`}
        onClick={() => onNavigate('usuarios')}
      >
        Perfiles de usuarios
      </button>
      
      <button 
        className={`nav-button ${currentView === 'logs' ? 'active' : ''}`}
        onClick={() => onNavigate('logs')}
      >
        Logger de eventos
      </button>

      {/* Botón de Logout al final */}
      <button className="nav-button logout-button" onClick={handleLogout}>
        Cerrar Sesión
      </button>
    </nav>
  );
};

export default Sidebar;