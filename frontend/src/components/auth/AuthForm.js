import React, { useState } from 'react';
import { useAuth } from '../../context/AuthContext'; // Ajusta la ruta si es necesario



const styles = {
    card: {
        padding: '2rem',
        backgroundColor: '#ffffff',
        borderRadius: '12px',
        boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
        width: '100%',
        maxWidth: '400px',
        margin: 'auto',
    },
    inputGroup: { marginBottom: '1.5rem' },
    label: { display: 'block', fontSize: '0.875rem', fontWeight: '500', color: '#374151', marginBottom: '0.5rem' },
    input: { width: '100%', padding: '0.75rem 1rem', border: '1px solid #d1d5db', borderRadius: '8px', fontSize: '1rem', boxSizing: 'border-box' },
    button: { width: '100%', padding: '0.75rem 1rem', fontSize: '1rem', fontWeight: '600', borderRadius: '8px', color: 'white', backgroundColor: '#4f46e5', border: 'none', cursor: 'pointer', transition: 'background-color 0.2s' },
    switchButton: { color: '#4f46e5', cursor: 'pointer', marginTop: '1.5rem', textAlign: 'center', fontSize: '0.9rem' },
    error: { backgroundColor: '#fee2e2', color: '#b91c1c', padding: '1rem', borderRadius: '8px', marginBottom: '1rem', fontSize: '0.9rem' },
    success: { backgroundColor: '#dcfce7', color: '#16a34a', padding: '1rem', borderRadius: '8px', marginBottom: '1rem', fontSize: '0.9rem' },
};


const AuthForm = () => {

    const [isRegisterMode, setIsRegisterMode] = useState(false);
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [nombre, setNombre] = useState('');
    const [apellido, setApellido] = useState('');
    const [cedula, setCedula] = useState(''); 

    const { login, register, loading, error, successMessage, setError, setSuccessMessage } = useAuth();

    const handleSwitchMode = () => {
        setIsRegisterMode(!isRegisterMode);
 
        setUsername('');
        setPassword('');
        setConfirmPassword('');
        setNombre('');
        setApellido('');
        setCedula(''); 
        setError('');
        setSuccessMessage('');
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        setSuccessMessage('');

        if (isRegisterMode) {
            // --- Modo Registro ---
            if (password !== confirmPassword) {
                setError('Las contraseñas no coinciden.');
                return;
            }

            await register(username, password, nombre, apellido, cedula);

        } else {
            // --- Modo Login ---
            await login(username, password);
        }
    };

    return (
        <div style={styles.card}>
            <h2 style={{ textAlign: 'center', color: '#1f2937', marginBottom: '2rem' }}>
                {isRegisterMode ? 'Crear Cuenta' : 'Iniciar Sesión'}
            </h2>

            {/* Muestra de Errores y Éxito */}
            {error && <div style={styles.error}>{error}</div>}
            {successMessage && <div style={styles.success}>{successMessage}</div>}

            <form onSubmit={handleSubmit}>
                {/* Campos de Registro Adicionales */}
                {isRegisterMode && (
                    <>
                        <div style={styles.inputGroup}>
                            <label htmlFor="nombre" style={styles.label}>Nombre</label>
                            <input id="nombre" type="text" value={nombre} onChange={(e) => setNombre(e.target.value)} required style={styles.input} />
                        </div>
                        <div style={styles.inputGroup}>
                            <label htmlFor="apellido" style={styles.label}>Apellido</label>
                            <input id="apellido" type="text" value={apellido} onChange={(e) => setApellido(e.target.value)} required style={styles.input} />
                        </div>
                        {/*  */}
                        <div style={styles.inputGroup}>
                            <label htmlFor="cedula" style={styles.label}>Cédula</label>
                            <input id="cedula" type="text" value={cedula} onChange={(e) => setCedula(e.target.value)} required style={styles.input} placeholder="V-12345678" />
                        </div>
                    </>
                )}

                {/* Campos Comunes (Username, Password) */}
                <div style={styles.inputGroup}>
                    <label htmlFor="username" style={styles.label}>Usuario</label>
                    <input id="username" type="text" value={username} onChange={(e) => setUsername(e.target.value)} required style={styles.input} />
                </div>
                <div style={styles.inputGroup}>
                    <label htmlFor="password" style={styles.label}>Contraseña</label>
                    <input id="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} required style={styles.input} />
                </div>

                {/* Campo de Confirmar Contraseña (Solo Registro) */}
                {isRegisterMode && (
                    <div style={styles.inputGroup}>
                        <label htmlFor="confirmPassword" style={styles.label}>Confirmar Contraseña</label>
                        <input id="confirmPassword" type="password" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} required style={styles.input} />
                    </div>
                )}

                <button
                    type="submit"
                    disabled={loading}
                    style={{ ...styles.button, opacity: loading ? 0.7 : 1 }}
                    onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#4338ca'}
                    onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#4f46e5'}
                >
                    {loading ? 'Cargando...' : (isRegisterMode ? 'Registrarse' : 'Iniciar Sesión')}
                </button>
            </form>

            <p onClick={handleSwitchMode} style={styles.switchButton}>
                {isRegisterMode ? '¿Ya tienes cuenta? Inicia Sesión' : '¿No tienes cuenta? Regístrate'}
            </p>
        </div>
    );
};

export default AuthForm;