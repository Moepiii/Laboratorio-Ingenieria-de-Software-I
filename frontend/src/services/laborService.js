import { apiCall } from './authService';

/**
 * Obtiene todas las labores de un proyecto.
 */
export const getLabores = (token, proyectoId, adminUsername) => {
    const body = {
        proyecto_id: proyectoId,
        admin_username: adminUsername
    };
    return apiCall('/admin/get-labores', 'POST', body, token);
};


export const createLabor = (token, laborData, adminUsername) => {

    const body = {
        ...laborData,
        admin_username: adminUsername
    };
    return apiCall('/admin/create-labor', 'POST', body, token);
};


export const updateLabor = (token, laborData, adminUsername) => {

    const body = {
        ...laborData,
        admin_username: adminUsername
    };
    return apiCall('/admin/update-labor', 'POST', body, token);
};

/**
 * Elimina una labor.
 */
export const deleteLabor = (token, laborId, adminUsername) => {
    const body = {
        id: laborId,
        admin_username: adminUsername
    };
    return apiCall('/admin/delete-labor', 'POST', body, token);
};