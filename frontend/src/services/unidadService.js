import { apiCall } from './authService';

export const getUnidades = (token) => {
    // Se asume que GET no requiere body, o si tu backend valida adminUsername en GET, ajustalo
    return apiCall('/admin/get-unidades', 'GET', null, token); 
    // NOTA: Si tu apiCall no soporta GET sin body, usa POST con {admin_username} como hiciste en otros.
    // Para mantener consistencia con tu proyecto, probablemente prefieras POST:
    // return apiCall('/admin/get-unidades', 'POST', {}, token);
};

export const createUnidad = (token, unidadData, adminUsername) => {
    const body = { ...unidadData, admin_username: adminUsername };
    return apiCall('/admin/create-unidad', 'POST', body, token);
};

export const updateUnidad = (token, unidadData, adminUsername) => {
    const body = { ...unidadData, admin_username: adminUsername };
    return apiCall('/admin/update-unidad', 'POST', body, token);
};

export const deleteUnidad = (token, id, adminUsername) => {
    const body = { id, admin_username: adminUsername };
    return apiCall('/admin/delete-unidad', 'POST', body, token);
};