import React from 'react';
import { useParams } from 'react-router-dom';

const EquiposEImplementos = () => {
    const { id } = useParams(); // Obtiene el ID del proyecto de la URL
    return (
        <div style={{ padding: '2rem', color: '#333' }}>
            <h2>Equipos e Implementos</h2>
            <p>Mostrando equipos para el Proyecto ID: <strong>{id}</strong></p>
            {/* Aquí iría tu futura tabla o lógica para los equipos */}
        </div>
    );
};

export default EquiposEImplementos;