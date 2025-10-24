import React, { useState, useEffect, useCallback } from 'react';

const styles = {
  adminFormContainer: {
    padding: '1.5rem',
    backgroundColor: '#f9fafb',
    borderRadius: '8px',
    marginBottom: '2rem',
    border: '1px solid #e5e7eb',
  },
  input: {
    width: '100%',
    padding: '0.75rem 1rem',
    border: '1px solid #d1d5db',
    borderRadius: '8px',
    fontSize: '1rem',
    transition: 'border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out',
    boxSizing: 'border-box',
  },
  button: {
    width: '100%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    padding: '0.75rem 1rem',
    fontSize: '1rem',
    fontWeight: '600',
    borderRadius: '8px',
    color: 'white',
    backgroundColor: '#4f46e5',
    border: 'none',
    cursor: 'pointer',
    transition: 'background-color 0.2s, transform 0.1s',
    boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
  },
  error: {
    fontSize: '0.875rem',
    color: '#dc2626',
    fontWeight: '500',
    backgroundColor: '#fef2f2',
    padding: '0.75rem',
    borderRadius: '8px',
    border: '1px solid #fecaca',
  },
  success: {
    fontSize: '0.875rem',
    color: '#059669',
    fontWeight: '500',
    backgroundColor: '#ecfdf5',
    padding: '0.75rem',
    borderRadius: '8px',
    border: '1px solid #a7f3d0',
  },
  selectAssign: {
    padding: '0.5rem',
    borderRadius: '6px',
    border: '1px solid #d1d5db',
    marginRight: '0.5rem',
    fontSize: '0.875rem',
    backgroundColor: 'white',
  },
  buttonAssign: {
    padding: '0.5rem 1rem',
    borderRadius: '6px',
    fontSize: '0.875rem',
    fontWeight: '600',
    backgroundColor: '#3b82f6',
    color: 'white',
    border: 'none',
    cursor: 'pointer',
    transition: 'background-color 0.2s',
  },
  label: { display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.25rem' },
  inputGroup: { marginBottom: '1.5rem' },
  card: { padding: '2rem', backgroundColor: '#ffffff', borderRadius: '12px', boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.25)', width: '100%', maxWidth: '400px', margin: 'auto' },
};

const PerfilesUsuarios = ({ apiCall, currentUser, userRole }) => {
  const [users, setUsers] = useState([]);
  const [loadingUsers, setLoadingUsers] = useState(false);
  const [adminError, setAdminError] = useState('');
  const [adminSuccess, setAdminSuccess] = useState('');
  const [newUsername, setNewUsername] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [newNombre, setNewNombre] = useState('');
  const [newApellido, setNewApellido] = useState('');
  const [loadingAdd, setLoadingAdd] = useState(false);
  const [proyectosList, setProyectosList] = useState([]);
  const [selectedProyectos, setSelectedProyectos] = useState({});
  const [selectedRoles, setSelectedRoles] = useState({});

  const fetchUsersAndProyectos = useCallback(async () => {
    setLoadingUsers(true);
    setAdminError('');
    try {
      const usersResult = await apiCall('admin/users', { admin_username: currentUser }, 'POST');
      if (usersResult.success) {
        const fetchedUsers = usersResult.data.users || [];
        setUsers(fetchedUsers);
        const initialSelections = {};
        const initialRoles = {};
        fetchedUsers.forEach(user => {
          initialSelections[user.id] = user.proyecto_id || '0';
          initialRoles[user.id] = user.role;
        });
        setSelectedProyectos(initialSelections);
        setSelectedRoles(initialRoles);
      } else {
        setAdminError('No se pudo cargar la lista de usuarios: ' + (usersResult.data.error || 'Desconocido'));
      }
      const proyectosResult = await apiCall('admin/get-proyectos', { admin_username: currentUser }, 'POST');
      if (proyectosResult.success) {
        setProyectosList(proyectosResult.data.proyectos || []);
      } else {
        setAdminError('No se pudo cargar la lista de proyectos.');
      }
    } catch (e) {
      setAdminError(`Error de conexión: ${e.message}`);
    } finally {
      setLoadingUsers(false);
    }
  }, [currentUser, apiCall]);

  useEffect(() => {
    fetchUsersAndProyectos();
  }, [fetchUsersAndProyectos]);

  const handleAdminAddUser = async (e) => {
    e.preventDefault();
    setAdminError(''); setAdminSuccess(''); setLoadingAdd(true);
    if (newPassword.length < 6) {
      setAdminError('La contraseña debe tener al menos 6 caracteres.');
      setLoadingAdd(false);
      return;
    }
    try {
      const result = await apiCall('admin/add-user', {
        username: newUsername, password: newPassword, nombre: newNombre, apellido: newApellido, admin_username: currentUser
      }, 'POST');
      if (result.success) {
        setAdminSuccess(result.data.mensaje);
        setNewUsername(''); setNewPassword(''); setNewNombre(''); setNewApellido('');
        fetchUsersAndProyectos();
      } else {
        setAdminError(result.data.error || 'Fallo al añadir el usuario.');
      }
    } catch (e) {
      setAdminError(`Error de conexión al añadir usuario: ${e.message}`);
    } finally {
      setLoadingAdd(false);
    }
  };

  const handleAdminDeleteUser = async (userId, username) => {
    setAdminError(''); setAdminSuccess('');
    if (!window.confirm(`¿Estás seguro de borrar permanentemente a ${username} (ID: ${userId})?`)) { return; }
    try {
      const result = await apiCall('admin/delete-user', { id: userId, admin_username: currentUser }, 'POST');
      if (result.success) {
        setAdminSuccess(result.data.mensaje);
        fetchUsersAndProyectos();
      } else {
        setAdminError(result.data.error || 'Fallo al borrar el usuario.');
      }
    } catch (e) {
      setAdminError(`Error de conexión al borrar usuario: ${e.message}`);
    }
  };

  const handleAdminUpdateRole = async (userId, newRole) => {
    setAdminError(''); setAdminSuccess('');
    if (!newRole) {
      setAdminError('Rol no seleccionado.');
      return;
    }
    try {
      const result = await apiCall('admin/update-user', {
        id: userId,
        new_role: newRole,
        admin_username: currentUser
      }, 'POST');
      if (result.success) {
        setAdminSuccess(`Rol actualizado con éxito.`);
        fetchUsersAndProyectos();
      } else {
        setAdminError(result.data.error || 'Fallo al actualizar el rol.');
      }
    } catch (e) {
      setAdminError(`Error de conexión al actualizar rol: ${e.message}`);
    }
  };

  const handleSelectRoleChange = (userId, newRole) => {
    setSelectedRoles(prev => ({
      ...prev,
      [userId]: newRole
    }));
  };

  const handleAssignProyecto = async (userId, proyectoId) => {
    setAdminError(''); setAdminSuccess('');
    const idAsignar = parseInt(proyectoId, 10);
    try {
      const result = await apiCall('admin/assign-proyecto', {
        user_id: userId,
        proyecto_id: idAsignar,
        admin_username: currentUser
      }, 'POST');
      if (result.success) {
        setAdminSuccess(result.data.mensaje);
        fetchUsersAndProyectos();
      } else {
        setAdminError(result.data.error || 'Fallo al asignar el proyecto.');
      }
    } catch (e) {
      setAdminError(`Error de conexión al asignar: ${e.message}`);
    }
  };

  const handleSelectProyectoChange = (userId, newProyectoId) => {
    setSelectedProyectos(prev => ({
      ...prev,
      [userId]: newProyectoId
    }));
  };

  return (
    <div>
      <h2 style={{ fontSize: '1.875rem', fontWeight: '700', color: '#1f2937' }}>
        Perfiles de Usuarios
      </h2>
      <p style={{ fontSize: '1.125rem', color: '#4b5563', marginBottom: '1.5rem' }}>
        Logueado como: **{currentUser}** (Rol: {userRole})
      </p>

      {adminSuccess && <p style={{ ...styles.success, marginBottom: '1rem' }}>{adminSuccess}</p>}
      {adminError && <p style={{ ...styles.error, marginBottom: '1rem' }}>{adminError}</p>}

      {userRole === 'admin' && (
        <div style={styles.adminFormContainer}>
          <h3 style={{ fontSize: '1.25rem', fontWeight: '700', color: '#1f2937', marginBottom: '1rem' }}>➕ Crear Nuevo Usuario</h3>
          <form onSubmit={handleAdminAddUser} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            <div style={{ display: 'flex', flexDirection: window.innerWidth > 600 ? 'row' : 'column', gap: '1rem' }}>
              <input type="text" placeholder="Nombre de Usuario" value={newUsername} onChange={(e) => setNewUsername(e.target.value)} required style={{ ...styles.input, flex: 1 }} disabled={loadingAdd} />
              <input type="password" placeholder="Contraseña (mín. 6 caracteres)" value={newPassword} onChange={(e) => setNewPassword(e.target.value)} required style={{ ...styles.input, flex: 1 }} disabled={loadingAdd} />
            </div>
            <div style={{ display: 'flex', flexDirection: window.innerWidth > 600 ? 'row' : 'column', gap: '1rem' }}>
              <input type="text" placeholder="Nombre" value={newNombre} onChange={(e) => setNewNombre(e.target.value)} required style={{ ...styles.input, flex: 1 }} disabled={loadingAdd} />
              <input type="text" placeholder="Apellido" value={newApellido} onChange={(e) => setNewApellido(e.target.value)} required style={{ ...styles.input, flex: 1 }} disabled={loadingAdd} />
            </div>
            <button type="submit" style={{ ...styles.button, width: '100%', backgroundColor: '#10b81' }} onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#059669'} onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#10b81'} disabled={loadingAdd}>
              {loadingAdd ? 'Creando...' : 'Crear Usuario'}
            </button>
          </form>
        </div>
      )}

      <h3 style={{ fontSize: '1.25rem', fontWeight: '700', color: '#1f2937', marginBottom: '1rem' }}>
        Lista de Usuarios (Total: {users.length})
      </h3>

      {loadingUsers ? (
        <p style={{ textAlign: 'center', padding: '2rem', color: '#4f46e5' }}>Cargando usuarios y proyectos...</p>
      ) : (
        <div style={{ maxHeight: '600px', overflowY: 'auto', border: '1px solid #e5e7eb', borderRadius: '8px' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left' }}>
            <thead>
              <tr style={{ backgroundColor: '#f9fafb', borderBottom: '1px solid #e5e7eb' }}>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>ID</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Usuario</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Nombre</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Apellido</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Rol</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Proyecto Actual</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Asignar Proyecto</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280', textAlign: 'center' }}>Acciones (Solo Admin)</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user, index) => {
                const isSelf = user.username === currentUser;
                return (
                  <tr key={user.id} style={{ borderBottom: index < users.length - 1 ? '1px solid #f3f4f6' : 'none', backgroundColor: isSelf ? '#fffbeb' : 'white' }}>

                    <td style={{ padding: '0.75rem', fontWeight: '500' }}>{user.id}</td>
                    <td style={{ padding: '0.75rem', fontWeight: '600' }}>{user.username} {isSelf && <span style={{ fontSize: '0.75rem', color: '#f59e0b' }}>(Tú)</span>}</td>
                    <td style={{ padding: '0.75rem' }}>{user.nombre}</td>
                    <td style={{ padding: '0.75rem' }}>{user.apellido}</td>

                    <td style={{ padding: '0.75rem' }}>
                      <span style={{
                        padding: '0.25rem 0.75rem',
                        borderRadius: '9999px',
                        fontSize: '0.75rem',
                        fontWeight: '700',
                        backgroundColor: user.role === 'admin' ? '#e0e7ff' : (user.role === 'gerente' ? '#fef3c7' : '#d1fae5'),
                        color: user.role === 'admin' ? '#3730a3' : (user.role === 'gerente' ? '#92400e' : '#065f46'),
                      }}>
                        {user.role.toUpperCase()}
                      </span>
                    </td>

                    <td style={{ padding: '0.75rem', fontWeight: '500', color: user.proyecto_nombre ? '#1d4ed8' : '#6b7280' }}>
                      {user.proyecto_nombre || <span style={{ fontStyle: 'italic' }}>No asignado</span>}
                    </td>

                    <td style={{ padding: '0.75rem' }}>
                      {isSelf ? (
                        <span style={{ fontSize: '0.875rem', color: '#6b7280' }}>N/A</span>
                      ) : (
                        <div style={{ display: 'flex', alignItems: 'center' }}>
                          <select
                            style={styles.selectAssign}
                            value={selectedProyectos[user.id] || '0'}
                            onChange={(e) => handleSelectProyectoChange(user.id, e.target.value)}
                          >
                            <option value="0">-- No asignado --</option>
                            {proyectosList.map(p => (
                              <option key={p.id} value={p.id}>{p.nombre}</option>
                            ))}
                          </select>
                          <button
                            style={styles.buttonAssign}
                            onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#2563eb'}
                            onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#3b82f6'}
                            onClick={() => handleAssignProyecto(user.id, selectedProyectos[user.id])}
                          >
                            Asignar
                          </button>
                        </div>
                      )}
                    </td>

                    <td style={{ padding: '0.75rem', textAlign: 'center' }}>
                      {(userRole === 'admin' && !isSelf) ? (
                        <div style={{ display: 'flex', gap: '0.5rem', justifyContent: 'center' }}>
                          <select
                            style={styles.selectAssign}
                            value={selectedRoles[user.id] || user.role}
                            onChange={(e) => handleSelectRoleChange(user.id, e.target.value)}
                          >
                            <option value="user">User</option>
                            <option value="gerente">Gerente</option>
                            <option value="admin">Admin</option>
                          </select>
                          <button
                            style={{ ...styles.buttonAssign, backgroundColor: '#f59e0b' }}
                            onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#d97706'}
                            onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#f59e0b'}
                            onClick={() => handleAdminUpdateRole(user.id, selectedRoles[user.id] || user.role)}
                          >
                            Guardar Rol
                          </button>
                          <button
                            onClick={() => handleAdminDeleteUser(user.id, user.username)}
                            style={{ ...styles.buttonAssign, backgroundColor: '#ef4444' }}
                            onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#dc2626'}
                            onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#ef4444'}
                          >
                            Borrar
                          </button>
                        </div>
                      ) : (
                        <span style={{ fontSize: '0.875rem', color: '#6b7280' }}>{isSelf ? '(Tú)' : 'N/A'}</span>
                      )}
                    </td>

                  </tr>
                )
              })}
              {users.length === 0 && !loadingUsers && (
                <tr><td colSpan="9" style={{ textAlign: 'center', padding: '1rem', color: '#6b7280' }}>No hay usuarios para mostrar.</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default PerfilesUsuarios;