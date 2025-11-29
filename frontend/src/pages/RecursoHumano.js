import React, { useState } from 'react';
import Modal from '../components/auth/Modal';

// Estilos
const styles = {
    container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
    header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '2px solid #e5e7eb', paddingBottom: '1rem', marginBottom: '2rem' },
    h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', margin: 0 },
    addButton: { padding: '0.75rem 1.5rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#2563eb', border: 'none', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '0.5rem' },
    tableContainer: { overflowX: 'auto', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
    table: { width: '100%', borderCollapse: 'collapse' },
    th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '2px solid #e5e7eb', color: '#374151', fontWeight: '600', fontSize: '0.875rem' },
    td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', color: '#4b5563', fontSize: '0.875rem' },
    actionButton: { padding: '0.4rem 0.8rem', borderRadius: '4px', border: 'none', cursor: 'pointer', fontSize: '0.875rem', marginLeft: '0.5rem', color: 'white' },
    editButton: { backgroundColor: '#f59e0b' },
    deleteButton: { backgroundColor: '#ef4444' },
    noResults: { textAlign: 'center', padding: '2rem', color: '#6b7280' },

    // Estilos del Formulario
    form: { display: 'flex', flexDirection: 'column', gap: '1rem' },
    formGroup: { display: 'flex', flexDirection: 'column', gap: '0.5rem' },
    rowGroup: { display: 'flex', gap: '1rem' }, 
    label: { fontSize: '0.875rem', fontWeight: '500', color: '#374151' },
    input: { padding: '0.6rem', border: '1px solid #d1d5db', borderRadius: '6px', fontSize: '1rem' },
    formActions: { display: 'flex', justifyContent: 'flex-end', gap: '1rem', marginTop: '1.5rem' },
    cancelButton: { padding: '0.6rem 1.2rem', border: '1px solid #d1d5db', borderRadius: '6px', backgroundColor: 'white', color: '#374151', cursor: 'pointer' },
    saveButton: { padding: '0.6rem 1.2rem', border: 'none', borderRadius: '6px', backgroundColor: '#2563eb', color: 'white', cursor: 'pointer' },
};

const RecursoHumano = () => {
    // Datos de prueba (Tabla)
    const [recursos, setRecursos] = useState([
        { 
            id: 1, 
            actividad: 'Cosecha', 
            accion: 'Recolección Manual', 
            tiempo: '8 horas', 
            cantidad: 5, 
            costo: 80.00, // Nuevo campo interno
            responsable: 'Juan Pérez', 
            monto: 400.00 
        },
        { 
            id: 2, 
            actividad: 'Poda', 
            accion: 'Poda de Formación', 
            tiempo: '4 horas', 
            cantidad: 2, 
            costo: 75.25,
            responsable: 'Carlos López', 
            monto: 150.50 
        },
    ]);

    // Estado del Formulario
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [formData, setFormData] = useState({
        actividad: '',
        accion: '',
        tiempo: '',
        cantidad: '',
        costo: '',      // ⭐️ Nuevo campo
        responsable: '',
        monto: ''
    });

    const handleOpenModal = () => {
        // Limpiar formulario al abrir
        setFormData({ actividad: '', accion: '', tiempo: '', cantidad: '', costo: '', responsable: '', monto: '' });
        setIsModalOpen(true);
    };

    const handleCloseModal = () => setIsModalOpen(false);

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        
        // Lógica opcional: Si el usuario cambia Cantidad o Costo, calcular Monto automáticamente
        let newData = { ...formData, [name]: value };

        if (name === 'cantidad' || name === 'costo') {
            const cant = parseFloat(name === 'cantidad' ? value : formData.cantidad) || 0;
            const cost = parseFloat(name === 'costo' ? value : formData.costo) || 0;
            if (cant > 0 && cost > 0) {
                newData.monto = (cant * cost).toFixed(2);
            }
        }

        setFormData(newData);
    };

    const handleSubmit = (e) => {
        e.preventDefault();
        const nuevoRecurso = {
            id: Date.now(),
            ...formData,
            cantidad: parseInt(formData.cantidad) || 0,
            costo: parseFloat(formData.costo) || 0,
            monto: parseFloat(formData.monto) || 0
        };
        setRecursos([...recursos, nuevoRecurso]);
        handleCloseModal();
    };

    const handleDeleteClick = (id) => {
        if (window.confirm("¿Estás seguro de eliminar este registro?")) {
            setRecursos(recursos.filter(r => r.id !== id));
        }
    };

    return (
        <div style={styles.container}>
            <div style={styles.header}>
                <h2 style={styles.h2}>Recurso Humano</h2>
                <button style={styles.addButton} onClick={handleOpenModal}>
                    <span>+</span> Añadir
                </button>
            </div>

            {/* Tabla con las columnas solicitadas */}
            <div style={styles.tableContainer}>
                <table style={styles.table}>
                    <thead>
                        <tr>
                            <th style={styles.th}>ID</th>
                            <th style={styles.th}>Actividad</th>
                            <th style={styles.th}>Acción</th>
                            <th style={styles.th}>Tiempo</th>
                            <th style={styles.th}>Cantidad</th>
                            <th style={styles.th}>Responsable</th>
                            <th style={styles.th}>Monto ($)</th>
                            <th style={styles.th}>Acciones</th>
                        </tr>
                    </thead>
                    <tbody>
                        {recursos.length > 0 ? (
                            recursos.map((rec) => (
                                <tr key={rec.id}>
                                    <td style={styles.td}>{rec.id}</td>
                                    <td style={styles.td}>{rec.actividad}</td>
                                    <td style={styles.td}>{rec.accion}</td>
                                    <td style={styles.td}>{rec.tiempo}</td>
                                    <td style={styles.td}>{rec.cantidad}</td>
                                    <td style={styles.td}>{rec.responsable}</td>
                                    <td style={styles.td}>${rec.monto.toFixed(2)}</td>
                                    <td style={styles.td}>
                                        <button style={{...styles.actionButton, ...styles.editButton}} onClick={() => alert("Editar: " + rec.id)}>Editar</button>
                                        <button style={{...styles.actionButton, ...styles.deleteButton}} onClick={() => handleDeleteClick(rec.id)}>Borrar</button>
                                    </td>
                                </tr>
                            ))
                        ) : (
                            <tr><td colSpan="8" style={styles.noResults}>No hay registros.</td></tr>
                        )}
                    </tbody>
                </table>
            </div>

            {/* Modal de Añadir */}
            <Modal isOpen={isModalOpen} onClose={handleCloseModal} title="Añadir Recurso Humano">
                <form onSubmit={handleSubmit} style={styles.form}>
                    
                    <div style={styles.formGroup}>
                        <label style={styles.label}>Actividad</label>
                        <input name="actividad" value={formData.actividad} onChange={handleInputChange} required style={styles.input} placeholder="Ej. Cosecha" />
                    </div>

                    <div style={styles.formGroup}>
                        <label style={styles.label}>Acción</label>
                        <input name="accion" value={formData.accion} onChange={handleInputChange} required style={styles.input} placeholder="Ej. Recolección" />
                    </div>

                    <div style={styles.rowGroup}>
                        <div style={{...styles.formGroup, flex: 1}}>
                            <label style={styles.label}>Tiempo</label>
                            <input name="tiempo" value={formData.tiempo} onChange={handleInputChange} required style={styles.input} placeholder="Ej. 8 horas" />
                        </div>
                        <div style={{...styles.formGroup, flex: 1}}>
                            <label style={styles.label}>Cantidad</label>
                            <input type="number" name="cantidad" value={formData.cantidad} onChange={handleInputChange} required style={styles.input} placeholder="0" />
                        </div>
                    </div>

                    <div style={styles.rowGroup}>
                        {/* ⭐️ NUEVO CAMPO: Costo ($) */}
                        <div style={{...styles.formGroup, flex: 1}}>
                            <label style={styles.label}>Costo Unitario ($)</label>
                            <input type="number" step="0.01" name="costo" value={formData.costo} onChange={handleInputChange} required style={styles.input} placeholder="0.00" />
                        </div>
                        <div style={{...styles.formGroup, flex: 1}}>
                            <label style={styles.label}>Responsable</label>
                            <input name="responsable" value={formData.responsable} onChange={handleInputChange} required style={styles.input} placeholder="Nombre" />
                        </div>
                    </div>

                    <div style={styles.formGroup}>
                        <label style={styles.label}>Monto Total ($)</label>
                        <input type="number" step="0.01" name="monto" value={formData.monto} onChange={handleInputChange} required style={styles.input} placeholder="0.00" />
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

export default RecursoHumano;