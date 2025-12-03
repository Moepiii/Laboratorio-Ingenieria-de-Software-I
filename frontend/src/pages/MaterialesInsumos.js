import React, { useState, useEffect, useMemo } from 'react';
import { useParams } from 'react-router-dom';
import Modal from '../components/auth/Modal';
import { useAuth } from '../context/AuthContext';
// Servicios
import { getDatosProyecto } from '../services/actividadService';
import { getMateriales, createMaterial, updateMaterial, deleteMaterial } from '../services/materialService';

const styles = {
    container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
    header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '2px solid #e5e7eb', paddingBottom: '1rem', marginBottom: '2rem' },
    h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', margin: 0 },
    addButton: { padding: '0.75rem 1.5rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#2563eb', border: 'none', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '0.5rem' },
    tableContainer: { overflowX: 'auto', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
    table: { width: '100%', borderCollapse: 'collapse', minWidth: '1000px' },
    th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '2px solid #e5e7eb', color: '#374151', fontWeight: '600', fontSize: '0.875rem', whiteSpace: 'nowrap' },
    td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', color: '#4b5563', fontSize: '0.875rem' },

    actionButton: { padding: '0.4rem 0.8rem', fontSize: '0.85rem', fontWeight: '600', borderRadius: '6px', border: 'none', cursor: 'pointer', marginLeft: '0.5rem' },
    editButton: { backgroundColor: '#f59e0b', color: 'white' },
    deleteButton: { backgroundColor: '#ef4444', color: 'white' },

    formGroup: { marginBottom: '1rem' },
    label: { display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' },
    input: { width: '100%', padding: '0.6rem', border: '1px solid #d1d5db', borderRadius: '6px', fontSize: '0.9rem', boxSizing: 'border-box' },
    select: { width: '100%', padding: '0.6rem', border: '1px solid #d1d5db', borderRadius: '6px', fontSize: '0.9rem', backgroundColor: 'white' },
    rowGroup: { display: 'flex', gap: '1rem', marginBottom: '0.5rem' },
    formActions: { display: 'flex', justifyContent: 'flex-end', gap: '1rem', marginTop: '1.5rem' },
    cancelButton: { padding: '0.5rem 1rem', backgroundColor: '#e5e7eb', color: '#374151', border: 'none', borderRadius: '6px', cursor: 'pointer' },
    saveButton: { padding: '0.5rem 1rem', backgroundColor: '#2563eb', color: 'white', border: 'none', borderRadius: '6px', cursor: 'pointer' }
};

const MaterialesInsumos = () => {
    const { id } = useParams();
    const { token, currentUser } = useAuth();

    const [materiales, setMateriales] = useState([]);

    // Listas para los Selects
    const [listaActividades, setListaActividades] = useState([]);
    const [listaLabores, setListaLabores] = useState([]);
    // Lista de Responsables
    const [listaResponsables, setListaResponsables] = useState([]);

    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingId, setEditingId] = useState(null);

    const [formData, setFormData] = useState({
        actividad: '',
        accion: '',
        categoria: '',
        responsable: '',
        nombre: '',
        unidad: '',
        cantidad: '',
        costo_unitario: '',
        monto: ''
    });

    const refreshMateriales = () => {
        if (id && token && currentUser) {
            getMateriales(token, id, currentUser.username)
                .then(res => {
                    if (res && res.materiales) setMateriales(res.materiales);
                })
                .catch(err => console.error("Error cargando materiales:", err));
        }
    };

    const totalGeneral = useMemo(() => {
        if (!materiales) return 0;
        return materiales.reduce((acc, item) => {
            return acc + (parseFloat(item.monto) || 0);
        }, 0);
    }, [materiales]);

    useEffect(() => {
        if (id && token && currentUser?.username) {
            const pid = parseInt(id, 10);

            getDatosProyecto(token, pid, currentUser.username)
                .then(res => {
                    if (res) {
                        if (res.actividades) setListaActividades(res.actividades);
                        if (res.labores) setListaLabores(res.labores);
                        // Cargar Encargados como Responsables
                        if (res.encargados) setListaResponsables(res.encargados);
                    }
                })
                .catch(err => console.error("Error cargando datos:", err));

            refreshMateriales();
        }
    }, [id, token, currentUser]);

    const handleOpenModal = () => {
        setEditingId(null);
        // Reseteamos el formulario incluyendo responsable
        setFormData({ actividad: '', accion: '', categoria: '', responsable: '', nombre: '', unidad: '', cantidad: '', costo_unitario: '', monto: '' });
        setIsModalOpen(true);
    };

    const handleCloseModal = () => setIsModalOpen(false);

    const handleEditClick = (mat) => {
        setEditingId(mat.id);
        setFormData({
            actividad: mat.actividad,
            accion: mat.accion,
            categoria: mat.categoria,
            responsable: mat.responsable || '',
            nombre: mat.nombre,
            unidad: mat.unidad,
            cantidad: mat.cantidad,
            costo_unitario: mat.costo_unitario,
            monto: mat.monto
        });
        setIsModalOpen(true);
    };

    const handleDeleteClick = async (matId) => {
        if (window.confirm("¿Borrar este material?")) {
            try {
                await deleteMaterial(token, matId, currentUser.username);
                refreshMateriales();
            } catch (e) { alert(e.message); }
        }
    };

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        let newData = { ...formData, [name]: value };

        // FÓRMULA: Monto = Cantidad * Costo
        if (name === 'cantidad' || name === 'costo_unitario') {
            const c = parseFloat(name === 'cantidad' ? value : formData.cantidad) || 0;
            const p = parseFloat(name === 'costo_unitario' ? value : formData.costo_unitario) || 0;
            newData.monto = (c * p).toFixed(2);
        }

        setFormData(newData);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            // CORRECCIÓN CRÍTICA: Convertir tipos antes de enviar
            // Esto evita el Error 400 en Go
            const payload = {
                ...formData,
                proyecto_id: parseInt(id), // ID del proyecto como entero
                cantidad: parseFloat(formData.cantidad) || 0,
                costo_unitario: parseFloat(formData.costo_unitario) || 0,
                monto: parseFloat(formData.monto) || 0
            };

            if (editingId) {
                await updateMaterial(token, { ...payload, id: editingId }, currentUser.username);
            } else {
                await createMaterial(token, payload, currentUser.username);
            }
            refreshMateriales();
            handleCloseModal();
        } catch (e) {
            console.error(e);
            alert("Error al guardar: " + e.message);
        }
    };

    return (
        <div style={styles.container}>
            <div style={styles.header}>
                <h2 style={styles.h2}>Materiales e Insumos</h2>
                <button style={styles.addButton} onClick={handleOpenModal}>+ Añadir</button>
            </div>

            <div style={styles.tableContainer}>
                <table style={styles.table}>
                    <thead>
                        <tr>
                            <th style={styles.th}>ID</th>
                            <th style={styles.th}>Actividad</th>
                            <th style={styles.th}>Acción</th>
                            <th style={styles.th}>Categoría</th>
                            <th style={styles.th}>Responsable</th>
                            <th style={styles.th}>Descripción</th>
                            <th style={styles.th}>Medida</th>
                            <th style={styles.th}>Cantidad</th>
                            <th style={styles.th}>Monto ($)</th>
                            <th style={styles.th}>Acciones</th>
                        </tr>
                    </thead>
                    <tbody>
                        {materiales.map((mat) => (
                            <tr key={mat.id}>
                                <td style={styles.td}>{mat.id}</td>
                                <td style={styles.td}>{mat.actividad}</td>
                                <td style={styles.td}>{mat.accion}</td>
                                <td style={styles.td}>{mat.categoria}</td>
                                <td style={styles.td}>{mat.responsable}</td>
                                <td style={styles.td}>{mat.nombre}</td>
                                <td style={styles.td}>{mat.unidad}</td>
                                <td style={styles.td}>{mat.cantidad}</td>
                                <td style={{ ...styles.td, fontWeight: 'bold' }}>{mat.monto}</td>
                                <td style={styles.td}>
                                    <button style={{ ...styles.actionButton, ...styles.editButton }} onClick={() => handleEditClick(mat)}>Editar</button>
                                    <button style={{ ...styles.actionButton, ...styles.deleteButton }} onClick={() => handleDeleteClick(mat.id)}>Borrar</button>
                                </td>
                            </tr>
                        ))}
                        {materiales.length === 0 && (
                            <tr><td colSpan="10" style={{ padding: '2rem', textAlign: 'center' }}>No hay materiales registrados.</td></tr>
                        )}
                    </tbody>
                    <tfoot>
                        <tr style={{ backgroundColor: '#f3f4f6', borderTop: '2px solid #e5e7eb' }}>
                            <td colSpan={8} style={{ ...styles.td, textAlign: 'right', fontWeight: 'bold', fontSize: '1.1rem' }}>
                                Monto Total Materiales e Insumos ($):
                            </td>
                            <td style={{ ...styles.td, fontWeight: 'bold', fontSize: '1.1rem', color: '#2563eb' }}>
                                ${totalGeneral.toFixed(2)}
                            </td>
                            <td></td>
                        </tr>
                    </tfoot>
                </table>
            </div>

            <Modal isOpen={isModalOpen} onClose={handleCloseModal} title={editingId ? "Editar Material" : "Nuevo Material / Insumo"}>
                <form onSubmit={handleSubmit}>

                    {/* SELECTOR ACTIVIDAD */}
                    <div style={styles.formGroup}>
                        <label style={styles.label}>Actividad (Del Proyecto)</label>
                        <select name="actividad" value={formData.actividad} onChange={handleInputChange} required style={styles.select}>
                            <option value="">-- Seleccione Actividad --</option>
                            {listaActividades.map(a => (
                                <option key={a.id} value={a.actividad}>{a.actividad}</option>
                            ))}
                        </select>
                    </div>

                    {/* SELECTOR ACCIÓN (LABOR) */}
                    <div style={styles.formGroup}>
                        <label style={styles.label}>Acción (Labor Agronómica)</label>
                        <select name="accion" value={formData.accion} onChange={handleInputChange} required style={styles.select}>
                            <option value="">-- Seleccione Labor --</option>
                            {listaLabores.map(l => (
                                <option key={l.id} value={l.descripcion}>{l.descripcion}</option>
                            ))}
                        </select>
                    </div>

                    {/* SELECTOR CATEGORÍA */}
                    <div style={styles.formGroup}>
                        <label style={styles.label}>Categoría</label>
                        <select name="categoria" value={formData.categoria} onChange={handleInputChange} required style={styles.select}>
                            <option value="">-- Seleccione Categoría --</option>
                            <option value="Ninguno">Ninguno</option>
                            <option value="Materiales">Materiales</option>
                            <option value="Insumos">Insumos</option>
                            <option value="Equipos">Equipos</option>
                        </select>
                    </div>

                    {/* SELECTOR RESPONSABLE */}
                    <div style={styles.formGroup}>
                        <label style={styles.label}>Responsable (Asignado)</label>
                        <select name="responsable" value={formData.responsable} onChange={handleInputChange} required style={styles.select}>
                            <option value="">-- Seleccione Responsable --</option>
                            {listaResponsables.map((r, index) => (
                                <option key={index} value={r.nombre}>{r.nombre}</option>
                            ))}
                        </select>
                    </div>

                    <div style={styles.rowGroup}>
                        <div style={{ ...styles.formGroup, flex: 2 }}>
                            <label style={styles.label}>Descripción</label>
                            <input name="nombre" value={formData.nombre} onChange={handleInputChange} required style={styles.input} placeholder="Ej. Urea, Glifosato" />
                        </div>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Medida</label>
                            <input name="unidad" value={formData.unidad} onChange={handleInputChange} required style={styles.input} placeholder="Kg, Lts" />
                        </div>
                    </div>

                    <div style={styles.rowGroup}>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Cantidad</label>
                            <input type="number" name="cantidad" value={formData.cantidad} onChange={handleInputChange} required style={styles.input} placeholder="0" />
                        </div>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Costo ($)</label>
                            <input type="number" step="0.01" name="costo_unitario" value={formData.costo_unitario} onChange={handleInputChange} required style={styles.input} placeholder="0.00" />
                        </div>
                    </div>

                    <div style={styles.formGroup}>
                        <label style={styles.label}>Monto ($)</label>
                        <input type="number" name="monto" value={formData.monto} readOnly style={{ ...styles.input, backgroundColor: '#f3f4f6' }} />
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

export default MaterialesInsumos;