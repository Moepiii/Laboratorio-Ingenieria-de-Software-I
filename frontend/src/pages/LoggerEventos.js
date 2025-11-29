import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../context/AuthContext';
// ⭐️ Importamos la nueva función deleteLogsByRange
import { getLogs, deleteLogs, deleteLogsByRange } from '../services/loggerService';
// ⭐️ Importamos el Modal
import Modal from '../components/auth/Modal';

const styles = {
    container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
    headerRow: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem', borderBottom: '2px solid #e5e7eb', paddingBottom: '0.75rem' },
    h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', margin: 0 },

    filterContainer: { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1rem', padding: '1.5rem', backgroundColor: '#f9fafb', borderRadius: '8px', marginBottom: '2rem', border: '1px solid #e5e7eb' },
    filterGroup: { display: 'flex', flexDirection: 'column' },
    label: { fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' },
    input: { padding: '0.6rem', border: '1px solid #d1d5db', borderRadius: '6px', fontSize: '0.9rem' },
    select: { padding: '0.6rem', border: '1px solid #d1d5db', borderRadius: '6px', fontSize: '0.9rem', backgroundColor: 'white' },

    buttonGroup: { display: 'flex', justifyContent: 'flex-end', gap: '1rem', marginTop: '1rem' },
    button: { padding: '0.6rem 1.2rem', backgroundColor: '#2563eb', color: 'white', border: 'none', borderRadius: '6px', cursor: 'pointer', fontWeight: '600', transition: 'background 0.2s' },

    tableContainer: { overflowX: 'auto', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 2px 4px rgba(0,0,0,0.05)' },
    table: { width: '100%', borderCollapse: 'collapse', fontSize: '0.9rem' },
    th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '1px solid #e5e7eb', color: '#4b5563', fontWeight: '600' },
    td: { padding: '0.75rem 1rem', borderBottom: '1px solid #f3f4f6', color: '#1f2937' },
    checkbox: { cursor: 'pointer', width: '16px', height: '16px' },
    noResults: { textAlign: 'center', padding: '2rem', color: '#6b7280' },

    // Estilos para el Modal de Borrado
    cancelButton: { padding: '0.6rem 1.2rem', backgroundColor: '#9ca3af', color: 'white', border: 'none', borderRadius: '6px', cursor: 'pointer', marginRight: '10px' },
    deleteButton: { padding: '0.6rem 1.2rem', backgroundColor: '#dc2626', color: 'white', border: 'none', borderRadius: '6px', cursor: 'pointer' }
};

const LoggerEventos = () => {
    const { token, currentUser } = useAuth();
    const [logs, setLogs] = useState([]);
    const [filters, setFilters] = useState({
        fecha_inicio: '',
        fecha_cierre: '',
        usuario_username: '',
        accion: '',
        entidad: ''
    });
    const [selectedIds, setSelectedIds] = useState([]);
    const [loading, setLoading] = useState(false);

    // ⭐️ ESTADOS PARA EL BORRADO MASIVO
    const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
    const [deleteMode, setDeleteMode] = useState('month'); // 'month', 'quarter', 'year'
    const [selectedYear, setSelectedYear] = useState(new Date().getFullYear());
    const [selectedMonth, setSelectedMonth] = useState(new Date().getMonth() + 1);
    const [selectedQuarter, setSelectedQuarter] = useState(1);

    const fetchLogs = useCallback(async () => {
        if (!token || !currentUser) return;
        setLoading(true);
        try {
            const data = await getLogs(token, currentUser.username, filters);
            setLogs(data || []);
        } catch (error) {
            console.error(error);
            alert("Error cargando logs: " + error.message);
        } finally {
            setLoading(false);
        }
    }, [token, currentUser, filters]);

    useEffect(() => {
        fetchLogs();
    }, [fetchLogs]);

    const handleFilterChange = (e) => {
        const { name, value } = e.target;
        setFilters(prev => ({ ...prev, [name]: value }));
    };

    const handleSelectOne = (id) => {
        if (selectedIds.includes(id)) {
            setSelectedIds(selectedIds.filter(sid => sid !== id));
        } else {
            setSelectedIds([...selectedIds, id]);
        }
    };

    const handleSelectAll = (e) => {
        if (e.target.checked) {
            setSelectedIds(logs.map(l => l.id));
        } else {
            setSelectedIds([]);
        }
    };

    const handleDeleteSelected = async () => {
        if (!window.confirm(`¿Estás seguro de eliminar ${selectedIds.length} eventos?`)) return;

        try {
            await deleteLogs(token, selectedIds, currentUser.username);
            alert("Logs eliminados correctamente.");
            setSelectedIds([]);
            fetchLogs();
        } catch (error) {
            alert("Error eliminando logs: " + error.message);
        }
    };

    // ⭐️ LÓGICA DE BORRADO MASIVO POR FECHAS
    const handleBulkDelete = async (e) => {
        e.preventDefault();

        const confirmMsg = "⚠️ ¿Estás seguro? Esta acción borrará PERMANENTEMENTE todo el historial del periodo seleccionado. No se puede deshacer.";
        if (!window.confirm(confirmMsg)) return;

        let startDate = '';
        let endDate = '';

        // Cálculo de fechas
        if (deleteMode === 'month') {
            // Ejemplo: 2024-02-01
            startDate = `${selectedYear}-${String(selectedMonth).padStart(2, '0')}-01`;
            // Último día del mes (el día 0 del mes siguiente)
            const lastDay = new Date(selectedYear, selectedMonth, 0).getDate();
            endDate = `${selectedYear}-${String(selectedMonth).padStart(2, '0')}-${lastDay}`;
        }
        else if (deleteMode === 'quarter') {
            const startMonth = (selectedQuarter - 1) * 3 + 1; // 1, 4, 7, 10
            const endMonth = startMonth + 2; // 3, 6, 9, 12

            startDate = `${selectedYear}-${String(startMonth).padStart(2, '0')}-01`;
            const lastDay = new Date(selectedYear, endMonth, 0).getDate();
            endDate = `${selectedYear}-${String(endMonth).padStart(2, '0')}-${lastDay}`;
        }
        else if (deleteMode === 'year') {
            startDate = `${selectedYear}-01-01`;
            endDate = `${selectedYear}-12-31`;
        }

        try {
            await deleteLogsByRange(token, startDate, endDate, currentUser.username);
            alert(`Historial del periodo eliminado correctamente.`);
            setIsDeleteModalOpen(false);
            fetchLogs(); // Recargar tabla
        } catch (err) {
            alert("Error al eliminar historial: " + err.message);
        }
    };

    return (
        <div style={styles.container}>
            <div style={styles.headerRow}>
                <h2 style={styles.h2}>Auditoría de Eventos</h2>
                <div style={{ display: 'flex', gap: '10px' }}>
                    <button
                        onClick={fetchLogs}
                        style={styles.button}
                        disabled={loading}
                    >
                        {loading ? 'Cargando...' : '🔄 Refrescar'}
                    </button>
                </div>
            </div>

            {/* Filtros */}
            <div style={styles.filterContainer}>
                <div style={styles.filterGroup}>
                    <label style={styles.label}>Fecha Inicio</label>
                    <input type="date" name="fecha_inicio" value={filters.fecha_inicio} onChange={handleFilterChange} style={styles.input} />
                </div>
                <div style={styles.filterGroup}>
                    <label style={styles.label}>Fecha Fin</label>
                    <input type="date" name="fecha_cierre" value={filters.fecha_cierre} onChange={handleFilterChange} style={styles.input} />
                </div>
                <div style={styles.filterGroup}>
                    <label style={styles.label}>Usuario (Username)</label>
                    <input type="text" name="usuario_username" value={filters.usuario_username} onChange={handleFilterChange} placeholder="Ej. admin" style={styles.input} />
                </div>
                <div style={styles.filterGroup}>
                    <label style={styles.label}>Acción</label>
                    <select name="accion" value={filters.accion} onChange={handleFilterChange} style={styles.select}>
                        <option value="">-- Todas --</option>
                        <option value="INICIO DE SESIÓN">INICIO DE SESIÓN</option>
                        <option value="CREACIÓN">CREACIÓN</option>
                        <option value="MODIFICACIÓN">MODIFICACIÓN</option>
                        <option value="ELIMINACIÓN">ELIMINACIÓN</option>
                        <option value="CAMBIO ESTADO">CAMBIO ESTADO</option>
                        <option value="ASIGNACIÓN PROYECTO">ASIGNACIÓN PROYECTO</option>
                    </select>
                </div>
                <div style={styles.filterGroup}>
                    <label style={styles.label}>Entidad</label>
                    <select name="entidad" value={filters.entidad} onChange={handleFilterChange} style={styles.select}>
                        <option value="">-- Todas --</option>
                        <option value="Auth">Auth (Login)</option>
                        <option value="Usuarios">Usuarios</option>
                        <option value="Proyectos">Proyectos</option>
                        <option value="Labores">Labores</option>
                        <option value="Equipos/Implementos">Equipos</option>
                        <option value="Unidades Medida">Unidades Medida</option>
                        <option value="Actividades">Actividades</option>
                        <option value="Logs">Logs (Auditoría)</option>
                    </select>
                </div>
            </div>

            {/* Botones de Acción */}
            <div style={{ marginBottom: '1rem', display: 'flex', gap: '10px' }}>
                <button
                    onClick={handleDeleteSelected}
                    disabled={selectedIds.length === 0}
                    style={{
                        ...styles.button,
                        backgroundColor: selectedIds.length > 0 ? '#ef4444' : '#9ca3af',
                        cursor: selectedIds.length > 0 ? 'pointer' : 'not-allowed'
                    }}
                >
                    🗑️ Eliminar Seleccionados ({selectedIds.length})
                </button>

                {/* ⭐️ BOTÓN NUEVO PARA BORRADO MASIVO */}
                <button
                    onClick={() => setIsDeleteModalOpen(true)}
                    style={{ ...styles.button, backgroundColor: '#dc2626' }}
                >
                    📅 Borrar Historial Antiguo
                </button>
            </div>

            {/* Tabla */}
            <div style={styles.tableContainer}>
                <table style={styles.table}>
                    <thead>
                        <tr>
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
                            <th style={styles.th}>ID Entidad</th>
                        </tr>
                    </thead>
                    <tbody>
                        {logs.length > 0 ? (
                            logs.map(log => (
                                <tr key={log.id}>
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

            {/* ⭐️ MODAL DE BORRADO MASIVO POR FECHAS ⭐️ */}
            <Modal isOpen={isDeleteModalOpen} onClose={() => setIsDeleteModalOpen(false)} title="Eliminación de Historial por Periodo">
                <form onSubmit={handleBulkDelete} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    <div style={{ backgroundColor: '#fff3cd', color: '#856404', padding: '10px', borderRadius: '4px', fontSize: '0.9rem' }}>
                        <strong>Atención:</strong> Estás a punto de borrar grandes cantidades de registros. Esta acción no se puede deshacer.
                    </div>

                    {/* Selector de TIPO */}
                    <div style={styles.filterGroup}>
                        <label style={styles.label}>Modo de Borrado:</label>
                        <select
                            value={deleteMode}
                            onChange={(e) => setDeleteMode(e.target.value)}
                            style={styles.select}
                        >
                            <option value="month">Borrar un Mes entero</option>
                            <option value="quarter">Borrar un Trimestre</option>
                            <option value="year">Borrar todo el Año</option>
                        </select>
                    </div>

                    {/* Selector de AÑO */}
                    <div style={styles.filterGroup}>
                        <label style={styles.label}>Año:</label>
                        <input
                            type="number"
                            value={selectedYear}
                            onChange={(e) => setSelectedYear(e.target.value)}
                            style={styles.input}
                        />
                    </div>

                    {/* Selector condicional: MES */}
                    {deleteMode === 'month' && (
                        <div style={styles.filterGroup}>
                            <label style={styles.label}>Mes:</label>
                            <select
                                value={selectedMonth}
                                onChange={(e) => setSelectedMonth(Number(e.target.value))}
                                style={styles.select}
                            >
                                <option value={1}>Enero</option>
                                <option value={2}>Febrero</option>
                                <option value={3}>Marzo</option>
                                <option value={4}>Abril</option>
                                <option value={5}>Mayo</option>
                                <option value={6}>Junio</option>
                                <option value={7}>Julio</option>
                                <option value={8}>Agosto</option>
                                <option value={9}>Septiembre</option>
                                <option value={10}>Octubre</option>
                                <option value={11}>Noviembre</option>
                                <option value={12}>Diciembre</option>
                            </select>
                        </div>
                    )}

                    {/* Selector condicional: TRIMESTRE */}
                    {deleteMode === 'quarter' && (
                        <div style={styles.filterGroup}>
                            <label style={styles.label}>Trimestre:</label>
                            <select
                                value={selectedQuarter}
                                onChange={(e) => setSelectedQuarter(Number(e.target.value))}
                                style={styles.select}
                            >
                                <option value={1}>Trimestre 1 (Ene - Mar)</option>
                                <option value={2}>Trimestre 2 (Abr - Jun)</option>
                                <option value={3}>Trimestre 3 (Jul - Sep)</option>
                                <option value={4}>Trimestre 4 (Oct - Dic)</option>
                            </select>
                        </div>
                    )}

                    <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '1rem', borderTop: '1px solid #eee', paddingTop: '1rem' }}>
                        <button type="button" onClick={() => setIsDeleteModalOpen(false)} style={styles.cancelButton}>Cancelar</button>
                        <button type="submit" style={styles.deleteButton}>Eliminar Definitivamente</button>
                    </div>
                </form>
            </Modal>
        </div>
    );
};

export default LoggerEventos;