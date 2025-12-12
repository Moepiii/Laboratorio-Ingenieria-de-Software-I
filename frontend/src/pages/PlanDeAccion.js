import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import Modal from '../components/auth/Modal';
import { useAuth } from '../context/AuthContext';
import { getDatosProyecto } from '../services/actividadService';
// Importamos update y delete
import { getPlanes, createPlan, updatePlan, deletePlan } from '../services/planService';

const styles = {
    container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
    header: { display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '2px solid #e5e7eb', paddingBottom: '1rem', marginBottom: '2rem' },
    h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', margin: 0 },
    addButton: { padding: '0.75rem 1.5rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#2563eb', border: 'none', cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '0.5rem' },
    tableContainer: { overflowX: 'auto', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
    table: { width: '100%', borderCollapse: 'collapse', minWidth: '900px' },
    th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '2px solid #e5e7eb', color: '#374151', fontWeight: '600', fontSize: '0.875rem', whiteSpace: 'nowrap' },
    td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', color: '#4b5563', fontSize: '0.875rem' },

    // Estilo para el Total
    footerTd: { padding: '1rem', borderTop: '2px solid #e5e7eb', color: '#111827', fontSize: '1rem', fontWeight: '700', backgroundColor: '#f9fafb' },

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

const PlanDeAccion = () => {
    const { id } = useParams();
    const { token, currentUser } = useAuth();

    const [planes, setPlanes] = useState([]);
    const [listaActividadesOrigen, setListaActividadesOrigen] = useState([]);
    const [listaLabores, setListaLabores] = useState([]);
    const [listaEncargados, setListaEncargados] = useState([]);

    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingPlanId, setEditingPlanId] = useState(null);

    const [formData, setFormData] = useState({
        actividad: '',
        accion: '',
        fecha_inicio: '',
        fecha_cierre: '',
        horas: '',
        responsable: '',
        costo_unitario: '',
        monto: ''
    });

    // FUNCIÓN PARA FORMATEAR FECHA (YYYY-MM-DD -> DD-MM-YYYY)
    const formatearFecha = (fechaISO) => {
        if (!fechaISO) return '';
        const soloFecha = fechaISO.split('T')[0];
        const [year, month, day] = soloFecha.split('-');
        return `${day}-${month}-${year}`;
    };

    // CÁLCULO DEL TOTAL GENERAL
    const totalMonto = planes.reduce((acc, plan) => {
        return acc + (parseFloat(plan.monto) || 0);
    }, 0);

    const refreshPlanes = () => {
        if (id && token && currentUser) {
            getPlanes(token, id, currentUser.username)
                .then(res => { if (res && res.planes) setPlanes(res.planes); })
                .catch(err => console.error(err));
        }
    };

    useEffect(() => {
        if (id && token && currentUser?.username) {
            const pid = parseInt(id, 10);
            getDatosProyecto(token, pid, currentUser.username)
                .then(res => {
                    if (res) {
                        if (res.actividades) setListaActividadesOrigen(res.actividades);
                        if (res.labores) setListaLabores(res.labores);
                        if (res.encargados) setListaEncargados(res.encargados);
                    }
                });
            refreshPlanes();
        }
    }, [id, token, currentUser]);

    const handleOpenModal = () => {
        setEditingPlanId(null);
        setFormData({ actividad: '', accion: '', fecha_inicio: '', fecha_cierre: '', horas: '', responsable: '', costo_unitario: '', monto: '' });
        setIsModalOpen(true);
    };

    const handleCloseModal = () => setIsModalOpen(false);

    const handleEditClick = (plan) => {
        setEditingPlanId(plan.id);
        setFormData({
            actividad: plan.actividad,
            accion: plan.accion,
            fecha_inicio: plan.fecha_inicio.split('T')[0],
            fecha_cierre: plan.fecha_cierre.split('T')[0],
            horas: plan.horas,
            responsable: plan.responsable,
            costo_unitario: plan.costo_unitario,
            monto: plan.monto
        });
        setIsModalOpen(true);
    };

    const handleDeleteClick = async (planId) => {
        if (window.confirm("¿Seguro que quieres borrar este plan?")) {
            try {
                await deletePlan(token, planId, currentUser.username);
                refreshPlanes();
            } catch (error) {
                alert("Error al borrar: " + error.message);
            }
        }
    };

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        let newData = { ...formData, [name]: value };

        if (name === 'horas' || name === 'costo_unitario') {
            const hrs = parseFloat(name === 'horas' ? value : formData.horas) || 0;
            const unit = parseFloat(name === 'costo_unitario' ? value : formData.costo_unitario) || 0;
            newData.monto = (hrs * unit).toFixed(2);
        }
        setFormData(newData);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const planData = { proyecto_id: id, ...formData };
            if (editingPlanId) {
                await updatePlan(token, { ...planData, id: editingPlanId }, currentUser.username);
            } else {
                await createPlan(token, planData, currentUser.username);
            }
            refreshPlanes();
            handleCloseModal();
        } catch (error) {
            alert("Error al guardar: " + error.message);
        }
    };

    return (
        <div style={styles.container}>
            <div style={styles.header}>
                <h2 style={styles.h2}>Plan de Acción (General)</h2>
                <button style={styles.addButton} onClick={handleOpenModal}>+ Añadir</button>
            </div>

            <div style={styles.tableContainer}>
                <table style={styles.table}>
                    <thead>
                        <tr>
                            <th style={styles.th}>ID</th>
                            <th style={styles.th}>Actividad</th>
                            {/*  */}
                            <th style={styles.th}>Acción</th>
                            <th style={styles.th}>Fecha Inicio</th>
                            <th style={styles.th}>Fecha Cierre</th>
                            {/*  */}
                            <th style={styles.th}>Cantidad Horas</th>
                            <th style={styles.th}>Responsable</th>
                            {/*  */}
                            <th style={styles.th}>Monto ($)</th>
                            <th style={styles.th}>Acciones</th>
                        </tr>
                    </thead>
                    <tbody>
                        {planes.map((plan) => (
                            <tr key={plan.id}>
                                <td style={styles.td}>{plan.id}</td>
                                <td style={styles.td}>{plan.actividad}</td>
                                <td style={styles.td}>{plan.accion}</td>
                                <td style={styles.td}>{formatearFecha(plan.fecha_inicio)}</td>
                                <td style={styles.td}>{formatearFecha(plan.fecha_cierre)}</td>
                                <td style={styles.td}>{plan.horas}</td>
                                <td style={styles.td}>{plan.responsable}</td>
                                <td style={{ ...styles.td, fontWeight: 'bold' }}>{plan.monto}</td>
                                <td style={styles.td}>
                                    <button style={{ ...styles.actionButton, ...styles.editButton }} onClick={() => handleEditClick(plan)}>Editar</button>
                                    <button style={{ ...styles.actionButton, ...styles.deleteButton }} onClick={() => handleDeleteClick(plan.id)}>Borrar</button>
                                </td>
                            </tr>
                        ))}
                        {planes.length === 0 && (
                            <tr><td colSpan="9" style={{ padding: '2rem', textAlign: 'center', color: '#6b7280' }}>No hay planes registrados. Añade uno nuevo.</td></tr>
                        )}
                    </tbody>

                    {planes.length > 0 && (
                        <tfoot>
                            <tr>
                                {/* ColSpan ajustado para alineación */}
                                <td colSpan="7" style={{ ...styles.footerTd, textAlign: 'right' }}>
                                    Monto Total Invertido ($):
                                </td>
                                <td style={{ ...styles.footerTd, color: '#2563eb' }}>
                                    ${totalMonto.toFixed(2)}
                                </td>
                                <td style={styles.footerTd}></td>
                            </tr>
                        </tfoot>
                    )}

                </table>
            </div>

            <Modal isOpen={isModalOpen} onClose={handleCloseModal} title={editingPlanId ? "Editar Plan" : "Nuevo Plan de Acción"}>
                <form onSubmit={handleSubmit}>

                    <div style={styles.formGroup}>
                        <label style={styles.label}>Actividad</label>
                        <select name="actividad" value={formData.actividad} onChange={handleInputChange} required style={styles.select}>
                            <option value="">-- Seleccione --</option>
                            {listaActividadesOrigen.map((act) => (
                                <option key={act.id} value={act.actividad}>{act.actividad}</option>
                            ))}
                        </select>
                    </div>

                    <div style={styles.formGroup}>
                        {/*  */}
                        <label style={styles.label}>Acción</label>
                        <select name="accion" value={formData.accion} onChange={handleInputChange} required style={styles.select}>
                            <option value="">-- Seleccione --</option>
                            {listaLabores.map((labor) => (
                                <option key={labor.id} value={labor.descripcion}>{labor.descripcion}</option>
                            ))}
                        </select>
                    </div>

                    <div style={styles.rowGroup}>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Fecha Inicio</label>
                            <input type="date" name="fecha_inicio" value={formData.fecha_inicio} onChange={handleInputChange} required style={styles.input} />
                        </div>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Fecha Cierre</label>
                            <input type="date" name="fecha_cierre" value={formData.fecha_cierre} onChange={handleInputChange} required style={styles.input} />
                        </div>
                    </div>

                    <div style={styles.rowGroup}>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            {/*  */}
                            <label style={styles.label}>Cantidad Horas</label>
                            <input type="number" name="horas" value={formData.horas} onChange={handleInputChange} required style={styles.input} placeholder="0" />
                        </div>

                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Responsable</label>
                            <select name="responsable" value={formData.responsable} onChange={handleInputChange} required style={styles.select}>
                                <option value="">-- Seleccione --</option>
                                {listaEncargados.map((user) => (
                                    <option key={user.id} value={`${user.nombre} ${user.apellido}`}>
                                        {user.nombre} {user.apellido}
                                    </option>
                                ))}
                            </select>
                        </div>
                    </div>

                    <div style={styles.rowGroup}>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Dinero ($)</label>
                            <input type="number" step="0.01" name="costo_unitario" value={formData.costo_unitario} onChange={handleInputChange} required style={styles.input} placeholder="0.00" />
                        </div>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Monto ($)</label>
                            <input type="number" name="monto" value={formData.monto} readOnly style={{ ...styles.input, backgroundColor: '#f3f4f6' }} placeholder="0.00" />
                        </div>
                    </div>

                    <div style={styles.formActions}>
                        <button type="button" onClick={handleCloseModal} style={styles.cancelButton}>Cancelar</button>
                        <button type="submit" style={styles.saveButton}>{editingPlanId ? "Actualizar" : "Guardar"}</button>
                    </div>
                </form>
            </Modal>
        </div>
    );
};

export default PlanDeAccion;