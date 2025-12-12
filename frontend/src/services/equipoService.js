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


export const createEquipo = (token, equipoData, adminUsername) => {

    const body = {
        ...equipoData,
        admin_username: adminUsername
    };
    return apiCall('/admin/create-equipo', 'POST', body, token);
};


export const updateEquipo = (token, equipoData, adminUsername) => {

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