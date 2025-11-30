import { apiCall } from './authService';

export const getPlanes = (token, proyectoId, adminUsername) => {
    return apiCall('/admin/get-planes', 'POST', {
        proyecto_id: parseInt(proyectoId),
        admin_username: adminUsername
    }, token);
};

export const createPlan = (token, planData, adminUsername) => {
    const body = {
        ...planData,
        proyecto_id: parseInt(planData.proyecto_id),
        horas: parseFloat(planData.horas),
        costo_unitario: parseFloat(planData.costo_unitario),
        monto: parseFloat(planData.monto),
        admin_username: adminUsername
    };
    return apiCall('/admin/create-plan', 'POST', body, token);
};

// ⭐️ NUEVO: Actualizar
export const updatePlan = (token, planData, adminUsername) => {
    const body = {
        ...planData,
        id: parseInt(planData.id), // Importante el ID
        horas: parseFloat(planData.horas),
        costo_unitario: parseFloat(planData.costo_unitario),
        monto: parseFloat(planData.monto),
        admin_username: adminUsername
    };
    return apiCall('/admin/update-plan', 'POST', body, token);
};

// ⭐️ NUEVO: Borrar
export const deletePlan = (token, planId, adminUsername) => {
    const body = {
        id: parseInt(planId),
        admin_username: adminUsername
    };
    return apiCall('/admin/delete-plan', 'POST', body, token);
};