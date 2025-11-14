import React, { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
// Importamos el servicio de equipos que acabamos de modificar
import { getEquipos, createEquipo, updateEquipo, deleteEquipo } from '../services/equipoService';

// (Todos tus estilos 'styles' van aquí...)
const styles = {
    container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
    h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', marginBottom: '0.5rem' },
    p: { fontSize: '1rem', color: '#4b5563', marginBottom: '2rem' },
    formContainer: { padding: '1.5rem', backgroundColor: '#f9fafb', borderRadius: '8px', marginBottom: '2rem', border: '1px solid #e5e7eb' },
    h3: { fontSize: '1.25rem', fontWeight: '600', color: '#111827', marginTop: '0', marginBottom: '1rem' },
    input: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', boxSizing: 'border-box' },
    select: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', boxSizing: 'border-box', backgroundColor: 'white' },
    // El grid se adaptará solo a 3 columnas al quitar un elemento
    grid: { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1rem' },
    button: { padding: '0.75rem 1.5rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#4f46e5', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s', marginTop: '1rem' },
    tableContainer: { overflowX: 'auto', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
    table: { width: '100%', borderCollapse: 'collapse' },
    th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '2px solid #e5e7eb', color: '#374151', fontWeight: '600', fontSize: '0.875rem' },
    td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', color: '#111827', verticalAlign: 'middle' },
    error: { color: 'red', marginTop: '1rem', backgroundColor: '#fee2e2', padding: '1rem', borderRadius: '8px', border: '1px solid #f87171' },
    loading: { textAlign: 'center', padding: '2rem', fontSize: '1.2rem', color: '#6b7280' },
    actionButton: { padding: '0.4rem 0.8rem', fontSize: '0.875rem', fontWeight: '500', borderRadius: '6px', border: 'none', cursor: 'pointer', marginRight: '0.5rem', transition: 'background-color 0.2s' },
    editButton: { backgroundColor: '#3b82f6', color: 'white' },
    deleteButton: { backgroundColor: '#ef4444', color: 'white' },
    saveButton: { backgroundColor: '#10b981', color: 'white' },
    cancelButton: { backgroundColor: '#6b7280', color: 'white' }
};


const EquiposEImplementos = () => {
    const { id } = useParams(); // ID del proyecto
    const { token, currentUser } = useAuth();

    const [equipos, setEquipos] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    // ⭐️ --- (INICIO) CAMBIO EN EL ESTADO --- ⭐️
    // Estado para el formulario de nuevo equipo
    const [newEquipoData, setNewEquipoData] = useState({
        // 'codigo_equipo' eliminado
        nombre: '',
        tipo: 'Equipo',
        estado: 'Activo'
    });
    // ⭐️ --- (FIN) CAMBIO EN EL ESTADO --- ⭐️

    // Estado para la edición
    const [editingId, setEditingId] = useState(null);
    const [editFormData, setEditFormData] = useState({});

    // Cargar Equipos
    const fetchEquipos = useCallback(async () => {
        try {
            setLoading(true);
            const data = await getEquipos(token, parseInt(id), currentUser.username);
            setEquipos(data.equipos || []);
        } catch (err) {
            setError(err.message || "Error al cargar equipos");
        } finally {
            setLoading(false);
        }
    }, [id, token, currentUser]);

    useEffect(() => {
        fetchEquipos();
    }, [fetchEquipos]);

    // Manejador del formulario de nuevo equipo
    const handleNewDataChange = (e) => {
        const { name, value } = e.target;
        setNewEquipoData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    // ⭐️ --- (INICIO) CAMBIO EN EL GUARDADO --- ⭐️
    // Guardar nuevo equipo
    const handleSaveNewEquipo = async (e) => {
        e.preventDefault();
        setError('');
        try {
            // El servicio 'createEquipo' (archivo anterior) ya espera
            // que 'newEquipoData' NO tenga 'codigo_equipo'.
            const nuevoEquipo = await createEquipo(
                token,
                { proyecto_id: parseInt(id), ...newEquipoData },
                currentUser.username
            );
            setEquipos(prev => [nuevoEquipo, ...prev]);

            // Limpiar formulario (ya no se resetea el código)
            setNewEquipoData({
                nombre: '',
                tipo: 'Equipo', // Resetear al valor por defecto
                estado: 'Activo'
            });
        } catch (err) {
            setError(err.message || "Error al guardar el equipo");
        }
    };
    // ⭐️ --- (FIN) CAMBIO EN EL GUARDADO --- ⭐️

    // Manejadores de edición (sin cambios)
    const handleEditClick = (equipo) => {
        setEditingId(equipo.id);
        setEditFormData({
            codigo_equipo: equipo.codigo_equipo,
            nombre: equipo.nombre,
            tipo: equipo.tipo,
            estado: equipo.estado
        });
    };

    const handleEditDataChange = (e) => {
        const { name, value } = e.target;
        setEditFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleCancelClick = () => {
        setEditingId(null);
    };

    const handleUpdateEquipo = async (equipoId) => {
        setError('');
        try {
            await updateEquipo(
                token,
                { id: equipoId, ...editFormData },
                currentUser.username
            );
            setEditingId(null);
            fetchEquipos(); // Recargar la lista
        } catch (err) {
            setError(err.message || "Error al actualizar el equipo");
        }
    };

    // Borrar equipo (sin cambios)
    const handleDeleteEquipo = async (equipoId) => {
        if (window.confirm('¿Estás seguro de que quieres borrar este equipo?')) {
            setError('');
            try {
                await deleteEquipo(token, equipoId, currentUser.username);
                fetchEquipos(); // Recargar la lista
            } catch (err) {
                setError(err.message || "Error al borrar el equipo");
            }
        }
    };

    return (
        <div style={styles.container}>
            <h2 style={styles.h2}>Equipos e Implementos</h2>
            <p style={styles.p}>Configuración de equipos e implementos para el Proyecto (ID: {id})</p>

            {/* ⭐️ --- (INICIO) CAMBIO EN EL FORMULARIO --- ⭐️ */}
            <form onSubmit={handleSaveNewEquipo} style={styles.formContainer}>
                <h3 style={styles.h3}>Añadir Nuevo Equipo o Implemento</h3>
                <div style={styles.grid}>

                    {/* Campo "Codigo" ELIMINADO */}

                    {/* Nombre */}
                    <input
                        name="nombre"
                        value={newEquipoData.nombre}
                        onChange={handleNewDataChange}
                        placeholder="Nombre descriptivo"
                        required
                        style={styles.input}
                    />
                    {/* Tipo */}
                    <select
                        name="tipo"
                        value={newEquipoData.tipo}
                        onChange={handleNewDataChange}
                        required
                        style={styles.select}
                    >
                        <option value="Equipo">Equipo</option>
                        <option value="Implemento">Implemento</option>
                    </select>
                    {/* Estado */}
                    <select
                        name="estado"
                        value={newEquipoData.estado}
                        onChange={handleNewDataChange}
                        required
                        style={styles.select}
                    >
                        <option value="Activo">Activo</option>
                        <option value="Inactivo">Inactivo</option>
                    </select>
                </div>
                <button type="submit" style={styles.button}>Guardar Nuevo</button>
            </form>
            {/* ⭐️ --- (FIN) CAMBIO EN EL FORMULARIO --- ⭐️ */}

            {error && <p style={styles.error}>{error}</p>}

            {/* --- Tabla de Equipos --- */}
            {loading ? (
                <p style={styles.loading}>Cargando...</p>
            ) : (
                <table style={{ ...styles.table, marginTop: '2rem' }}>
                    <thead>
                        <tr>
                            <th style={styles.th}>Código</th>
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
                                <td colSpan="6" style={{ ...styles.td, textAlign: 'center', color: '#6b7280' }}>
                                    No hay equipos o implementos definidos para este proyecto.
                                </td>
                            </tr>
                        ) : (
                            equipos.map(equipo => (
                                <tr key={equipo.id}>
                                    {editingId === equipo.id ? (
                                        // --- Fila en Modo Edición ---
                                        <>
                                            <td style={styles.td}><input name="codigo_equipo" value={editFormData.codigo_equipo} onChange={handleEditDataChange} style={styles.input} /></td>
                                            <td style={styles.td}><input name="nombre" value={editFormData.nombre} onChange={handleEditDataChange} style={styles.input} /></td>
                                            <td style={styles.td}>
                                                <select name="tipo" value={editFormData.tipo} onChange={handleEditDataChange} style={styles.select}>
                                                    <option value="Equipo">Equipo</option>
                                                    <option value="Implemento">Implemento</option>
                                                </select>
                                            </td>
                                            <td style={styles.td}>
                                                <select name="estado" value={editFormData.estado} onChange={handleEditDataChange} style={styles.select}>
                                                    <option value="Activo">Activo</option>
                                                    <option value="Inactivo">Inactivo</option>
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
                                            {/* La tabla SÍ muestra el código que generó el backend */}
                                            <td style={styles.td}>{equipo.codigo_equipo}</td>
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