import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

import Sidebar from '../Sidebar';
import Portafolio from './Portafolio';
import Usuarios from './Usuarios';
import '../AdminDashboard.css';

import LaboresAgronomicas from './LaboresAgronomicas';
import EquiposEImplementos from './EquiposEImplementos';
import DatosProyecto from './DatosProyecto';
import LoggerEventos from './LoggerEventos';
import UnidadesMedida from './UnidadesMedida';

// Importamos los componentes de Planes de Acción
import PlanDeAccion from './PlanDeAccion';
import RecursoHumano from './RecursoHumano';
// ⭐️ 1. IMPORTAR EL NUEVO COMPONENTE
import MaterialesInsumos from './MaterialesInsumos';

const AdminDashboard = () => {

  const { userRole } = useAuth();

  return (
    <div className="admin-container">
      <Sidebar />

      <main className="main-content">
        <Routes>
          {/* ... Rutas anteriores ... */}
          <Route path="proyectos" element={<Portafolio />} />
          <Route path="usuarios" element={<Usuarios />} />
          <Route path="configuraciones" element={<Portafolio />} />

          <Route path="configuraciones/proyecto/:id/labores" element={<LaboresAgronomicas />} />
          <Route path="configuraciones/proyecto/:id/equipos" element={<EquiposEImplementos />} />
          <Route path="configuraciones/proyecto/:id/unidades" element={<UnidadesMedida />} />
          <Route path="proyectos/datos/:id" element={<DatosProyecto />} />

          {/* Rutas Planes de Acción */}
          <Route path="planes-accion/general" element={<PlanDeAccion />} />
          <Route path="planes-accion/recursos" element={<RecursoHumano />} />
          
          {/* ⭐️ 2. RUTA ACTUALIZADA: MATERIALES E INSUMOS ⭐️ */}
          <Route path="planes-accion/materiales" element={<MaterialesInsumos />} />

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