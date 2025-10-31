import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import { AuthProvider } from './context/AuthContext'; // 1. Importa el provider



const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    {/* 2. Envuelve <App /> con el <AuthProvider /> */}
    <AuthProvider>
      <App />
    </AuthProvider>
  </React.StrictMode>
);