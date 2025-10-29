import React from 'react';
import AuthForm from '../components/auth/AuthForm';

// Estilos que antes estaban en App.js para centrar el formulario
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
};

const LoginPage = () => {
    return (
        <div style={styles.container}>
            <AuthForm />
        </div>
    );
};

export default LoginPage;