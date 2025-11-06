import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
// Importaremos esto en el futuro para obtener los datos del proyecto
// import { getProjectById } from '../services/projectService'; 

// Estilos para la nueva página
const styles = {
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
        backgroundColor: '#6b7280', // gray-500
        border: 'none',
        cursor: 'pointer',
        textDecoration: 'none',
    },
    section: {
        padding: '1.5rem',
        backgroundColor: '#f9fafb', // gray-50
        borderRadius: '8px',
        border: '1px solid #e5e7eb',
        marginBottom: '1.5rem',
        boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)'
    },
    sectionTitle: {
        fontSize: '1.25rem',
        fontWeight: '600',
        color: '#111827',
        marginTop: 0,
        marginBottom: '1rem',
        borderBottom: '1px solid #d1d5db',
        paddingBottom: '0.5rem'
    },
    placeholder: {
        fontSize: '1rem',
        color: '#6b7280',
        fontStyle: 'italic',
    }
};

const DatosProyecto = () => {
    const { id } = useParams(); // Obtiene el ID del proyecto de la URL
    // const [proyecto, setProyecto] = useState(null); // Lo usaremos en el futuro
    // const [loading, setLoading] = useState(true);
    // const { token, currentUser } = useAuth();

    // useEffect(() => {
    //     // En el futuro, aquí llamaremos a una función para cargar los datos del proyecto
    //     // const fetchDatos = async () => { ... }
    //     // fetchDatos();
    // }, [id, token]);

    return (
        <div style={styles.container}>
            <div style={styles.header}>
                <h2 style={styles.h2}>Datos del Proyecto (ID: {id})</h2>
                {/* Este botón te regresa a la lista de proyectos */}
                <Link to="/admin/proyectos" style={styles.backButton}>
                    &larr; Volver al Portafolio
                </Link>
            </div>

            {/* Sección de Resumen (Placeholder) */}
            <section style={styles.section}>
                <h3 style={styles.sectionTitle}>Resumen</h3>
                <p style={styles.placeholder}>
                    (Aquí se mostrará el resumen y estado del proyecto)
                </p>
            </section>

            {/* Sección de Labores (Placeholder) */}
            <section style={styles.section}>
                <h3 style={styles.sectionTitle}>Labores Agronómicas</h3>
                <p style={styles.placeholder}>
                    (Aquí se mostrará un resumen o dashboard de las labores)
                </p>
            </section>

            {/* Sección de Equipos (Placeholder) */}
            <section style={styles.section}>
                <h3 style={styles.sectionTitle}>Equipos e Implementos</h3>
                <p style={styles.placeholder}>
                    (Aquí se mostrará un resumen o dashboard de los equipos)
                </p>
            </section>
        </div>
    );
};

export default DatosProyecto;