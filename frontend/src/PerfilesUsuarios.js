import React, { useState, useEffect, useCallback } from 'react';

// --- ¡IMPORTANTE! ---
// Copia y pega el objeto 'styles' completo desde tu archivo 'App.js' aquí.
// Lo necesitas para que el formulario y la tabla se vean bien.

// Estilos básicos en línea (inline styles) usando un diseño responsivo.
const styles = {
  container: {
    minHeight: '100vh',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: '#f3f4f6', // gray-100
    padding: '1rem',
    fontFamily: 'Inter, sans-serif',
  },
  card: {
    padding: '2rem',
    backgroundColor: '#ffffff',
    borderRadius: '12px',
    boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.25)', // shadow-2xl
    width: '100%',
    maxWidth: '400px',
    margin: 'auto', // Asegura el centrado en móvil
  },
  inputGroup: {
    marginBottom: '1.5rem',
  },
  label: {
    display: 'block',
    fontSize: '0.875rem',
    fontWeight: '500',
    color: '#374151', // gray-700
    marginBottom: '0.25rem',
  },
  input: {
    width: '100%',
    padding: '0.75rem 1rem',
    border: '1px solid #d1d5db',
    borderRadius: '8px',
    fontSize: '1rem',
    transition: 'border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out',
    boxSizing: 'border-box', // Importante para responsividad
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
    backgroundColor: '#4f46e5', // indigo-600
    border: 'none',
    cursor: 'pointer',
    transition: 'background-color 0.2s, transform 0.1s',
    boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)',
  },
  error: {
    fontSize: '0.875rem',
    color: '#dc2626', // red-600
    fontWeight: '500',
    backgroundColor: '#fef2f2', // red-50
    padding: '0.75rem',
    borderRadius: '8px',
    border: '1px solid #fecaca', // red-200
  },
  success: {
    fontSize: '0.875rem',
    color: '#059669', // emerald-600
    fontWeight: '500',
    backgroundColor: '#ecfdf5', // emerald-50
    padding: '0.75rem',
    borderRadius: '8px',
    border: '1px solid #a7f3d0', // emerald-200
  },
  welcomeButton: {
    display: 'inline-flex',
    alignItems: 'center',
    padding: '0.75rem 1.5rem',
    fontSize: '1rem',
    fontWeight: '600',
    borderRadius: '8px',
    color: 'white',
    backgroundColor: '#ef4444', // red-500
    border: 'none',
    cursor: 'pointer',
    transition: 'background-color 0.2s, transform 0.1s',
  },
  switchText: {
    fontSize: '0.875rem',
    textAlign: 'center',
    color: '#6b7280',
    marginTop: '1rem',
    cursor: 'pointer',
    textDecoration: 'underline',
    fontWeight: '600',
  },
  // Nuevos estilos para el dashboard
  adminFormContainer: {
    padding: '1.5rem',
    backgroundColor: '#f9fafb', // gray-50
    borderRadius: '8px',
    marginBottom: '2rem',
    border: '1px solid #e5e7eb',
  }
};

const PerfilesUsuarios = ({ currentUser, apiCall }) => {
  // --- Comienzo del código que me pasaste ---
  const [users, setUsers] = useState([]);
  const [loadingUsers, setLoadingUsers] = useState(false);
  const [adminError, setAdminError] = useState('');
  const [adminSuccess, setAdminSuccess] = useState('');

  // NUEVOS ESTADOS para el formulario de Añadir Usuario
  const [newUsername, setNewUsername] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [newNombre, setNewNombre] = useState(''); // NUEVO
  const [newApellido, setNewApellido] = useState(''); // NUEVO
  const [loadingAdd, setLoadingAdd] = useState(false);


  // Función para obtener la lista de usuarios
  const fetchUsers = useCallback(async () => {
    setLoadingUsers(true);
    setAdminError('');
    try {
      const result = await apiCall('admin/users', { admin_username: currentUser }, 'POST');
      if (result.success) {
        // Filtramos el usuario actual y lo ponemos de primero
        const filteredUsers = (result.data.users || []).filter(u => u.username !== currentUser);
        const selfUser = (result.data.users || []).find(u => u.username === currentUser);

        // Los usuarios ahora tienen .nombre y .apellido
        setUsers(selfUser ? [selfUser, ...filteredUsers] : filteredUsers);
      } else {
        setAdminError('No se pudo cargar la lista de usuarios: ' + (result.data.error || 'Desconocido'));
      }
    } catch (e) {
      setAdminError(e.message.includes('Acceso denegado') ? 'Error: Acceso Denegado (Go no te reconoce como Admin)' : `Error de conexión: ${e.message}`);
    } finally {
      setLoadingUsers(false);
    }
  }, [currentUser, apiCall]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);


  // --- FUNCIÓN: Añadir Usuario (ACTUALIZADA) ---
  const handleAdminAddUser = async (e) => {
    e.preventDefault();
    setAdminError('');
    setAdminSuccess('');
    setLoadingAdd(true);

    if (newPassword.length < 6) {
      setAdminError('La contraseña debe tener al menos 6 caracteres.');
      setLoadingAdd(false);
      return;
    }

    try {
      // ACTUALIZADO: Enviamos nombre y apellido
      const result = await apiCall('admin/add-user', {
        username: newUsername,
        password: newPassword,
        nombre: newNombre,
        apellido: newApellido,
        admin_username: currentUser
      }, 'POST');

      if (result.success) {
        setAdminSuccess(result.data.mensaje);
        setNewUsername('');
        setNewPassword('');
        setNewNombre(''); // Limpiamos
        setNewApellido(''); // Limpiamos
        fetchUsers(); // Volver a cargar la lista
      } else {
        setAdminError(result.data.error || 'Fallo al añadir el usuario.');
      }
    } catch (e) {
      setAdminError(`Error de conexión al añadir usuario: ${e.message}`);
    } finally {
      setLoadingAdd(false);
    }
  };


  // --- FUNCIÓN: Borrar Usuario ---
  const handleAdminDeleteUser = async (userId, username) => {
    setAdminError('');
    setAdminSuccess('');

    if (!window.confirm(`¿Estás seguro de borrar permanentemente a ${username} (ID: ${userId})?`)) {
      return;
    }

    try {
      const result = await apiCall('admin/delete-user', {
        id: userId,
        admin_username: currentUser
      }, 'POST');

      if (result.success) {
        setAdminSuccess(result.data.mensaje);
        fetchUsers(); // Volver a cargar la lista
      } else {
        setAdminError(result.data.error || 'Fallo al borrar el usuario.');
      }
    } catch (e) {
      setAdminError(`Error de conexión al borrar usuario: ${e.message}`);
    }
  };


  // Función para cambiar el rol de un usuario
  const handleAdminUpdateRole = async (userId, username, currentRole) => {
    const newRole = currentRole === 'admin' ? 'user' : 'admin';
    setAdminError('');
    setAdminSuccess('');

    if (window.confirm(`¿Estás seguro de cambiar el rol de ${username} a "${newRole}"?`)) {
      try {
        const result = await apiCall('admin/update-user', {
          id: userId,
          new_role: newRole,
          admin_username: currentUser
        }, 'POST');

        if (result.success) {
          setAdminSuccess(`Rol de ${username} actualizado a ${newRole.toUpperCase()} con éxito.`);
          fetchUsers(); // Volver a cargar la lista
        } else {
          setAdminError(result.data.error || 'Fallo al actualizar el rol.');
        }
      } catch (e) {
        setAdminError(`Error de conexión al actualizar rol: ${e.message}`);
      }
    }
  };
  // --- Fin del código que me pasaste ---


  return (
    <div style={{ ...styles.card, maxWidth: '900px', width: '100%', padding: '2.5rem' }}>
      
      {/* Título de la vista */}
      <h2 style={{ fontSize: '2rem', fontWeight: '800', color: '#4f46e5', marginBottom: '1.5rem', borderBottom: '2px solid #e5e7eb', paddingBottom: '0.5rem' }}>
        🔑 Perfiles de usuarios
      </h2>
      <p style={{ fontSize: '1.125rem', color: '#4b5563', marginBottom: '1.5rem' }}>
        Administrador: <strong>{currentUser}</strong>.
      </p>

      {adminSuccess && <p style={{ ...styles.success, marginBottom: '1rem' }}>{adminSuccess}</p>}
      {adminError && <p style={{ ...styles.error, marginBottom: '1rem' }}>{adminError}</p>}

      {/* --- Formulario para Añadir Usuario (ACTUALIZADO) --- */}
      <div style={styles.adminFormContainer}>
        <h3 style={{ fontSize: '1.25rem', fontWeight: '700', color: '#1f2937', marginBottom: '1rem' }}>➕ Crear Nuevo Usuario</h3>
        <form onSubmit={handleAdminAddUser} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>

          {/* Fila 1: Usuario y Contraseña */}
          <div style={{ display: 'flex', flexDirection: window.innerWidth > 600 ? 'row' : 'column', gap: '1rem' }}>
            <input
              type="text"
              placeholder="Nombre de Usuario"
              value={newUsername}
              onChange={(e) => setNewUsername(e.target.value)}
              required
              style={{ ...styles.input, flex: 1 }}
              disabled={loadingAdd}
            />
            <input
              type="password"
              placeholder="Contraseña (mín. 6 caracteres)"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              required
              style={{ ...styles.input, flex: 1 }}
              disabled={loadingAdd}
            />
          </div>

          {/* Fila 2: Nombre y Apellido (NUEVO) */}
          <div style={{ display: 'flex', flexDirection: window.innerWidth > 600 ? 'row' : 'column', gap: '1rem' }}>
            <input
              type="text"
              placeholder="Nombre"
              value={newNombre}
              onChange={(e) => setNewNombre(e.target.value)}
              required
              style={{ ...styles.input, flex: 1 }}
              disabled={loadingAdd}
            />
            <input
              type="text"
              placeholder="Apellido"
              value={newApellido}
              onChange={(e) => setNewApellido(e.target.value)}
              required
              style={{ ...styles.input, flex: 1 }}
              disabled={loadingAdd}
            />
          </div>

          {/* Fila 3: Botón */}
          <button
            type="submit"
            style={{ ...styles.button, width: '100%', backgroundColor: '#10b981' }} // emerald-500
            onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#059669'} // emerald-600
            onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#10b981'}
            disabled={loadingAdd}
          >
            {loadingAdd ? 'Creando...' : 'Crear Usuario'}
          </button>
        </form>
      </div>

      {/* --- Tabla de Usuarios (ACTUALIZADA) --- */}
      <h3 style={{ fontSize: '1.25rem', fontWeight: '700', color: '#1f2937', marginBottom: '1rem' }}>
        Lista de Usuarios (Total: {users.length})
      </h3>

      {loadingUsers ? (
        <p style={{ textAlign: 'center', padding: '2rem', color: '#4f46e5' }}>Cargando usuarios...</p>
      ) : (
        <div style={{ maxHeight: '400px', overflowY: 'auto', border: '1px solid #e5e7eb', borderRadius: '8px' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left' }}>
            <thead>
              <tr style={{ backgroundColor: '#f9fafb', borderBottom: '1px solid #e5e7eb' }}>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>ID</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Usuario</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Nombre</th>{/* NUEVO */}
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Apellido</th>{/* NUEVO */}
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280' }}>Rol Actual</th>
                <th style={{ padding: '0.75rem', fontSize: '0.875rem', color: '#6b7280', textAlign: 'center' }}>Acciones</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user, index) => {
                const isSelf = user.username === currentUser;
                return (
                  <tr key={user.id} style={{ borderBottom: index < users.length - 1 ? '1px solid #f3f4f6' : 'none', backgroundColor: isSelf ? '#fffbeb' : (user.role === 'admin' ? '#f5f3ff' : 'white') }}>
                    <td style={{ padding: '0.75rem', fontWeight: '500' }}>{user.id}</td>
                    <td style={{ padding: '0.75rem', fontWeight: '600' }}>{user.username} {isSelf && <span style={{ fontSize: '0.75rem', color: '#f59e0b' }}>(Tú)</span>}</td>
                    <td style={{ padding: '0.75rem' }}>{user.nombre}</td>{/* NUEVO */}
                    <td style={{ padding: '0.75rem' }}>{user.apellido}</td>{/* NUEVO */}
                    <td style={{ padding: '0.75rem' }}>
                      <span style={{
                        padding: '0.25rem 0.75rem',
                        borderRadius: '9999px',
                        fontSize: '0.75rem',
                        fontWeight: '700',
                        backgroundColor: user.role === 'admin' ? '#eef2ff' : '#d1fae5',
                        color: user.role === 'admin' ? '#4f46e5' : '#065f46',
                      }}>
                        {user.role.toUpperCase()}
                      </span>
                    </td>
                    <td style={{ padding: '0.75rem', textAlign: 'center' }}>

                      {/* Botón de CAMBIAR ROL */}
                      <button
                        onClick={() => handleAdminUpdateRole(user.id, user.username, user.role)}
                        style={{
                          padding: '0.5rem 1rem',
                          borderRadius: '6px',
                          fontSize: '0.875rem',
                          fontWeight: '600',
                          marginRight: '0.5rem',
                          backgroundColor: user.role === 'admin' ? '#f87171' : '#34d399', // Red for demote, Green for promote
                          color: 'white',
                          border: 'none',
                          cursor: isSelf ? 'not-allowed' : 'pointer',
                          opacity: isSelf ? 0.5 : 1,
                          transition: 'background-color 0.2s',
                        }}
                        onMouseOver={(e) => { if (!isSelf) e.currentTarget.style.backgroundColor = user.role === 'admin' ? '#ef4444' : '#059669' }}
                        onMouseOut={(e) => { if (!isSelf) e.currentTarget.style.backgroundColor = user.role === 'admin' ? '#f87171' : '#34d399' }}
                        disabled={isSelf}
                      >
                        {user.role === 'admin' ? 'Degradar a USER' : 'Promover a ADMIN'}
                      </button>

                      {/* Botón de BORRAR */}
                      <button
                        onClick={() => handleAdminDeleteUser(user.id, user.username)}
                        style={{
                          padding: '0.5rem 1rem',
                          borderRadius: '6px',
                          fontSize: '0.875rem',
                          fontWeight: '600',
                          backgroundColor: '#ef4444', // Red-500
                          color: 'white',
                          border: 'none',
                          cursor: isSelf ? 'not-allowed' : 'pointer',
                          opacity: isSelf ? 0.5 : 1,
                          transition: 'background-color 0.2s',
                        }}
                        onMouseOver={(e) => { if (!isSelf) e.currentTarget.style.backgroundColor = '#dc2626' }}
                        onMouseOut={(e) => { if (!isSelf) e.currentTarget.style.backgroundColor = '#ef4444' }}
                        disabled={isSelf}
                      >
                        Borrar
                      </button>
                    </td>
                  </tr>
                )
              })}
              {users.length === 0 && !loadingUsers && (
                /* ACTUALIZADO: colspan es 6 ahora */
                <tr><td colSpan="6" style={{ textAlign: 'center', padding: '1rem', color: '#6b7280' }}>No hay usuarios para mostrar.</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      {/* BOTÓN DE LOGOUT ELIMINADO: 
        Ya no es necesario aquí, porque está en el Sidebar.
      */}
      
    </div>
  );
};

export default PerfilesUsuarios;