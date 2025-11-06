import { apiCall } from './authService';

/**
 * Obtiene todos los equipos de un proyecto.
 * Llama a: /api/admin/get-equipos
 */
export const getEquipos = (token, proyectoId, adminUsername) => {
    // El body que espera el GetEquiposHandler en Go
    const body = {
        proyecto_id: proyectoId,
        admin_username: adminUsername
    };
    return apiCall('/admin/get-equipos', 'POST', body, token);
};

/**
 * Crea un nuevo equipo.
 * Llama a: /api/admin/create-equipo
 */
export const createEquipo = (token, equipoData, adminUsername) => {
    // equipoData = { proyecto_id, nombre, tipo, estado }
    const body = {
        ...equipoData,
        admin_username: adminUsername
    };
    return apiCall('/admin/create-equipo', 'POST', body, token);
};

/**
 * Actualiza un equipo existente.
 * Llama a: /api/admin/update-equipo
 */
export const updateEquipo = (token, equipoData, adminUsername) => {
    // equipoData = { id, nombre, tipo, estado }
    const body = {
        ...equipoData,
        admin_username: adminUsername
    };
    return apiCall('/admin/update-equipo', 'POST', body, token);
};

/**
 * Elimina un equipo.
 * Llama a: /api/admin/delete-equipo
 */
export const deleteEquipo = (token, equipoId, adminUsername) => {
    const body = {
        id: equipoId,
        admin_username: adminUsername
    };
    return apiCall('/admin/delete-equipo', 'POST', body, token);
};