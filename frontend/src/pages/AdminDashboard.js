import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';


// Importa los componentes
import Sidebar from '../Sidebar';
import Portafolio from './Portafolio';
import Usuarios from './Usuarios';
import '../AdminDashboard.css'; // Importa el CSS

const AdminDashboard = () => {


  return (
    <div className="admin-container">
      {/* El Sidebar ahora es totalmente independiente.
        Obtiene 'logout' y 'userRole' desde useAuth() por sí mismo.
      */}
      <Sidebar />

      <main className="main-content">
        <Routes>
          {/* Rutas anidadas (ej. /admin/proyectos) */}
          <Route path="proyectos" element={<Portafolio />} />
          <Route path="usuarios" element={<Usuarios />} />

          {/* Ruta de Logs (comentada).
            Si la descomentamos, tendríamos que volver a importar useAuth
            y obtener userRole aquí.
          */}
          {/*
          {userRole === 'admin' && (
            <Route path="logs" element={<Logs />} />
          )}
          */}

          {/* Ruta por defecto (ej. /admin/) */}
          <Route
            index
            element={<Navigate to="proyectos" replace />} // Redirige a proyectos
          />
        </Routes>
      </main>
    </div>
  );
};

export default AdminDashboard;