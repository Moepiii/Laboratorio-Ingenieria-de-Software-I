import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../context/AuthContext';
import { getUserDashboardData } from '../services/userService';


const styles = {
    container: { padding: '2rem', backgroundColor: '#ffffff', borderRadius: '12px', boxShadow: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)', width: '100%', maxWidth: '800px', margin: '2rem auto', fontFamily: 'Inter, sans-serif' },
    header: { fontSize: '1.875rem', fontWeight: '700', color: '#1f2937', marginBottom: '1.5rem', borderBottom: '1px solid #e5e7eb', paddingBottom: '0.75rem' },
    sectionTitle: { fontSize: '1.25rem', fontWeight: '600', color: '#374151', marginTop: '1.5rem', marginBottom: '0.75rem' },
    projectInfo: { backgroundColor: '#f9fafb', padding: '1rem', borderRadius: '8px', border: '1px solid #e5e7eb' },
    infoGrid: { display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1rem' },
    infoItem: { fontSize: '1rem' },
    infoLabel: { display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#6b7280' },
    badge: { padding: '0.25rem 0.75rem', borderRadius: '9999px', fontSize: '0.875rem', fontWeight: '500' },
    badgeSuccess: { backgroundColor: '#d1fae5', color: '#065f46' },
    badgeError: { backgroundColor: '#fee2e2', color: '#991b1b' },
    memberList: { listStyle: 'none', padding: 0 },
    memberItem: { display: 'flex', justifyContent: 'space-between', padding: '0.5rem 0', borderBottom: '1px solid #e5e7eb' },
    memberName: { fontWeight: '500' },
    memberRole: { fontSize: '0.875rem', color: '#6b7280' },
    logoutButton: { width: '100%', padding: '0.75rem 1rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#ef4444', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s', marginTop: '2rem' },
    error: { color: '#dc2626', backgroundColor: '#fee2e2', padding: '1rem', borderRadius: '8px' },
    noDataMessage: { fontStyle: 'italic', color: '#6b7280', marginTop: '1rem' }
};


const formatDate = (dateString) => {
    if (!dateString || dateString.startsWith("0001-01-01")) return "N/A";
    try {
        const date = new Date(dateString);
        // Verifica si la fecha es válida después de crearla
        if (isNaN(date.getTime())) return "Fecha inválida";
        return new Intl.DateTimeFormat('es-ES', { day: '2-digit', month: '2-digit', year: 'numeric' }).format(date);
    } catch (e) {
        return "Fecha inválida";
    }
};

const UserDashboard = () => {

    const { currentUser, userId, token, logout } = useAuth();
    const [projectData, setProjectData] = useState(null);
    const [loading, setLoading] = useState(!userId || !token);
    const [error, setError] = useState('');

    const fetchUserDashboard = useCallback(async () => {
        if (!userId || typeof userId !== 'number' || userId === 0 || !token) {
            setLoading(false);
            if (userId === 0) setError("ID de usuario no válido.");
            else setLoading(false);
            return;
        }
        setLoading(true);
        setError('');
        try {
            const data = await getUserDashboardData(token, userId);
            setProjectData(data);
            console.log("Datos recibidos del backend:", data);

        } catch (err) {
            if (err.message.includes("404") && (err.message.includes("vacía") || err.message.includes("Usuario no encontrado") || err.message.includes("Not Found"))) {
                setProjectData(null);
                setError('');
                console.log("Detectado 404 manejado, interpretado como 'sin proyecto'.");
            } else {
                setError(`Error de conexión: ${err.message}`);
            }
        } finally {
            setLoading(false);
        }
    }, [userId, token]);

    useEffect(() => {
        if (token && userId) {
            fetchUserDashboard();
        } else {
            setLoading(!token || !userId);
        }
    }, [fetchUserDashboard, token, userId]);


    // --- Renderizado Condicional ---
    if (loading || !currentUser) { // Muestra cargando si falta user o loading está activo
        return <div style={styles.container}><p>Cargando dashboard...</p></div>;
    }
    if (error) { /*  */
        return (
            <div style={styles.container}>
                <h1 style={styles.header}> Bienvenido, {currentUser?.nombre || 'Usuario'} </h1>
                <p style={styles.error}>{error}</p>
                <button onClick={logout} style={styles.logoutButton}> 🚪 Cerrar Sesión </button>
            </div>
        );
    }
    if (!projectData) { /*  */
        return (
            <div style={styles.container}>
                <h1 style={styles.header}> Bienvenido, {currentUser?.nombre || 'Usuario'} </h1>
                <p style={styles.noDataMessage}>Actualmente no estás asignado a ningún proyecto.</p>
                <button onClick={logout} style={styles.logoutButton}> 🚪 Cerrar Sesión </button>
            </div>
        );
    }

    const { proyecto, gerentes, miembros } = projectData;

    return (
        <div style={styles.container}>
            <h1 style={styles.header}> Bienvenido, {currentUser.nombre} </h1>

            {/* Detalles del Proyecto */}
            {proyecto && proyecto.id > 0 ? (
                <div style={styles.projectInfo}>
                    <h2 style={styles.sectionTitle}>Detalles de tu Proyecto </h2>
                    <div style={styles.infoGrid}>
                        <div style={styles.infoItem}><span style={styles.infoLabel}>Nombre</span> {proyecto.nombre || 'N/A'}</div>
                        {/*  */}
                        <div style={styles.infoItem}><span style={styles.infoLabel}>Fecha Inicio</span> {formatDate(proyecto.fecha_inicio)}</div>
                        {/*  */}
                        <div style={styles.infoItem}><span style={styles.infoLabel}>Fecha Cierre</span> {formatDate(proyecto.fecha_cierre)}</div>
                        <div style={styles.infoItem}>
                            <span style={styles.infoLabel}>Estado</span>
                            {proyecto.estado && (
                                <span style={{ ...styles.badge, ...(proyecto.estado === 'habilitado' ? styles.badgeSuccess : styles.badgeError) }}>
                                    {proyecto.estado}
                                </span>
                            )}
                        </div>
                    </div>
                </div>
            ) : (
                <p style={styles.noDataMessage}>Actualmente no estás asignado a ningún proyecto válido.</p>
            )}

            {/* Gerentes (sin cambios) */}
            {gerentes && Array.isArray(gerentes) && gerentes.length > 0 ? (
                <>
                    <h2 style={styles.sectionTitle}>Gerente(s) del Proyecto </h2>
                    <ul style={styles.memberList}>
                        {gerentes.map(gerente => (
                            <li key={`gerente-${gerente.id}`} style={styles.memberItem}>
                                <span style={styles.memberName}>{gerente.nombre || ''} {gerente.apellido || ''}</span>
                                <span style={styles.memberRole}>({gerente.username || 'N/A'})</span>
                            </li>
                        ))}
                    </ul>
                </>
            ) : (
                proyecto && proyecto.id > 0 && <p style={styles.noDataMessage}>No hay gerentes asignados a este proyecto.</p>
            )}

            {/* Miembros (sin cambios) */}
            {miembros && Array.isArray(miembros) && miembros.length > 0 ? (
                <>
                    <h2 style={styles.sectionTitle}>Compañeros en el Proyecto </h2>
                    <ul style={styles.memberList}>
                        {miembros.map(miembro => (
                            <li key={`miembro-${miembro.id}`} style={styles.memberItem}>
                                <span style={styles.memberName}>{miembro.nombre || ''} {miembro.apellido || ''}</span>
                                <span style={styles.memberRole}>({miembro.username || 'N/A'})</span>
                            </li>
                        ))}
                    </ul>
                </>
            ) : (
                proyecto && proyecto.id > 0 && <p style={styles.noDataMessage}>No hay otros compañeros asignados a este proyecto.</p>
            )}

            {/* Botón Logout (sin cambios) */}
            <button onClick={logout} style={styles.logoutButton}> 🚪 Cerrar Sesión </button>
        </div>
    );
};

export default UserDashboard;