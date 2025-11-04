import React from 'react';
// ⭐️ YA NO SE NECESITA 'useParams' aquí
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

import Sidebar from '../Sidebar';
import Portafolio from './Portafolio';
import Usuarios from './Usuarios';
import '../AdminDashboard.css';

// ⭐️ 1. IMPORTA LOS NUEVOS COMPONENTES DESDE SUS ARCHIVOS
import LaboresAgronomicas from './LaboresAgronomicas';
import EquiposEImplementos from './EquiposEImplementos';

// ⭐️ 2. LOS COMPONENTES PLACEHOLDER FUERON ELIMINADOS DE AQUÍ

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

          {/* Ruta de Configuraciones (Muestra la lista de proyectos) */}
          <Route path="configuraciones" element={<Portafolio />} />

          {/* ⭐️ 3. LAS RUTAS AHORA USAN LOS COMPONENTES IMPORTADOS */}

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