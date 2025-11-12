import { apiCall } from './authService';

// --- Funciones para el Dashboard de Usuario ---

export const getUserDashboardData = (token, userId) => {
    // Coincide con json:"user_id" en models.go
    return apiCall('/user/project-details', 'POST', { user_id: userId }, token);
};


// --- Funciones para el Panel de Admin ---

export const getAdminUsers = (token, adminUsername) => {
    return apiCall('/admin/users', 'POST', { admin_username: adminUsername }, token);
};

export const getGerentes = async (token, adminUsername) => {
    const data = await getAdminUsers(token, adminUsername);
    const gerentes = data.users.filter(user => user.role === 'gerente');
    return { users: gerentes };
};


// ⭐️ INICIO DE LA CORRECCIÓN ⭐️
export const adminAddUser = (token, userData, adminUsername) => {

    // 'userData' es el objeto {username, password, ...}
    // El backend espera { user: {...}, admin_username: "..." }
    const body = {
        user: userData, // ⬅️ En lugar de "...userData"
        admin_username: adminUsername
    };

    return apiCall('/admin/add-user', 'POST', body, token);
};
// ⭐️ FIN DE LA CORRECCIÓN ⭐️


export const adminDeleteUser = (token, userId, adminUsername) => {
    const body = {
        id: userId,
        admin_username: adminUsername
    };
    return apiCall('/admin/delete-user', 'POST', body, token);
};

export const adminUpdateUserRole = (token, userId, newRole, adminUsername) => {
    const body = {
        id: userId,
        new_role: newRole,
        admin_username: adminUsername
    };
    return apiCall('/admin/update-user', 'POST', body, token);
};

export const adminAssignProjectToUser = (token, userId, projectId, adminUsername) => {
    const body = {
        user_id: userId,
        proyecto_id: projectId,
        admin_username: adminUsername
    };
    return apiCall('/admin/assign-project', 'POST', body, token);
};

// Obtiene la data del dashboard de un usuario (para el rol 'user')
export const getProjectDetailsForUser = (token, userId) => {
    const body = {
        user_id: userId
    };
    return apiCall('/user/project-details', 'POST', body, token);
};