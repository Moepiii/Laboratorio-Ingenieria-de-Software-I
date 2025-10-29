import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { useAuth } from '../context/AuthContext';
import {
  getAllProjects,
  createProject,
  updateProject,
  deleteProject,
  setProjectState // Importa la función de estado
} from '../services/projectService';

// Objeto de estilos (incluye estilos para botones de estado)
const styles = {
  adminFormContainer: { padding: '1.5rem', backgroundColor: '#f9fafb', borderRadius: '8px', marginBottom: '2rem', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
  actionToolbar: { display: 'flex', flexWrap: 'wrap', gap: '0.75rem', marginBottom: '1rem', alignItems: 'center' },
  input: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', boxSizing: 'border-box' },
  searchInput: { padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', minWidth: '250px', marginLeft: 'auto' },
  inputGroup: { marginBottom: '1.25rem' },
  label: { display: 'block', marginBottom: '0.5rem', fontSize: '0.875rem', fontWeight: '500', color: '#374151' },
  buttonPrimary: { padding: '0.6rem 1.25rem', backgroundColor: '#4f46e5', color: 'white', border: 'none', borderRadius: '8px', cursor: 'pointer', fontWeight: '500', fontSize: '0.9rem', display: 'inline-flex', alignItems: 'center', gap: '0.5rem', transition: 'backgroundColor 0.2s' },
  buttonSecondary: { padding: '0.6rem 1.25rem', backgroundColor: '#ffffff', color: '#374151', border: '1px solid #d1d5db', borderRadius: '8px', cursor: 'pointer', fontWeight: '500', fontSize: '0.9rem', display: 'inline-flex', alignItems: 'center', gap: '0.5rem', transition: 'backgroundColor 0.2s, borderColor 0.2s' },
  buttonDanger: { padding: '0.6rem 1.25rem', backgroundColor: '#ef4444', color: 'white', border: 'none', borderRadius: '8px', cursor: 'pointer', fontWeight: '500', fontSize: '0.9rem', display: 'inline-flex', alignItems: 'center', gap: '0.5rem', transition: 'backgroundColor 0.2s' },
  // Estilos para botones de estado
  buttonSuccess: { padding: '0.6rem 1.25rem', backgroundColor: '#10b981', color: 'white', border: 'none', borderRadius: '8px', cursor: 'pointer', fontWeight: '500', fontSize: '0.9rem', display: 'inline-flex', alignItems: 'center', gap: '0.5rem', transition: 'backgroundColor 0.2s' }, // verde
  buttonWarning: { padding: '0.6rem 1.25rem', backgroundColor: '#f59e0b', color: 'white', border: 'none', borderRadius: '8px', cursor: 'pointer', fontWeight: '500', fontSize: '0.9rem', display: 'inline-flex', alignItems: 'center', gap: '0.5rem', transition: 'backgroundColor 0.2s' }, // ambar

  tableContainer: { overflowX: 'auto', backgroundColor: 'white', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06)' },
  table: { width: '100%', borderCollapse: 'collapse', fontSize: '0.9rem' },
  th: { padding: '0.75rem 1rem', textAlign: 'left', borderBottom: '2px solid #e5e7eb', backgroundColor: '#f9fafb', color: '#6b7280', fontWeight: '600', textTransform: 'uppercase', letterSpacing: '0.05em' },
  td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', color: '#374151' },
  selectedRow: { backgroundColor: '#eef2ff' },
  badgeSuccess: { backgroundColor: '#d1fae5', color: '#065f46', padding: '0.25rem 0.75rem', borderRadius: '9999px', fontSize: '0.875rem', fontWeight: '500' },
  badgeError: { backgroundColor: '#fee2e2', color: '#991b1b', padding: '0.25rem 0.75rem', borderRadius: '9999px', fontSize: '0.875rem', fontWeight: '500' },
  tableIconButton: { padding: '0.5rem', backgroundColor: '#ffffff', border: '1px solid #d1d5db', borderRadius: '6px', cursor: 'pointer', color: '#4f46e5', display: 'inline-flex', alignItems: 'center', justifyContent: 'center', fontSize: '0.9rem', lineHeight: 1, transition: 'background-color 0.2s, color 0.2s' },
  error: { fontSize: '0.875rem', color: '#dc2626', fontWeight: '500', backgroundColor: '#fef2f2', padding: '0.75rem', borderRadius: '8px', border: '1px solid #fecaca' },
  success: { fontSize: '0.875rem', color: '#059669', fontWeight: '500', backgroundColor: '#ecfdf5', padding: '0.75rem', borderRadius: '8px', border: '1px solid #a7f3d0' },
};

// Sub-componente ProyectoForm (sin cambios)
const ProyectoForm = ({ mode, formData, setFormData, handleSubmit, handleCancel, selectedProyecto }) => {
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };
  return (
    <div style={styles.adminFormContainer}>
      <h3 style={{ marginTop: 0, marginBottom: '1.5rem', fontSize: '1.25rem', fontWeight: '600' }}>
        {mode === 'create' ? 'Crear Nuevo Proyecto' : 'Modificar Proyecto'}
      </h3>
      <form onSubmit={handleSubmit}>
        <div style={styles.inputGroup}>
          <label htmlFor="nombre" style={styles.label}>Nombre del Proyecto</label>
          <input type="text" id="nombre" name="nombre" value={formData.nombre} onChange={handleChange} required style={styles.input} placeholder="Ej: Proyecto Titán" />
        </div>
        <div style={styles.inputGroup}>
          <label htmlFor="fecha_inicio" style={styles.label}>Fecha de Inicio</label>
          <input type="date" id="fecha_inicio" name="fecha_inicio" value={formData.fecha_inicio} onChange={handleChange} required style={styles.input} />
        </div>
        <div style={styles.inputGroup}>
          <label htmlFor="fecha_cierre" style={styles.label}>Fecha de Cierre (Opcional)</label>
          <input type="date" id="fecha_cierre" name="fecha_cierre" value={formData.fecha_cierre} onChange={handleChange} style={styles.input} />
        </div>
        {mode === 'edit' && selectedProyecto && (
          <div style={styles.inputGroup}>
            <label style={styles.label}>Estado Actual</label>
            <input type="text" value={selectedProyecto.estado === 'habilitado' ? 'Habilitado' : 'Cerrado'} readOnly style={{ ...styles.input, backgroundColor: '#f3f4f6', cursor: 'not-allowed', color: selectedProyecto.estado === 'habilitado' ? '#065f46' : '#991b1b', fontWeight: '500' }} />
          </div>
        )}
        <div style={{ display: 'flex', gap: '0.75rem', justifyContent: 'flex-end' }}>
          <button type="button" onClick={handleCancel} style={styles.buttonSecondary}> Cancelar </button>
          <button type="submit" style={styles.buttonPrimary}> {mode === 'create' ? 'Crear Proyecto' : 'Guardar Cambios'} </button>
        </div>
      </form>
    </div>
  );
};

// Sub-componente ListaProyectos (sin cambios)
const ListaProyectos = ({ proyectos, selectedProyectoId, setSelectedProyectoId, searchTerm }) => {
  const filteredProyectos = proyectos.filter(p =>
    p.nombre.toLowerCase().includes(searchTerm.toLowerCase())
  );
  const formatDateForDisplay = (dateString) => {
    if (!dateString || dateString === "0001-01-01T00:00:00Z" || dateString.startsWith("0001-01-01")) return null;
    try {
      let date = new Date(dateString);
      if (isNaN(date.getTime())) { date = new Date(dateString + 'T00:00:00'); }
      if (isNaN(date.getTime())) { return dateString; }
      const day = String(date.getDate()).padStart(2, '0');
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const year = date.getFullYear();
      return `${day}/${month}/${year}`;
    } catch (error) { return dateString; }
  };
  return (
    <div style={styles.tableContainer}>
      <table style={styles.table}>
        <thead>
          <tr>
            <th style={styles.th}>ID</th>
            <th style={styles.th}>Nombre</th>
            <th style={styles.th}>Fecha Inicio</th>
            <th style={styles.th}>Fecha Cierre</th>
            <th style={styles.th}>Estado</th>
            <th style={styles.th}>Generar</th>
          </tr>
        </thead>
        <tbody>
          {filteredProyectos.map((proyecto) => {
            const isSelected = proyecto.id === selectedProyectoId;
            return (
              <tr key={proyecto.id} style={isSelected ? styles.selectedRow : { cursor: 'pointer' }}
                onClick={() => setSelectedProyectoId(proyecto.id)}
              >
                <td style={styles.td}>{proyecto.id}</td>
                <td style={{ ...styles.td, fontWeight: '600' }}>{proyecto.nombre}</td>
                <td style={styles.td}>{formatDateForDisplay(proyecto.fecha_inicio)}</td>
                <td style={styles.td}>{formatDateForDisplay(proyecto.fecha_cierre) || <span style={{ color: '#9ca3af' }}>N/A</span>}</td>
                <td style={styles.td}>
                  {proyecto.estado === 'habilitado' ? (
                    <span style={styles.badgeSuccess}>Habilitado</span>
                  ) : (
                    <span style={styles.badgeError}>Cerrado</span>
                  )}
                </td>
                <td style={styles.td}>
                  <button style={styles.tableIconButton} onClick={(e) => { e.stopPropagation(); alert(`Mostrar detalles del proyecto ID: ${proyecto.id}`); }}
                    onMouseOver={(e) => { e.currentTarget.style.backgroundColor = '#eef2ff'; e.currentTarget.style.color = '#4338ca'; }}
                    onMouseOut={(e) => { e.currentTarget.style.backgroundColor = '#ffffff'; e.currentTarget.style.color = '#4f46e5'; }}
                    title="Ver detalles"
                  >
                    <i className="fas fa-search"></i>
                  </button>
                </td>
              </tr>
            )
          })}
          {filteredProyectos.length === 0 && (
            <tr><td colSpan="6" style={{ ...styles.td, textAlign: 'center', color: '#6b7280', padding: '1.5rem' }}>
              {searchTerm ? 'No se encontraron proyectos.' : 'No hay proyectos creados.'}
            </td></tr>
          )}
        </tbody>
      </table>
    </div>
  );
};


// Componente principal Portafolio (Refactorizado)
const Portafolio = () => {
  const { token, currentUser } = useAuth();

  const [proyectos, setProyectos] = useState([]);
  const [mode, setMode] = useState('list');
  const [selectedProyectoId, setSelectedProyectoId] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [formData, setFormData] = useState({
    nombre: '',
    fecha_inicio: '',
    fecha_cierre: '',
  });
  const [editingId, setEditingId] = useState(null);

  const selectedProyecto = useMemo(() => {
    return proyectos.find(p => p.id === selectedProyectoId);
  }, [proyectos, selectedProyectoId]);

  // Lógica de carga
  const fetchProyectos = useCallback(async () => {
    // Si no está cargando ya, ponlo a cargar
    // Evita cargas múltiples si currentUser cambia rápido
    if (!loading) setLoading(true);
    setError('');
    // Asegura que currentUser exista antes de llamar
    if (!currentUser) {
      setLoading(false);
      return;
    }
    try {
      const data = await getAllProjects(token, currentUser.username);
      setProyectos(data.proyectos || []);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [token, currentUser, loading]); // Añade 'loading' a dependencias

  useEffect(() => {
    fetchProyectos();
    // Quita fetchProyectos de las dependencias, ya está dentro del useCallback
    // eslint-disable-next-line react-hooks/exhaustive-deps 
  }, [token, currentUser]); // Llama solo cuando token o currentUser cambian

  // Lógica de CRUD
  const handleCrearProyecto = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(''); setSuccess('');
    try {
      const body = {
        nombre: formData.nombre,
        fecha_inicio: formData.fecha_inicio,
        fecha_cierre: formData.fecha_cierre || null,
        admin_username: currentUser.username
      };
      await createProject(token, body);
      setSuccess('Proyecto creado con éxito');
      setMode('list');
      fetchProyectos();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleActualizarProyecto = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(''); setSuccess('');
    try {
      const body = {
        id: editingId,
        nombre: formData.nombre,
        fecha_inicio: formData.fecha_inicio,
        fecha_cierre: formData.fecha_cierre || null,
        admin_username: currentUser.username
      };
      await updateProject(token, body);
      setSuccess('Proyecto actualizado con éxito');
      setMode('list');
      fetchProyectos();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleEliminarProyecto = async (id) => {
    if (!id) return; // No hacer nada si no hay ID
    if (window.confirm('¿Seguro que quieres eliminar este proyecto?')) {
      setError(''); setSuccess('');
      try {
        await deleteProject(token, id, currentUser.username);
        setSuccess('Proyecto eliminado');
        setSelectedProyectoId(null); // Deselecciona al borrar
        fetchProyectos();
      } catch (err) {
        setError(err.message);
      }
    }
  };

  // ⭐️⭐️ Lógica para Cambiar Estado ⭐️⭐️
  const handleChangeEstado = async (newState) => {
    if (!selectedProyecto) {
      alert("Por favor, selecciona un proyecto.");
      return;
    }
    const confirmMessage = `¿Seguro que quieres ${newState === 'habilitado' ? 'habilitar' : 'cerrar'} el proyecto "${selectedProyecto.nombre}"?`;
    if (window.confirm(confirmMessage)) {
      setError(''); setSuccess('');
      try {
        const result = await setProjectState(token, selectedProyecto.id, newState, currentUser.username);
        setSuccess(result.mensaje || `Proyecto ${newState}.`);
        fetchProyectos(); // Recarga la lista
      } catch (err) {
        setError(err.message);
      }
    }
  };

  // Handlers específicos para los botones
  const handleSetEstadoHabilitar = () => handleChangeEstado('habilitado');
  const handleSetEstadoCerrar = () => handleChangeEstado('cerrado');


  // Funciones auxiliares
  const handleEditClick = () => {
    if (!selectedProyecto) {
      alert("Por favor, selecciona un proyecto para modificar.");
      return;
    }
    if (selectedProyecto.estado === 'cerrado') {
      alert("No se puede modificar un proyecto cerrado.");
      return;
    }
    setMode('edit');
    setEditingId(selectedProyecto.id);
    setFormData({
      nombre: selectedProyecto.nombre,
      fecha_inicio: selectedProyecto.fecha_inicio.split('T')[0],
      fecha_cierre: selectedProyecto.fecha_cierre ? selectedProyecto.fecha_cierre.split('T')[0] : '',
    });
  };
  const handleCancel = () => {
    setMode('list');
    setEditingId(null);
    setFormData({ nombre: '', fecha_inicio: '', fecha_cierre: '' });
  };
  const handleNewProjectClick = () => {
    setMode('create');
    setEditingId(null);
    setFormData({ nombre: '', fecha_inicio: '', fecha_cierre: '' });
  };

  // JSX
  return (
    <div style={{ width: '100%' }}>
      {error && <p style={{ ...styles.error, marginBottom: '1rem' }}>{error}</p>}
      {success && <p style={{ ...styles.success, marginBottom: '1rem' }}>{success}</p>}

      {mode !== 'list' ? (
        <ProyectoForm
          mode={mode}
          formData={formData}
          setFormData={setFormData}
          handleSubmit={mode === 'create' ? handleCrearProyecto : handleActualizarProyecto}
          handleCancel={handleCancel}
          selectedProyecto={selectedProyecto}
        />
      ) : (
        <>
          {/* ⭐️⭐️ TOOLBAR CON TODOS LOS BOTONES ⭐️⭐️ */}
          <div style={styles.actionToolbar}>
            {/* Botones CRUD */}
            <button style={styles.buttonPrimary} onClick={handleNewProjectClick}>
              <i className="fas fa-plus" style={{ marginRight: '8px' }}></i>
              Crear
            </button>
            <button
              style={styles.buttonSecondary}
              onClick={handleEditClick}
              // Deshabilita si no hay selección O si está cerrado
              disabled={!selectedProyecto || selectedProyecto.estado === 'cerrado'}
            >
              <i className="fas fa-edit" style={{ marginRight: '8px' }}></i>
              Modificar
            </button>
            <button
              style={styles.buttonDanger}
              onClick={() => handleEliminarProyecto(selectedProyectoId)}
              disabled={!selectedProyecto}
            >
              <i className="fas fa-trash" style={{ marginRight: '8px' }}></i>
              Borrar
            </button>

            {/* Botones de Estado */}
            <button
              style={styles.buttonSuccess}
              onClick={handleSetEstadoHabilitar}
              disabled={!selectedProyecto || selectedProyecto.estado === 'habilitado'}
            >
              <i className="fas fa-check-circle" style={{ marginRight: '8px' }}></i>
              Habilitar
            </button>
            <button
              style={styles.buttonWarning}
              onClick={handleSetEstadoCerrar}
              disabled={!selectedProyecto || selectedProyecto.estado === 'cerrado'}
            >
              <i className="fas fa-times-circle" style={{ marginRight: '8px' }}></i>
              Cerrar
            </button>

            {/* Botón Imprimir (sin funcionalidad por ahora) */}
            <button style={{ ...styles.buttonPrimary, backgroundColor: '#1d4ed8' }}>
              Imprimir
            </button>

            {/* Búsqueda */}
            <input
              type="text"
              placeholder="Buscar por nombre..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              style={styles.searchInput}
            />
          </div>

          {/* Lista de Proyectos */}
          {loading ? (
            <p>Cargando proyectos...</p>
          ) : (
            <ListaProyectos
              proyectos={proyectos}
              selectedProyectoId={selectedProyectoId}
              setSelectedProyectoId={setSelectedProyectoId}
              searchTerm={searchTerm}
            />
          )}
        </>
      )}
    </div>
  );
};

export default Portafolio;