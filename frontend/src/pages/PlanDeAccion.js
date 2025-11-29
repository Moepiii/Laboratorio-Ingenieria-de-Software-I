import React, { useState } from 'react';
// 1. Importamos el componente Modal que ya tienes
import Modal from '../components/auth/Modal';

const styles = {
    container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
    header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '2px solid #e5e7eb', paddingBottom: '1rem', marginBottom: '2rem' },
    h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', margin: 0 },
    addButton: { padding: '0.75rem 1.5rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#2563eb', border: 'none', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '0.5rem' },
    tableContainer: { overflowX: 'auto', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
    table: { width: '100%', borderCollapse: 'collapse' },
    th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '2px solid #e5e7eb', color: '#374151', fontWeight: '600', fontSize: '0.875rem' },
    td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', color: '#4b5563', fontSize: '0.875rem' },
    
    // Estilos de botones de acción
    actionButton: { padding: '0.4rem 0.8rem', borderRadius: '4px', border: 'none', cursor: 'pointer', fontSize: '0.875rem', marginLeft: '0.5rem', color: 'white' },
    editButton: { backgroundColor: '#f59e0b' },
    deleteButton: { backgroundColor: '#ef4444' },
    noResults: { textAlign: 'center', padding: '2rem', color: '#6b7280' },

    // --- Estilos del Formulario en el Modal ---
    form: { display: 'flex', flexDirection: 'column', gap: '1rem' },
    formGroup: { display: 'flex', flexDirection: 'column', gap: '0.5rem' },
    // Para poner dos campos en una misma fila (opcional, se ve mejor)
    rowGroup: { display: 'flex', gap: '1rem' }, 
    label: { fontSize: '0.875rem', fontWeight: '500', color: '#374151' },
    input: { padding: '0.6rem', border: '1px solid #d1d5db', borderRadius: '6px', fontSize: '1rem' },
    
    formActions: { display: 'flex', justifyContent: 'flex-end', gap: '1rem', marginTop: '1.5rem' },
    cancelButton: { padding: '0.6rem 1.2rem', border: '1px solid #d1d5db', borderRadius: '6px', backgroundColor: 'white', color: '#374151', cursor: 'pointer' },
    saveButton: { padding: '0.6rem 1.2rem', border: 'none', borderRadius: '6px', backgroundColor: '#2563eb', color: 'white', cursor: 'pointer' },
};

const PlanDeAccion = () => {
    // Estado de la tabla (Datos iniciales)
    const [planes, setPlanes] = useState([
        { id: 1, actividad: 'Preparación de Suelo', accion: 'Arado Profundo', fecha_inicio: '2023-12-01', fecha_cierre: '2023-12-05', horas: 40, responsable: 'Juan Pérez', monto: 1500.00 },
        { id: 2, actividad: 'Siembra', accion: 'Siembra de Maíz', fecha_inicio: '2023-12-10', fecha_cierre: '2023-12-12', horas: 16, responsable: 'María Gómez', monto: 800.50 },
    ]);

    // Estados para el Modal
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [formData, setFormData] = useState({
        actividad: '',
        accion: '',
        fecha_inicio: '',
        fecha_cierre: '',
        horas: '',
        responsable: '',
        monto: ''
    });

    // Abrir modal y limpiar formulario
    const handleOpenModal = () => {
        setFormData({ actividad: '', accion: '', fecha_inicio: '', fecha_cierre: '', horas: '', responsable: '', monto: '' });
        setIsModalOpen(true);
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
    };

    // Manejar cambios en los inputs
    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData({ ...formData, [name]: value });
    };

    // Guardar (Agregar a la lista localmente)
    const handleSubmit = (e) => {
        e.preventDefault();
        
        // Creamos un nuevo objeto con un ID temporal
        const nuevoPlan = {
            id: Date.now(), // ID único basado en fecha
            ...formData,
            horas: Number(formData.horas),
            monto: Number(formData.monto)
        };

        setPlanes([...planes, nuevoPlan]); // Agregamos a la tabla
        handleCloseModal(); // Cerramos modal
    };

    const handleDeleteClick = (id) => {
        if(window.confirm("¿Borrar este elemento?")) {
            setPlanes(planes.filter(p => p.id !== id));
        }
    };

    return (
        <div style={styles.container}>
            <div style={styles.header}>
                <h2 style={styles.h2}>Plan de Acción</h2>
                <button style={styles.addButton} onClick={handleOpenModal}>
                    <span>+</span> Añadir
                </button>
            </div>

            {/* Tabla */}
            <div style={styles.tableContainer}>
                <table style={styles.table}>
                    <thead>
                        <tr>
                            <th style={styles.th}>ID</th>
                            <th style={styles.th}>Actividad</th>
                            <th style={styles.th}>Acción</th>
                            <th style={styles.th}>Fecha Inicio</th>
                            <th style={styles.th}>Fecha Cierre</th>
                            <th style={styles.th}>Horas</th>
                            <th style={styles.th}>Responsable</th>
                            <th style={styles.th}>Monto ($)</th>
                            <th style={styles.th}>Acciones</th>
                        </tr>
                    </thead>
                    <tbody>
                        {planes.length > 0 ? (
                            planes.map((plan) => (
                                <tr key={plan.id}>
                                    <td style={styles.td}>{plan.id}</td>
                                    <td style={styles.td}>{plan.actividad}</td>
                                    <td style={styles.td}>{plan.accion}</td>
                                    <td style={styles.td}>{plan.fecha_inicio}</td>
                                    <td style={styles.td}>{plan.fecha_cierre}</td>
                                    <td style={styles.td}>{plan.horas}</td>
                                    <td style={styles.td}>{plan.responsable}</td>
                                    <td style={styles.td}>${parseFloat(plan.monto).toFixed(2)}</td>
                                    <td style={styles.td}>
                                        <button style={{...styles.actionButton, ...styles.editButton}} onClick={() => alert("Editar: " + plan.id)}>Editar</button>
                                        <button style={{...styles.actionButton, ...styles.deleteButton}} onClick={() => handleDeleteClick(plan.id)}>Borrar</button>
                                    </td>
                                </tr>
                            ))
                        ) : (
                            <tr><td colSpan="9" style={styles.noResults}>No hay planes registrados.</td></tr>
                        )}
                    </tbody>
                </table>
            </div>

            {/* ⭐️ MODAL DE AÑADIR ⭐️ */}
            <Modal isOpen={isModalOpen} onClose={handleCloseModal} title="Añadir Plan de Acción">
                <form onSubmit={handleSubmit} style={styles.form}>
                    
                    <div style={styles.formGroup}>
                        <label style={styles.label}>Actividad</label>
                        <input name="actividad" value={formData.actividad} onChange={handleInputChange} required style={styles.input} placeholder="Ej. Preparación de Suelo" />
                    </div>

                    <div style={styles.formGroup}>
                        <label style={styles.label}>Acción</label>
                        <input name="accion" value={formData.accion} onChange={handleInputChange} required style={styles.input} placeholder="Ej. Arado" />
                    </div>

                    <div style={styles.rowGroup}>
                        <div style={{...styles.formGroup, flex: 1}}>
                            <label style={styles.label}>Fecha de Inicio</label>
                            <input type="date" name="fecha_inicio" value={formData.fecha_inicio} onChange={handleInputChange} required style={styles.input} />
                        </div>
                        <div style={{...styles.formGroup, flex: 1}}>
                            <label style={styles.label}>Fecha de Cierre</label>
                            <input type="date" name="fecha_cierre" value={formData.fecha_cierre} onChange={handleInputChange} required style={styles.input} />
                        </div>
                    </div>

                    <div style={styles.rowGroup}>
                        <div style={{...styles.formGroup, flex: 1}}>
                            <label style={styles.label}>Cantidad de Horas</label>
                            <input type="number" name="horas" value={formData.horas} onChange={handleInputChange} required style={styles.input} placeholder="0" />
                        </div>
                        <div style={{...styles.formGroup, flex: 1}}>
                            <label style={styles.label}>Monto ($)</label>
                            <input type="number" step="0.01" name="monto" value={formData.monto} onChange={handleInputChange} required style={styles.input} placeholder="0.00" />
                        </div>
                    </div>

                    <div style={styles.formGroup}>
                        <label style={styles.label}>Responsable</label>
                        <input name="responsable" value={formData.responsable} onChange={handleInputChange} required style={styles.input} placeholder="Nombre del responsable" />
                    </div>

                    <div style={styles.formActions}>
                        <button type="button" onClick={handleCloseModal} style={styles.cancelButton}>Cancelar</button>
                        <button type="submit" style={styles.saveButton}>Guardar</button>
                    </div>
                </form>
            </Modal>
        </div>
    );
};

export default PlanDeAccion;