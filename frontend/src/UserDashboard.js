import React, { useState, useEffect, useCallback } from 'react';

// Estilos básicos (puedes personalizarlos o usar los de App.js si los exportas)
const styles = {
    container: {
        padding: '2rem',
        backgroundColor: '#ffffff',
        borderRadius: '12px',
        boxShadow: '0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05)',
        width: '100%',
        maxWidth: '800px', // Más ancho para mostrar listas
        margin: '2rem auto',
        fontFamily: 'Inter, sans-serif',
    },
    header: {
        fontSize: '1.875rem',
        fontWeight: '700',
        color: '#1f2937',
        marginBottom: '1.5rem',
        borderBottom: '1px solid #e5e7eb',
        paddingBottom: '0.75rem',
    },
    sectionTitle: {
        fontSize: '1.25rem',
        fontWeight: '600',
        color: '#374151',
        marginTop: '1.5rem',
        marginBottom: '0.75rem',
    },
    projectInfo: {
        backgroundColor: '#f9fafb',
        padding: '1rem',
        borderRadius: '8px',
        border: '1px solid #e5e7eb',
        marginBottom: '1.5rem',
    },
    projectDetail: {
        marginBottom: '0.5rem',
        fontSize: '1rem',
    },
    label: {
        fontWeight: '600',
        color: '#4b5563',
        marginRight: '0.5rem',
    },
    memberList: {
        listStyle: 'none',
        padding: 0,
        margin: 0,
    },
    memberItem: {
        padding: '0.5rem 0',
        borderBottom: '1px solid #f3f4f6',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
    },
    memberName: {
        fontWeight: '500',
        color: '#1f2937',
    },
    memberRole: {
        fontSize: '0.875rem',
        color: '#6b7280',
        fontStyle: 'italic',
    },
    logoutButton: {
        display: 'inline-flex',
        alignItems: 'center',
        padding: '0.6rem 1.2rem',
        fontSize: '0.9rem',
        fontWeight: '600',
        borderRadius: '8px',
        color: 'white',
        backgroundColor: '#ef4444', // red-500
        border: 'none',
        cursor: 'pointer',
        transition: 'background-color 0.2s',
        marginTop: '2rem',
        float: 'right', // Alinea a la derecha
    },
    loadingText: {
        textAlign: 'center',
        color: '#4f46e5',
        padding: '2rem',
    },
    errorText: {
        color: '#dc2626', // red-600
        backgroundColor: '#fef2f2', // red-50
        padding: '1rem',
        borderRadius: '8px',
        border: '1px solid #fecaca', // red-200
        textAlign: 'center',
    },
    noProjectText: {
        textAlign: 'center',
        color: '#6b7280',
        padding: '2rem',
        fontStyle: 'italic',
    }
};

// ⭐️ --- NUEVO COMPONENTE --- ⭐️
const UserDashboard = ({ currentUser, userId, apiCall, handleLogout }) => {
    const [projectData, setProjectData] = useState(null); // { proyecto: {}, miembros: [], gerentes: [] }
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    // Función para formatear fecha (igual que en PortafolioProyectos)
    const formatDate = (dateString) => {
        if (!dateString) return 'N/A';
        try {
            const [year, month, day] = dateString.split('-');
            return `${day}/${month}/${year}`;
        } catch { return dateString; }
    };

    // Carga los detalles del proyecto al montar
    useEffect(() => {
        const fetchDetails = async () => {
            if (!userId) {
                setError("No se pudo obtener el ID del usuario.");
                setLoading(false);
                return;
            }
            setLoading(true);
            setError('');
            try {
                const result = await apiCall('user/project-details', { user_id: userId }, 'POST');
                if (result.success) {
                    setProjectData(result.data);
                } else {
                    setError(result.data.error || 'Error al cargar los datos del proyecto.');
                }
            } catch (e) {
                setError(`Error de conexión: ${e.message}`);
            } finally {
                setLoading(false);
            }
        };

        fetchDetails();
    }, [apiCall, userId]); // Depende de apiCall y userId

    // Renderizado condicional
    if (loading) {
        return <div style={styles.container}><p style={styles.loadingText}>Cargando información del proyecto...</p></div>;
    }

    if (error) {
        return <div style={styles.container}><p style={styles.errorText}>Error: {error}</p></div>;
    }

    // Si no hay proyecto asignado
    if (!projectData || !projectData.proyecto) {
        return (
            <div style={styles.container}>
                <h1 style={styles.header}>Bienvenido, {currentUser}</h1>
                <p style={styles.noProjectText}>Actualmente no estás asignado a ningún proyecto.</p>
                <button
                    onClick={handleLogout}
                    style={styles.logoutButton}
                    onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#dc2626'} // red-600
                    onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#ef4444'} // red-500
                >
                    🚪 Cerrar Sesión
                </button>
            </div>
        );
    }

    // Si hay proyecto, muestra los detalles
    const { proyecto, miembros, gerentes } = projectData;

    return (
        <div style={styles.container}>
            <h1 style={styles.header}>Bienvenido, {currentUser}</h1>

            {/* Información del Proyecto */}
            <h2 style={styles.sectionTitle}>Proyecto Asignado 📂</h2>
            <div style={styles.projectInfo}>
                <p style={styles.projectDetail}><span style={styles.label}>Nombre:</span> {proyecto.nombre}</p>
                <p style={styles.projectDetail}><span style={styles.label}>Inicio:</span> {formatDate(proyecto.fecha_inicio)}</p>
                <p style={styles.projectDetail}><span style={styles.label}>Cierre:</span> {formatDate(proyecto.fecha_cierre)}</p>
            </div>

            {/* Gerentes del Proyecto */}
            <h2 style={styles.sectionTitle}>Gerentes Asignados al Proyecto 🧑‍💼</h2>
            {gerentes.length > 0 ? (
                <ul style={styles.memberList}>
                    {gerentes.map(gerente => (
                        <li key={gerente.id} style={styles.memberItem}>
                            <span style={styles.memberName}>{gerente.nombre} {gerente.apellido}</span>
                            <span style={styles.memberRole}>({gerente.username})</span>
                        </li>
                    ))}
                </ul>
            ) : (
                <p style={{ fontStyle: 'italic', color: '#6b7280' }}>No hay gerentes asignados a este proyecto.</p>
            )}

            {/* Compañeros del Proyecto */}
            <h2 style={styles.sectionTitle}>Compañeros en el Proyecto 👥</h2>
            {miembros.length > 0 ? (
                <ul style={styles.memberList}>
                    {miembros.map(miembro => (
                        <li key={miembro.id} style={styles.memberItem}>
                            <span style={styles.memberName}>{miembro.nombre} {miembro.apellido}</span>
                            <span style={styles.memberRole}>({miembro.username})</span>
                        </li>
                    ))}
                </ul>
            ) : (
                <p style={{ fontStyle: 'italic', color: '#6b7280' }}>No hay otros compañeros asignados a este proyecto.</p>
            )}

            {/* Botón de Logout */}
            <button
                onClick={handleLogout}
                style={styles.logoutButton}
                onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#dc2626'} // red-600
                onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#ef4444'} // red-500
            >
                🚪 Cerrar Sesión
            </button>
        </div>
    );
};

export default UserDashboard;