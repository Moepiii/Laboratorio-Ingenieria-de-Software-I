import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

import UnidadesMedida from './UnidadesMedida';

import Sidebar from '../Sidebar';
import Portafolio from './Portafolio';
import Usuarios from './Usuarios';
import '../AdminDashboard.css';

// Importamos los componentes de las sub-páginas
import LaboresAgronomicas from './LaboresAgronomicas';
import EquiposEImplementos from './EquiposEImplementos';
import DatosProyecto from './DatosProyecto';
import LoggerEventos from './LoggerEventos'; // ⭐️ 1. IMPORTAMOS EL NUEVO COMPONENTE

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

          {/* Ruta para "Labores": /admin/configuraciones/proyecto/:id/labores */}
          <Route
            path="configuraciones/proyecto/:id/labores"
            element={<LaboresAgronomicas />}
          />

          {/* Ruta para "Equipos\": /admin/configuraciones/proyecto/:id/equipos */}
          <Route
            path="configuraciones/proyecto/:id/equipos"
            element={<EquiposEImplementos />}
          />

          <Route path="unidades" element={<UnidadesMedida />} />

          {/* Ruta para la página de Datos del Proyecto */}
          <Route
            path="proyectos/datos/:id"
            element={<DatosProyecto />}
          />

          {/* Ruta de Logs (condicional) */}
          {userRole === 'admin' && (
            // ⭐️ 2. REEMPLAZAMOS EL DIV DE RELLENO
            <Route path="logs" element={<LoggerEventos />} />
          )}

          {/* Ruta por defecto (redirige a /admin/proyectos) */}
          <Route
            path="/"
            element={<Navigate to="/admin/proyectos" replace />}
          />

          {/* Ruta "catch-all" (opcional) */}
          <Route
            path="*"
            element={<Navigate to="/admin/proyectos" replace />}
          />
        </Routes>
      </main>
    </div>
  );
};

export default AdminDashboard;