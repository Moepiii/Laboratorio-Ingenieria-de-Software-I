import React, { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { getEquipos, createEquipo, updateEquipo, deleteEquipo } from '../services/equipoService';

// Estilos (similares a los que ya usas)
const styles = {
    container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
    h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', marginBottom: '0.5rem' },
    p: { fontSize: '1rem', color: '#4b5563', marginBottom: '2rem' },
    formContainer: { padding: '1.5rem', backgroundColor: '#f9fafb', borderRadius: '8px', marginBottom: '2rem', border: '1px solid #e5e7eb' },
    h3: { fontSize: '1.25rem', fontWeight: '600', color: '#111827', marginTop: '0', marginBottom: '1rem' },
    input: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', boxSizing: 'border-box' },
    select: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', boxSizing: 'border-box', backgroundColor: 'white' },
    formGrid: { display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' },
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

const EquiposEImplementos = () => {
    // --- 1. Hooks y Estado ---
    const { id } = useParams(); // ID del proyecto de la URL
    const { token, currentUser } = useAuth(); // Token y usuario del contexto

    const [equipos, setEquipos] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    // Estado para el formulario de crear
    const [newEquipoNombre, setNewEquipoNombre] = useState('');
    const [newEquipoTipo, setNewEquipoTipo] = useState('implemento'); // Valor por defecto

    // Estado para la edición en línea
    const [editingId, setEditingId] = useState(null);
    const [editFormData, setEditFormData] = useState({ nombre: '', tipo: '', estado: '' });

    // Variables de utilidad
    const adminUsername = currentUser?.username;
    const proyectoIdNum = parseInt(id, 10);

    // --- 2. Función para Cargar Datos ---
    const fetchEquipos = useCallback(async () => {
        if (!token || !adminUsername || !proyectoIdNum) return;

        setLoading(true);
        setError('');
        try {
            // Llama al servicio que creamos en el paso 2.1
            const data = await getEquipos(token, proyectoIdNum, adminUsername);
            setEquipos(data.equipos || []);
        } catch (err) {
            setError(err.message || 'Error al cargar equipos.');
        } finally {
            setLoading(false);
        }
    }, [token, adminUsername, proyectoIdNum]);

    // Carga inicial de datos
    useEffect(() => {
        fetchEquipos();
    }, [fetchEquipos]);

    // --- 3. Handlers (Manejadores de eventos) CRUD ---

    // CREAR un nuevo equipo
    const handleCreateEquipo = async (e) => {
        e.preventDefault();
        if (newEquipoNombre.trim() === '') {
            setError('El nombre no puede estar vacío.');
            return;
        }
        setError('');

        const equipoData = {
            proyecto_id: proyectoIdNum,
            nombre: newEquipoNombre,
            tipo: newEquipoTipo,
            estado: 'disponible' // Estado por defecto al crear
        };

        try {
            const nuevoEquipo = await createEquipo(token, equipoData, adminUsername);
            setEquipos([nuevoEquipo, ...equipos]); // Añade al inicio de la lista
            setNewEquipoNombre(''); // Limpia el formulario
            setNewEquipoTipo('implemento');
        } catch (err) {
            setError(err.message || 'Error al crear el equipo.');
        }
    };

    // BORRAR un equipo
    const handleDeleteEquipo = async (equipoId) => {
        if (!window.confirm('¿Estás seguro de que quieres borrar este equipo?')) {
            return;
        }

        try {
            await deleteEquipo(token, equipoId, adminUsername);
            setEquipos(equipos.filter(equipo => equipo.id !== equipoId));
        } catch (err) {
            setError(err.message || 'Error al borrar el equipo.');
        }
    };

    // --- 4. Handlers para Edición En Línea ---

    // Clic en "Editar"
    const handleEditClick = (equipo) => {
        setEditingId(equipo.id);
        setEditFormData({ nombre: equipo.nombre, tipo: equipo.tipo, estado: equipo.estado });
    };

    // Clic en "Cancelar"
    const handleCancelClick = () => {
        setEditingId(null);
    };

    // Maneja el cambio en los inputs de edición
    const handleEditFormChange = (e) => {
        const { name, value } = e.target;
        setEditFormData({ ...editFormData, [name]: value });
    };

    // Clic en "Guardar": ACTUALIZA el equipo
    const handleUpdateEquipo = async (equipoId) => {
        const equipoData = {
            id: equipoId,
            ...editFormData // { nombre, tipo, estado }
        };

        try {
            await updateEquipo(token, equipoData, adminUsername);

            // Actualiza la lista local de equipos
            const updatedEquipos = equipos.map(equipo =>
                equipo.id === equipoId ? { ...equipo, ...editFormData } : equipo
            );
            setEquipos(updatedEquipos);
            setEditingId(null); // Desactiva el modo edición
        } catch (err) {
            setError(err.message || 'Error al actualizar el equipo.');
        }
    };


    // --- 5. Renderizado del Componente ---
    return (
        <div style={styles.container}>
            <h2 style={styles.h2}>Gestión de Equipos e Implementos</h2>
            <p style={styles.p}>Administrando equipos para el Proyecto ID: <strong>{id}</strong></p>

            {/* Formulario de Creación */}
            <div style={styles.formContainer}>
                <h3 style={styles.h3}>Nuevo Equipo o Implemento</h3>
                <form onSubmit={handleCreateEquipo}>
                    <div style={styles.formGrid}>
                        <div>
                            <label htmlFor="nombre" style={{ display: 'block', marginBottom: '0.5rem' }}>Nombre</label>
                            <input
                                id="nombre"
                                type="text"
                                value={newEquipoNombre}
                                onChange={(e) => setNewEquipoNombre(e.target.value)}
                                style={styles.input}
                                placeholder="Ej: Tractor John Deere"
                            />
                        </div>
                        <div>
                            <label htmlFor="tipo" style={{ display: 'block', marginBottom: '0.5rem' }}>Tipo</label>
                            <select
                                id="tipo"
                                value={newEquipoTipo}
                                onChange={(e) => setNewEquipoTipo(e.target.value)}
                                style={styles.select}
                            >
                                <option value="implemento">Implemento</option>
                                <option value="equipo">Equipo</option>
                            </select>
                        </div>
                    </div>
                    <button type="submit" style={styles.button}>
                        Crear Equipo
                    </button>
                </form>
            </div>

            {/* Muestra de Errores */}
            {error && <p style={styles.errorText}>{error}</p>}

            {/* Tabla de Equipos */}
            {loading ? (
                <p>Cargando equipos...</p>
            ) : (
                <table style={styles.table}>
                    <thead>
                        <tr>
                            <th style={styles.th}>Nombre</th>
                            <th style={styles.th}>Tipo</th>
                            <th style={styles.th}>Estado</th>
                            <th style={styles.th}>Fecha Creación</th>
                            <th style={styles.th}>Acciones</th>
                        </tr>
                    </thead>
                    <tbody>
                        {equipos.length === 0 ? (
                            <tr>
                                <td colSpan="5" style={{ ...styles.td, textAlign: 'center' }}>
                                    No hay equipos registrados para este proyecto.
                                </td>
                            </tr>
                        ) : (
                            equipos.map(equipo => (
                                <tr key={equipo.id}>
                                    {editingId === equipo.id ? (
                                        // --- Fila en Modo Edición ---
                                        <>
                                            <td style={styles.td}>
                                                <input
                                                    type="text"
                                                    name="nombre"
                                                    value={editFormData.nombre}
                                                    onChange={handleEditFormChange}
                                                    style={styles.input}
                                                />
                                            </td>
                                            <td style={styles.td}>
                                                <select
                                                    name="tipo"
                                                    value={editFormData.tipo}
                                                    onChange={handleEditFormChange}
                                                    style={styles.select}
                                                >
                                                    <option value="implemento">Implemento</option>
                                                    <option value="equipo">Equipo</option>
                                                </select>
                                            </td>
                                            <td style={styles.td}>
                                                <select
                                                    name="estado"
                                                    value={editFormData.estado}
                                                    onChange={handleEditFormChange}
                                                    style={styles.select}
                                                >
                                                    <option value="disponible">Disponible</option>
                                                    <option value="en uso">En Uso</option>
                                                    <option value="mantenimiento">Mantenimiento</option>
                                                </select>
                                            </td>
                                            <td style={styles.td}>{new Date(equipo.fecha_creacion).toLocaleDateString()}</td>
                                            <td style={styles.td}>
                                                <button style={{ ...styles.actionButton, ...styles.saveButton }} onClick={() => handleUpdateEquipo(equipo.id)}>Guardar</button>
                                                <button style={{ ...styles.actionButton, ...styles.cancelButton }} onClick={handleCancelClick}>Cancelar</button>
                                            </td>
                                        </>
                                    ) : (
                                        // --- Fila en Modo Lectura ---
                                        <>
                                            <td style={styles.td}>{equipo.nombre}</td>
                                            <td style={styles.td}>{equipo.tipo}</td>
                                            <td style={styles.td}>{equipo.estado}</td>
                                            <td style={styles.td}>{new Date(equipo.fecha_creacion).toLocaleDateString()}</td>
                                            <td style={styles.td}>
                                                <button style={{ ...styles.actionButton, ...styles.editButton }} onClick={() => handleEditClick(equipo)}>Editar</button>
                                                <button style={{ ...styles.actionButton, ...styles.deleteButton }} onClick={() => handleDeleteEquipo(equipo.id)}>Borrar</button>
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

export default EquiposEImplementos;