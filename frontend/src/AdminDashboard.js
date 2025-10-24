import React, { useState } from 'react';
import Sidebar from './Sidebar';
import PortafolioProyectos from './PortafolioProyectos';
import PerfilesUsuarios from './PerfilesUsuarios';
import LoggerEventos from './LoggerEventos';
import './AdminDashboard.css';

// ⭐️ Recibe 'userRole' de App.js
const AdminDashboard = ({ currentUser, userRole, handleLogout, apiCall }) => {

  const [currentView, setCurrentView] = useState('proyectos');

  // Función para renderizar el componente de la vista actual
  const renderView = () => {
    switch (currentView) {
      case 'proyectos':
        return <PortafolioProyectos apiCall={apiCall} currentUser={currentUser} />;
      case 'usuarios':
        // ⭐️ Pasa 'userRole' a PerfilesUsuarios
        return <PerfilesUsuarios apiCall={apiCall} currentUser={currentUser} userRole={userRole} />;
      case 'logs':
        return userRole === 'admin' ? <LoggerEventos apiCall={apiCall} /> : <PortafolioProyectos apiCall={apiCall} currentUser={currentUser} />;
      default:
        return <PortafolioProyectos apiCall={apiCall} currentUser={currentUser} />;
    }
  };

  return (
    <div className="admin-container">
      {/* El panel lateral (Sidebar) */}
      <Sidebar
        onNavigate={setCurrentView}
        handleLogout={handleLogout}
        currentView={currentView}
        userRole={userRole} // ⭐️ Pasa 'userRole' al Sidebar
      />

      {/* El contenido principal que cambia dinámicamente */}
      <main className="main-content">
        <div className="view-wrapper">
          {renderView()}
        </div>
      </main>
    </div>
  );
};

export default AdminDashboard;