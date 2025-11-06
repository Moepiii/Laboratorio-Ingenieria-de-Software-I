import React, { useState, useMemo } from 'react';
import Modal from './Modal'; // Importa el Modal genérico

// Estilos para este componente
const styles = {
    searchInput: {
        width: '100%',
        padding: '0.75rem 1rem',
        border: '1px solid #d1d5db',
        borderRadius: '8px',
        fontSize: '1rem',
        boxSizing: 'border-box',
        marginBottom: '1rem',
    },
    // Estilos para la nueva tabla
    tableWrapper: {
        maxHeight: '300px',
        overflowY: 'auto',
        border: '1px solid #e5e7eb',
        borderRadius: '8px',
    },
    table: {
        width: '100%',
        borderCollapse: 'collapse',
    },
    th: {
        padding: '0.75rem 1rem',
        textAlign: 'left',
        backgroundColor: '#f3f4f6',
        borderBottom: '2px solid #e5e7eb',
        color: '#374151',
        fontWeight: '600',
    },
    td: {
        padding: '0.75rem 1rem',
        borderBottom: '1px solid #e5e7eb',
        verticalAlign: 'middle',
    },
    // ⭐️ --- (INICIO) NUEVO ESTILO --- ⭐️
    colCedula: {
        minWidth: '130px',      // Le damos un ancho mínimo de 130px
        whiteSpace: 'nowrap',   // Evita que el texto se parta en dos líneas
    },
    // ⭐️ --- (FIN) NUEVO ESTILO --- ⭐️
    trHover: {
        cursor: 'pointer',
        backgroundColor: '#f9fafb', // Fondo suave al pasar el mouse
    },
    radioOuter: {
        width: '18px',
        height: '18px',
        borderRadius: '50%',
        border: '2px solid #6b7280', // Gris
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
    radioOuterSelected: {
        borderColor: '#4f46e5', // Azul
    },
    radioInner: {
        width: '10px',
        height: '10px',
        borderRadius: '50%',
        backgroundColor: '#4f46e5', // Azul
    },
    noResults: {
        padding: '1rem',
        textAlign: 'center',
        color: '#6b7280',
    }
};

/**
 * Componente para el Modal de Búsqueda de Encargados.
 * @param {boolean} isOpen - Si el modal está abierto.
 * @param {function} onClose - Función para cerrar el modal.
 * @param {function} onSelect - Función que se llama con el encargado seleccionado.
 * @param {Array} encargadosList - La lista completa de encargados.
 * @param {number} selectedId - El ID del encargado actualmente seleccionado.
 */
const EncargadoSearchModal = ({ isOpen, onClose, onSelect, encargadosList, selectedId }) => {
    const [searchTerm, setSearchTerm] = useState('');
    const [hoveredId, setHoveredId] = useState(null); // Para el efecto hover

    // Filtra la lista
    const filteredEncargados = useMemo(() => {
        const term = searchTerm.toLowerCase();
        if (!term) {
            return encargadosList;
        }
        return encargadosList.filter(enc =>
            enc.nombre.toLowerCase().includes(term) ||
            enc.apellido.toLowerCase().includes(term) ||
            enc.cedula.toLowerCase().includes(term)
        );
    }, [searchTerm, encargadosList]);

    // Manejador de selección
    const handleSelect = (encargado) => {
        onSelect(encargado);
        onClose();
        setSearchTerm('');
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Buscar Encargado">
            <input
                type="text"
                style={styles.searchInput}
                placeholder="Buscar por nombre, apellido o cédula..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
            />
            <div style={styles.tableWrapper}>
                <table style={styles.table}>
                    <thead style={styles.thead}>
                        <tr>
                            <th style={styles.th}>Opciones</th>
                            {/* ⭐️ CAMBIO: Aplicamos el nuevo estilo a la cabecera */}
                            <th style={{ ...styles.th, ...styles.colCedula }}>Cédula</th>
                            <th style={styles.th}>Nombre</th>
                            <th style={styles.th}>Apellido</th>
                        </tr>
                    </thead>
                    <tbody>
                        {filteredEncargados.length > 0 ? (
                            filteredEncargados.map(enc => {
                                const isSelected = enc.id === selectedId;
                                const isHovered = enc.id === hoveredId;

                                return (
                                    <tr
                                        key={enc.id}
                                        style={isHovered ? { ...styles.tr, ...styles.trHover } : styles.tr}
                                        onMouseEnter={() => setHoveredId(enc.id)}
                                        onMouseLeave={() => setHoveredId(null)}
                                        onClick={() => handleSelect(enc)}
                                    >
                                        <td style={styles.td}>
                                            <div style={{ ...styles.radioOuter, ...(isSelected ? styles.radioOuterSelected : {}) }}>
                                                {isSelected && <div style={styles.radioInner}></div>}
                                            </div>
                                        </td>
                                        {/* ⭐️ CAMBIO: Aplicamos el nuevo estilo a la celda */}
                                        <td style={{ ...styles.td, ...styles.colCedula }}>{enc.cedula}</td>
                                        <td style={styles.td}>{enc.nombre}</td>
                                        <td style={styles.td}>{enc.apellido}</td>
                                    </tr>
                                );
                            })
                        ) : (
                            <tr>
                                <td colSpan="4" style={{ ...styles.td, ...styles.noResults }}>
                                    No se encontraron encargados.
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </Modal>
    );
};

export default EncargadoSearchModal;