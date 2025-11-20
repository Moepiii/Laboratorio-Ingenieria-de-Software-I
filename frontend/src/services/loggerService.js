import { apiCall } from './authService';

/**
 * Llama al endpoint del backend para obtener la bitácora de eventos.
 *
 * @param {string} token - El token JWT del admin.
 * @param {string} adminUsername - El username del admin.
 * @param {object} filters - Un objeto con los filtros a aplicar.
 * @param {string} [filters.fecha_inicio] - 'YYYY-MM-DD'
 * @param {string} [filters.fecha_cierre] - 'YYYY-MM-DD'
 * @param {string} [filters.usuario_username] - Nombre de usuario a buscar
 * @param {string} [filters.accion] - Acción específica (ej. "CREACIÓN")
 * @param {string} [filters.entidad] - Entidad específica (ej. "Proyectos")
 * @returns {Promise<Array>} - Una promesa que resuelve a la lista de logs.
 */
export const getLogs = (token, adminUsername, filters = {}) => {

    // El cuerpo de la solicitud incluye el admin (para permisos)
    // y el objeto de filtros.
    const body = {
        admin_username: adminUsername,
        ...filters
    };

    // Llama al nuevo endpoint que creamos en Go
    return apiCall('/admin/get-logs', 'POST', body, token);
};
export const deleteLogs = (token, logIds, adminUsername) => {
    const body = {
        ids: logIds, // Array de números
        admin_username: adminUsername
    };
    return apiCall('/admin/delete-logs', 'POST', body, token);
};