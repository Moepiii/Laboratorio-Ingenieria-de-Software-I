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

/**
 * Crea una nueva labor.
 * ⭐️ MODIFICADO: El 'laborData' ya NO incluye 'codigo_labor'.
 */
export const createLabor = (token, laborData, adminUsername) => {
    // laborData ahora es { proyecto_id, descripcion, estado }
    // El backend generará el 'codigo_labor' automáticamente.
    const body = {
        ...laborData,
        admin_username: adminUsername
    };
    return apiCall('/admin/create-labor', 'POST', body, token);
};

/**
 * Actualiza una labor existente.
 * ⭐️ (Sin cambios, la actualización aún envía el código)
 */
export const updateLabor = (token, laborData, adminUsername) => {
    // laborData ahora es { id, codigo_labor, descripcion, estado }
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