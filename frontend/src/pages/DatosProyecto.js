import React, { useState, useEffect, useCallback } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
// ⭐️ 1. Importa el servicio de actividades que creamos
import {
    getDatosProyecto,
    createActividad,
    updateActividad,
    deleteActividad
} from '../services/actividadService';
// ⭐️ 2. Importa el Modal que acabamos de crear
import Modal from '../components/auth/Modal';

// ⭐️ 3. Todos los estilos están definidos aquí (sin .css)
const styles = {
    // Estilos de la página principal
    container: { padding: '2rem', fontFamily: 'Inter, sans-serif' },
    header: {
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        borderBottom: '2px solid #e5e7eb',
        paddingBottom: '1rem',
        marginBottom: '2rem'
    },
    h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', margin: 0 },
    backButton: {
        padding: '0.6rem 1.2rem',
        fontSize: '1rem',
        fontWeight: '600',
        borderRadius: '8px',
        color: 'white',
        backgroundColor: '#6b7280',
        border: 'none',
        cursor: 'pointer',
        textDecoration: 'none',
    },
    // Estilos de la sección de la tabla
    section: {
        padding: '1.5rem',
        backgroundColor: 'white',
        borderRadius: '8px',
        border: '1px solid #e5e7eb',
        marginBottom: '1.5rem',
        boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)'
    },
    sectionHeader: {
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: '1.5rem'
    },
    sectionTitle: {
        fontSize: '1.25rem',
        fontWeight: '600',
        color: '#111827',
        margin: 0,
    },
    addButton: {
        padding: '0.6rem 1.2rem',
        fontSize: '1rem',
        fontWeight: '600',
        borderRadius: '8px',
        color: 'white',
        backgroundColor: '#4f46e5',
        border: 'none',
        cursor: 'pointer',
    },
    // Estilos de la tabla
    tableWrapper: { overflowX: 'auto' },
    table: { width: '100%', borderCollapse: 'collapse', minWidth: '900px' },
    th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '2px solid #e5e7eb', color: '#374151', fontWeight: '600' },
    td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', verticalAlign: 'middle', whiteSpace: 'nowrap' },
    actionButton: { background: 'none', border: 'none', cursor: 'pointer', fontSize: '1.2rem', padding: '0.5rem' },
    // Estilos del Modal
    formGrid: {
        display: 'grid',
        gridTemplateColumns: '1fr 1fr',
        gap: '1rem',
        marginBottom: '1.5rem'
    },
    formGridTriple: { // Estilo para 3 columnas
        display: 'grid',
        gridTemplateColumns: '1fr 1fr 1fr',
        gap: '1rem',
        marginBottom: '1.5rem'
    },
    formGroup: { display: 'flex', flexDirection: 'column' },
    label: { display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' },
    input: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', boxSizing: 'border-box' },
    select: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', boxSizing: 'border-box', backgroundColor: 'white' },
    textarea: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', boxSizing: 'border-box', minHeight: '80px', fontFamily: 'inherit' },
    formActions: { display: 'flex', justifyContent: 'flex-end', gap: '1rem', marginTop: '1.5rem' },
    saveButton: { padding: '0.6rem 1.2rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#22c55e', border: 'none', cursor: 'pointer' },
    cancelButton: { padding: '0.6rem 1.2rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: '#374151', backgroundColor: '#f3f4f6', border: '1px solid #d1d5db', cursor: 'pointer' },
    errorText: { color: 'red', marginTop: '1rem' }
};


const DatosProyecto = () => {
    const { id } = useParams(); // ID del proyecto desde la URL
    const { token, currentUser } = useAuth();
    const adminUsername = currentUser?.username;
    const proyectoIdNum = parseInt(id, 10);

    // Estado para los datos de la página
    const [actividades, setActividades] = useState([]);
    const [labores, setLabores] = useState([]);
    const [equipos, setEquipos] = useState([]);
    const [encargados, setEncargados] = useState([]);

    // Estado de UI
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [isModalOpen, setIsModalOpen] = useState(false);

    // Estado del formulario del modal
    const [currentActividad, setCurrentActividad] = useState(null); // Si es != null, estamos editando
    const [formData, setFormData] = useState({
        actividad: '',
        labor_agronomica_id: 0,
        equipo_implemento_id: 0,
        encargado_id: 0,
        recurso_humano: 0,
        costo: 0,
        observaciones: ''
    });

    // Carga de datos inicial
    const loadPageData = useCallback(async () => {
        if (!token || !adminUsername || !proyectoIdNum) return;
        setLoading(true);
        setError('');
        try {
            const data = await getDatosProyecto(token, proyectoIdNum, adminUsername);
            setActividades(data.actividades || []);
            setLabores(data.labores || []);
            setEquipos(data.equipos || []);
            setEncargados(data.encargados || []);
        } catch (err) {
            setError(err.message || 'Error al cargar los datos del proyecto.');
        } finally {
            setLoading(false);
        }
    }, [token, adminUsername, proyectoIdNum]);

    useEffect(() => {
        loadPageData();
    }, [loadPageData]);

    // --- Manejo del Modal ---
    const handleOpenModal = () => {
        setCurrentActividad(null); // Limpia para "Crear"
        setFormData({
            actividad: '',
            labor_agronomica_id: 0,
            equipo_implemento_id: 0,
            encargado_id: 0,
            recurso_humano: 0,
            costo: 0,
            observaciones: ''
        });
        setIsModalOpen(true);
        setError('');
    };

    const handleOpenEditModal = (actividad) => {
        setCurrentActividad(actividad); // Guarda la actividad para "Editar"
        setFormData({
            actividad: actividad.actividad,
            labor_agronomica_id: actividad.labor_agronomica_id.Valid ? actividad.labor_agronomica_id.Int64 : 0,
            equipo_implemento_id: actividad.equipo_implemento_id.Valid ? actividad.equipo_implemento_id.Int64 : 0,
            encargado_id: actividad.encargado_id.Valid ? actividad.encargado_id.Int64 : 0,
            recurso_humano: actividad.recurso_humano,
            costo: actividad.costo,
            observaciones: actividad.observaciones.Valid ? actividad.observaciones.String : ''
        });
        setIsModalOpen(true);
        setError('');
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
        setCurrentActividad(null);
        setError('');
    };

    const handleFormChange = (e) => {
        const { name, value } = e.target;
        // Convierte a número si es necesario
        const val = (name === 'labor_agronomica_id' || name === 'equipo_implemento_id' || name === 'encargado_id' || name === 'recurso_humano' || name === 'costo')
            ? Number(value)
            : value;

        setFormData(prev => ({ ...prev, [name]: val }));
    };

    // --- Funciones CRUD (Llamadas a la API) ---

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        // Prepara los datos para enviar (convierte 0 a null para IDs opcionales)
        const dataToSend = {
            ...formData,
            proyecto_id: proyectoIdNum,
            labor_agronomica_id: formData.labor_agronomica_id > 0 ? formData.labor_agronomica_id : null,
            equipo_implemento_id: formData.equipo_implemento_id > 0 ? formData.equipo_implemento_id : null,
            encargado_id: formData.encargado_id > 0 ? formData.encargado_id : null,
        };

        try {
            let response;
            if (currentActividad) {
                // --- Actualizar ---
                response = await updateActividad(token, { ...dataToSend, id: currentActividad.id }, adminUsername);
            } else {
                // --- Crear ---
                response = await createActividad(token, dataToSend, adminUsername);
            }
            // El backend devuelve la lista actualizada de actividades
            setActividades(response.actividades || []);
            handleCloseModal();
        } catch (err) {
            setError(err.message || 'Error al guardar la actividad.');
        }
    };

    const handleDelete = async (actividadId) => {
        if (!window.confirm('¿Estás seguro de que quieres borrar esta actividad?')) return;
        setError('');
        try {
            await deleteActividad(token, actividadId, adminUsername);
            // Recarga todos los datos (o solo filtra localmente)
            setActividades(prev => prev.filter(act => act.id !== actividadId));
        } catch (err) {
            setError(err.message || 'Error al borrar la actividad.');
        }
    };

    return (
        <div style={styles.container}>
            <div style={styles.header}>
                <h2 style={styles.h2}>Datos del Proyecto (ID: {id})</h2>
                <Link to="/admin/proyectos" style={styles.backButton}>
                    &larr; Volver al Portafolio
                </Link>
            </div>

            {loading && <p>Cargando datos...</p>}
            {error && !isModalOpen && <p style={styles.errorText}>{error}</p>}

            {/* Sección de Actividades */}
            <section style={styles.section}>
                <div style={styles.sectionHeader}>
                    <h3 style={styles.sectionTitle}>Registro de Actividades</h3>
                    <button style={styles.addButton} onClick={handleOpenModal}>
                        + Nueva Actividad
                    </button>
                </div>
                <div style={styles.tableWrapper}>
                    <table style={styles.table}>
                        <thead>
                            <tr>
                                <th style={styles.th}>ID</th>
                                <th style={styles.th}>Actividad</th>
                                <th style={styles.th}>Labor</th>
                                <th style={styles.th}>Equipo</th>
                                <th style={styles.th}>Encargado</th>
                                <th style={styles.th}>Rec. Humano</th>
                                <th style={styles.th}>Costo</th>
                                <th style={styles.th}>Observaciones</th>
                                <th style={styles.th}>Acciones</th>
                            </tr>
                        </thead>
                        <tbody>
                            {actividades.length === 0 ? (
                                <tr>
                                    <td colSpan="9" style={{ ...styles.td, textAlign: 'center' }}>No hay actividades registradas.</td>
                                </tr>
                            ) : (
                                actividades.map(act => (
                                    <tr key={act.id}>
                                        <td style={styles.td}>{act.id}</td>
                                        <td style={styles.td}>{act.actividad}</td>
                                        <td style={styles.td}>{act.labor_descripcion.Valid ? act.labor_descripcion.String : 'N/A'}</td>
                                        <td style={styles.td}>{act.equipo_nombre.Valid ? act.equipo_nombre.String : 'N/A'}</td>
                                        <td style={styles.td}>{act.encargado_nombre.Valid ? act.encargado_nombre.String : 'N/A'}</td>
                                        <td style={styles.td}>{act.recurso_humano}</td>
                                        <td style={styles.td}>{act.costo.toFixed(2)}</td>
                                        <td style={styles.td}>{act.observaciones.Valid ? act.observaciones.String : ''}</td>
                                        <td style={styles.td}>
                                            <button onClick={() => handleOpenEditModal(act)} style={{ ...styles.actionButton, color: '#f59e0b' }}>✏️</button>
                                            <button onClick={() => handleDelete(act.id)} style={{ ...styles.actionButton, color: '#ef4444' }}>❌</button>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </section>

            {/* Modal para Crear/Editar Actividad */}
            <Modal isOpen={isModalOpen} onClose={handleCloseModal} title={currentActividad ? 'Editar Actividad' : 'Nueva Actividad'}>
                <form onSubmit={handleSubmit}>

                    {/* Fila 1 */}
                    <div style={styles.formGroup}>
                        <label htmlFor="actividad" style={styles.label}>Nombre de la Actividad</label>
                        <input
                            type="text"
                            name="actividad"
                            id="actividad"
                            value={formData.actividad}
                            onChange={handleFormChange}
                            style={styles.input}
                            required
                        />
                    </div>

                    {/* Fila 2 */}
                    <div style={styles.formGrid}>
                        <div style={styles.formGroup}>
                            <label htmlFor="labor_agronomica_id" style={styles.label}>Labor Agronómica</label>
                            <select
                                name="labor_agronomica_id"
                                id="labor_agronomica_id"
                                value={formData.labor_agronomica_id}
                                onChange={handleFormChange}
                                style={styles.select}
                            >
                                <option value="0">--- Ninguna ---</option>
                                {labores.map(lab => (
                                    <option key={lab.id} value={lab.id}>{lab.codigo_labor} - {lab.descripcion}</option>
                                ))}
                            </select>
                        </div>
                        <div style={styles.formGroup}>
                            <label htmlFor="equipo_implemento_id" style={styles.label}>Equipo / Implemento</label>
                            <select
                                name="equipo_implemento_id"
                                id="equipo_implemento_id"
                                value={formData.equipo_implemento_id}
                                onChange={handleFormChange}
                                style={styles.select}
                            >
                                <option value="0">--- Ninguno ---</option>
                                {equipos.map(eq => (
                                    <option key={eq.id} value={eq.id}>{eq.codigo_equipo} - {eq.nombre}</option>
                                ))}
                            </select>
                        </div>
                    </div>

                    {/* Fila 3 */}
                    <div style={styles.formGridTriple}> {/* ⭐️ Estilo de 3 columnas */}
                        <div style={styles.formGroup}>
                            <label htmlFor="encargado_id" style={styles.label}>Encargado</label>
                            <select
                                name="encargado_id"
                                id="encargado_id"
                                value={formData.encargado_id}
                                onChange={handleFormChange}
                                style={styles.select}
                            >
                                <option value="0">--- Ninguno ---</option>
                                {encargados.map(enc => (
                                    <option key={enc.id} value={enc.id}>{enc.nombre} {enc.apellido}</option>
                                ))}
                            </select>
                        </div>
                        <div style={styles.formGroup}>
                            <label htmlFor="recurso_humano" style={styles.label}>Recurso Humano (Nro.)</label>
                            <input
                                type="number"
                                name="recurso_humano"
                                id="recurso_humano"
                                value={formData.recurso_humano}
                                onChange={handleFormChange}
                                style={styles.input}
                                min="0"
                            />
                        </div>
                        <div style={styles.formGroup}>
                            <label htmlFor="costo" style={styles.label}>Costo ($)</label>
                            <input
                                type="number"
                                name="costo"
                                id="costo"
                                value={formData.costo}
                                onChange={handleFormChange}
                                style={styles.input}
                                min="0"
                                step="0.01"
                            />
                        </div>
                    </div>

                    {/* Fila 4 */}
                    <div style={styles.formGroup}>
                        <label htmlFor="observaciones" style={styles.label}>Observaciones</label>
                        <textarea
                            name="observaciones"
                            id="observaciones"
                            value={formData.observaciones}
                            onChange={handleFormChange}
                            style={styles.textarea}
                        />
                    </div>

                    {error && <p style={styles.errorText}>{error}</p>}

                    <div style={styles.formActions}>
                        <button type="button" style={styles.cancelButton} onClick={handleCloseModal}>Cancelar</button>
                        <button type="submit" style={styles.saveButton}>{currentActividad ? 'Actualizar' : 'Guardar'}</button>
                    </div>
                </form>
            </Modal>
        </div>
    );
};

export default DatosProyecto;