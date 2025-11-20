import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../context/AuthContext';
import { getUnidades, createUnidad, updateUnidad, deleteUnidad } from '../services/unidadService';
import Modal from '../components/auth/Modal';

// Estilos idénticos a Portafolio.js
const styles = {
    container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
    header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '2px solid #e5e7eb', paddingBottom: '1rem', marginBottom: '2rem' },
    h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', margin: 0 },
    addButton: { padding: '0.75rem 1.5rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#2563eb', border: 'none', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '0.5rem' },
    tableContainer: { overflowX: 'auto', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
    table: { width: '100%', borderCollapse: 'collapse' },
    th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '2px solid #e5e7eb', color: '#374151', fontWeight: '600', fontSize: '0.875rem' },
    td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', color: '#4b5563', fontSize: '0.875rem' },
    
    // Formulario Modal
    form: { display: 'flex', flexDirection: 'column', gap: '1rem' },
    formGroup: { display: 'flex', flexDirection: 'column', gap: '0.5rem' },
    label: { fontSize: '0.875rem', fontWeight: '500', color: '#374151' },
    input: { padding: '0.6rem', border: '1px solid #d1d5db', borderRadius: '6px', fontSize: '1rem' },
    select: { padding: '0.6rem', border: '1px solid #d1d5db', borderRadius: '6px', fontSize: '1rem', backgroundColor: 'white' },
    formActions: { display: 'flex', justifyContent: 'flex-end', gap: '1rem', marginTop: '1.5rem' },
    cancelButton: { padding: '0.6rem 1.2rem', border: '1px solid #d1d5db', borderRadius: '6px', backgroundColor: 'white', color: '#374151', cursor: 'pointer' },
    saveButton: { padding: '0.6rem 1.2rem', border: 'none', borderRadius: '6px', backgroundColor: '#2563eb', color: 'white', cursor: 'pointer' },
    
    actionButton: { padding: '0.4rem 0.8rem', borderRadius: '4px', border: 'none', cursor: 'pointer', fontSize: '0.875rem', marginLeft: '0.5rem', color: 'white' },
    editButton: { backgroundColor: '#f59e0b' },
    deleteButton: { backgroundColor: '#ef4444' },
};

const UnidadesMedida = () => {
    const { token, currentUser } = useAuth();
    const [unidades, setUnidades] = useState([]);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [currentUnidad, setCurrentUnidad] = useState(null);

    // ⭐️ NUEVO ESTADO: dimension
    const [formData, setFormData] = useState({ nombre: '', abreviatura: '', tipo: 'Peso', dimension: '' });

    const fetchUnidades = useCallback(async () => {
        try {
            const res = await getUnidades(token);
            setUnidades(res || []); 
        } catch (error) {
            console.error("Error cargando unidades:", error);
        }
    }, [token]);

    useEffect(() => { fetchUnidades(); }, [fetchUnidades]);

    const handleOpenModal = (unidad = null) => {
        if (unidad) {
            setCurrentUnidad(unidad);
            // Cargamos los datos incluyendo la dimensión
            setFormData({ 
                nombre: unidad.nombre, 
                abreviatura: unidad.abreviatura, 
                tipo: unidad.tipo,
                dimension: unidad.dimension // ⭐️ Carga valor existente
            });
        } else {
            setCurrentUnidad(null);
            setFormData({ nombre: '', abreviatura: '', tipo: 'Peso', dimension: '' });
        }
        setIsModalOpen(true);
    };

    const handleCloseModal = () => setIsModalOpen(false);

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        
        // Convertimos dimension a float para enviarlo correctamente
        const dataToSend = {
            ...formData,
            dimension: parseFloat(formData.dimension) || 0
        };

        try {
            if (currentUnidad) {
                await updateUnidad(token, { ...dataToSend, id: currentUnidad.id }, currentUser.username);
            } else {
                await createUnidad(token, dataToSend, currentUser.username);
            }
            fetchUnidades();
            handleCloseModal();
        } catch (error) {
            alert("Error: " + error.message);
        }
    };

    const handleDelete = async (id) => {
        if (window.confirm("¿Estás seguro de eliminar esta unidad?")) {
            try {
                await deleteUnidad(token, id, currentUser.username);
                fetchUnidades();
            } catch (error) {
                alert("Error: " + error.message);
            }
        }
    };

    return (
        <div style={styles.container}>
            <div style={styles.header}>
                <h2 style={styles.h2}>Unidades de Medida</h2>
                <button style={styles.addButton} onClick={() => handleOpenModal()}>
                    <span>+</span> Añadir Unidad
                </button>
            </div>

            <div style={styles.tableContainer}>
                <table style={styles.table}>
                    <thead>
                        <tr>
                            {/* ⭐️ NUEVA COLUMNA: ID */}
                            <th style={styles.th}>ID</th>
                            <th style={styles.th}>Nombre</th>
                            <th style={styles.th}>Abreviatura</th>
                            <th style={styles.th}>Tipo</th>
                            {/* ⭐️ NUEVA COLUMNA: DIMENSIÓN */}
                            <th style={styles.th}>Dimensión</th>
                            <th style={styles.th}>Acciones</th>
                        </tr>
                    </thead>
                    <tbody>
                        {unidades.length > 0 ? unidades.map((u) => (
                            <tr key={u.id}>
                                {/* ⭐️ MUESTRA ID */}
                                <td style={styles.td}>{u.id}</td>
                                <td style={styles.td}>{u.nombre}</td>
                                <td style={styles.td}>{u.abreviatura}</td>
                                <td style={styles.td}>{u.tipo}</td>
                                {/* ⭐️ MUESTRA DIMENSIÓN */}
                                <td style={styles.td}>{u.dimension}</td>
                                <td style={styles.td}>
                                    <button style={{...styles.actionButton, ...styles.editButton}} onClick={() => handleOpenModal(u)}>Editar</button>
                                    <button style={{...styles.actionButton, ...styles.deleteButton}} onClick={() => handleDelete(u.id)}>Borrar</button>
                                </td>
                            </tr>
                        )) : (
                            <tr><td colSpan="6" style={{...styles.td, textAlign:'center'}}>No hay unidades registradas.</td></tr>
                        )}
                    </tbody>
                </table>
            </div>

            <Modal isOpen={isModalOpen} onClose={handleCloseModal} title={currentUnidad ? "Editar Unidad" : "Añadir Unidad"}>
                <form onSubmit={handleSubmit} style={styles.form}>
                    <div style={styles.formGroup}>
                        <label htmlFor="nombre" style={styles.label}>Nombre</label>
                        <input 
                            type="text" id="nombre" name="nombre" 
                            value={formData.nombre} onChange={handleChange} required 
                            style={styles.input} placeholder="Ej. Kilogramo" 
                        />
                    </div>
                    <div style={styles.formGroup}>
                        <label htmlFor="abreviatura" style={styles.label}>Abreviatura</label>
                        <input 
                            type="text" id="abreviatura" name="abreviatura" 
                            value={formData.abreviatura} onChange={handleChange} required 
                            style={styles.input} placeholder="Ej. kg" 
                        />
                    </div>
                    
                    {/* ⭐️ NUEVO INPUT: DIMENSIÓN */}
                    <div style={styles.formGroup}>
                        <label htmlFor="dimension" style={styles.label}>Dimensión (Número Real)</label>
                        <input 
                            type="number" 
                            step="any" // Permite decimales
                            id="dimension" name="dimension" 
                            value={formData.dimension} onChange={handleChange} required 
                            style={styles.input} placeholder="Ej. 10.5" 
                        />
                    </div>

                    <div style={styles.formGroup}>
                        <label htmlFor="tipo" style={styles.label}>Tipo</label>
                        <select 
                            id="tipo" name="tipo" 
                            value={formData.tipo} onChange={handleChange} 
                            style={styles.select}
                        >
                            <option value="Peso">Peso</option>
                            <option value="Líquido">Líquido</option>
                            <option value="Longitud">Longitud</option>
                        </select>
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

export default UnidadesMedida;