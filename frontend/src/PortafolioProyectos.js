import React, { useState, useEffect, useCallback } from 'react';

// OBJETO DE ESTILOS (NECESARIO EN ESTE ARCHIVO)
const styles = {
  // Contenedor del formulario
  adminFormContainer: {
    padding: '1.5rem',
    backgroundColor: '#f9fafb', // gray-50
    borderRadius: '8px',
    marginBottom: '2rem',
    border: '1px solid #e5e7eb',
  },
  // Contenedor de botones de acción
  actionToolbar: {
    display: 'flex',
    flexWrap: 'wrap', // Para responsividad
    gap: '0.75rem',
    marginBottom: '1rem',
    alignItems: 'center',
  },
  // Estilo para inputs (texto y fechas)
  input: {
    width: '100%',
    padding: '0.75rem 1rem',
    border: '1px solid #d1d5db',
    borderRadius: '8px',
    fontSize: '1rem',
    boxSizing: 'border-box', // Importante
  },
  // Estilo para el input de búsqueda
  searchInput: {
    padding: '0.75rem 1rem',
    border: '1px solid #d1d5db',
    borderRadius: '8px',
    fontSize: '1rem',
    minWidth: '250px',
    marginLeft: 'auto', // Mueve la búsqueda a la derecha
  },
  // Estilo para botones genéricos
  button: {
    padding: '0.75rem 1rem',
    fontSize: '0.875rem',
    fontWeight: '600',
    borderRadius: '8px',
    color: 'white',
    border: 'none',
    cursor: 'pointer',
    transition: 'background-color 0.2s',
    display: 'flex',
    alignItems: 'center',
    gap: '0.5rem',
  },
  // Estilo para mensajes de error
  error: {
    fontSize: '0.875rem',
    color: '#dc2626', // red-600
    fontWeight: '500',
    backgroundColor: '#fef2f2', // red-50
    padding: '0.75rem',
    borderRadius: '8px',
    border: '1px solid #fecaca', // red-200
  },
  // Estilo para mensajes de éxito
  success: {
    fontSize: '0.875rem',
    color: '#059669', // emerald-600
    fontWeight: '500',
    backgroundColor: '#ecfdf5', // emerald-50
    padding: '0.75rem',
    borderRadius: '8px',
    border: '1px solid #a7f3d0', // emerald-200
  },
  // Contenedor para los campos del formulario
  formGrid: {
    display: 'grid',
    gridTemplateColumns: '1fr 1fr', // 2 columnas
    gap: '1rem',
  },
  formField: {
    display: 'flex',
    flexDirection: 'column',
  },
  label: {
    display: 'block',
    fontSize: '0.875rem',
    fontWeight: '500',
    color: '#374151',
    marginBottom: '0.25rem',
  },
  formFullWidth: {
    gridColumn: '1 / -1', // Ocupa todo el ancho
  },
  formActions: {
    gridColumn: '1 / -1',
    display: 'flex',
    gap: '1rem',
    marginTop: '1rem',
  }
};


const PortafolioProyectos = ({ apiCall, currentUser }) => {

  // Lista de todos los proyectos
  const [proyectos, setProyectos] = useState([]);

  // --- Estados de UI ---
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedProyectoId, setSelectedProyectoId] = useState(null);

  // --- Estados del Formulario ---
  const [isEditing, setIsEditing] = useState(false);
  const [currentProyectoId, setCurrentProyectoId] = useState(null); // Para saber qué ID editar
  const [nombre, setNombre] = useState('');
  const [fechaInicio, setFechaInicio] = useState('');
  const [fechaCierre, setFechaCierre] = useState('');

  // ⭐️ --- NUEVA FUNCIÓN PARA FORMATEAR FECHAS --- ⭐️
  /**
   * Convierte una fecha de 'AAAA-MM-DD' a 'DD/MM/AAAA'.
   * @param {string} dateString La fecha en formato 'AAAA-MM-DD'.
   * @returns {string | null} La fecha formateada o null si no hay fecha.
   */
  const formatDate = (dateString) => {
    if (!dateString || dateString.length === 0) {
      return null;
    }
    try {
      // Separa la fecha (ej: '2025-10-20')
      const [year, month, day] = dateString.split('-');
      // Devuelve la fecha en el nuevo formato (ej: '20/10/2025')
      if (day && month && year) {
        return `${day}/${month}/${year}`;
      }
      return dateString; // Devuelve el original si el formato falla
    } catch (e) {
      return dateString; // Devuelve el original en caso de error
    }
  };

  // 1. Cargar Proyectos (sin cambios)
  const fetchProyectos = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const result = await apiCall('admin/get-proyectos', { admin_username: currentUser }, 'POST');
      if (result.success) {
        setProyectos(result.data.proyectos || []);
      } else {
        setError(result.data.error || 'No se pudieron cargar los proyectos.');
      }
    } catch (e) {
      setError(`Error de conexión: ${e.message}`);
    } finally {
      setLoading(false);
    }
  }, [apiCall, currentUser]);

  useEffect(() => {
    fetchProyectos();
  }, [fetchProyectos]);

  // 2. Limpiar formulario y resetear modos (sin cambios)
  const resetForm = () => {
    setNombre('');
    setFechaInicio('');
    setFechaCierre('');
    setIsEditing(false);
    setCurrentProyectoId(null);
  };

  // 3. Botón "Agregar Nuevo Proyecto" (sin cambios)
  const handleShowAddForm = () => {
    resetForm();
    setSelectedProyectoId(null);
    setIsEditing(false);
  };

  // 4. Botón "Modificar Proyecto" (sin cambios)
  // (Nota: Esto funciona porque el <input type="date"> SÍ espera 'AAAA-MM-DD')
  const handleStartEdit = () => {
    if (!selectedProyectoId) return;

    const proyectoToEdit = proyectos.find(p => p.id === selectedProyectoId);
    if (proyectoToEdit) {
      setNombre(proyectoToEdit.nombre);
      setFechaInicio(proyectoToEdit.fecha_inicio || '');
      setFechaCierre(proyectoToEdit.fecha_cierre || '');
      setCurrentProyectoId(proyectoToEdit.id);
      setIsEditing(true);
    }
  };

  // 5. Botón "Eliminar Proyecto" (sin cambios)
  const handleDelete = async () => {
    if (!selectedProyectoId) return;

    const proyectoToDelete = proyectos.find(p => p.id === selectedProyectoId);
    if (!window.confirm(`¿Seguro que quieres borrar el proyecto "${proyectoToDelete.nombre}"?`)) {
      return;
    }

    setError('');
    setSuccess('');

    try {
      const result = await apiCall('admin/delete-proyecto', {
        id: selectedProyectoId,
        admin_username: currentUser
      }, 'POST');

      if (result.success) {
        setSuccess(result.data.mensaje);
        fetchProyectos();
        setSelectedProyectoId(null);
        resetForm();
      } else {
        setError(result.data.error || 'Fallo al borrar el proyecto.');
      }
    } catch (e) {
      setError(`Error de conexión: ${e.message}`);
    }
  };

  // 6. Submit del Formulario (Crea o Modifica) (sin cambios)
  // (Nota: Esto funciona porque el backend SÍ espera 'AAAA-MM-DD')
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    if (!nombre || !fechaInicio) {
      setError('El Nombre y la Fecha de Inicio son obligatorios.');
      return;
    }

    const proyectoData = {
      nombre,
      fecha_inicio: fechaInicio,
      fecha_cierre: fechaCierre,
      admin_username: currentUser
    };

    try {
      let result;
      if (isEditing) {
        result = await apiCall('admin/update-proyecto', {
          id: currentProyectoId,
          ...proyectoData
        }, 'POST');
      } else {
        result = await apiCall('admin/create-proyecto', proyectoData, 'POST');
      }

      if (result.success) {
        setSuccess(result.data.mensaje);
        resetForm();
        fetchProyectos();
      } else {
        setError(result.data.error || 'Fallo al guardar el proyecto.');
      }
    } catch (e) {
      setError(`Error de conexión: ${e.message}`);
    }
  };

  // 7. Filtrar proyectos basado en la búsqueda (sin cambios)
  const filteredProyectos = proyectos.filter(p =>
    p.nombre.toLowerCase().includes(searchTerm.toLowerCase())
  );


  return (
    <div>
      <h2 style={{ fontSize: '1.875rem', fontWeight: '700', color: '#1f2937' }}>
        Portafolio de Proyectos
      </h2>

      {success && <p style={{ ...styles.success, margin: '1rem 0' }}>{success}</p>}
      {error && <p style={{ ...styles.error, margin: '1rem 0' }}>{error}</p>}

      {/* --- Barra de Botones y Búsqueda (sin cambios) --- */}
      <div style={styles.actionToolbar}>
        <button
          onClick={handleShowAddForm}
          style={{ ...styles.button, backgroundColor: '#10b981' }} // emerald-500
          onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#059669'}
          onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#10b981'}
        >
          <span>➕</span> Agregar Nuevo Proyecto
        </button>
        <button
          onClick={handleStartEdit}
          disabled={!selectedProyectoId}
          style={{ ...styles.button, backgroundColor: '#3b82f6', opacity: !selectedProyectoId ? 0.5 : 1, cursor: !selectedProyectoId ? 'not-allowed' : 'pointer' }} // blue-500
          onMouseOver={(e) => !e.currentTarget.disabled && (e.currentTarget.style.backgroundColor = '#2563eb')}
          onMouseOut={(e) => !e.currentTarget.disabled && (e.currentTarget.style.backgroundColor = '#3b82f6')}
        >
          <span>✏️</span> Modificar Proyecto
        </button>
        <button
          onClick={handleDelete}
          disabled={!selectedProyectoId}
          style={{ ...styles.button, backgroundColor: '#ef4444', opacity: !selectedProyectoId ? 0.5 : 1, cursor: !selectedProyectoId ? 'not-allowed' : 'pointer' }} // red-500
          onMouseOver={(e) => !e.currentTarget.disabled && (e.currentTarget.style.backgroundColor = '#dc2626')}
          onMouseOut={(e) => !e.currentTarget.disabled && (e.currentTarget.style.backgroundColor = '#ef4444')}
        >
          <span>❌</span> Eliminar Proyecto
        </button>
        <input
          type="text"
          placeholder="Buscar Proyecto por nombre..."
          style={styles.searchInput}
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
      </div>

      {/* --- Formulario de Agregar / Modificar (sin cambios) --- */}
      <div style={{ ...styles.adminFormContainer, marginTop: '1rem' }}>
        <h3 style={{ fontSize: '1.25rem', fontWeight: '700', color: '#1f2937', marginBottom: '1rem' }}>
          {isEditing ? 'Modificar Proyecto ✏️' : 'Agregar Nuevo Proyecto ➕'}
        </h3>
        <form onSubmit={handleSubmit}>
          <div style={styles.formGrid}>
            {/* Campo Nombre */}
            <div style={{ ...styles.formField, ...styles.formFullWidth }}>
              <label htmlFor="nombre" style={styles.label}>Descripción (Nombre del Proyecto)</label>
              <input
                id="nombre"
                type="text"
                placeholder="Nombre del proyecto..."
                value={nombre}
                onChange={(e) => setNombre(e.target.value)}
                required
                style={styles.input}
              />
            </div>

            {/* Campo Fecha Inicio */}
            <div style={styles.formField}>
              <label htmlFor="fechaInicio" style={styles.label}>Fecha de Inicio</label>
              <input
                id="fechaInicio"
                type="date"
                value={fechaInicio}
                onChange={(e) => setFechaInicio(e.target.value)}
                required
                style={styles.input}
              />
            </div>

            {/* Campo Fecha Cierre */}
            <div style={styles.formField}>
              <label htmlFor="fechaCierre" style={styles.label}>Fecha de Cierre (Opcional)</label>
              <input
                id="fechaCierre"
                type="date"
                value={fechaCierre}
                onChange={(e) => setFechaCierre(e.target.value)}
                style={styles.input}
              />
            </div>

            {/* Botones del Formulario */}
            <div style={styles.formActions}>
              <button
                type="submit"
                style={{ ...styles.button, backgroundColor: '#4f46e5' }} // indigo-600
                onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#4338ca'}
                onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#4f46e5'}
              >
                {isEditing ? 'Guardar Cambios' : 'Agregar Proyecto'}
              </button>
              {isEditing && (
                <button
                  type="button"
                  onClick={resetForm}
                  style={{ ...styles.button, backgroundColor: '#6b7280' }} // gray-500
                  onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#4b5563'}
                  onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#6b7280'}
                >
                  Cancelar Edición
                </button>
              )}
            </div>
          </div>
        </form>
      </div>

      {/* --- Tabla de Proyectos Existentes --- */}
      <h3 style={{ fontSize: '1.25rem', fontWeight: '700', color: '#1f2937', marginBottom: '1rem' }}>
        Proyectos Existentes (Total: {filteredProyectos.length})
      </h3>
      {loading ? (
        <p>Cargando proyectos...</p>
      ) : (
        <div style={{ maxHeight: '600px', overflowY: 'auto', border: '1px solid #e5e7eb', borderRadius: '8px' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left' }}>
            <thead>
              <tr style={{ backgroundColor: '#f9fafb', borderBottom: '1px solid #e5e7eb' }}>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>ID</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Descripción</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Inicio</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Cierre</th>
              </tr>
            </thead>
            <tbody>
              {filteredProyectos.map((proyecto, index) => {
                const isSelected = proyecto.id === selectedProyectoId;
                return (
                  <tr
                    key={proyecto.id}
                    style={{
                      borderBottom: index < filteredProyectos.length - 1 ? '1px solid #f3f4f6' : 'none',
                      backgroundColor: isSelected ? '#eef2ff' : 'white', // Resaltar fila seleccionada
                      cursor: 'pointer'
                    }}
                    onClick={() => setSelectedProyectoId(proyecto.id)} // Seleccionar fila al hacer clic
                  >
                    <td style={{ padding: '0.75rem', fontWeight: '500' }}>{proyecto.id}</td>
                    <td style={{ padding: '0.75rem', fontWeight: '600' }}>{proyecto.nombre}</td>
                    {/* ⭐️ --- CAMBIO AQUÍ --- ⭐️ */}
                    <td style={{ padding: '0.75rem' }}>
                      {formatDate(proyecto.fecha_inicio) || <span style={{ color: '#9ca3af' }}>N/A</span>}
                    </td>
                    {/* ⭐️ --- CAMBIO AQUÍ --- ⭐️ */}
                    <td style={{ padding: '0.75rem' }}>
                      {formatDate(proyecto.fecha_cierre) || <span style={{ color: '#9ca3af' }}>N/A</span>}
                    </td>
                  </tr>
                )
              })}
              {filteredProyectos.length === 0 && (
                <tr><td colSpan="4" style={{ textAlign: 'center', padding: '1rem', color: '#6b7280' }}>
                  {searchTerm ? 'No se encontraron proyectos con ese nombre.' : 'No hay proyectos creados.'}
                </td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}

    </div>
  );
};

export default PortafolioProyectos;