import React, { useState, useEffect, useCallback } from 'react';
import AdminDashboard from './AdminDashboard';
import UserDashboard from './UserDashboard'; // ⭐️ 1. Importa el nuevo componente

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
  adminFormContainer: {
    padding: '1.5rem',
    backgroundColor: '#f9fafb', // gray-50
    borderRadius: '8px',
    marginBottom: '2rem',
    border: '1px solid #e5e7eb',
  }
};

// --- Componente de Formulario de Login/Registro (AuthForm) ---
const AuthForm = ({
  isRegisterMode,
  successMessage,
  error,
  loading,
  username,
  password,
  confirmPassword,
  nombre,
  apellido,
  setUsername,
  setPassword,
  setConfirmPassword,
  setNombre,
  setApellido,
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

      {error && (
        <p style={styles.error}>
          {error}
        </p>
      )}

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

      <p
        style={styles.switchText}
        onClick={handleSwitchMode}
      >
        {isRegisterMode ? '¿Ya tienes una cuenta? Inicia sesión' : '¿No tienes cuenta? Regístrate aquí'}
      </p>
    </form>
  </div>
);


// Componente principal de la aplicación
const App = () => {
  // === ESTADOS CLAVE ===
  const [userRole, setUserRole] = useState(''); // 'user', 'admin', o 'gerente'
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [currentUser, setCurrentUser] = useState('');
  const [userId, setUserId] = useState(null); // ⭐️ 2. Añade estado para userId
  // ======================
  // (Estados existentes - sin cambios)
  const [loading, setLoading] = useState(false);
  const [isRegisterMode, setIsRegisterMode] = useState(false);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [nombre, setNombre] = useState('');
  const [apellido, setApellido] = useState('');
  const [error, setError] = useState('');
  const [successMessage, setSuccessMessage] = useState('');

  // (handleSwitchMode sin cambios)
  const handleSwitchMode = () => {
    setIsRegisterMode(prev => !prev);
    setError('');
    setSuccessMessage('');
    setPassword('');
    setConfirmPassword('');
    setNombre('');
    setApellido('');
  };

  // (apiCall sin cambios)
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
  }, []);


  // ⭐️ 3. handleLogin AHORA GUARDA EL userId
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
        // Go devuelve el rol y el id en el login
        const role = result.data.role || 'user';
        const fetchedUserId = result.data.id; // Obtiene el ID del backend

        setIsLoggedIn(true);
        setCurrentUser(result.data.usuario);
        setUserRole(role);
        setUserId(fetchedUserId); // Guarda el ID en el estado
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

  // (handleRegister sin cambios)
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
    if (!nombre || !apellido) {
      setError('El nombre y el apellido son obligatorios.');
      setLoading(false);
      return;
    }

    try {
      const result = await apiCall('register', { username, password, nombre, apellido });
      if (result.success) {
        setSuccessMessage(result.data.mensaje + ". ¡Ya puedes iniciar sesión!");
        handleSwitchMode();
      }
    } catch (e) {
      setError(e.message.includes('El nombre de usuario ya existe') ? 'Ese nombre de usuario ya está registrado.' : e.message);
      setPassword('');
      setConfirmPassword('');
    } finally {
      setLoading(false);
    }
  };

  // ⭐️ 4. handleLogout AHORA LIMPIA userId
  const handleLogout = () => {
    setIsLoggedIn(false);
    setUsername('');
    setPassword('');
    setConfirmPassword('');
    setNombre('');
    setApellido('');
    setError('');
    setSuccessMessage('');
    setCurrentUser('');
    setUserRole('');
    setUserId(null); // Limpia el ID
  };


  // ⭐️ 5. --- CONDICIONAL DE RENDERIZADO ACTUALIZADO ---
  return (
    <div style={styles.container}>
      {isLoggedIn ? (
        // Si es admin o gerente -> AdminDashboard
        (userRole === 'admin' || userRole === 'gerente') ? (
          <AdminDashboard
            currentUser={currentUser}
            userRole={userRole} // Pasa el rol
            handleLogout={handleLogout}
            apiCall={apiCall}
          />
        ) : ( // Si es 'user' -> UserDashboard
          <UserDashboard
            currentUser={currentUser}
            userId={userId} // Pasa el ID del usuario
            apiCall={apiCall}
            handleLogout={handleLogout}
          />
        )
      ) : (
        // Si no está logueado -> AuthForm
        <AuthForm
          isRegisterMode={isRegisterMode}
          successMessage={successMessage}
          error={error}
          loading={loading}
          username={username}
          password={password}
          confirmPassword={confirmPassword}
          nombre={nombre}
          apellido={apellido}
          setUsername={setUsername}
          setPassword={setPassword}
          setConfirmPassword={setConfirmPassword}
          setNombre={setNombre}
          setApellido={setApellido}
          handleRegister={handleRegister}
          handleLogin={handleLogin}
          handleSwitchMode={handleSwitchMode}
        />
      )}
    </div>
  );
};

export default App;