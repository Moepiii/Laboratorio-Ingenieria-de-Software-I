import { apiCall } from './authService';

/**
 * Llama al endpoint del backend para obtener la bitácora de eventos.
 *
 * @param {string} token - El token JWT del admin.
 * @param {string} adminUsername - El username del admin.
 * @param {object} filters - Un objeto con los filtros a aplicar.
 * @returns {Promise<Array>} - Una promesa que resuelve a la lista de logs.
 */
export const getLogs = (token, adminUsername, filters = {}) => {
    // El cuerpo de la solicitud incluye el admin (para permisos)
    // y el objeto de filtros.
    const body = {
        admin_username: adminUsername,
        ...filters
    };

    // Llama al endpoint de Go
    return apiCall('/admin/get-logs', 'POST', body, token);
};

/**
 * Elimina logs específicos seleccionados por ID.
 */
export const deleteLogs = (token, logIds, adminUsername) => {
    const body = {
        ids: logIds, // Array de números
        admin_username: adminUsername
    };
    return apiCall('/admin/delete-logs', 'POST', body, token);
};

/**
 * ⭐️ NUEVA FUNCIÓN: Elimina logs por rango de fechas (Masivo)
 * Llama a: /api/admin/delete-logs-range
 */
export const deleteLogsByRange = (token, startDate, endDate, adminUsername) => {
    const body = {
        fecha_inicio: startDate, // YYYY-MM-DD
        fecha_fin: endDate,      // YYYY-MM-DD
        admin_username: adminUsername
    };
    return apiCall('/admin/delete-logs-range', 'POST', body, token);
};