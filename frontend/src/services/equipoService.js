import { apiCall } from './authService';

/**
 * Obtiene todos los equipos de un proyecto.
 */
export const getEquipos = (token, proyectoId, adminUsername) => {
    const body = {
        proyecto_id: proyectoId,
        admin_username: adminUsername
    };
    return apiCall('/admin/get-equipos', 'POST', body, token);
};

/**
 * Crea un nuevo equipo.
 * ⭐️ MODIFICADO: El 'equipoData' ahora incluye 'codigo_equipo'
 */
export const createEquipo = (token, equipoData, adminUsername) => {
    // equipoData ahora es { proyecto_id, codigo_equipo, nombre, tipo, estado }
    const body = {
        ...equipoData,
        admin_username: adminUsername
    };
    return apiCall('/admin/create-equipo', 'POST', body, token);
};

/**
 * Actualiza un equipo existente.
 * ⭐️ MODIFICADO: El 'equipoData' ahora incluye 'codigo_equipo'
 */
export const updateEquipo = (token, equipoData, adminUsername) => {
    // equipoData ahora es { id, codigo_equipo, nombre, tipo, estado }
    const body = {
        ...equipoData,
        admin_username: adminUsername
    };
    return apiCall('/admin/update-equipo', 'POST', body, token);
};

/**
 * Elimina un equipo.
 */
export const deleteEquipo = (token, equipoId, adminUsername) => {
    const body = {
        id: equipoId,
        admin_username: adminUsername
    };
    return apiCall('/admin/delete-equipo', 'POST', body, token);
};