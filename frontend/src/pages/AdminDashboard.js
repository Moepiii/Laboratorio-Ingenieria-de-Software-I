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

import PlanDeAccion from './PlanDeAccion';
import RecursoHumano from './RecursoHumano';
import MaterialesInsumos from './MaterialesInsumos';

const AdminDashboard = () => {
  const { userRole } = useAuth();

  return (
    <div className="admin-container">
      <Sidebar />

      <main className="main-content">
        <Routes>
          {/* --- Rutas Principales --- */}
          <Route path="proyectos" element={<Portafolio />} />
          <Route path="usuarios" element={<Usuarios />} />
          
          {/* 1. Configuraciones (Lista de Proyectos) */}
          <Route path="configuraciones" element={<Portafolio />} />

          {/* 2. Planes de Acción (AHORA USA PORTAFOLIO TAMBIÉN) */}
          <Route path="planes-accion" element={<Portafolio />} />


          {/* --- Sub-Rutas de Configuración de Proyecto --- */}
          <Route path="configuraciones/proyecto/:id/labores" element={<LaboresAgronomicas />} />
          <Route path="configuraciones/proyecto/:id/equipos" element={<EquiposEImplementos />} />
          <Route path="configuraciones/proyecto/:id/unidades" element={<UnidadesMedida />} />

          {/* --- Sub-Rutas de Planes de Acción (Con ID) --- */}
          <Route path="planes-accion/proyecto/:id/general" element={<PlanDeAccion />} />
          <Route path="planes-accion/proyecto/:id/recursos" element={<RecursoHumano />} />
          <Route path="planes-accion/proyecto/:id/materiales" element={<MaterialesInsumos />} />

          {/* Datos del Proyecto */}
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