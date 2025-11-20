import { apiCall } from './authService';

// ⭐️ Ahora requiere proyectoId
export const getUnidades = (token, proyectoId, adminUsername) => {
    const body = { 
        proyecto_id: parseInt(proyectoId), // Aseguramos que sea número
        admin_username: adminUsername 
    };
    return apiCall('/admin/get-unidades', 'POST', body, token);
};

export const createUnidad = (token, unidadData, adminUsername) => {
    const body = { 
        ...unidadData, 
        // unidadData ya debe traer proyecto_id desde el componente
        admin_username: adminUsername 
    };
    return apiCall('/admin/create-unidad', 'POST', body, token);
};

// Update y Delete no cambian (usan ID de la unidad)
export const updateUnidad = (token, unidadData, adminUsername) => {
    const body = { ...unidadData, admin_username: adminUsername };
    return apiCall('/admin/update-unidad', 'POST', body, token);
};

export const deleteUnidad = (token, id, adminUsername) => {
    const body = { id, admin_username: adminUsername };
    return apiCall('/admin/delete-unidad', 'POST', body, token);
};