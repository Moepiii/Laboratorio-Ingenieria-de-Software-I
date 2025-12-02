import { apiCall } from './authService';

export const getMateriales = (token, proyectoId, adminUsername) => {
    return apiCall('/admin/get-materiales', 'POST', {
        proyecto_id: parseInt(proyectoId),
        admin_username: adminUsername
    }, token);
};

export const createMaterial = (token, data, adminUsername) => {
    const body = {
        ...data,
        proyecto_id: parseInt(data.proyecto_id),
        cantidad: parseFloat(data.cantidad),
        costo_unitario: parseFloat(data.costo_unitario),
        monto: parseFloat(data.monto),
        admin_username: adminUsername
    };
    return apiCall('/admin/create-material', 'POST', body, token);
};

export const updateMaterial = (token, data, adminUsername) => {
    const body = {
        ...data,
        id: parseInt(data.id),
        cantidad: parseFloat(data.cantidad),
        costo_unitario: parseFloat(data.costo_unitario),
        monto: parseFloat(data.monto),
        admin_username: adminUsername
    };
    return apiCall('/admin/update-material', 'POST', body, token);
};

export const deleteMaterial = (token, id, adminUsername) => {
    return apiCall('/admin/delete-material', 'POST', { id: parseInt(id), admin_username: adminUsername }, token);
};