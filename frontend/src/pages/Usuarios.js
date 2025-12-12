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

// Estilos
const styles = {
  container: { padding: '2rem', color: '#333', fontFamily: 'Inter, sans-serif' },
  h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', marginBottom: '1.5rem', borderBottom: '2px solid #e5e7eb', paddingBottom: '0.75rem' },
  adminFormContainer: { padding: '1.5rem', backgroundColor: '#f9fafb', borderRadius: '8px', marginBottom: '2rem', border: '1px solid #e5e7eb' },
  h3: { fontSize: '1.25rem', fontWeight: '600', color: '#111827', marginTop: '0', marginBottom: '1rem' },
  input: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', transition: 'border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out', boxSizing: 'border-box' },
  select: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', boxSizing: 'border-box' },
  button: { width: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '0.75rem 1rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#4f46e5', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s, transform 0.1s', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
  error: { color: 'red', marginTop: '1rem' },
  success: { color: 'green', marginTop: '1rem' },
  tableContainer: { overflowX: 'auto', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
  table: { width: '100%', borderCollapse: 'collapse' },
  th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '2px solid #e5e7eb', color: '#374151', fontWeight: '600', fontSize: '0.875rem' },
  td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', verticalAlign: 'middle' },

  tdCedula: { minWidth: '150px' },
  // -----------------------------------------------------------
  roleSelect: { padding: '0.5rem', borderRadius: '4px', border: '1px solid #d1d5db', minWidth: '120px' },
  projectSelect: { padding: '0.5rem', borderRadius: '4px', border: '1px solid #d1d5db', minWidth: '150px' },
  buttonSave: { padding: '0.5rem 1rem', fontSize: '0.75rem', fontWeight: '600', borderRadius: '4px', color: 'white', backgroundColor: '#f59e0b', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s' },
  buttonAssign: { padding: '0.6rem 1.2rem', fontSize: '0.875rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#3b82f6', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s, transform 0.1s' }
};

const PerfilesUsuarios = () => {
  const { token, currentUser, userId, userRole, setError, setSuccessMessage } = useAuth();
  const [users, setUsers] = useState([]);
  const [projects, setProjects] = useState([]);
  const [newUser, setNewUser] = useState({
    username: '',
    password: '',
    nombre: '',
    apellido: '',
    cedula: '', 
    role: 'user',
    proyecto_id: '' 
  });
  const [selectedRoles, setSelectedRoles] = useState({});
  const [selectedProjects, setSelectedProjects] = useState({});
  const [loading, setLoading] = useState(false);
  const [generalError, setGeneralError] = useState('');
  const [generalSuccess, setGeneralSuccess] = useState('');

  const adminUsername = currentUser?.username;

  const fetchUsersAndProjects = useCallback(async () => {
    if (!token || !adminUsername) return;

    setLoading(true);
    setGeneralError('');
    try {
      // 1. Obtener Usuarios
      const usersData = await getAdminUsers(token, adminUsername);
      setUsers(usersData.users);

      // 2. Obtener Proyectos
      const projectsData = await getAdminProjects(token, adminUsername);
      setProjects(projectsData.proyectos || []);

      // 3. Inicializar estados de selección
      const initialRoles = {};
      const initialProjects = {};
      usersData.users.forEach(user => {
        initialRoles[user.id] = user.role;
        initialProjects[user.id] = user.proyecto_id || '';
      });
      setSelectedRoles(initialRoles);
      setSelectedProjects(initialProjects);

    } catch (err) {
      setGeneralError('Error al cargar datos: ' + err.message);
    } finally {
      setLoading(false);
    }
  }, [token, adminUsername]);

  useEffect(() => {
    fetchUsersAndProjects();
  }, [fetchUsersAndProjects]);


  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setNewUser(prev => ({ ...prev, [name]: value }));
  };

  const handleAdminAddUser = async (e) => {
    e.preventDefault();
    setGeneralError('');
    setGeneralSuccess('');

    if (newUser.username.trim() === '' || newUser.password.trim() === '' || newUser.cedula.trim() === '' || newUser.nombre.trim() === '' || newUser.apellido.trim() === '') {
      return setGeneralError('Todos los campos son obligatorios.');
    }

    setLoading(true);
    try {
      const userData = {
        username: newUser.username,
        password: newUser.password,
        role: newUser.role,
        nombre: newUser.nombre,
        apellido: newUser.apellido,
        cedula: newUser.cedula,
        proyecto_id: newUser.proyecto_id || null // Asegura que sea null si está vacío
      };

      const result = await adminAddUser(token, userData, adminUsername);
      setGeneralSuccess(result.mensaje || 'Usuario agregado con éxito.');
      setNewUser({ username: '', password: '', nombre: '', apellido: '', cedula: '', role: 'user', proyecto_id: '' });
      await fetchUsersAndProjects(); // Recargar la lista
    } catch (err) {
      setGeneralError('Error al agregar usuario: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleRoleChange = (userId, newRole) => {
    setSelectedRoles(prev => ({ ...prev, [userId]: newRole }));
  };

  const handleProjectChange = (userId, newProjectId) => {
    setSelectedProjects(prev => ({ ...prev, [userId]: newProjectId || null })); // Guarda null si es "No Asignar"
  };

  const handleAdminUpdateRole = async (userId, newRole) => {
    if (userId === currentUser.id) return setGeneralError("No puedes cambiar tu propio rol.");

    setGeneralError('');
    setGeneralSuccess('');
    setLoading(true);

    try {
      const result = await adminUpdateUserRole(token, userId, newRole, adminUsername);
      setGeneralSuccess(result.mensaje || 'Rol de usuario actualizado con éxito.');
      await fetchUsersAndProjects();
    } catch (err) {
      setGeneralError('Error al actualizar el rol: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleAdminAssignProject = async (userId, projectId) => {
    if (userId === currentUser.id && userRole !== 'admin') {
      return setGeneralError("Solo un administrador puede asignarse un proyecto a sí mismo.");
    }

    setGeneralError('');
    setGeneralSuccess('');
    setLoading(true);

    try {
      const actualProjectId = projectId === '' ? null : parseInt(projectId, 10);

      const result = await adminAssignProjectToUser(token, userId, actualProjectId, adminUsername);
      setGeneralSuccess(result.mensaje || 'Proyecto asignado/desasignado con éxito.');
      await fetchUsersAndProjects();
    } catch (err) {
      setGeneralError('Error al asignar proyecto: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleAdminDeleteUser = async (userId, username) => {
    if (userId === currentUser.id) return setGeneralError("No puedes borrar tu propia cuenta.");
    if (!window.confirm(`¿Estás seguro de que quieres borrar al usuario ${username}?`)) return;

    setGeneralError('');
    setGeneralSuccess('');
    setLoading(true);

    try {
      const result = await adminDeleteUser(token, userId, adminUsername);
      setGeneralSuccess(result.mensaje || 'Usuario borrado con éxito.');
      await fetchUsersAndProjects();
    } catch (err) {
      setGeneralError('Error al borrar usuario: ' + err.message);
    } finally {
      setLoading(false);
    }
  };


  return (
    <div style={styles.container}>
      <h2 style={styles.h2}>Perfiles de Usuarios (Admin/Gerente)</h2>

      {/* Formulario para Agregar Usuario (Solo Admin) */}
      {userRole === 'admin' && (
        <div style={styles.adminFormContainer}>
          <h3 style={styles.h3}>Agregar Nuevo Usuario</h3>
          <form onSubmit={handleAdminAddUser}>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1rem', marginBottom: '1rem' }}>
              <input type="text" name="username" placeholder="Usuario" value={newUser.username} onChange={handleInputChange} style={styles.input} required />
              <input type="password" name="password" placeholder="Contraseña" value={newUser.password} onChange={handleInputChange} style={styles.input} required />
              <input type="text" name="nombre" placeholder="Nombre" value={newUser.nombre} onChange={handleInputChange} style={styles.input} required />
              <input type="text" name="apellido" placeholder="Apellido" value={newUser.apellido} onChange={handleInputChange} style={styles.input} required />
              <input type="text" name="cedula" placeholder="Cédula" value={newUser.cedula} onChange={handleInputChange} style={styles.input} required />
              <select name="role" value={newUser.role} onChange={handleInputChange} style={styles.select}>
                <option value="user">User</option>
                <option value="encargado">Encargado</option>
                <option value="gerente">Gerente</option>
                <option value="admin">Admin</option>
              </select>
              <select name="proyecto_id" value={newUser.proyecto_id} onChange={handleInputChange} style={styles.select}>
                <option value="">No Asignar Proyecto</option>
                {projects.map(p => (
                  <option key={p.id} value={p.id}>{p.nombre}</option>
                ))}
              </select>
              <button type="submit" style={styles.button} disabled={loading}>
                {loading ? 'Cargando...' : 'Agregar Usuario'}
              </button>
            </div>
          </form>
          {generalError && <p style={styles.error}>{generalError}</p>}
          {generalSuccess && <p style={styles.success}>{generalSuccess}</p>}
        </div>
      )}


      {/* Tabla de Usuarios */}
      {loading ? (
        <p>Cargando usuarios...</p>
      ) : (
        <div style={styles.tableContainer}>
          <table style={styles.table}>
            <thead>
              <tr>
                <th style={styles.th}>ID</th>
                <th style={styles.th}>Usuario</th>
                <th style={styles.th}>Nombre</th>
                <th style={styles.th}>Apellido</th>
                {/*  APLICACIÓN DEL ESTILO AL ENCABEZADO  */}
                <th style={{ ...styles.th, ...styles.tdCedula }}>Cédula</th>
                <th style={styles.th}>Rol</th>
                <th style={styles.th}>Proyecto Asignado</th>
                <th style={styles.th} colSpan="2">Acciones</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => {
                const isSelf = user.id === userId;
                return (
                  <tr key={user.id}>
                    <td style={styles.td}>{user.id}</td>
                    <td style={styles.td}>{user.username}</td>
                    <td style={styles.td}>{user.nombre}</td>
                    <td style={styles.td}>{user.apellido}</td>
                    {/*  APLICACIÓN DEL ESTILO A LA CELDA DE DATOS  */}
                    <td style={{ ...styles.td, ...styles.tdCedula }}>{user.cedula}</td>
                    <td style={styles.td}>
                      <select
                        style={styles.roleSelect}
                        value={selectedRoles[user.id] || user.role}
                        onChange={(e) => handleRoleChange(user.id, e.target.value)}
                        disabled={isSelf || userRole !== 'admin'} // Solo Admin puede cambiar roles
                      >
                        <option value="user">User</option>
                        <option value="encargado">Encargado</option>
                        <option value="gerente">Gerente</option>
                        <option value="admin" disabled={isSelf}>Admin</option>
                      </select>
                    </td>
                    <td style={styles.td}>
                      <select
                        style={styles.projectSelect}
                        value={selectedProjects[user.id] || ''}
                        onChange={(e) => handleProjectChange(user.id, e.target.value)}
                        disabled={userRole !== 'admin' && userRole !== 'gerente'} // Admin/Gerente pueden asignar
                      >
                        <option value="">No Asignar</option>
                        {projects.map(p => (
                          <option key={p.id} value={p.id}>{p.nombre}</option>
                        ))}
                      </select>
                    </td>
                    <td style={{ ...styles.td, width: '150px' }}>
                      {(userRole === 'admin' || userRole === 'gerente') && !isSelf && (
                        <button
                          style={{ ...styles.buttonSave, margin: '0 0.5rem 0 0' }}
                          onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#d97706'}
                          onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#f59e0b'}
                          onClick={() => handleAdminAssignProject(user.id, selectedProjects[user.id] || '')}
                        >
                          Guardar Proy.
                        </button>
                      )}
                    </td>
                    <td style={{ ...styles.td, width: '150px' }}>
                      {(userRole === 'admin' && !isSelf) ? (
                        <div style={{ display: 'flex', gap: '0.5rem' }}>
                          <button
                            style={{ ...styles.buttonSave }}
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