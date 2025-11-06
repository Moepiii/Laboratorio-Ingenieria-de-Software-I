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
  formGrid: { display: 'grid', gridTemplateColumns: 'repeat(1, minmax(0, 1fr))', gap: '1rem' }, // Grid por defecto
  h2: { fontSize: '1.75rem', fontWeight: '700', color: '#1f2937', marginBottom: '1.5rem' },
  tableContainer: { overflowX: 'auto', borderRadius: '8px', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' },
  table: { width: '100%', borderCollapse: 'collapse', minWidth: '800px' },
  th: { padding: '0.75rem 1rem', textAlign: 'left', backgroundColor: '#f3f4f6', borderBottom: '2px solid #e5e7eb', color: '#374151', fontWeight: '600', fontSize: '0.875rem' },
  td: { padding: '0.75rem 1rem', borderBottom: '1px solid #e5e7eb', verticalAlign: 'middle' },
  select: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', backgroundColor: 'white' },
  buttonAssign: { padding: '0.4rem 0.8rem', fontSize: '0.875rem', fontWeight: '500', borderRadius: '6px', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s', width: '100px', marginTop: '0.25rem' },
  error: { backgroundColor: '#fee2e2', color: '#b91c1c', padding: '1rem', borderRadius: '8px', marginBottom: '1rem', fontSize: '0.9rem' },
  success: { backgroundColor: '#dcfce7', color: '#16a34a', padding: '1rem', borderRadius: '8px', marginBottom: '1rem', fontSize: '0.9rem' },
  inputGroup: { marginBottom: '0.5rem' },
  label: { display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' },
};

// Media query para el grid
if (window.innerWidth >= 768) {
  styles.formGrid = {
    display: 'grid',
    gridTemplateColumns: 'repeat(3, minmax(0, 1fr))', // 3 columnas en desktop
    gap: '1rem'
  };
}


const PerfilesUsuarios = () => {
  const { token, currentUser } = useAuth();
  const adminUsername = currentUser?.username;

  // Estados del componente
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [generalError, setGeneralError] = useState('');

  // Estados del formulario "Añadir"
  const [addingUser, setAddingUser] = useState(false);
  const [newUsername, setNewUsername] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [newNombre, setNewNombre] = useState('');
  const [newApellido, setNewApellido] = useState('');
  const [newCedula, setNewCedula] = useState('');
  const [addError, setAddError] = useState('');
  const [addSuccess, setAddSuccess] = useState('');

  // Estados para la asignación de roles y proyectos
  const [proyectos, setProyectos] = useState([]);
  const [selectedRoles, setSelectedRoles] = useState({});
  const [selectedProyectos, setSelectedProyectos] = useState({});

  // Carga inicial de datos
  const loadData = useCallback(async () => {
    if (!token || !adminUsername) return;
    setLoading(true);
    setGeneralError('');
    try {
      const [usersData, proyectosData] = await Promise.all([
        getAdminUsers(token, adminUsername),
        getAdminProjects(token, adminUsername)
      ]);

      const usersList = usersData.users || [];
      setUsers(usersList);
      setProyectos(proyectosData.proyectos || []);

      const initialRoles = {};
      const initialProyectos = {};
      usersList.forEach(user => {
        initialRoles[user.id] = user.role;
        initialProyectos[user.id] = user.proyecto_id || 0;
      });
      setSelectedRoles(initialRoles);
      setSelectedProyectos(initialProyectos);

    } catch (err) {
      setGeneralError(err.message || 'Error al cargar los datos.');
    } finally {
      setLoading(false);
    }
  }, [token, adminUsername]);

  useEffect(() => {
    loadData();
  }, [loadData]);


  // Helper para limpiar el formulario de añadir
  const clearForm = () => {
    setNewUsername('');
    setNewPassword('');
    setNewNombre('');
    setNewApellido('');
    setNewCedula('');
    setAddingUser(false);
    setAddError('');
    setAddSuccess('');
  };

  // --- Handlers ---

  const handleAdminAddUser = async (e) => {
    e.preventDefault();
    setAddError('');
    setAddSuccess('');

    const newUser = {
      username: newUsername,
      password: newPassword,
      nombre: newNombre,
      apellido: newApellido,
      cedula: newCedula
    };

    if (!newUser.username || !newUser.password || !newUser.nombre || !newUser.apellido || !newUser.cedula) {
      setAddError('Todos los campos (username, password, nombre, apellido y cédula) son requeridos.');
      return;
    }

    try {
      const response = await adminAddUser(token, newUser, adminUsername);
      setAddSuccess(response.mensaje || 'Usuario añadido con éxito.');
      loadData();
      clearForm();

    } catch (err) {
      setAddError(err.message || 'Error al añadir usuario.');
    }
  };

  const handleAdminDeleteUser = async (userId, username) => {
    if (username === currentUser.username) {
      setGeneralError('No puedes borrar tu propia cuenta.');
      return;
    }
    if (!window.confirm(`¿Estás seguro de que quieres borrar al usuario ${username}?`)) return;

    setGeneralError('');
    try {
      await adminDeleteUser(token, userId, adminUsername);
      loadData();
    } catch (err) {
      setGeneralError(err.message || 'Error al borrar usuario.');
    }
  };

  const handleAdminUpdateRole = async (userId, newRole) => {
    setGeneralError('');
    try {
      await adminUpdateUserRole(token, userId, newRole, adminUsername);
      setUsers(prevUsers =>
        prevUsers.map(user =>
          user.id === userId ? { ...user, role: newRole } : user
        )
      );
    } catch (err) {
      setGeneralError(err.message || 'Error al actualizar el rol.');
    }
  };

  const handleAdminAssignProject = async (userId, proyectoId) => {
    setGeneralError('');
    const pId = parseInt(proyectoId, 10);
    try {
      await adminAssignProjectToUser(token, userId, pId, adminUsername);
      setUsers(prevUsers =>
        prevUsers.map(user =>
          user.id === userId ? { ...user, proyecto_id: pId || null, proyecto_nombre: proyectos.find(p => p.id === pId)?.nombre || null } : user
        )
      );
    } catch (err) {
      setGeneralError(err.message || 'Error al asignar el proyecto.');
    }
  };

  // --- Renderizado ---

  if (loading) return <div style={{ padding: '2rem' }}>Cargando perfiles...</div>;

  return (
    <div style={{ padding: '2rem', fontFamily: 'Inter, sans-serif' }}>
      <h2 style={styles.h2}>Perfiles de Usuarios</h2>

      {generalError && <div style={styles.error}>{generalError}</div>}

      {!addingUser && (
        <button
          onClick={() => setAddingUser(true)}
          style={{ ...styles.button, width: 'auto', marginBottom: '2rem', backgroundColor: '#22c55e' }}
        >
          <i className="fas fa-plus" style={{ marginRight: '8px' }}></i>
          Añadir Nuevo Usuario
        </button>
      )}

      {addingUser && (
        <div style={styles.adminFormContainer}>
          <form onSubmit={handleAdminAddUser}>
            <div style={styles.formGrid}>
              <div style={styles.inputGroup}>
                <label style={styles.label}>Nombre</label>
                <input type="text" style={styles.input} value={newNombre} onChange={(e) => setNewNombre(e.target.value)} />
              </div>
              <div style={styles.inputGroup}>
                <label style={styles.label}>Apellido</label>
                <input type="text" style={styles.input} value={newApellido} onChange={(e) => setNewApellido(e.target.value)} />
              </div>
              <div style={styles.inputGroup}>
                <label style={styles.label}>Cédula</label>
                <input type="text" style={styles.input} value={newCedula} onChange={(e) => setNewCedula(e.target.value)} placeholder="V-12345678" />
              </div>
              <div style={styles.inputGroup}>
                <label style={styles.label}>Username (Login)</label>
                <input type="text" style={styles.input} value={newUsername} onChange={(e) => setNewUsername(e.target.value)} />
              </div>
              <div style={styles.inputGroup}>
                <label style={styles.label}>Password Temporal</label>
                <input type="password" style={styles.input} value={newPassword} onChange={(e) => setNewPassword(e.target.value)} />
              </div>
            </div>

            {addError && <p style={styles.error}>{addError}</p>}
            {addSuccess && <p style={styles.success}>{addSuccess}</p>}

            <div style={{ display: 'flex', gap: '1rem', marginTop: '1rem' }}>
              <button type="submit" style={{ ...styles.button, width: '150px' }}>Crear Usuario</button>
              <button type="button" onClick={clearForm} style={{ ...styles.button, width: '150px', backgroundColor: '#6b7280' }}>Cancelar</button>
            </div>
          </form>
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
                {/* ⭐️ LÍNEA CORREGIDA: ⭐️ */}
                <th style={styles.th}>Username</th>
                <th style={styles.th}>Nombre</th>
                <th style={styles.th}>Cédula</th>
                <th style={styles.th}>Proyecto Asignado</th>
                <th style={styles.th}>Rol</th>
                <th style={styles.th}>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => {
                const isSelf = user.username === currentUser.username;
                return (
                  <tr key={user.id}>
                    <td style={styles.td}>{user.id}</td>
                    <td style={styles.td}>{user.username}</td>
                    <td style={styles.td}>{user.nombre} {user.apellido}</td>
                    <td style={styles.td}>{user.cedula}</td>
                    <td style={styles.td}>
                      {isSelf ? (
                        <span style={{ fontSize: '0.875rem', color: '#6b7280' }}>N/A</span>
                      ) : (
                        <select
                          style={styles.select}
                          value={selectedProyectos[user.id] || 0}
                          onChange={(e) => {
                            const newProjId = parseInt(e.target.value, 10);
                            setSelectedProyectos(prev => ({ ...prev, [user.id]: newProjId }));
                            handleAdminAssignProject(user.id, newProjId);
                          }}
                        >
                          <option value="0">--- No Asignado ---</option>
                          {proyectos.map(p => (
                            <option key={p.id} value={p.id}>{p.nombre}</option>
                          ))}
                        </select>
                      )}
                    </td>
                    <td style={styles.td}>
                      {isSelf ? (
                        <span style={{ fontSize: '0.875rem', fontWeight: '500', color: '#1f2937' }}>{user.role}</span>
                      ) : (
                        <select
                          style={styles.select}
                          value={selectedRoles[user.id] || user.role}
                          onChange={(e) => setSelectedRoles(prev => ({ ...prev, [user.id]: e.target.value }))}
                        >
                          <option value="user">User</option>
                          <option value="encargado">Encargado</option>
                          <option value="gerente">Gerente</option>
                          <option value="admin">Admin</option>
                        </select>
                      )}
                    </td>
                    <td style={styles.td}>
                      {!isSelf ? (
                        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.25rem' }}>
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