import { apiCall } from './authService';

/**
 * Obtiene todos los proyectos.
 */
export const getAllProjects = (token, adminUsername) => {
    return apiCall('/admin/get-proyectos', 'POST', { admin_username: adminUsername }, token);
};

/**
 * Crea un nuevo proyecto.
 */
export const createProject = (token, proyectoData) => {
    // El body se construye en Portafolio.js { nombre, fecha_inicio, ..., admin_username }
    return apiCall('/admin/create-proyecto', 'POST', proyectoData, token);
};

/**
 * Actualiza un proyecto existente.
 */
export const updateProject = (token, proyectoData) => {
    // El body se construye en Portafolio.js { id, nombre, fecha_inicio, ..., admin_username }
    return apiCall('/admin/update-proyecto', 'POST', proyectoData, token);
};

/**
 * Elimina un proyecto.
 */
export const deleteProject = (token, projectId, adminUsername) => {
    const body = {
        id: projectId,
        admin_username: adminUsername
    };
    return apiCall('/admin/delete-proyecto', 'POST', body, token);
};

/**
 * Obtiene la lista de proyectos para el panel de admin
 */
export const getAdminProjects = (token, adminUsername) => {
    return apiCall('/admin/get-proyectos', 'POST', { admin_username: adminUsername }, token);
};

/**
 * ⭐️ FUNCIÓN PARA CAMBIAR ESTADO ⭐️
 * Cambia el estado de un proyecto (habilitado/cerrado).
 * (Corresponde a tu 'AdminSetProyectoEstadoHandler')
 */
export const setProjectState = (token, projectId, newState, adminUsername) => {
    // Tu handler espera { id, estado, admin_username }
    const body = {
        id: projectId,
        estado: newState, // 'habilitado' o 'cerrado'
        admin_username: adminUsername
    };
    return apiCall('/admin/set-proyecto-estado', 'POST', body, token);
};