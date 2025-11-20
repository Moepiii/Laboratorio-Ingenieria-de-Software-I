import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

import Sidebar from '../Sidebar';
import Portafolio from './Portafolio';
import Usuarios from './Usuarios';
import '../AdminDashboard.css';

// Importamos los componentes de las sub-páginas
import LaboresAgronomicas from './LaboresAgronomicas';
import EquiposEImplementos from './EquiposEImplementos';
import DatosProyecto from './DatosProyecto';
import LoggerEventos from './LoggerEventos';
import UnidadesMedida from './UnidadesMedida'; // Importamos

const AdminDashboard = () => {

  const { userRole } = useAuth();

  return (
    <div className="admin-container">
      <Sidebar />

      <main className="main-content">
        <Routes>
          <Route path="proyectos" element={<Portafolio />} />
          <Route path="usuarios" element={<Usuarios />} />
          <Route path="configuraciones" element={<Portafolio />} />

          {/* Configuración: Labores */}
          <Route
            path="configuraciones/proyecto/:id/labores"
            element={<LaboresAgronomicas />}
          />

          {/* Configuración: Equipos */}
          <Route
            path="configuraciones/proyecto/:id/equipos"
            element={<EquiposEImplementos />}
          />

          {/* ⭐️ NUEVA RUTA ANIDADA: Unidades de Medida ⭐️ */}
          {/* NOTA: Ya no es /admin/unidades, ahora es específica del proyecto */}
          <Route
            path="configuraciones/proyecto/:id/unidades"
            element={<UnidadesMedida />}
          />

          <Route path="proyectos/datos/:id" element={<DatosProyecto />} />

          {userRole === 'admin' && (
            <Route path="logs" element={<LoggerEventos />} />
          )}

          <Route path="/" element={<Navigate to="/admin/proyectos" replace />} />
          <Route path="*" element={<Navigate to="/admin/proyectos" replace />} />
        </Routes>
      </main>
    </div>
  );
};

export default AdminDashboard;