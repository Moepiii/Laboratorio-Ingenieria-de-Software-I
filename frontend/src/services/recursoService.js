import { apiCall } from './authService';

/**
 * Obtiene la lista de recursos humanos de un proyecto.
 */
export const getRecursos = (token, proyectoId, adminUsername) => {
    return apiCall('/admin/get-recursos', 'POST', {
        proyecto_id: parseInt(proyectoId),
        admin_username: adminUsername
    }, token);
};

/**
 * Crea un nuevo recurso humano.
 */
export const createRecurso = (token, data, adminUsername) => {
    const body = {
        ...data,
        proyecto_id: parseInt(data.proyecto_id),
        // Aseguramos que los números sean enviados correctamente
        tiempo: parseFloat(data.tiempo),
        cantidad: parseFloat(data.cantidad),
        costo_unitario: parseFloat(data.costo_unitario),
        monto: parseFloat(data.monto),
        admin_username: adminUsername
    };
    return apiCall('/admin/create-recurso', 'POST', body, token);
};

/**
 * Actualiza un recurso existente.
 */
export const updateRecurso = (token, data, adminUsername) => {
    const body = {
        ...data,
        id: parseInt(data.id),
        tiempo: parseFloat(data.tiempo),
        cantidad: parseFloat(data.cantidad),
        costo_unitario: parseFloat(data.costo_unitario),
        monto: parseFloat(data.monto),
        admin_username: adminUsername
    };
    return apiCall('/admin/update-recurso', 'POST', body, token);
};

/**
 * Elimina un recurso humano.
 */
export const deleteRecurso = (token, id, adminUsername) => {
    return apiCall('/admin/delete-recurso', 'POST', {
        id: parseInt(id),
        admin_username: adminUsername
    }, token);
};