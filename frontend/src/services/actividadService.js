import { apiCall } from './authService';


export const getDatosProyecto = (token, proyectoId, adminUsername) => {
    const body = {
        proyecto_id: proyectoId,
        admin_username: adminUsername
    };
    return apiCall('/admin/get-datos-proyecto', 'POST', body, token);
};


export const createActividad = (token, actividadData, adminUsername) => {

    const body = {
        ...actividadData,
        admin_username: adminUsername
    };

    return apiCall('/admin/create-actividad', 'POST', body, token);
};


export const updateActividad = (token, actividadData, adminUsername) => {

    const body = {
        ...actividadData,
        admin_username: adminUsername
    };

    return apiCall('/admin/update-actividad', 'POST', body, token);
};


export const deleteActividad = (token, actividadId, adminUsername) => {
    const body = {
        id: actividadId,
        admin_username: adminUsername
    };
    return apiCall('/admin/delete-actividad', 'POST', body, token);
};