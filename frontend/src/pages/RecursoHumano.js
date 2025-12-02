import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import Modal from '../components/auth/Modal';
import { useAuth } from '../context/AuthContext';
// Servicios
import { getDatosProyecto } from '../services/actividadService';
import { getRecursos, createRecurso, updateRecurso, deleteRecurso } from '../services/recursoService';

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

const RecursoHumano = () => {
    const { id } = useParams();
    const { token, currentUser } = useAuth();

    const [recursos, setRecursos] = useState([]);

    // Listas para los Selects
    const [listaActividades, setListaActividades] = useState([]);
    const [listaLabores, setListaLabores] = useState([]);
    const [listaEncargados, setListaEncargados] = useState([]);

    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingId, setEditingId] = useState(null);

    const [formData, setFormData] = useState({
        actividad: '',
        accion: '',
        nombre: '',
        cedula: '',
        tiempo: '',
        cantidad: '',
        costo_unitario: '',
        monto: ''
    });

    // Función para recargar la tabla de recursos guardados
    const refreshRecursos = () => {
        if (id && token && currentUser) {
            getRecursos(token, id, currentUser.username)
                .then(res => {
                    if (res && res.recursos) setRecursos(res.recursos);
                })
                .catch(err => console.error("Error cargando recursos:", err));
        }
    };

    useEffect(() => {
        if (id && token && currentUser?.username) {
            const pid = parseInt(id, 10);

            console.log("📡 Solicitando Datos del Proyecto ID:", pid);

            // 1. Cargar Configuración (Actividades, Labores, Encargados)
            getDatosProyecto(token, pid, currentUser.username)
                .then(res => {
                    console.log("✅ Datos recibidos:", res); // MIRA LA CONSOLA (F12) SI ESTO NO FUNCIONA
                    if (res) {
                        if (res.actividades) setListaActividades(res.actividades);
                        if (res.labores) setListaLabores(res.labores);
                        if (res.encargados) setListaEncargados(res.encargados);
                    }
                })
                .catch(err => console.error("Error cargando datos:", err));

            // 2. Cargar la tabla
            refreshRecursos();
        }
    }, [id, token, currentUser]);

    const handleOpenModal = () => {
        setEditingId(null);
        setFormData({ actividad: '', accion: '', nombre: '', cedula: '', tiempo: '', cantidad: '', costo_unitario: '', monto: '' });
        setIsModalOpen(true);
    };

    const handleCloseModal = () => setIsModalOpen(false);

    const handleEditClick = (rec) => {
        setEditingId(rec.id);
        setFormData({
            actividad: rec.actividad,
            accion: rec.accion,
            nombre: rec.nombre,
            cedula: rec.cedula,
            tiempo: rec.tiempo,
            cantidad: rec.cantidad,
            costo_unitario: rec.costo_unitario,
            monto: rec.monto
        });
        setIsModalOpen(true);
    };

    const handleDeleteClick = async (recId) => {
        if (window.confirm("¿Borrar este recurso?")) {
            try {
                await deleteRecurso(token, recId, currentUser.username);
                refreshRecursos();
            } catch (e) { alert(e.message); }
        }
    };

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        let newData = { ...formData, [name]: value };

        // ⭐️ Lógica para rellenar Cédula si seleccionan un Encargado de la lista
        if (name === 'nombre') {
            const encargadoSeleccionado = listaEncargados.find(enc => `${enc.nombre} ${enc.apellido}` === value);
            if (encargadoSeleccionado) {
                newData.cedula = encargadoSeleccionado.cedula;
            }
        }

        // ⭐️ FÓRMULA: Monto = Tiempo * Cantidad * Costo
        if (name === 'tiempo' || name === 'cantidad' || name === 'costo_unitario') {
            const t = parseFloat(name === 'tiempo' ? value : formData.tiempo) || 0;
            const c = parseFloat(name === 'cantidad' ? value : formData.cantidad) || 0;
            const p = parseFloat(name === 'costo_unitario' ? value : formData.costo_unitario) || 0;

            newData.monto = (t * c * p).toFixed(2);
        }

        setFormData(newData);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const dataToSend = { proyecto_id: id, ...formData };

            if (editingId) {
                await updateRecurso(token, { ...dataToSend, id: editingId }, currentUser.username);
            } else {
                await createRecurso(token, dataToSend, currentUser.username);
            }
            refreshRecursos();
            handleCloseModal();
        } catch (e) {
            alert("Error al guardar: " + e.message);
        }
    };

    return (
        <div style={styles.container}>
            <div style={styles.header}>
                <h2 style={styles.h2}>Recurso Humano</h2>
                <button style={styles.addButton} onClick={handleOpenModal}>+ Añadir</button>
            </div>

            <div style={styles.tableContainer}>
                <table style={styles.table}>
                    <thead>
                        <tr>
                            <th style={styles.th}>ID</th>
                            <th style={styles.th}>Actividad</th>
                            <th style={styles.th}>Acción</th>
                            <th style={styles.th}>Nombre</th>
                            <th style={styles.th}>Cédula</th>
                            <th style={styles.th}>Tiempo</th>
                            <th style={styles.th}>Cant.</th>
                            <th style={styles.th}>Costo Unit</th>
                            <th style={styles.th}>Total ($)</th>
                            <th style={styles.th}>Acciones</th>
                        </tr>
                    </thead>
                    <tbody>
                        {recursos.map((rec) => (
                            <tr key={rec.id}>
                                <td style={styles.td}>{rec.id}</td>
                                <td style={styles.td}>{rec.actividad}</td>
                                <td style={styles.td}>{rec.accion}</td>
                                <td style={styles.td}>{rec.nombre}</td>
                                <td style={styles.td}>{rec.cedula}</td>
                                <td style={styles.td}>{rec.tiempo}</td>
                                <td style={styles.td}>{rec.cantidad}</td>
                                <td style={styles.td}>{rec.costo_unitario}</td>
                                <td style={{ ...styles.td, fontWeight: 'bold' }}>{rec.monto}</td>
                                <td style={styles.td}>
                                    <button style={{ ...styles.actionButton, ...styles.editButton }} onClick={() => handleEditClick(rec)}>Editar</button>
                                    <button style={{ ...styles.actionButton, ...styles.deleteButton }} onClick={() => handleDeleteClick(rec.id)}>Borrar</button>
                                </td>
                            </tr>
                        ))}
                        {recursos.length === 0 && (
                            <tr><td colSpan="10" style={{ padding: '2rem', textAlign: 'center' }}>No hay recursos registrados.</td></tr>
                        )}
                    </tbody>
                </table>
            </div>

            <Modal isOpen={isModalOpen} onClose={handleCloseModal} title={editingId ? "Editar Recurso" : "Nuevo Recurso Humano"}>
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
                        {listaActividades.length === 0 && <small style={{ color: 'red' }}>No hay actividades registradas</small>}
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
                        {listaLabores.length === 0 && <small style={{ color: 'red' }}>No hay labores configuradas</small>}
                    </div>

                    <div style={styles.rowGroup}>
                        {/* ⭐️ CAMBIO: Ahora es un SELECT para elegir Encargados */}
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Nombre / Responsable</label>
                            <select name="nombre" value={formData.nombre} onChange={handleInputChange} required style={styles.select}>
                                <option value="">-- Seleccione Persona --</option>
                                {listaEncargados.map(e => (
                                    <option key={e.id} value={`${e.nombre} ${e.apellido}`}>
                                        {e.nombre} {e.apellido}
                                    </option>
                                ))}
                            </select>
                            {listaEncargados.length === 0 && <small style={{ color: 'orange' }}>No hay encargados registrados</small>}
                        </div>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Cédula (Auto)</label>
                            <input name="cedula" value={formData.cedula} onChange={handleInputChange} style={{ ...styles.input, backgroundColor: '#f9fafb' }} readOnly />
                        </div>
                    </div>

                    <div style={styles.rowGroup}>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Tiempo (Días/Horas)</label>
                            <input type="number" name="tiempo" value={formData.tiempo} onChange={handleInputChange} required style={styles.input} placeholder="0" />
                        </div>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Cantidad (Personas)</label>
                            <input type="number" name="cantidad" value={formData.cantidad} onChange={handleInputChange} required style={styles.input} placeholder="0" />
                        </div>
                    </div>

                    <div style={styles.rowGroup}>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Costo Unitario ($)</label>
                            <input type="number" step="0.01" name="costo_unitario" value={formData.costo_unitario} onChange={handleInputChange} required style={styles.input} placeholder="0.00" />
                        </div>
                        <div style={{ ...styles.formGroup, flex: 1 }}>
                            <label style={styles.label}>Monto Total ($)</label>
                            <input type="number" name="monto" value={formData.monto} readOnly style={{ ...styles.input, backgroundColor: '#f3f4f6' }} />
                        </div>
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