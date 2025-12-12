import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { useLocation, useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import {
  getAllProjects,
  createProject,
  updateProject,
  deleteProject,
  setProjectState
} from '../services/projectService';


const ListaProyectos = ({ proyectos, selectedProyectoId, setSelectedProyectoId, searchTerm, showToolbar, navigate, location }) => {

  const filteredProyectos = useMemo(() => {
    return proyectos.filter(p => p.nombre.toLowerCase().includes(searchTerm.toLowerCase()));
  }, [proyectos, searchTerm]);

  const handleRowClick = (proyecto) => {
    setSelectedProyectoId(proyecto.id);
    
    // Si no es modo toolbar (es decir, estamos en la lista principal)
    if (!showToolbar) {
        const currentPath = location.pathname;

        
        if (currentPath.includes('/configuraciones')) {
            // Si venimos de Configuraciones -> Vamos a Labores
            navigate(`/admin/configuraciones/proyecto/${proyecto.id}/labores`);
        } else if (currentPath.includes('/planes-accion')) {
            // Si venimos de Planes de Acción -> Vamos a Plan General
            navigate(`/admin/planes-accion/proyecto/${proyecto.id}/general`);
        } else {
            // Si venimos de Portafolio -> Vamos a Datos del Proyecto (Comportamiento por defecto)
            navigate(`/admin/proyectos/datos/${proyecto.id}`);
        }
    }
  };

  return (
    <div style={styles.tableContainer}>
      <table style={styles.table}>
        <thead>
          <tr>
            <th style={styles.th}></th>
            <th style={styles.th}>Nombre del Proyecto</th>
            <th style={styles.th}>Fecha Inicio</th>
            <th style={styles.th}>Fecha Cierre</th>
            <th style={styles.th}>Estado</th>
          </tr>
        </thead>
        <tbody>
          {filteredProyectos.length > 0 ? (
            filteredProyectos.map(proyecto => (
              <tr
                key={proyecto.id}
                onClick={() => handleRowClick(proyecto)}
                style={selectedProyectoId === proyecto.id ? styles.trSelected : styles.tr}
              >
                <td style={styles.td}>
                  <input
                    type="radio"
                    checked={selectedProyectoId === proyecto.id}
                    onChange={() => handleRowClick(proyecto)}
                  />
                </td>
                <td style={styles.td}>{proyecto.nombre}</td>
                <td style={styles.td}>{proyecto.fecha_inicio}</td>
                <td style={styles.td}>{proyecto.fecha_cierre}</td>
                <td style={styles.td}>
                  <span style={proyecto.estado === 'habilitado' ? styles.statusEnabled : styles.statusDisabled}>
                    {proyecto.estado}
                  </span>
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="5" style={{ ...styles.td, textAlign: 'center' }}>No hay proyectos para mostrar.</td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};



// --- Componente Principal Portafolio ---
const Portafolio = () => {
  const [proyectos, setProyectos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedProyectoId, setSelectedProyectoId] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [showForm, setShowForm] = useState(false);

  const [nombre, setNombre] = useState('');
  const [fechaInicio, setFechaInicio] = useState('');
  const [fechaCierre, setFechaCierre] = useState('');

  const { token, currentUser } = useAuth();
  const adminUsername = currentUser?.username;

  const navigate = useNavigate();
  const location = useLocation();
  const showToolbar = location.pathname.startsWith('/admin/proyectos');

  const selectedProyecto = useMemo(() => {
    return proyectos.find(p => p.id === selectedProyectoId) || null;
  }, [proyectos, selectedProyectoId]);

  const fetchProyectos = useCallback(async () => {
    if (!token || !adminUsername) return;
    setLoading(true);
    try {
      const data = await getAllProjects(token, adminUsername);
      setProyectos(data.proyectos || []);
    } catch (err) {
      setError(err.message || 'Error al cargar proyectos');
    } finally {
      setLoading(false);
    }
  }, [token, adminUsername]);

  useEffect(() => {
    fetchProyectos();
  }, [fetchProyectos]);

  useEffect(() => {
    if (selectedProyecto && showForm) {
      setNombre(selectedProyecto.nombre);
      setFechaInicio(selectedProyecto.fecha_inicio);
      setFechaCierre(selectedProyecto.fecha_cierre);
    } else {
      setNombre('');
      setFechaInicio('');
      setFechaCierre('');
    }
  }, [selectedProyecto, showForm]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    const proyectoData = {
      nombre,
      fecha_inicio: fechaInicio,
      fecha_cierre: fechaCierre,
      admin_username: adminUsername
    };

    try {
      if (selectedProyectoId) { 
        const updated = await updateProject(token, { ...proyectoData, id: selectedProyectoId });
        setProyectos(proyectos.map(p => p.id === selectedProyectoId ? updated : p));
      } else { // Crear
        const newProject = await createProject(token, proyectoData);
        setProyectos([newProject, ...proyectos]);
      }
      setShowForm(false);
      setSelectedProyectoId(null);
    } catch (err) {
      setError(err.message || 'Error al guardar el proyecto');
    }
  };

  const handleBorrarProyecto = async (id) => {
    if (!id || !window.confirm('¿Estás seguro de que quieres borrar este proyecto?')) return;
    setError('');
    try {
      await deleteProject(token, id, adminUsername);
      setProyectos(proyectos.filter(p => p.id !== id));
      setSelectedProyectoId(null);
    } catch (err) {
      setError(err.message || 'Error al borrar el proyecto');
    }
  };

  const handleSetEstado = async (estado) => {
    if (!selectedProyectoId) return;
    setError('');
    try {
      await setProjectState(token, selectedProyectoId, estado, adminUsername);
      setProyectos(proyectos.map(p =>
        p.id === selectedProyectoId ? { ...p, estado: estado } : p
      ));
    } catch (err) {
      setError(err.message || 'Error al cambiar el estado');
    }
  };

  const handleSetEstadoHabilitar = () => handleSetEstado('habilitado');
  const handleSetEstadoCerrar = () => handleSetEstado('cerrado');

  return (
    <div style={styles.container}>
      {showForm ? (
        // 
        <div style={styles.adminFormContainer}>
          <form onSubmit={handleSubmit}>

            {/* 1. CABECERA (SOLO TÍTULO) */}
            <div style={styles.formHeader}>
              <h3 style={styles.formTitle}>
                {selectedProyectoId ? 'Actualizar Proyecto' : 'Crear Nuevo Proyecto'}
              </h3>
            </div>

            {/* 2. BOTONES MOVIDOS AQUÍ */}
            <div style={styles.buttonGroup}>
              <button type="submit" style={styles.button}>
                {selectedProyectoId ? 'Actualizar' : 'Guardar'}
              </button>
              <button
                type="button"
                style={styles.buttonWarning}
                onClick={() => { setShowForm(false); setSelectedProyectoId(null); }}
              >
                Cancelar
              </button>
            </div>

            {/* 3. CAMPOS DEL FORMULARIO */}
            <div style={styles.inputGroup}>
              <label style={styles.label}>Nombre del Proyecto</label>
              <input type="text" style={styles.input} value={nombre} onChange={(e) => setNombre(e.target.value)} required />
            </div>
            <div style={styles.inputGroup}>
              <label style={styles.label}>Fecha Inicio</label>
              <input type="date" style={styles.input} value={fechaInicio} onChange={(e) => setFechaInicio(e.target.value)} required />
            </div>
            <div style={styles.inputGroup}>
              <label style={styles.label}>Fecha Cierre</label>
              <input type="date" style={styles.input} value={fechaCierre} onChange={(e) => setFechaCierre(e.target.value)} required />
            </div>

            {error && <p style={styles.error}>{error}</p>}

          </form>
        </div>
        // 
      ) : (
        <>
          <h2 style={styles.h2}>{showToolbar ? 'Portafolio de Proyectos' : 'Configuraciones'}</h2>
          {showToolbar && (
            <div style={styles.actionToolbar}>
              <button style={styles.button} onClick={() => { setShowForm(true); setSelectedProyectoId(null); }}>
                <i className="fas fa-plus" style={{ marginRight: '8px' }}></i>
                Crear
              </button>

              <Link
                to={`/admin/proyectos/datos/${selectedProyectoId}`}
                style={{
                  ...styles.buttonInfo,
                  ...(!selectedProyecto ?
                    { backgroundColor: '#9ca3af', cursor: 'not-allowed', opacity: 0.7 }
                    :
                    {}
                  )
                }}
                onClick={(e) => {
                  if (!selectedProyecto) {
                    e.preventDefault();
                  }
                }}
              >
                <i className="fas fa-database" style={{ marginRight: '8px' }}></i>
                Datos
              </Link>

              <button style={styles.buttonDanger} onClick={() => handleBorrarProyecto(selectedProyectoId)} disabled={!selectedProyecto} >
                <i className="fas fa-trash" style={{ marginRight: '8px' }}></i>
                Borrar
              </button>
              <button style={styles.buttonSuccess} onClick={handleSetEstadoHabilitar} disabled={!selectedProyecto || selectedProyecto.estado === 'habilitado'} >
                <i className="fas fa-check-circle" style={{ marginRight: '8px' }}></i>
                Habilitar
              </button>
              <button style={styles.buttonWarning} onClick={handleSetEstadoCerrar} disabled={!selectedProyecto || selectedProyecto.estado === 'cerrado'} >
                <i className="fas fa-times-circle" style={{ marginRight: '8px' }}></i>
                Cerrar
              </button>
              <button style={{ ...styles.buttonPrimary, backgroundColor: '#1d4ed8' }}> Imprimir </button>
              <input type="text" placeholder="Buscar por nombre..." value={searchTerm} onChange={(e) => setSearchTerm(e.target.value)} style={styles.searchInput} />
            </div>
          )}
          {error && <p style={styles.error}>{error}</p>}
          {loading ? (
            <p>Cargando proyectos...</p>
          ) : (
            <ListaProyectos
              proyectos={proyectos}
              selectedProyectoId={selectedProyectoId}
              setSelectedProyectoId={setSelectedProyectoId}
              searchTerm={searchTerm}
              showToolbar={showToolbar}
              navigate={navigate}
              location={location}
            />
          )}
        </>
      )}
    </div>
  );
};

// --- ESTILOS (CON NUEVOS ESTILOS AÑADIDOS) ---
const styles = {
  container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
  h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', marginBottom: '1.5rem' },
  adminFormContainer: { padding: '1.5rem', backgroundColor: '#f9fafb', borderRadius: '8px', marginBottom: '2rem', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
  actionToolbar: { display: 'flex', flexWrap: 'wrap', gap: '0.75rem', marginBottom: '1rem', alignItems: 'center' },
  input: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', boxSizing: 'border-box' },
  searchInput: { padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', minWidth: '250px', marginLeft: 'auto' },
  inputGroup: { marginBottom: '1.25rem' },
  label: { display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' },

  // 
  formHeader: {
    // Ya no es 'flex', solo un bloque normal
    marginBottom: '1.5rem', // Espacio después del título
    borderBottom: '1px solid #e5e7eb',
    paddingBottom: '1rem'
  },
  formTitle: {
    fontSize: '1.25rem',
    fontWeight: '600',
    color: '#111827',
    margin: 0
  },
  buttonGroup: {
    display: 'flex',
    gap: '0.75rem', // Espacio entre botones
    marginBottom: '1.5rem' // Espacio ANTES del primer campo "Nombre del Proyecto"
  },
  // 

  button: { padding: '0.6rem 1.2rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#4f46e5', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s', display: 'inline-flex', alignItems: 'center' },
  buttonSuccess: { padding: '0.6rem 1.2rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#22c55e', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s', display: 'inline-flex', alignItems: 'center' },
  buttonWarning: { padding: '0.6rem 1.2rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#f59e0b', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s', display: 'inline-flex', alignItems: 'center' },
  buttonDanger: { padding: '0.6rem 1.2rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#ef4444', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s', display: 'inline-flex', alignItems: 'center' },
  buttonPrimary: { padding: '0.6rem 1.2rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#3b82f6', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s', display: 'inline-flex', alignItems: 'center' },
  buttonInfo: {
    padding: '0.6rem 1.2rem',
    fontSize: '1rem',
    fontWeight: '600',
    borderRadius: '8px',
    color: 'white',
    backgroundColor: '#3b82f6', // blue-500
    border: 'none',
    cursor: 'pointer',
    transition: 'background-color 0.2s',
    textDecoration: 'none',
    display: 'inline-flex',
    alignItems: 'center',
  },
  error: { color: 'red', marginTop: '1rem' },
  tableContainer: { overflowX: 'auto', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
  table: { width: '100%', borderCollapse: 'collapse' },
  th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '2px solid #e5e7eb', color: '#374151', fontWeight: '600', fontSize: '0.875rem' },
  td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', verticalAlign: 'middle' },
  tr: { cursor: 'pointer', transition: 'background-color 0.15s' },
  trSelected: { cursor: 'pointer', backgroundColor: '#eef2ff' },
  statusEnabled: { color: '#16a34a', backgroundColor: '#dcfce7', padding: '0.25rem 0.5rem', borderRadius: '99px', fontSize: '0.875rem', fontWeight: '500' },
  statusDisabled: { color: '#b91c1c', backgroundColor: '#fee2e2', padding: '0.25rem 0.5rem', borderRadius: '99px', fontSize: '0.875rem', fontWeight: '500' },
};

export default Portafolio;