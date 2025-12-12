import React from 'react';

// Estilos para el Modal (método inline)
const styles = {
    overlay: {
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: 'rgba(0, 0, 0, 0.7)', // Fondo oscuro semitransparente
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: 1000,
    },
    modal: {
        backgroundColor: 'white',
        padding: '2rem',
        borderRadius: '8px',
        width: '90%',
        maxWidth: '600px',
        maxHeight: '90vh', // Alto máximo
        overflowY: 'auto',  // Scroll si el contenido es muy largo
        boxShadow: '0 10px 25px rgba(0, 0, 0, 0.1)',
        position: 'relative',
    },
    closeButton: {
        position: 'absolute',
        top: '1rem',
        right: '1rem',
        background: 'transparent',
        border: 'none',
        fontSize: '1.5rem',
        cursor: 'pointer',
        color: '#6b7280',
    },
    header: {
        fontSize: '1.5rem',
        fontWeight: '600',
        color: '#1f2937',
        marginTop: 0,
        marginBottom: '1.5rem',
    }
};


const Modal = ({ isOpen, onClose, title, children }) => {
    if (!isOpen) {
        return null;
    }

    const handleOverlayClick = (e) => {
        if (e.target === e.currentTarget) {
            onClose();
        }
    };

    return (
        <div style={styles.overlay} onClick={handleOverlayClick}>
            <div style={styles.modal}>
                <button style={styles.closeButton} onClick={onClose}>&times;</button>
                {title && <h2 style={styles.header}>{title}</h2>}
                {children}
            </div>
        </div>
    );
};

export default Modal;