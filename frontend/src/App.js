import React, { useState, useEffect, useCallback } from 'react';

// Base URL para tu backend de Go
const API_BASE_URL = 'http://localhost:8080/api';

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

// --- Componente de Formulario de Login/Registro ---
const AuthForm = ({
  isRegisterMode,
  successMessage,
  error,
  loading,
  username,
  password,
  confirmPassword,
  nombre, // NUEVO
  apellido, // NUEVO
  setUsername,
  setPassword,
  setConfirmPassword,
  setNombre, // NUEVO
  setApellido, // NUEVO
  handleRegister,
  handleLogin,
  handleSwitchMode
}) => (
  <div style={styles.card}>
    <h2 style={{ fontSize: '1.875rem', fontWeight: '800', color: '#1f2937', textAlign: 'center', marginBottom: '1.5rem' }}>
      {isRegisterMode ? 'Registrar Nuevo Usuario' : 'Iniciar Sesión'}
    </h2>

    {successMessage && <p style={styles.success}>{successMessage}</p>}

    <form onSubmit={isRegisterMode ? handleRegister : handleLogin} style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>

      {/* Campo de Usuario */}
      <div style={styles.inputGroup}>
        <label htmlFor="username" style={styles.label}>
          Usuario
        </label>
        <div style={{ position: 'relative' }}>
          <span style={{ position: 'absolute', left: '0.75rem', top: '50%', transform: 'translateY(-50%)', color: '#818cf8', fontSize: '1.25rem' }}>👤</span>
          <input
            id="username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="Escribe tu usuario aquí"
            style={{ ...styles.input, paddingLeft: '2.5rem' }}
            required
            disabled={loading}
          />
        </div>
      </div>

      {/* --- NUEVOS CAMPOS DE NOMBRE Y APELLIDO (SOLO REGISTRO) --- */}
      {isRegisterMode && (
        <>
          <div style={styles.inputGroup}>
            <label htmlFor="nombre" style={styles.label}>
              Nombre
            </label>
            <input
              id="nombre"
              type="text"
              value={nombre}
              onChange={(e) => setNombre(e.target.value)}
              placeholder="Tu nombre"
              style={styles.input}
              required
              disabled={loading}
            />
          </div>
          <div style={styles.inputGroup}>
            <label htmlFor="apellido" style={styles.label}>
              Apellido
            </label>
            <input
              id="apellido"
              type="text"
              value={apellido}
              onChange={(e) => setApellido(e.target.value)}
              placeholder="Tu apellido"
              style={styles.input}
              required
              disabled={loading}
            />
          </div>
        </>
      )}
      {/* --- FIN DE NUEVOS CAMPOS --- */}


      {/* Campo de Contraseña */}
      <div style={styles.inputGroup}>
        <label htmlFor="password" style={styles.label}>
          Contraseña
        </label>
        <div style={{ position: 'relative' }}>
          <span style={{ position: 'absolute', left: '0.75rem', top: '50%', transform: 'translateY(-50%)', color: '#818cf8', fontSize: '1.25rem' }}>🔒</span>
          <input
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
            style={{ ...styles.input, paddingLeft: '2.5rem' }}
            required
            disabled={loading}
          />
        </div>
      </div>

      {/* Campo de Confirmar Contraseña (Solo en modo Registro) */}
      {isRegisterMode && (
        <div style={styles.inputGroup}>
          <label htmlFor="confirmPassword" style={styles.label}>
            Confirmar Contraseña
          </label>
          <div style={{ position: 'relative' }}>
            <span style={{ position: 'absolute', left: '0.75rem', top: '50%', transform: 'translateY(-50%)', color: '#818cf8', fontSize: '1.25rem' }}>🔑</span>
            <input
              id="confirmPassword"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              placeholder="Repite tu contraseña"
              style={{ ...styles.input, paddingLeft: '2.5rem' }}
              required
              disabled={loading}
            />
          </div>
        </div>
      )}

      {/* Mensaje de Error */}
      {error && (
        <p style={styles.error}>
          {error}
        </p>
      )}

      {/* Botón Principal */}
      <button
        type="submit"
        style={styles.button}
        onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#4338ca'}
        onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#4f46e5'}
        disabled={loading}
      >
        {loading ? (
          'Cargando...'
        ) : (
          <>
            <span style={{ marginRight: '0.5rem', fontSize: '1rem' }}>{isRegisterMode ? '✨' : '➡️'}</span>
            {isRegisterMode ? 'Registrarse' : 'Ingresar'}
          </>
        )}
      </button>

      {/* Cambiar modo (Login/Registro) */}
      <p
        style={styles.switchText}
        onClick={handleSwitchMode}
      >
        {isRegisterMode ? '¿Ya tienes una cuenta? Inicia sesión' : '¿No tienes cuenta? Regístrate aquí'}
      </p>
    </form>
  </div>
);

// --- Componente de Bienvenida para Usuarios Normales ---
const WelcomeMessage = ({ currentUser, handleLogout, loading, fetchError, backendMessage, fetchGoGreeting }) => {
  useEffect(() => {
    fetchGoGreeting();
  }, [fetchGoGreeting]);

  return (
    <div style={{ ...styles.card, maxWidth: '500px', textAlign: 'center', padding: '2.5rem' }}>
      <span style={{ fontSize: '3rem', color: '#10b981', display: 'block', margin: '0 auto 1rem' }}>✅</span>
      <h2 style={{ fontSize: '2.25rem', fontWeight: '800', color: '#047857', marginBottom: '1rem' }}>
        ¡Bienvenido, Usuario!
      </h2>
      <p style={{ fontSize: '1.125rem', color: '#4b5563', marginBottom: '1.5rem' }}>
        Has iniciado sesión como **{currentUser}**. Eres un usuario normal.
      </p>

      {/* Sección del Mensaje del Backend de Go */}
      <div style={{ padding: '1rem', border: '1px solid #6366f1', borderRadius: '8px', marginBottom: '1.5rem', backgroundColor: '#eef2ff' }}>
        <h3 style={{ fontSize: '1rem', fontWeight: '700', color: '#4f46e5', marginBottom: '0.5rem' }}>
          Mensaje del Endpoint de Saludo de Go:
        </h3>
        {loading ? (
          <p style={{ color: '#4f46e5' }}>Cargando...</p>
        ) : fetchError ? (
          <p style={{ color: '#dc2626', fontWeight: 'bold' }}>❌ Error: {backendMessage}</p>
        ) : (
          <p style={{ fontSize: '1.1rem', fontWeight: 'bold', color: '#1f2937' }}>{backendMessage}</p>
        )}
      </div>

      <button
        onClick={handleLogout}
        style={styles.welcomeButton}
        onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#dc2626'}
        onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#ef4444'}
      >
        <span style={{ marginRight: '0.5rem', fontSize: '1rem' }}>🚪</span>
        Cerrar Sesión
      </button>
    </div>
  );
};

// --- Componente de Dashboard para Administradores ---
const AdminDashboard = ({ currentUser, handleLogout, apiCall }) => {
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


  return (
    <div style={{ ...styles.card, maxWidth: '900px', width: '100%', padding: '2.5rem' }}>
      <h2 style={{ fontSize: '2rem', fontWeight: '800', color: '#4f46e5', marginBottom: '1.5rem', borderBottom: '2px solid #e5e7eb', paddingBottom: '0.5rem' }}>
        🔑 Panel de Administración
      </h2>
      <p style={{ fontSize: '1.125rem', color: '#4b5563', marginBottom: '1.5rem' }}>
        Administrador: **{currentUser}**.
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

      <button
        onClick={handleLogout}
        style={{ ...styles.welcomeButton, marginTop: '2rem', float: 'right', width: 'auto' }}
        onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#dc2626'}
        onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#ef4444'}
      >
        <span style={{ marginRight: '0.5rem', fontSize: '1rem' }}>🚪</span>
        Cerrar Sesión
      </button>
    </div>
  );
};


// Componente principal de la aplicación
const App = () => {
  // === ESTADOS CLAVE ===
  const [userRole, setUserRole] = useState(''); // 'user' o 'admin'
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [currentUser, setCurrentUser] = useState(''); // Almacena el usuario autenticado
  // ======================
  const [backendMessage, setBackendMessage] = useState('Cargando mensaje de Go...');
  const [loading, setLoading] = useState(false);
  const [fetchError, setFetchError] = useState(null);
  const [isRegisterMode, setIsRegisterMode] = useState(false);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState(''); // Para el registro
  const [nombre, setNombre] = useState(''); // NUEVO
  const [apellido, setApellido] = useState(''); // NUEVO
  const [error, setError] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  // Función de utilidad para alternar entre Login y Registro y limpiar estados de error/éxito
  const handleSwitchMode = () => {
    setIsRegisterMode(prev => !prev);
    setError('');
    setSuccessMessage('');
    setPassword('');
    setConfirmPassword('');
    setNombre(''); // Limpiamos
    setApellido(''); // Limpiamos
  };

  // FIX: Usamos useCallback para estabilizar apiCall, ya que es usada en otras funciones con useCallback.
  const apiCall = useCallback(async (endpoint, data, method = 'POST') => {
    let attempts = 0;
    const maxAttempts = 3;
    const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

    while (attempts < maxAttempts) {
      try {
        const config = {
          method: method,
          headers: {
            'Content-Type': 'application/json',
          },
        };

        if (method !== 'GET' && data) {
          config.body = JSON.stringify(data);
        }

        const response = await fetch(`${API_BASE_URL}/${endpoint}`, config);

        const responseData = (response.headers.get('content-type')?.includes('application/json'))
          ? await response.json()
          : { error: await response.text() || 'Respuesta vacía o no JSON' };


        if (response.ok) {
          return { success: true, data: responseData };
        } else {
          // Si el servidor responde con un error HTTP (4xx, 5xx), lo lanzamos
          throw new Error(responseData.error || `Error HTTP ${response.status}: ${response.statusText}`);
        }

      } catch (e) {
        attempts++;
        if (attempts >= maxAttempts) {
          throw new Error(`Fallo de conexión tras ${maxAttempts} intentos: ${e.message}`);
        }
        await delay(1000 * attempts);
      }
    }
    throw new Error('Máximo de reintentos alcanzado sin éxito.');
  }, []); // apiCall no depende de ningún estado o prop, así que se define una sola vez.


  /**
   * FIX: Usamos useCallback para estabilizar fetchGoGreeting.
   */
  const fetchGoGreeting = useCallback(async () => {
    setLoading(true);
    setFetchError(null);
    try {
      // Usamos el método GET para el saludo simple
      const response = await fetch(`${API_BASE_URL}/saludo`);

      if (!response.ok) {
        throw new Error(`Error HTTP: ${response.status}`);
      }

      const data = await response.json();
      setBackendMessage(data.mensaje);

    } catch (error) {
      console.error("Error al obtener saludo de Go:", error);
      setBackendMessage('No se pudo conectar con el backend de Go. Revisa la consola y CORS.');
      setFetchError(error.message);
    } finally {
      setLoading(false);
    }
  }, []); // No tiene dependencias internas, se crea una vez.

  /**
   * Maneja la lógica de inicio de sesión.
   */
  const handleLogin = async (e) => {
    e.preventDefault();
    setError('');
    setSuccessMessage('');
    setLoading(true);

    if (!username || !password) {
      setError('Por favor, ingresa usuario y contraseña.');
      setLoading(false);
      return;
    }

    try {
      const result = await apiCall('login', { username, password });

      if (result.success) {
        // Go devuelve el rol en el login
        const role = result.data.role || 'user';

        setIsLoggedIn(true);
        setCurrentUser(result.data.usuario);
        setUserRole(role);
        setError('');
        setPassword('');
      }
    } catch (e) {
      setError(e.message.includes('Credenciales inválidas') ? 'Usuario o contraseña incorrectos.' : e.message);
      setPassword('');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Maneja la lógica de registro de usuario. (ACTUALIZADO)
   */
  const handleRegister = async (e) => {
    e.preventDefault();
    setError('');
    setSuccessMessage('');
    setLoading(true);

    if (password !== confirmPassword) {
      setError('Las contraseñas no coinciden.');
      setLoading(false);
      return;
    }

    if (password.length < 6) {
      setError('La contraseña debe tener al menos 6 caracteres.');
      setLoading(false);
      return;
    }

    // NUEVA VALIDACIÓN
    if (!nombre || !apellido) {
      setError('El nombre y el apellido son obligatorios.');
      setLoading(false);
      return;
    }

    try {
      // ACTUALIZADO: Enviamos nombre y apellido
      const result = await apiCall('register', { username, password, nombre, apellido });

      if (result.success) {
        setSuccessMessage(result.data.mensaje + ". ¡Ya puedes iniciar sesión!");
        handleSwitchMode(); // Vuelve a la vista de login y limpia los campos
      }
    } catch (e) {
      setError(e.message.includes('El nombre de usuario ya existe') ? 'Ese nombre de usuario ya está registrado.' : e.message);
      setPassword('');
      setConfirmPassword('');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Cierra la sesión del usuario.
   */
  const handleLogout = () => {
    setIsLoggedIn(false);
    setUsername('');
    setPassword('');
    setConfirmPassword('');
    setNombre(''); // Limpiamos
    setApellido(''); // Limpiamos
    setError('');
    setSuccessMessage('');
    setCurrentUser('');
    setUserRole(''); // Limpiamos el rol
  };


  // Condicional Renderizado basado en el Rol
  return (
    <div style={styles.container}>
      {isLoggedIn ? (
        userRole === 'admin' ? (
          // Vista para Administradores
          <AdminDashboard
            currentUser={currentUser}
            handleLogout={handleLogout}
            apiCall={apiCall} // apiCall es estable gracias a useCallback
          />
        ) : (
          // Vista para Usuarios Normales
          <WelcomeMessage
            currentUser={currentUser}
            handleLogout={handleLogout}
            loading={loading}
            fetchError={fetchError}
            backendMessage={backendMessage}
            fetchGoGreeting={fetchGoGreeting} // fetchGoGreeting es estable gracias a useCallback
          />
        )
      ) : (
        // Vista de Login/Registro (ACTUALIZADO)
        <AuthForm
          isRegisterMode={isRegisterMode}
          successMessage={successMessage}
          error={error}
          loading={loading}
          username={username}
          password={password}
          confirmPassword={confirmPassword}
          nombre={nombre}           // NUEVO
          apellido={apellido}     // NUEVO
          setUsername={setUsername}
          setPassword={setPassword}
          setConfirmPassword={setConfirmPassword}
          setNombre={setNombre}       // NUEVO
          setApellido={setApellido}   // NUEVO
          handleRegister={handleRegister}
          handleLogin={handleLogin}
          handleSwitchMode={handleSwitchMode}
        />
      )}
    </div>
  );
};

export default App;