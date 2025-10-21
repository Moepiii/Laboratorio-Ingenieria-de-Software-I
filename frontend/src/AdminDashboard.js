import React, { useState } from 'react';
import Sidebar from './Sidebar'; // Importamos el menú lateral
import PortafolioProyectos from './PortafolioProyectos'; // Importamos la vista de Proyectos
import PerfilesUsuarios from './PerfilesUsuarios'; // Importamos la vista de Usuarios
import LoggerEventos from './LoggerEventos'; // Importamos la vista de Logs
import './AdminDashboard.css'; // Importaremos los estilos para el layout

const AdminDashboard = ({ currentUser, handleLogout, apiCall }) => {
  // 'currentView' controla qué componente se muestra en el área principal
  const [currentView, setCurrentView] = useState('proyectos');

  // Función para renderizar el componente de la vista actual
  const renderView = () => {
    switch (currentView) {
      case 'proyectos':
        // ⭐️ AQUÍ ESTABA EL ERROR: Faltaba pasar currentUser
        return <PortafolioProyectos apiCall={apiCall} currentUser={currentUser} />;
      case 'usuarios':
        return <PerfilesUsuarios apiCall={apiCall} currentUser={currentUser} />;
      case 'logs':
        return <LoggerEventos apiCall={apiCall} />;
      default:
        // ⭐️ Y LO AÑADIMOS AQUÍ TAMBIÉN (en el 'default')
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
      />

      {/* El contenido principal que cambia dinámicamente */}
      <main className="main-content">
        {/* Aquí puedes añadir un header si quieres */}
        {/* <header className="main-header">...</header> */}

        <div className="view-wrapper">
          {renderView()}
        </div>
      </main>
    </div>
  );
};

export default AdminDashboard;