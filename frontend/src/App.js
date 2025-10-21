import React, { useState, useEffect, useCallback } from 'react';
import AdminDashboard from './AdminDashboard';

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