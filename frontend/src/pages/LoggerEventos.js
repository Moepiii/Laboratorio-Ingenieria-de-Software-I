import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../context/AuthContext';
import { getLogs } from '../services/loggerService'; // 1. Importamos el servicio

// Estilos (los mismos que usamos en otras páginas)
const styles = {
    container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
    h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', marginBottom: '1.5rem', borderBottom: '2px solid #e5e7eb', paddingBottom: '0.75rem' },
    
    // Estilos para los filtros (basado en Screenshot_10.jpg)
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
    filterGroup: {
        display: 'flex',
        flexDirection: 'column',
    },
    label: { 
        fontSize: '0.875rem', 
        fontWeight: '500', 
        color: '#374151', 
        marginBottom: '0.5rem' 
    },
    input: { 
        width: '100%', 
        padding: '0.75rem 1rem', 
        border: '1px solid #d1d5db', 
        borderRadius: '8px', 
        fontSize: '1rem', 
        boxSizing: 'border-box' 
    },
    select: { 
        width: '100%', 
        padding: '0.75rem 1rem', 
        border: '1px solid #d1d5db', 
        borderRadius: '8px', 
        fontSize: '1rem', 
        boxSizing: 'border-box',
        backgroundColor: 'white'
    },
    filterActions: {
        gridColumn: '1 / -1', // Ocupa todo el ancho
        display: 'flex',
        justifyContent: 'flex-end',
        gap: '1rem',
        marginTop: '1rem'
    },
    button: { 
        padding: '0.75rem 1.5rem', 
        fontSize: '1rem', 
        fontWeight: '600', 
        borderRadius: '8px', 
        color: 'white', 
        border: 'none', 
        cursor: 'pointer' 
    },

    // Estilos para la tabla (basado en Screenshot_9.jpg)
    tableContainer: { 
        overflowX: 'auto', 
        borderRadius: '8px', 
        border: '1px solid #e5e7eb', 
        boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' 
    },
    table: { width: '100%', borderCollapse: 'collapse' },
    th: { 
        padding: '0.75rem 1rem', 
        textAlign: 'left', 
        backgroundColor: '#f3f4f6', 
        borderBottom: '2px solid #e5e7eb', 
        color: '#374151', 
        fontWeight: '600', 
        fontSize: '0.875rem' 
    },
    td: { 
        padding: '0.75rem 1rem', 
        borderBottom: '1px solid #e5e7eb', 
        color: '#111827',
        fontSize: '0.9rem' // Un poco más pequeño para más datos
    },
    error: { color: 'red', marginTop: '1rem', textAlign: 'center' },
    loading: { textAlign: 'center', padding: '2rem', fontSize: '1.2rem', color: '#6b7280' },
    noResults: { textAlign: 'center', padding: '2rem', color: '#6b7280' }
};

const LoggerEventos = () => {
    const { token, currentUser } = useAuth();
    
    const [logs, setLogs] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    
    // Estado para los filtros
    const [filters, setFilters] = useState({
        usuario_username: '',
        accion: '',
        entidad: '',
        fecha_inicio: '',
        fecha_cierre: ''
    });

    // 2. Función para cargar los logs
    const fetchLogs = useCallback(async (currentFilters) => {
        setLoading(true);
        setError('');
        try {
            if (!currentUser?.username) {
                throw new Error("No se pudo identificar al administrador");
            }
            // Pasa los filtros al servicio
            const data = await getLogs(token, currentUser.username, currentFilters);
            setLogs(data || []);
        } catch (err) {
            setError(err.message || "Error al cargar la bitácora");
        } finally {
            setLoading(false);
        }
    }, [token, currentUser]); // Depende de token y currentUser

    // 3. Cargar logs al montar el componente
    useEffect(() => {
        fetchLogs(filters); // Carga inicial con filtros vacíos
    }, [fetchLogs]); // 'filters' no se incluye aquí para evitar recarga en cada tecleo

    // 4. Manejadores de filtros
    const handleFilterChange = (e) => {
        const { name, value } = e.target;
        setFilters(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleFilterSubmit = (e) => {
        e.preventDefault();
        fetchLogs(filters); // Llama al fetch con los filtros actuales
    };
    
    const clearFilters = () => {
        const clearedFilters = {
            usuario_username: '',
            accion: '',
            entidad: '',
            fecha_inicio: '',
            fecha_cierre: ''
        };
        setFilters(clearedFilters);
        fetchLogs(clearedFilters); // Carga con filtros limpios
    };

    // Opciones para los <select>
    const accionesOptions = [
        "CREACIÓN", "MODIFICACIÓN", "ELIMINACIÓN", "CAMBIO DE ROL", 
        "CAMBIO DE ESTADO", "ASIGNACIÓN DE PROYECTO", "DESASIGNACIÓN DE PROYECTO",
        "INICIO DE SESIÓN", "REGISTRO"
    ];
    
    const entidadesOptions = [
        "Proyectos", "Usuarios", "Labores", 
        "Equipos/Implementos", "Actividades", "Auth"
    ];

    return (
        <div style={styles.container}>
            <h2 style={styles.h2}>Logger de Eventos (Bitácora)</h2>

            {/* --- Contenedor de Filtros --- */}
            <form onSubmit={handleFilterSubmit}>
                <div style={styles.filterContainer}>
                    {/* Filtro Usuario */}
                    <div style={styles.filterGroup}>
                        <label htmlFor="usuario_username" style={styles.label}>Usuario</label>
                        <input
                            type="text"
                            id="usuario_username"
                            name="usuario_username"
                            value={filters.usuario_username}
                            onChange={handleFilterChange}
                            style={styles.input}
                            placeholder="Buscar por username..."
                        />
                    </div>

                    {/* Filtro Acción */}
                    <div style={styles.filterGroup}>
                        <label htmlFor="accion" style={styles.label}>Acción</label>
                        <select
                            id="accion"
                            name="accion"
                            value={filters.accion}
                            onChange={handleFilterChange}
                            style={styles.select}
                        >
                            <option value="">Todas</option>
                            {accionesOptions.map(opt => <option key={opt} value={opt}>{opt}</option>)}
                        </select>
                    </div>
                    
                    {/* Filtro Entidad */}
                    <div style={styles.filterGroup}>
                        <label htmlFor="entidad" style={styles.label}>Entidad</label>
                        <select
                            id="entidad"
                            name="entidad"
                            value={filters.entidad}
                            onChange={handleFilterChange}
                            style={styles.select}
                        >
                            <option value="">Todas</option>
                            {entidadesOptions.map(opt => <option key={opt} value={opt}>{opt}</option>)}
                        </select>
                    </div>

                    {/* Filtro Fecha Inicio */}
                    <div style={styles.filterGroup}>
                        <label htmlFor="fecha_inicio" style={styles.label}>Fecha Inicio</label>
                        <input
                            type="date"
                            id="fecha_inicio"
                            name="fecha_inicio"
                            value={filters.fecha_inicio}
                            onChange={handleFilterChange}
                            style={styles.input}
                        />
                    </div>
                    
                    {/* Filtro Fecha Cierre */}
                    <div style={styles.filterGroup}>
                        <label htmlFor="fecha_cierre" style={styles.label}>Fecha Cierre</label>
                        <input
                            type="date"
                            id="fecha_cierre"
                            name="fecha_cierre"
                            value={filters.fecha_cierre}
                            onChange={handleFilterChange}
                            style={styles.input}
                        />
                    </div>
                    
                    {/* Botones */}
                    <div style={styles.filterActions}>
                        <button 
                            type="button" 
                            style={{ ...styles.button, backgroundColor: '#6b7280' }} // gray-500
                            onClick={clearFilters}
                        >
                            Limpiar
                        </button>
                        <button 
                            type="submit" 
                            style={{ ...styles.button, backgroundColor: '#4f46e5' }} // indigo-600
                        >
                            Filtrar
                        </button>
                    </div>
                </div>
            </form>

            {/* --- Contenedor de la Tabla --- */}
            {loading ? (
                <p style={styles.loading}>Cargando bitácora...</p>
            ) : error ? (
                <p style={styles.error}>{error}</p>
            ) : (
                <div style={styles.tableContainer}>
                    <table style={styles.table}>
                        <thead>
                            <tr>
                                <th style={styles.th}>ID</th>
                                <th style={styles.th}>Timestamp</th>
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
                                    <td colSpan="7" style={styles.noResults}>No se encontraron eventos.</td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
};

export default LoggerEventos;