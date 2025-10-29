import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
// ⭐️ 'useAuth' ya no es necesario aquí, porque el Sidebar lo llama por su cuenta
// import { useAuth } from '../context/AuthContext'; 

// Importa los componentes
import Sidebar from '../Sidebar';
import Portafolio from './Portafolio';
import Usuarios from './Usuarios';
// import Logs from './Logs'; // No tenemos este archivo, lo dejamos comentado
import '../AdminDashboard.css'; // Importa el CSS

const AdminDashboard = () => {

  // ⭐️ LÍNEA CORREGIDA:
  // Ya no extraemos 'userRole' porque no se usa en ESTE componente.
  // El Sidebar ahora obtiene 'userRole' él mismo desde el contexto.

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