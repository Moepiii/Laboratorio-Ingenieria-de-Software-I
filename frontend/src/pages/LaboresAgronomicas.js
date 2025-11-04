import React, { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { getLabores, createLabor, updateLabor, deleteLabor } from '../services/laborService';

// Estilos (similares a los que ya usas en Portafolio.js y Usuarios.js)
const styles = {
    container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
    h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', marginBottom: '0.5rem' },
    p: { fontSize: '1rem', color: '#4b5563', marginBottom: '2rem' },
    formContainer: { padding: '1.5rem', backgroundColor: '#f9fafb', borderRadius: '8px', marginBottom: '2rem', border: '1px solid #e5e7eb' },
    h3: { fontSize: '1.25rem', fontWeight: '600', color: '#111827', marginTop: '0', marginBottom: '1rem' },
    input: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', boxSizing: 'border-box' },
    button: { padding: '0.75rem 1.5rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#4f46e5', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s', marginTop: '1rem' },
    table: { width: '100%', borderCollapse: 'collapse', marginTop: '1.5rem', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
    th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '2px solid #e5e7eb', color: '#374151', fontWeight: '600' },
    td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', verticalAlign: 'middle' },
    actionButton: { padding: '0.4rem 0.8rem', fontSize: '0.875rem', fontWeight: '500', borderRadius: '6px', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s', marginRight: '0.5rem' },
    editButton: { backgroundColor: '#f59e0b', color: 'white' },
    deleteButton: { backgroundColor: '#ef4444', color: 'white' },
    saveButton: { backgroundColor: '#22c55e', color: 'white' },
    cancelButton: { backgroundColor: '#6b7280', color: 'white' },
    errorText: { color: 'red', marginTop: '1rem' }
};

const LaboresAgronomicas = () => {
    // --- 1. Hooks y Estado ---
    const { id } = useParams(); // Obtiene el ID del proyecto de la URL
    const { token, currentUser } = useAuth(); // Obtiene el token y usuario del contexto

    const [labores, setLabores] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    // Estado para el formulario de crear
    const [newLaborDesc, setNewLaborDesc] = useState('');

    // Estado para la edición en línea
    const [editingId, setEditingId] = useState(null); // ID de la labor que se está editando
    const [editFormData, setEditFormData] = useState({ descripcion: '', estado: '' }); // Datos del formulario de edición

    // Variables de utilidad
    const adminUsername = currentUser?.username;
    const proyectoIdNum = parseInt(id, 10); // Convierte el ID de la URL a número

    // --- 2. Función para Cargar Datos ---
    const fetchLabores = useCallback(async () => {
        if (!token || !adminUsername || !proyectoIdNum) return;

        setLoading(true);
        setError('');
        try {
            // Llama al servicio que creamos en el paso 2.1
            const data = await getLabores(token, proyectoIdNum, adminUsername);
            setLabores(data.labores || []);
        } catch (err) {
            setError(err.message || 'Error al cargar labores.');
        } finally {
            setLoading(false);
        }
    }, [token, adminUsername, proyectoIdNum]);

    // Carga inicial de datos cuando el componente se monta
    useEffect(() => {
        fetchLabores();
    }, [fetchLabores]);

    // --- 3. Handlers (Manejadores de eventos) CRUD ---

    // CREAR una nueva labor
    const handleCreateLabor = async (e) => {
        e.preventDefault();
        if (newLaborDesc.trim() === '') {
            setError('La descripción no puede estar vacía.');
            return;
        }
        setError('');

        const laborData = {
            proyecto_id: proyectoIdNum,
            descripcion: newLaborDesc,
            estado: 'activa' // Estado por defecto al crear
        };

        try {
            // Llama al servicio
            const nuevaLabor = await createLabor(token, laborData, adminUsername);
            setLabores([nuevaLabor, ...labores]); // Añade la nueva labor al inicio de la lista
            setNewLaborDesc(''); // Limpia el input
        } catch (err) {
            setError(err.message || 'Error al crear la labor.');
        }
    };

    // BORRAR una labor
    const handleDeleteLabor = async (laborId) => {
        if (!window.confirm('¿Estás seguro de que quieres borrar esta labor?')) {
            return;
        }

        try {
            await deleteLabor(token, laborId, adminUsername);
            // Filtra la labor borrada del estado local
            setLabores(labores.filter(labor => labor.id !== laborId));
        } catch (err) {
            setError(err.message || 'Error al borrar la labor.');
        }
    };

    // --- 4. Handlers para Edición En Línea ---

    // Clic en "Editar": activa el modo edición para esa fila
    const handleEditClick = (labor) => {
        setEditingId(labor.id);
        setEditFormData({ descripcion: labor.descripcion, estado: labor.estado });
    };

    // Clic en "Cancelar": desactiva el modo edición
    const handleCancelClick = () => {
        setEditingId(null);
    };

    // Maneja el cambio en los inputs de edición
    const handleEditFormChange = (e) => {
        const { name, value } = e.target;
        setEditFormData({ ...editFormData, [name]: value });
    };

    // Clic en "Guardar": ACTUALIZA la labor
    const handleUpdateLabor = async (laborId) => {
        const laborData = {
            id: laborId,
            ...editFormData // { descripcion, estado }
        };

        try {
            await updateLabor(token, laborData, adminUsername);

            // Actualiza la lista local de labores con los nuevos datos
            const updatedLabores = labores.map(labor =>
                labor.id === laborId ? { ...labor, ...editFormData } : labor
            );
            setLabores(updatedLabores);
            setEditingId(null); // Desactiva el modo edición
        } catch (err) {
            setError(err.message || 'Error al actualizar la labor.');
        }
    };


    // --- 5. Renderizado del Componente ---
    return (
        <div style={styles.container}>
            <h2 style={styles.h2}>Gestión de Labores Agronómicas</h2>
            <p style={styles.p}>Administrando labores para el Proyecto ID: <strong>{id}</strong></p>

            {/* Formulario de Creación */}
            <div style={styles.formContainer}>
                <h3 style={styles.h3}>Nueva Labor</h3>
                <form onSubmit={handleCreateLabor}>
                    <div style={{ marginBottom: '1rem' }}>
                        <label htmlFor="descripcion" style={{ display: 'block', marginBottom: '0.5rem' }}>Descripción</label>
                        <input
                            id="descripcion"
                            type="text"
                            value={newLaborDesc}
                            onChange={(e) => setNewLaborDesc(e.target.value)}
                            style={styles.input}
                            placeholder="Ej: Preparación de suelo"
                        />
                    </div>
                    <button type="submit" style={styles.button}>
                        Crear Labor
                    </button>
                </form>
            </div>

            {/* Muestra de Errores */}
            {error && <p style={styles.errorText}>{error}</p>}

            {/* Tabla de Labores */}
            {loading ? (
                <p>Cargando labores...</p>
            ) : (
                <table style={styles.table}>
                    <thead>
                        <tr>
                            <th style={styles.th}>Descripción</th>
                            <th style={styles.th}>Estado</th>
                            <th style={styles.th}>Fecha Creación</th>
                            <th style={styles.th}>Acciones</th>
                        </tr>
                    </thead>
                    <tbody>
                        {labores.length === 0 ? (
                            <tr>
                                <td colSpan="4" style={{ ...styles.td, textAlign: 'center' }}>
                                    No hay labores registradas para este proyecto.
                                </td>
                            </tr>
                        ) : (
                            labores.map(labor => (
                                <tr key={labor.id}>
                                    {editingId === labor.id ? (
                                        // --- Fila en Modo Edición ---
                                        <>
                                            <td style={styles.td}>
                                                <input
                                                    type="text"
                                                    name="descripcion"
                                                    value={editFormData.descripcion}
                                                    onChange={handleEditFormChange}
                                                    style={styles.input}
                                                />
                                            </td>
                                            <td style={styles.td}>
                                                <select
                                                    name="estado"
                                                    value={editFormData.estado}
                                                    onChange={handleEditFormChange}
                                                    style={styles.input}
                                                >
                                                    <option value="activa">Activa</option>
                                                    <option value="completada">Completada</option>
                                                    <option value="archivada">Archivada</option>
                                                </select>
                                            </td>
                                            <td style={styles.td}>{new Date(labor.fecha_creacion).toLocaleDateString()}</td>
                                            <td style={styles.td}>
                                                <button style={{ ...styles.actionButton, ...styles.saveButton }} onClick={() => handleUpdateLabor(labor.id)}>Guardar</button>
                                                <button style={{ ...styles.actionButton, ...styles.cancelButton }} onClick={handleCancelClick}>Cancelar</button>
                                            </td>
                                        </>
                                    ) : (
                                        // --- Fila en Modo Lectura ---
                                        <>
                                            <td style={styles.td}>{labor.descripcion}</td>
                                            <td style={styles.td}>{labor.estado}</td>
                                            <td style={styles.td}>{new Date(labor.fecha_creacion).toLocaleDateString()}</td>
                                            <td style={styles.td}>
                                                <button style={{ ...styles.actionButton, ...styles.editButton }} onClick={() => handleEditClick(labor)}>Editar</button>
                                                <button style={{ ...styles.actionButton, ...styles.deleteButton }} onClick={() => handleDeleteLabor(labor.id)}>Borrar</button>
                                            </td>
                                        </>
                                    )}
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            )}
        </div>
    );
};

export default LaboresAgronomicas;