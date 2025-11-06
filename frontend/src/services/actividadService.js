import { apiCall } from './authService';

/**
 * Obtiene TODOS los datos para la página de Datos del Proyecto
 * (Labores, Equipos, Encargados y Actividades)
 * Llama a: /api/admin/get-datos-proyecto
 */
export const getDatosProyecto = (token, proyectoId, adminUsername) => {
    const body = {
        proyecto_id: proyectoId,
        admin_username: adminUsername
    };
    return apiCall('/admin/get-datos-proyecto', 'POST', body, token);
};

/**
 * Crea una nueva actividad.
 * Llama a: /api/admin/create-actividad
 */
export const createActividad = (token, actividadData, adminUsername) => {
    // actividadData = { proyecto_id, actividad, labor_agronomica_id, ... }
    const body = {
        ...actividadData,
        admin_username: adminUsername
    };
    // Devuelve la lista actualizada de actividades
    return apiCall('/admin/create-actividad', 'POST', body, token);
};

/**
 * Actualiza una actividad existente.
 * Llama a: /api/admin/update-actividad
 */
export const updateActividad = (token, actividadData, adminUsername) => {
    // actividadData = { id, proyecto_id, actividad, labor_agronomica_id, ... }
    const body = {
        ...actividadData,
        admin_username: adminUsername
    };
    // Devuelve la lista actualizada de actividades
    return apiCall('/admin/update-actividad', 'POST', body, token);
};

/**
 * Elimina una actividad.
 * Llama a: /api/admin/delete-actividad
 */
export const deleteActividad = (token, actividadId, adminUsername) => {
    const body = {
        id: actividadId,
        admin_username: adminUsername
    };
    return apiCall('/admin/delete-actividad', 'POST', body, token);
};