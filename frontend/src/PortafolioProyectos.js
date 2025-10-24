import React, { useState, useEffect, useCallback, useMemo } from 'react';

// OBJETO DE ESTILOS (ACTUALIZADO CON NUEVOS ESTILOS)
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
    marginLeft: 'auto', // Mueve el input al final
  },
  // Contenedor para un grupo de input + label
  inputGroup: {
    marginBottom: '1.25rem',
  },
  // Estilo para labels
  label: {
    display: 'block',
    marginBottom: '0.5rem',
    fontSize: '0.875rem',
    fontWeight: '500',
    color: '#374151', // gray-700
  },
  // Botón Primario (Crear, Guardar)
  buttonPrimary: {
    padding: '0.6rem 1.25rem',
    backgroundColor: '#4f46e5', // indigo-600
    color: 'white',
    border: 'none',
    borderRadius: '8px',
    cursor: 'pointer',
    fontWeight: '500',
    fontSize: '0.9rem',
    display: 'inline-flex',
    alignItems: 'center',
    gap: '0.5rem',
    transition: 'backgroundColor 0.2s',
  },
  // Botón Secundario (Modificar, Cancelar, etc.)
  buttonSecondary: {
    padding: '0.6rem 1.25rem',
    backgroundColor: '#ffffff',
    color: '#374151', // gray-700
    border: '1px solid #d1d5db', // gray-300
    borderRadius: '8px',
    cursor: 'pointer',
    fontWeight: '500',
    fontSize: '0.9rem',
    display: 'inline-flex',
    alignItems: 'center',
    gap: '0.5rem',
    transition: 'backgroundColor 0.2s, borderColor 0.2s',
  },
  // Botón de Peligro (Borrar)
  buttonDanger: {
    padding: '0.6rem 1.25rem',
    backgroundColor: '#ef4444', // red-500
    color: 'white',
    border: 'none',
    borderRadius: '8px',
    cursor: 'pointer',
    fontWeight: '500',
    fontSize: '0.9rem',
    display: 'inline-flex',
    alignItems: 'center',
    gap: '0.5rem',
    transition: 'backgroundColor 0.2s',
  },
  // NUEVO: Botón de Advertencia (Cerrar)
  buttonWarning: {
    padding: '0.6rem 1.25rem',
    backgroundColor: '#f59e0b', // amber-500
    color: 'white',
    border: 'none',
    borderRadius: '8px',
    cursor: 'pointer',
    fontWeight: '500',
    fontSize: '0.9rem',
    display: 'inline-flex',
    alignItems: 'center',
    gap: '0.5rem',
    transition: 'backgroundColor 0.2s',
  },
  // Contenedor de la tabla
  tableContainer: {
    overflowX: 'auto',
    backgroundColor: 'white',
    borderRadius: '8px',
    border: '1px solid #e5e7eb',
    boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06)',
  },
  modalBackdrop: {
    position: 'fixed',
    top: 0,
    left: 0,
    width: '100%',
    height: '100%',
    backgroundColor: 'rgba(0, 0, 0, 0.5)', // Fondo oscuro semitransparente
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000,
  },
  // NUEVO: Contenedor del modal
  modalContent: {
    backgroundColor: 'white',
    borderRadius: '8px',
    padding: '2rem',
    width: '100%',
    maxWidth: '500px',
    boxShadow: '0 10px 25px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
    zIndex: 1001,
  },
  // MEJORA: Sombra para el formulario (si decides no usar el modal)
  adminFormContainer: {
    padding: '1.5rem',
    backgroundColor: '#f9fafb',
    borderRadius: '8px',
    marginBottom: '2rem',
    border: '1px solid #e5e7eb',
    boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)', // Sombra añadida
  },
  // Tabla
  table: {
    width: '100%',
    borderCollapse: 'collapse',
    fontSize: '0.9rem',
  },
  // Encabezado de la tabla
  th: {
    padding: '0.75rem 1rem',
    textAlign: 'left',
    borderBottom: '2px solid #e5e7eb',
    backgroundColor: '#f9fafb',
    color: '#6b7280',
    fontWeight: '600',
    textTransform: 'uppercase',
    letterSpacing: '0.05em',
  },
  // Celda de la tabla
  td: {
    padding: '0.75rem 1rem',
    borderBottom: '1px solid #e5e7eb',
    color: '#374151',
  },
  // Fila seleccionada
  selectedRow: {
    backgroundColor: '#eef2ff', // indigo-50
  },
  // NUEVO: Pastilla de estado Habilitado
  badgeSuccess: {
    backgroundColor: '#d1fae5', // green-100
    color: '#065f46', // green-800
    padding: '0.25rem 0.75rem',
    borderRadius: '9999px',
    fontSize: '0.875rem',
    fontWeight: '500',
  },
  // NUEVO: Pastilla de estado Cerrado
  badgeError: {
    backgroundColor: '#fee2e2', // red-100
    color: '#991b1b', // red-800
    padding: '0.25rem 0.75rem',
    borderRadius: '9999px',
    fontSize: '0.875rem',
    fontWeight: '500',
  },
  // ⭐️ 1. NUEVO: Estilo para el botón de icono en la tabla
  tableIconButton: {
    padding: '0.5rem',
    backgroundColor: '#ffffff',
    border: '1px solid #d1d5db', // gray-300
    borderRadius: '6px',
    cursor: 'pointer',
    color: '#4f46e5', // indigo-600
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'center',
    fontSize: '0.9rem',
    lineHeight: 1,
    transition: 'background-color 0.2s, color 0.2s',
  },
};

// Componente de Formulario (ACTUALIZADO)
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
          <input
            type="text"
            id="nombre"
            name="nombre"
            value={formData.nombre}
            onChange={handleChange}
            required
            style={styles.input}
            placeholder="Ej: Proyecto Titán"
          />
        </div>
        <div style={styles.inputGroup}>
          <label htmlFor="fecha_inicio" style={styles.label}>Fecha de Inicio</label>
          <input
            type="date"
            id="fecha_inicio"
            name="fecha_inicio"
            value={formData.fecha_inicio}
            onChange={handleChange}
            required
            style={styles.input}
          />
        </div>
        <div style={styles.inputGroup}>
          <label htmlFor="fecha_cierre" style={styles.label}>Fecha de Cierre (Opcional)</label>
          <input
            type="date"
            id="fecha_cierre"
            name="fecha_cierre"
            value={formData.fecha_cierre}
            onChange={handleChange}
            style={styles.input}
          />
        </div>

        {mode === 'edit' && selectedProyecto && (
          <div style={styles.inputGroup}>
            <label style={styles.label}>Estado Actual</label>
            <input
              type="text"
              value={selectedProyecto.estado === 'habilitado' ? 'Habilitado' : 'Cerrado'}
              readOnly
              style={{ ...styles.input, backgroundColor: '#f3f4f6', cursor: 'not-allowed', color: selectedProyecto.estado === 'habilitado' ? '#065f46' : '#991b1b', fontWeight: '500' }}
            />
          </div>
        )}

        <div style={{ display: 'flex', gap: '0.75rem', justifyContent: 'flex-end' }}>
          <button type="button" onClick={handleCancel} style={styles.buttonSecondary}>
            Cancelar
          </button>
          <button type="submit" style={styles.buttonPrimary}>
            {mode === 'create' ? 'Crear Proyecto' : 'Guardar Cambios'}
          </button>
        </div>
      </form>
    </div>
  );
};


// Componente de Tabla (ACTUALIZADO)
const ListaProyectos = ({ proyectos, selectedProyectoId, setSelectedProyectoId, searchTerm }) => {

  const filteredProyectos = proyectos.filter(p =>
    p.nombre.toLowerCase().includes(searchTerm.toLowerCase())
  );

  // Formateador de fechas para mostrar en la tabla (DD/MM/YYYY)
  const formatDateForDisplay = (dateString) => {
    if (!dateString || dateString === "0001-01-01T00:00:00Z" || dateString.startsWith("0001-01-01")) return null;
    try {
      // Intenta parsear como AAAA-MM-DD primero (viene del input)
      let date;
      if (dateString.includes('-')) {
        date = new Date(dateString + 'T00:00:00'); // Añade hora para evitar problemas de zona horaria
      } else {
        date = new Date(dateString); // Intenta parsear formato ISO si viene del backend
      }

      if (isNaN(date.getTime())) { // Verifica si la fecha es válida
        return dateString; // Devuelve original si no es válida
      }

      const day = String(date.getDate()).padStart(2, '0');
      const month = String(date.getMonth() + 1).padStart(2, '0'); // Meses son 0-indexados
      const year = date.getFullYear();
      return `${day}/${month}/${year}`;
    } catch (error) {
      console.warn("Error formateando fecha:", dateString, error);
      return dateString; // Devuelve el string original si falla
    }
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
            <th style={styles.th}>Generar</th> {/* ⭐️ 2. NUEVO ENCABEZADO DE COLUMNA */}
          </tr>
        </thead>
        <tbody>
          {filteredProyectos.map((proyecto) => {
            const isSelected = proyecto.id === selectedProyectoId;
            return (
              <tr
                key={proyecto.id}
                style={isSelected ? styles.selectedRow : { cursor: 'pointer' }} // Añade cursor pointer
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
                {/* ⭐️ 3. NUEVA CELDA CON BOTÓN DE LUPA */}
                <td style={styles.td}>
                  <button
                    style={styles.tableIconButton}
                    onClick={(e) => {
                      e.stopPropagation(); // Evita que se seleccione la fila
                      alert(`Mostrar detalles del proyecto ID: ${proyecto.id}`);
                    }}
                    onMouseOver={(e) => {
                      e.currentTarget.style.backgroundColor = '#eef2ff'; // indigo-50
                      e.currentTarget.style.color = '#4338ca'; // indigo-700
                    }}
                    onMouseOut={(e) => {
                      e.currentTarget.style.backgroundColor = '#ffffff';
                      e.currentTarget.style.color = '#4f46e5'; // indigo-600
                    }}
                    title="Ver detalles"
                  >
                    <i className="fas fa-search"></i> {/* Icono de lupa */}
                  </button>
                </td>
              </tr>
            )
          })}
          {filteredProyectos.length === 0 && (
            <tr>
              {/* ⭐️ 4. COLSPAN ACTUALIZADO A 6 */}
              <td colSpan="6" style={{ ...styles.td, textAlign: 'center', color: '#6b7280', padding: '1.5rem' }}>
                {searchTerm ? 'No se encontraron proyectos con ese nombre.' : 'No hay proyectos creados.'}
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};


// --- COMPONENTE PRINCIPAL (ACTUALIZADO) ---
const PortafolioProyectos = ({ apiCall, currentUser }) => { // Recibe currentUser completo ahora
  const [proyectos, setProyectos] = useState([]);
  const [mode, setMode] = useState('view'); // 'view', 'create', 'edit'
  const [selectedProyectoId, setSelectedProyectoId] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    nombre: '',
    fecha_inicio: '',
    fecha_cierre: '',
  });

  const selectedProyecto = useMemo(() => {
    return proyectos.find(p => p.id === selectedProyectoId);
  }, [proyectos, selectedProyectoId]);

  const fetchProyectos = useCallback(async () => {
    setLoading(true);
    try {
      // Asegúrate de que currentUser exista y tenga la propiedad 'usuario'
      const username = currentUser; // Ajustado aquí
      if (!username) {
        throw new Error("Nombre de usuario no disponible para la solicitud.");
      }
      const result = await apiCall('admin/get-proyectos', { admin_username: username }, 'POST');
      if (result.success) {
        setProyectos(result.data.proyectos || []);
      } else {
        throw new Error(result.data.error || "Error desconocido al cargar proyectos");
      }
    } catch (error) {
      console.error("Error al cargar proyectos:", error);
      alert(`Error al cargar proyectos: ${error.message}`);
    } finally {
      setLoading(false);
    }
  }, [apiCall, currentUser]); // Dependencia ajustada

  useEffect(() => {
    fetchProyectos();
  }, [fetchProyectos]);

  // Formatear fechas para el formulario (YYYY-MM-DD)
  const formatDateForInput = (dateString) => {
    if (!dateString || dateString.startsWith("0001-01-01")) return '';
    try {
      // Intenta parsear como AAAA-MM-DD primero (viene del input)
      let date;
      if (dateString.includes('-') && !dateString.includes('T')) {
        date = new Date(dateString + 'T00:00:00'); // Añade hora para evitar problemas de zona horaria
      } else {
        date = new Date(dateString); // Intenta parsear formato ISO si viene del backend
      }

      if (isNaN(date.getTime())) { // Verifica si la fecha es válida
        console.warn("Fecha inválida recibida:", dateString);
        return ''; // Devuelve vacío si no es válida
      }
      return date.toISOString().split('T')[0];
    } catch (error) {
      console.error("Error formateando fecha para input:", dateString, error);
      return ''; // Devuelve vacío en caso de error
    }
  };


  const resetForm = () => {
    setFormData({ nombre: '', fecha_inicio: '', fecha_cierre: '' });
  };

  const handleCreateNew = () => {
    resetForm();
    setSelectedProyectoId(null);
    setMode('create');
  };

  const handleEditClick = () => {
    if (selectedProyecto) {
      if (selectedProyecto.estado === 'cerrado') {
        alert('No se puede modificar un proyecto que está "Cerrado".');
        return;
      }
      setFormData({
        nombre: selectedProyecto.nombre,
        fecha_inicio: formatDateForInput(selectedProyecto.fecha_inicio),
        fecha_cierre: formatDateForInput(selectedProyecto.fecha_cierre),
      });
      setMode('edit');
    } else {
      alert('Por favor, selecciona un proyecto para modificar.');
    }
  };

  const handleCancel = () => {
    resetForm();
    setMode('view');
  };

  const handleDelete = async () => {
    if (!selectedProyecto) {
      alert('Por favor, selecciona un proyecto para borrar.');
      return;
    }
    if (window.confirm(`¿Estás seguro de que quieres borrar el proyecto "${selectedProyecto.nombre}"? Esta acción no se puede deshacer.`)) {
      try {
        const username = currentUser; // Ajustado aquí
        if (!username) { throw new Error("Nombre de usuario no disponible."); }
        await apiCall('admin/delete-proyecto', {
          id: selectedProyecto.id,
          admin_username: username
        }, 'POST');
        alert('Proyecto borrado exitosamente.');
        fetchProyectos();
        setSelectedProyectoId(null);
        setMode('view');
      } catch (error) {
        console.error("Error al borrar proyecto:", error);
        alert(`Error al borrar proyecto: ${error.message}`);
      }
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const endpoint = mode === 'create' ? 'admin/create-proyecto' : 'admin/update-proyecto';
    const username = currentUser; // Ajustado aquí
    if (!username) { alert("Error: Nombre de usuario no disponible."); return; }
    const payload = {
      ...formData,
      admin_username: username,
    };
    if (mode === 'edit') {
      payload.id = selectedProyectoId;
    }

    // Validación simple de fecha de cierre
    if (formData.fecha_cierre && formData.fecha_inicio && formData.fecha_cierre < formData.fecha_inicio) {
      alert("La fecha de cierre no puede ser anterior a la fecha de inicio.");
      return;
    }

    try {
      const result = await apiCall(endpoint, payload, 'POST');
      alert(result.data.mensaje || 'Operación exitosa.');
      fetchProyectos();
      setMode('view');
      resetForm();
      if (mode === 'create') {
        setSelectedProyectoId(null);
      }
    } catch (error) {
      console.error("Error al guardar proyecto:", error);
      alert(`Error al guardar: ${error.message}`);
    }
  };

  const handleSetEstado = async (nuevoEstado) => {
    if (!selectedProyecto) {
      alert('Por favor, selecciona un proyecto primero.');
      return;
    }
    if (!window.confirm(`¿Estás seguro de que quieres cambiar el estado de "${selectedProyecto.nombre}" a "${nuevoEstado}"?`)) {
      return;
    }

    try {
      const endpoint = 'admin/set-proyecto-estado'; // Endpoint ajustado
      const username = currentUser; // Ajustado aquí
      if (!username) { throw new Error("Nombre de usuario no disponible."); }
      const payload = {
        id: selectedProyecto.id,
        estado: nuevoEstado,
        admin_username: username,
      };
      const result = await apiCall(endpoint, payload, 'POST'); // Método POST
      alert(result.data.mensaje || 'Estado actualizado con éxito.');
      fetchProyectos();
      if (nuevoEstado === 'cerrado' && mode === 'edit') {
        setMode('view');
      }
    } catch (error) {
      console.error('Error al cambiar estado:', error);
      alert(`Error al cambiar estado: ${error.message}`);
    }
  };


  return (
    <div style={{ padding: '2rem' }}>
      <h2 style={{ fontSize: '1.875rem', fontWeight: 'bold', color: '#111827', marginBottom: '1.5rem' }}>
        Gestión de Proyectos
      </h2>

      {mode !== 'view' ? (
        <ProyectoForm
          mode={mode}
          formData={formData}
          setFormData={setFormData}
          handleSubmit={handleSubmit}
          handleCancel={handleCancel}
          selectedProyecto={selectedProyecto}
        />
      ) : (
        <>
          <div style={styles.actionToolbar}>
            <button onClick={handleCreateNew} style={styles.buttonPrimary}>
              <i className="fas fa-plus" style={{ marginRight: '8px' }}></i>
              Crear Proyecto
            </button>
            <button
              onClick={handleEditClick}
              disabled={!selectedProyecto || selectedProyecto.estado === 'cerrado'}
              style={{
                ...styles.buttonSecondary,
                cursor: (!selectedProyecto || selectedProyecto.estado === 'cerrado') ? 'not-allowed' : 'pointer',
                opacity: (!selectedProyecto || selectedProyecto.estado === 'cerrado') ? 0.6 : 1
              }}
              title={
                !selectedProyecto
                  ? 'Selecciona un proyecto para modificar'
                  : selectedProyecto.estado === 'cerrado'
                    ? 'No se puede modificar un proyecto "Cerrado"'
                    : 'Modificar proyecto seleccionado'
              }
            >
              <i className="fas fa-pencil-alt" style={{ marginRight: '8px' }}></i>
              Modificar
            </button>

            <button
              onClick={handleDelete}
              disabled={!selectedProyecto}
              style={{
                ...styles.buttonDanger,
                cursor: !selectedProyecto ? 'not-allowed' : 'pointer',
                opacity: !selectedProyecto ? 0.6 : 1
              }}
              title={!selectedProyecto ? 'Selecciona un proyecto para borrar' : 'Borrar proyecto seleccionado'}
            >
              <i className="fas fa-trash-alt" style={{ marginRight: '8px' }}></i>
              Borrar
            </button>
            <button
              onClick={() => handleSetEstado('habilitado')}
              disabled={!selectedProyecto || selectedProyecto.estado === 'habilitado'}
              style={{
                ...styles.buttonSecondary,
                cursor: (!selectedProyecto || selectedProyecto.estado === 'habilitado') ? 'not-allowed' : 'pointer',
                opacity: (!selectedProyecto || selectedProyecto.estado === 'habilitado') ? 0.6 : 1,
                color: (!selectedProyecto || selectedProyecto.estado === 'habilitado') ? '#9ca3af' : '#10b981' // Verde
              }}
              title={
                !selectedProyecto
                  ? 'Selecciona un proyecto'
                  : selectedProyecto.estado === 'habilitado'
                    ? 'El proyecto ya está habilitado'
                    : 'Habilitar proyecto para edición'
              }
            >
              <i className="fas fa-check-circle" style={{ marginRight: '8px' }}></i>
              Habilitar
            </button>
            <button
              onClick={() => handleSetEstado('cerrado')}
              disabled={!selectedProyecto || selectedProyecto.estado === 'cerrado'}
              style={{
                ...styles.buttonWarning,
                cursor: (!selectedProyecto || selectedProyecto.estado === 'cerrado') ? 'not-allowed' : 'pointer',
                opacity: (!selectedProyecto || selectedProyecto.estado === 'cerrado') ? 0.6 : 1
              }}
              title={
                !selectedProyecto
                  ? 'Selecciona un proyecto'
                  : selectedProyecto.estado === 'cerrado'
                    ? 'El proyecto ya está cerrado'
                    : 'Cerrar proyecto para bloquear edición'
              }
            >
              <i className="fas fa-lock" style={{ marginRight: '8px' }}></i>
              Cerrar
            </button>
                        <button
              onClick={() => alert('Función de imprimir no implementada aún.')}
              style={{
                display: 'inline-flex',
                alignItems: 'center',
                justifyContent: 'center',
                padding: '0.5rem 1.5rem',
                fontSize: '1rem',
                fontWeight: '600',
                borderRadius: '8px',
                color: 'white',
                backgroundColor: '#4f46e5', // Color índigo (igual que el de login)
                border: 'none',
                cursor: 'pointer',
                transition: 'background-color 0.2s',
                boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
                marginTop: '0rem', // Añade un margen para separarlo
              }}
              onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#4338ca'} // Efecto hover
              onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#4f46e5'} // Vuelve al color original
            >
              
              Imprimir Proyecto
            </button>
            <input
              type="text"
              placeholder="Buscar por nombre..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              style={styles.searchInput}
            />
          </div>
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

export default PortafolioProyectos;