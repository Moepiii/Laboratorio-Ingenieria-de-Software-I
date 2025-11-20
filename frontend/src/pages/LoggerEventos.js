import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../context/AuthContext';
import { getLogs, deleteLogs } from '../services/loggerService';

const styles = {
    container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
    // Encabezado flex para poner el botón a la derecha
    headerRow: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem', borderBottom: '2px solid #e5e7eb', paddingBottom: '0.75rem' },
    h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', margin: 0 },
    
    filterContainer: {
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
        gap: '1rem',
        padding: '1.5rem',
        backgroundColor: '#f9fafb',
        borderRadius: '8px',
        marginBottom: '2rem',
        border: '1px solid #e5e7eb'
    },
    filterGroup: { display: 'flex', flexDirection: 'column' },
    label: { fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' },
    input: { padding: '0.6rem', border: '1px solid #d1d5db', borderRadius: '6px', fontSize: '1rem' },
    select: { padding: '0.6rem', border: '1px solid #d1d5db', borderRadius: '6px', fontSize: '1rem', backgroundColor: 'white' },
    
    // Botones
    buttonGroup: { display: 'flex', gap: '1rem', alignItems: 'flex-end' },
    searchButton: { padding: '0.6rem 1.2rem', border: 'none', borderRadius: '6px', backgroundColor: '#2563eb', color: 'white', cursor: 'pointer', fontWeight: '600' },
    clearButton: { padding: '0.6rem 1.2rem', border: '1px solid #d1d5db', borderRadius: '6px', backgroundColor: 'white', color: '#374151', cursor: 'pointer', fontWeight: '600' },
    
    // Botón de Eliminar (Rojo)
    deleteButton: { 
        padding: '0.6rem 1.2rem', 
        border: 'none', 
        borderRadius: '6px', 
        backgroundColor: '#ef4444', 
        color: 'white', 
        cursor: 'pointer', 
        fontWeight: '600',
        display: 'flex',
        alignItems: 'center',
        gap: '0.5rem'
    },

    tableContainer: { overflowX: 'auto', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
    table: { width: '100%', borderCollapse: 'collapse' },
    th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '2px solid #e5e7eb', color: '#374151', fontWeight: '600', fontSize: '0.875rem' },
    td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', color: '#4b5563', fontSize: '0.875rem' },
    noResults: { textAlign: 'center', padding: '2rem', color: '#6b7280' },
    
    // Checkbox styles
    checkbox: { transform: 'scale(1.2)', cursor: 'pointer' }
};

const LoggerEventos = () => {
    const { token, currentUser } = useAuth();
    const [logs, setLogs] = useState([]);
    const [selectedIds, setSelectedIds] = useState([]); // ⭐️ Estado para selección

    // Filtros
    const [filters, setFilters] = useState({
        fecha_inicio: '',
        fecha_cierre: '',
        usuario_username: '',
        accion: '',
        entidad: ''
    });

    const fetchLogs = useCallback(async () => {
        try {
            const data = await getLogs(token, currentUser.username, filters);
            setLogs(data || []);
            setSelectedIds([]); // Limpiar selección al buscar de nuevo
        } catch (error) {
            console.error("Error obteniendo logs:", error);
        }
    }, [token, currentUser.username, filters]);

    // Cargar al inicio
    useEffect(() => {
        // Una carga inicial sin filtros estrictos o con los defaults
        getLogs(token, currentUser.username, {}).then(data => setLogs(data || []));
        // eslint-disable-next-line
    }, []); 

    const handleFilterChange = (e) => {
        setFilters({ ...filters, [e.target.name]: e.target.value });
    };

    const handleSearch = (e) => {
        e.preventDefault();
        fetchLogs();
    };

    const handleClearFilters = () => {
        setFilters({ fecha_inicio: '', fecha_cierre: '', usuario_username: '', accion: '', entidad: '' });
        // Opcional: recargar logs sin filtros automáticamente
        getLogs(token, currentUser.username, {}).then(data => setLogs(data || []));
    };

    // ⭐️ MANEJO DE SELECCIÓN ⭐️
    const handleSelectAll = (e) => {
        if (e.target.checked) {
            // Seleccionar todos los IDs visibles
            const allIds = logs.map(log => log.id);
            setSelectedIds(allIds);
        } else {
            // Deseleccionar todo
            setSelectedIds([]);
        }
    };

    const handleSelectOne = (id) => {
        if (selectedIds.includes(id)) {
            setSelectedIds(selectedIds.filter(itemId => itemId !== id));
        } else {
            setSelectedIds([...selectedIds, id]);
        }
    };

    // ⭐️ ELIMINAR SELECCIONADOS ⭐️
    const handleDeleteSelected = async () => {
        if (selectedIds.length === 0) return;

        const confirmMessage = `¿Estás seguro de eliminar ${selectedIds.length} evento(s)? Esta acción no se puede deshacer.`;
        if (window.confirm(confirmMessage)) {
            try {
                await deleteLogs(token, selectedIds, currentUser.username);
                alert("Eventos eliminados correctamente.");
                fetchLogs(); // Recargar tabla
            } catch (error) {
                alert("Error: " + error.message);
            }
        }
    };

    return (
        <div style={styles.container}>
            <div style={styles.headerRow}>
                <h2 style={styles.h2}>Logger de Eventos</h2>
                
                {/* ⭐️ Botón de Eliminar (Solo visible si hay selección) */}
                {selectedIds.length > 0 && (
                    <button style={styles.deleteButton} onClick={handleDeleteSelected}>
                        🗑️ Eliminar Seleccionados ({selectedIds.length})
                    </button>
                )}
            </div>

            {/* Filtros (Igual que antes) */}
            <form style={styles.filterContainer} onSubmit={handleSearch}>
                <div style={styles.filterGroup}>
                    <label style={styles.label}>Desde</label>
                    <input type="date" name="fecha_inicio" value={filters.fecha_inicio} onChange={handleFilterChange} style={styles.input} />
                </div>
                <div style={styles.filterGroup}>
                    <label style={styles.label}>Hasta</label>
                    <input type="date" name="fecha_cierre" value={filters.fecha_cierre} onChange={handleFilterChange} style={styles.input} />
                </div>
                <div style={styles.filterGroup}>
                    <label style={styles.label}>Usuario</label>
                    <input type="text" name="usuario_username" placeholder="Username..." value={filters.usuario_username} onChange={handleFilterChange} style={styles.input} />
                </div>
                <div style={styles.filterGroup}>
                    <label style={styles.label}>Acción</label>
                    <select name="accion" value={filters.accion} onChange={handleFilterChange} style={styles.select}>
                        <option value="">Todas</option>
                        <option value="CREACIÓN">Creación</option>
                        <option value="MODIFICACIÓN">Modificación</option>
                        <option value="ELIMINACIÓN">Eliminación</option>
                        <option value="REGISTRO">Registro</option>
                        <option value="LOGIN">Login</option>
                    </select>
                </div>
                <div style={{ ...styles.filterGroup, justifyContent: 'flex-end' }}>
                    <div style={styles.buttonGroup}>
                        <button type="button" onClick={handleClearFilters} style={styles.clearButton}>Limpiar</button>
                        <button type="submit" style={styles.searchButton}>Buscar</button>
                    </div>
                </div>
            </form>

            {/* Tabla */}
            <div style={styles.tableContainer}>
                <table style={styles.table}>
                    <thead>
                        <tr>
                            {/* ⭐️ Checkbox Maestro */}
                            <th style={styles.th}>
                                <input 
                                    type="checkbox" 
                                    style={styles.checkbox}
                                    onChange={handleSelectAll}
                                    checked={logs.length > 0 && selectedIds.length === logs.length}
                                />
                            </th>
                            <th style={styles.th}>ID</th>
                            <th style={styles.th}>Fecha/Hora</th>
                            <th style={styles.th}>Usuario</th>
                            <th style={styles.th}>Rol</th>
                            <th style={styles.th}>Acción</th>
                            <th style={styles.th}>Entidad</th>
                            <th style={styles.th}>Entidad ID</th>
                        </tr>
                    </thead>
                    <tbody>
                        {logs.length > 0 ? (
                            logs.map(log => (
                                <tr key={log.id}>
                                    {/* ⭐️ Checkbox Individual */}
                                    <td style={styles.td}>
                                        <input 
                                            type="checkbox" 
                                            style={styles.checkbox}
                                            checked={selectedIds.includes(log.id)}
                                            onChange={() => handleSelectOne(log.id)}
                                        />
                                    </td>
                                    <td style={styles.td}>{log.id}</td>
                                    <td style={styles.td}>{new Date(log.timestamp).toLocaleString()}</td>
                                    <td style={styles.td}>{log.usuario_username}</td>
                                    <td style={styles.td}>{log.usuario_rol}</td>
                                    <td style={styles.td}>{log.accion}</td>
                                    <td style={styles.td}>{log.entidad}</td>
                                    <td style={styles.td}>{log.entidad_id}</td>
                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td colSpan="8" style={styles.noResults}>No se encontraron eventos.</td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default LoggerEventos;