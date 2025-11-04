import React from 'react';
import { useParams } from 'react-router-dom';

const LaboresAgronomicas = () => {
    const { id } = useParams(); // Obtiene el ID del proyecto de la URL
    return (
        <div style={{ padding: '2rem', color: '#333' }}>
            <h2>Labores Agronómicas</h2>
            <p>Mostrando labores para el Proyecto ID: <strong>{id}</strong></p>
            {/* Aquí iría tu futura tabla o lógica para las labores */}
        </div>
    );
};

export default LaboresAgronomicas;