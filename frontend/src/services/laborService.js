import { apiCall } from './authService';

/**
 * Obtiene todas las labores de un proyecto.
 * Llama a: /api/admin/get-labores
 */
export const getLabores = (token, proyectoId, adminUsername) => {
    // El body que espera el GetLaboresHandler en Go
    const body = {
        proyecto_id: proyectoId,
        admin_username: adminUsername
    };
    return apiCall('/admin/get-labores', 'POST', body, token);
};

/**
 * Crea una nueva labor.
 * Llama a: /api/admin/create-labor
 */
export const createLabor = (token, laborData, adminUsername) => {
    // laborData = { proyecto_id, descripcion, estado }
    const body = {
        ...laborData,
        admin_username: adminUsername
    };
    return apiCall('/admin/create-labor', 'POST', body, token);
};

/**
 * Actualiza una labor existente.
 * Llama a: /api/admin/update-labor
 */
export const updateLabor = (token, laborData, adminUsername) => {
    // laborData = { id, descripcion, estado }
    const body = {
        ...laborData,
        admin_username: adminUsername
    };
    return apiCall('/admin/update-labor', 'POST', body, token);
};

/**
 * Elimina una labor.
 * Llama a: /api/admin/delete-labor
 */
export const deleteLabor = (token, laborId, adminUsername) => {
    const body = {
        id: laborId,
        admin_username: adminUsername
    };
    return apiCall('/admin/delete-labor', 'POST', body, token);
};