import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from '../context/AuthContext';
import {
  getAdminUsers,
  adminAddUser,
  adminDeleteUser,
  adminUpdateUserRole,
  adminAssignProjectToUser
} from '../services/userService';
import { getAdminProjects } from '../services/projectService';

// (Estilos - sin cambios)
const styles = {
  adminFormContainer: { padding: '1.5rem', backgroundColor: '#f9fafb', borderRadius: '8px', marginBottom: '2rem', border: '1px solid #e5e7eb' },
  input: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', transition: 'border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out', boxSizing: 'border-box' },
  button: { width: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '0.75rem 1rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#4f46e5', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s, transform 0.1s', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
  error: { fontSize: '0.875rem', color: '#dc2626', fontWeight: '500', backgroundColor: '#fef2f2', padding: '0.75rem', borderRadius: '8px', border: '1px solid #fecaca' },
  success: { fontSize: '0.875rem', color: '#059669', fontWeight: '500', backgroundColor: '#ecfdf5', padding: '0.75rem', borderRadius: '8px', border: '1px solid #a7f3d0' },
  selectAssign: { padding: '0.5rem', borderRadius: '6px', border: '1px solid #d1d5db', marginRight: '0.5rem', fontSize: '0.875rem', backgroundColor: 'white' },
  buttonAssign: { padding: '0.5rem 1rem', borderRadius: '6px', fontSize: '0.875rem', fontWeight: '600', backgroundColor: '#3b82f6', color: 'white', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s' },
  label: { display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.25rem' },
  inputGroup: { marginBottom: '1.5rem' },
  card: { padding: '2rem', backgroundColor: '#ffffff', borderRadius: '12px', boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.25)', width: '100%', maxWidth: '400px', margin: 'auto' },
  tableContainer: { overflowX: 'auto', backgroundColor: 'white', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06)' },
  table: { width: '100%', borderCollapse: 'collapse', fontSize: '0.9rem' },
  th: { padding: '0.75rem 1rem', textAlign: 'left', borderBottom: '2px solid #e5e7eb', backgroundColor: '#f9fafb', color: '#6b7280', fontWeight: '600', textTransform: 'uppercase', letterSpacing: '0.05em' },
  td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', color: '#374151', verticalAlign: 'middle' },
};


const PerfilesUsuarios = () => {

  // 1. Hooks siempre al principio
  const { token, currentUser, userRole } = useAuth();
  const [users, setUsers] = useState([]);
  const [loadingUsers, setLoadingUsers] = useState(true);
  const [selectedProyectos, setSelectedProyectos] = useState({});
  const [selectedRoles, setSelectedRoles] = useState({});
  const [proyectosList, setProyectosList] = useState([]);
  const [adminError, setAdminError] = useState('');
  const [adminSuccess, setAdminSuccess] = useState('');
  const [loadingAdd, setLoadingAdd] = useState(false);
  const [newUsername, setNewUsername] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [newNombre, setNewNombre] = useState('');
  const [newApellido, setNewApellido] = useState('');

  // 2. useCallback con comprobación interna y dependencia completa
  const fetchUsersAndProyectos = useCallback(async () => {
    // Guarda: Salir si los datos requeridos no están listos
    if (!token || !currentUser?.username) {
      // console.warn("fetchUsersAndProyectos llamado antes de que user/token estuvieran listos.");
      setLoadingUsers(false);
      if (!currentUser && token) setAdminError("Error: No se pudieron cargar los datos del usuario actual.");
      return;
    }

    setLoadingUsers(true);
    setAdminError('');
    try {
      // Ahora es seguro usar currentUser.username
      const usersResult = await getAdminUsers(token, currentUser.username);
      const fetchedUsers = usersResult.users || [];
      setUsers(fetchedUsers);

      const initialRoles = {};
      const initialProyectos = {};
      fetchedUsers.forEach(user => {
        initialRoles[user.id] = user.role;
        initialProyectos[user.id] = user.proyecto_id || '0';
      });
      setSelectedRoles(initialRoles);
      setSelectedProyectos(initialProyectos);

      const proyectosResult = await getAdminProjects(token, currentUser.username);
      setProyectosList(proyectosResult.proyectos || []);

    } catch (e) {
      setAdminError(`Error al cargar datos: ${e.message}`);
    } finally {
      setLoadingUsers(false);
    }
    
  }, [token, currentUser]);

  // 3. useEffect llama a fetch solo cuando las dependencias están listas
  useEffect(() => {
    // Usa ?. en la guarda
    if (token && currentUser?.username) {
      fetchUsersAndProyectos();
    } else {
      setLoadingUsers(!token);
    }
   
  }, [fetchUsersAndProyectos, token, currentUser]);


  // 4. Handlers CRUD con guardas usando ?. (sin cambios en esta parte)
  const handleAdminAddUser = async (e) => {
    e.preventDefault();
    if (!token || !currentUser?.username) return;
    setLoadingAdd(true);
    setAdminError(''); setAdminSuccess('');
    try {
      const userData = { username: newUsername, password: newPassword, nombre: newNombre, apellido: newApellido };
      const result = await adminAddUser(token, userData, currentUser.username);
      setAdminSuccess(result.mensaje || 'Usuario creado');
      setNewUsername(''); setNewPassword(''); setNewNombre(''); setNewApellido('');
      fetchUsersAndProyectos();
    } catch (e) { setAdminError(`Error de conexión: ${e.message}`); }
    finally { setLoadingAdd(false); }
  };
  // ... (resto de handlers CRUD sin cambios) ...
  const handleAdminDeleteUser = async (userId, username) => {
    if (!token || !currentUser?.username || !window.confirm(`Borrar ${username}?`)) return;
    setAdminError(''); setAdminSuccess('');
    try {
      const result = await adminDeleteUser(token, userId, currentUser.username);
      setAdminSuccess(result.mensaje || 'Usuario borrado');
      fetchUsersAndProyectos();
    } catch (e) { setAdminError(`Error de conexión: ${e.message}`); }
  };
  const handleAdminUpdateRole = async (userId, newRole) => {
    if (!token || !currentUser?.username) return;
    setAdminError(''); setAdminSuccess('');
    try {
      await adminUpdateUserRole(token, userId, newRole, currentUser.username);
      setAdminSuccess(`Rol actualizado.`);
      fetchUsersAndProyectos();
    } catch (e) { setAdminError(`Error de conexión: ${e.message}`); }
  };
  const handleAssignProyecto = async (userId, proyectoId) => {
    if (!token || !currentUser?.username || !proyectoId) return;
    setAdminError(''); setAdminSuccess('');
    try {
      const idAsignar = parseInt(proyectoId, 10);
      const result = await adminAssignProjectToUser(token, userId, idAsignar, currentUser.username);
      setAdminSuccess(result.mensaje || 'Proyecto asignado/quitado');
      fetchUsersAndProyectos();
    } catch (e) { setAdminError(`Error de conexión: ${e.message}`); }
  };


  // Funciones auxiliares (sin cambios)
  const handleSelectRoleChange = (userId, newRole) => setSelectedRoles(prev => ({ ...prev, [userId]: newRole }));
  const handleSelectProyectoChange = (userId, newProyectoId) => setSelectedProyectos(prev => ({ ...prev, [userId]: newProyectoId }));

  // 5. Renderizado condicional principal
  if (loadingUsers || !currentUser?.username) {
    return <div style={{ padding: '2rem' }}>Cargando datos...</div>;
  }

  // --- Renderiza el componente principal ---
  return (
    <div style={{ width: '100%' }}>
      <h2 style={{ fontSize: '1.875rem', fontWeight: '700', color: '#1f2937' }}> Perfiles de Usuarios </h2>
      <p style={{ fontSize: '1.125rem', color: '#4b5563', marginBottom: '1.5rem' }}>
        Logueado como: **{currentUser?.username || 'Usuario Desconocido'}** (Rol: {userRole || 'N/A'})
      </p>

      {adminSuccess && <p style={{ ...styles.success, marginBottom: '1rem' }}>{adminSuccess}</p>}
      {!loadingUsers && adminError && <p style={{ ...styles.error, marginBottom: '1rem' }}>{adminError}</p>}

      {/* Formulario Admin */}
      {userRole === 'admin' && (
        <div style={styles.adminFormContainer}>
          {/* ... (contenido del formulario sin cambios) ... */}
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
            <button type="submit" style={{ ...styles.button, width: '100%', backgroundColor: loadingAdd ? '#9ca3af' : '#4f46e5' }} disabled={loadingAdd}>
              {loadingAdd ? 'Creando...' : 'Crear Usuario'}
            </button>
          </form>
        </div>
      )}

      {/* Tabla de Usuarios */}
      {!loadingUsers && !adminError && (
        <div style={styles.tableContainer}>
          <table style={styles.table}>
            <thead>
              <tr>
                <th style={styles.th}>ID</th>
                <th style={styles.th}>Usuario</th>
                <th style={styles.th}>Nombre</th>
                <th style={styles.th}>Apellido</th>
                <th style={styles.th}>Rol Actual</th>
                <th style={styles.th}>Asignar Proyecto</th>
                <th style={styles.th}>Acciones (Admin)</th>
              </tr>
            </thead>
            <tbody>
              {users.map(user => {
                const isSelf = user.username === currentUser?.username;
                const proyectosDisponibles = proyectosList.filter(p =>
                  user.role === 'user' ? p.estado === 'habilitado' : true
                );

                return (
                  <tr key={user.id}>
                    {/* ... (Celdas de la tabla sin cambios internos) ... */}
                    <td style={styles.td}>{user.id}</td>
                    <td style={styles.td}>{user.username} {isSelf && '(Tú)'}</td>
                    <td style={styles.td}>{user.nombre}</td>
                    <td style={styles.td}>{user.apellido}</td>
                    <td style={styles.td}>
                      {userRole === 'admin' && !isSelf ? (
                        <select
                          style={styles.selectAssign}
                          value={selectedRoles[user.id] || user.role}
                          onChange={(e) => handleSelectRoleChange(user.id, e.target.value)}
                        >
                          <option value="user">User</option>
                          <option value="gerente">Gerente</option>
                          <option value="admin">Admin</option>
                        </select>
                      ) : (
                        user.role
                      )}
                    </td>
                    <td style={styles.td}>
                      {(userRole === 'admin' || userRole === 'gerente') && !isSelf && (user.role === 'user' || user.role === 'gerente') ? (
                        <div style={{ display: 'flex' }}>
                          <select
                            style={styles.selectAssign}
                            value={selectedProyectos[user.id] || '0'}
                            onChange={(e) => handleSelectProyectoChange(user.id, e.target.value)}
                          >
                            <option value="0">Quitar / Ninguno</option>
                            {proyectosDisponibles.map(p => (
                              <option key={p.id} value={p.id}>{p.nombre}</option>
                            ))}
                          </select>
                          <button
                            style={styles.buttonAssign}
                            onClick={() => handleAssignProyecto(user.id, selectedProyectos[user.id])}
                          >
                            Asignar
                          </button>
                        </div>
                      ) : (
                        user.proyecto_nombre || <span style={{ color: '#9ca3af' }}>N/A</span>
                      )}
                    </td>
                    <td style={styles.td}>
                      {userRole === 'admin' && !isSelf ? (
                        <div style={{ display: 'flex', gap: '0.5rem' }}>
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
                );
              })}
              {users.length === 0 && (
                <tr><td colSpan="7" style={{ textAlign: 'center', padding: '1rem', color: '#6b7280' }}>No hay usuarios para mostrar.</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default PerfilesUsuarios;