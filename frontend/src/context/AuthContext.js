import React, { createContext, useContext, useState, useEffect } from 'react';
// ⭐️ MODIFICADO: Importamos la función actualizada de authService
import { loginUser, registerUser } from '../services/authService';

const AuthContext = createContext();

export const useAuth = () => {
    return useContext(AuthContext);
};

export const AuthProvider = ({ children }) => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [currentUser, setCurrentUser] = useState(null); // Esto guardará {username, nombre, apellido, cedula}
    const [userRole, setUserRole] = useState(null);
    const [userId, setUserId] = useState(null);
    const [token, setToken] = useState(null);

    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [successMessage, setSuccessMessage] = useState('');

    // Función de Logout
    const logout = () => {
        setIsLoggedIn(false);
        setCurrentUser(null);
        setUserRole(null);
        setUserId(null);
        setToken(null);

        localStorage.removeItem('token');
        localStorage.removeItem('user');
        localStorage.removeItem('role');
        localStorage.removeItem('userId');
    };

    // Cargar datos desde localStorage al iniciar
    useEffect(() => {
        setLoading(true);
        try {
            const storedToken = localStorage.getItem('token');
            const storedUser = localStorage.getItem('user'); // Esto ya contiene la cédula si el usuario inició sesión
            const storedRole = localStorage.getItem('role');
            const storedUserId = localStorage.getItem('userId');

            if (storedToken && storedUser && storedRole && storedUserId) {
                setIsLoggedIn(true);
                setToken(storedToken);
                setCurrentUser(JSON.parse(storedUser));
                setUserRole(storedRole);
                setUserId(parseInt(storedUserId, 10)); // Asegura que sea número
            }
        } catch (e) {
            console.error("Error al cargar datos de localStorage:", e);
            logout(); // Limpia el estado si hay un error
        } finally {
            setLoading(false);
        }
    }, []);

    // Función de Login
    const login = async (username, password) => {
        setLoading(true);
        setError('');
        setSuccessMessage('');
        try {
            // data = { token, user: {username, nombre, apellido, cedula}, role, userId }
            const data = await loginUser(username, password);

            // 1. Poner datos en el Estado
            setIsLoggedIn(true);
            setToken(data.token);
            setCurrentUser(data.user); // data.user ya incluye la cédula desde el backend
            setUserRole(data.role);
            setUserId(data.userId);

            // 2. Guarda en localStorage
            localStorage.setItem('token', data.token);
            localStorage.setItem('user', JSON.stringify(data.user)); // Se guarda el objeto user completo
            localStorage.setItem('role', data.role);
            localStorage.setItem('userId', data.userId); // Se guarda como string

        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    // ⭐️ MODIFICADO: La función 'register' ahora acepta y pasa la 'cedula'
    const register = async (username, password, nombre, apellido, cedula) => {
        setLoading(true);
        setError('');
        setSuccessMessage('');
        try {
            // Llama al servicio (que se encarga del PascalCase)
            const response = await registerUser(username, password, nombre, apellido, cedula);
            setSuccessMessage(response.mensaje || '¡Registro exitoso! Ahora puedes iniciar sesión.');

        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    // Define el "valor" que se pasará a los componentes hijos
    const value = {
        isLoggedIn,
        currentUser,
        userRole,
        userId,
        token,
        loading,
        error,
        successMessage,
        login,
        register,
        logout,
        setError,
        setSuccessMessage
    };

    return (
        <AuthContext.Provider value={value}>
            {!loading && children}
        </AuthContext.Provider>
    );
};