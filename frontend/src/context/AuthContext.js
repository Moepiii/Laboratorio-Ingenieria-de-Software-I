import React, { createContext, useContext, useState, useEffect } from 'react';
import { loginUser, registerUser } from '../services/authService';

const AuthContext = createContext();

export const useAuth = () => {
    return useContext(AuthContext);
};

export const AuthProvider = ({ children }) => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [currentUser, setCurrentUser] = useState(null);
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
            const storedUser = localStorage.getItem('user');
            const storedRole = localStorage.getItem('role');
            const storedUserId = localStorage.getItem('userId');

            if (storedToken && storedUser && storedUser !== "undefined" && storedRole && storedUserId) {

                setToken(storedToken);
                setCurrentUser(JSON.parse(storedUser));
                setUserRole(storedRole);
                // ⭐️ ARREGLO: Parsea el 'userId' (que es un string) a un número
                setUserId(parseInt(storedUserId, 10));
                setIsLoggedIn(true);

            } else if (storedToken || storedUser || storedRole || storedUserId) {
                logout(); // Limpia datos parciales o corruptos
            }

        } catch (e) {
            console.error("Error al cargar localStorage, limpiando sesión.", e);
            logout();
        } finally {
            setLoading(false);
        }
    }, []); // Se ejecuta solo una vez

    // Función de Login
    const login = async (username, password) => {
        setLoading(true);
        setError('');
        setSuccessMessage('');
        try {
            // Llama al servicio (que se encarga del PascalCase)
            const data = await loginUser(username, password);

            // Verifica la respuesta del backend (que ya envía token)
            if (!data.token || !data.user || !data.role || !data.userId) {
                throw new Error('Respuesta de login incompleta (faltan token, user, role, o userId).');
            }

            setToken(data.token);
            setCurrentUser(data.user);
            setUserRole(data.role);
            setUserId(data.userId); // data.userId ya es un número
            setIsLoggedIn(true);

            // Guarda en localStorage
            localStorage.setItem('token', data.token);
            localStorage.setItem('user', JSON.stringify(data.user));
            localStorage.setItem('role', data.role);
            localStorage.setItem('userId', data.userId); // Se guarda como string

        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    // Función de Registro
    const register = async (username, password, nombre, apellido) => {
        setLoading(true);
        setError('');
        setSuccessMessage('');
        try {
            // Llama al servicio (que se encarga del PascalCase)
            const response = await registerUser(username, password, nombre, apellido);
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
        userId, // Este ahora será un NÚMERO gracias al useEffect
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