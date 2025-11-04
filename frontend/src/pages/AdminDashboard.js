import React from 'react';
// 1. IMPORTA useParams PARA LOS COMPONENTES HIJOS
import { Routes, Route, Navigate, useParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

import Sidebar from '../Sidebar';
import Portafolio from './Portafolio';
import Usuarios from './Usuarios';
// Asumiendo que AdminDashboard.css está un nivel arriba de 'pages'
import '../AdminDashboard.css';

// 2. CREA LOS NUEVOS COMPONENTES (PLACEHOLDERS)
//    (Puedes moverlos a sus propios archivos .js si crecen mucho)

const LaboresAgronomicas = () => {
  const { id } = useParams(); // Obtiene el ID del proyecto de la URL
  return (
    <div style={{ padding: '2rem', color: '#333' }}>
      <h2>Labores Agronómicas</h2>
      <p>Mostrando labores para el Proyecto ID: <strong>{id}</strong></p>
      {/* Aquí iría tu futura tabla o lógica para las labores */}
    </div>
  );
};

const EquiposEImplementos = () => {
  const { id } = useParams(); // Obtiene el ID del proyecto de la URL
  return (
    <div style={{ padding: '2rem', color: '#333' }}>
      <h2>Equipos e Implementos</h2>
      <p>Mostrando equipos para el Proyecto ID: <strong>{id}</strong></p>
      {/* Aquí iría tu futura tabla o lógica para los equipos */}
    </div>
  );
};

// Componente principal del Dashboard
const AdminDashboard = () => {

  const { userRole } = useAuth();

  return (
    <div className="admin-container">
      <Sidebar />

      <main className="main-content">
        <Routes>
          {/* Rutas existentes */}
          <Route path="proyectos" element={<Portafolio />} />
          <Route path="usuarios" element={<Usuarios />} />

          {/* 3. RUTA DE CONFIGURACIONES (Muestra la lista de proyectos) */}
          <Route path="configuraciones" element={<Portafolio />} />

          {/* 4. AÑADE LAS NUEVAS RUTAS PARA EL SUB-MENÚ */}

          {/* Ruta para "Labores": /admin/configuraciones/proyecto/:id/labores */}
          <Route
            path="configuraciones/proyecto/:id/labores"
            element={<LaboresAgronomicas />}
          />

          {/* Ruta para "Equipos": /admin/configuraciones/proyecto/:id/equipos */}
          <Route
            path="configuraciones/proyecto/:id/equipos"
            element={<EquiposEImplementos />}
          />

          {/* Ruta de Logs (condicional) */}
          {userRole === 'admin' && (
            <Route path="logs" element={<div style={{ padding: '2rem', color: '#333' }}>Página de Logs (Aún no implementada)</div>} />
          )}

          {/* Ruta por defecto (redirige a /admin/proyectos) */}
          <Route
            index
            element={<Navigate to="proyectos" replace />}
          />
        </Routes>
      </main>
    </div>
  );
};

export default AdminDashboard;