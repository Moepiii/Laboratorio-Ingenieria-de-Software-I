// La URL base de tu API
const API_BASE_URL = 'http://localhost:8080/api';

/**
 * Función apiCall robusta (Versión Final)
 */
export const apiCall = async (endpoint, method, body = null, token = null) => {
    const url = `${API_BASE_URL}${endpoint}`;

    const headers = {
        'Content-Type': 'application/json',
    };

    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const config = {
        method: method,
        headers: headers,
    };

    if (body) {
        config.body = JSON.stringify(body);
    }

    let response; // Mueve la declaración aquí para usarla en el catch
    try {
        response = await fetch(url, config);
        const responseText = await response.text();

        // Si la respuesta está vacía Y el status es OK (ej. 204 No Content), devuelve null o un objeto vacío
        if (!responseText && response.ok) {
            // Devolver null es más explícito que undefined
            return null;
        }
        // Si la respuesta está vacía Y el status NO es OK, lanza error
        if (!responseText && !response.ok) {
            throw new Error(`Error ${response.status} ${response.statusText}: Respuesta vacía del servidor`);
        }

        // Intenta parsear como JSON
        let data;
        try {
            data = JSON.parse(responseText);
        } catch (e) {
            // Si falla el parseo Y el status es OK, algo raro pasó. Lanza error.
            if (response.ok) {
                throw new Error(`Respuesta inesperada no-JSON del servidor: ${responseText.substring(0, 100)}`); // Muestra parte del texto
            } else {
                // Si falla el parseo Y el status NO es OK, usa el texto como error.
                throw new Error(`Error ${response.status} ${response.statusText}: ${responseText.substring(0, 100)}`);
            }
        }

        // Si el parseo tiene éxito pero !response.ok, usa el mensaje de error del JSON
        if (!response.ok) {
            // Busca 'error' o 'mensaje' en la respuesta JSON
            throw new Error(data.error || data.mensaje || `Error ${response.status} ${response.statusText}`);
        }

        // Si todo está bien, devuelve los datos
        return data;

    } catch (error) {
        // Asegura que siempre lancemos un objeto Error
        const errorMessage = error instanceof Error ? error.message : String(error);
        console.error(`Error en apiCall a ${endpoint}:`, errorMessage);
        // Si tenemos el 'response', añade el status al mensaje
        if (response) {
            throw new Error(`Error ${response.status}: ${errorMessage}`);
        } else {
            throw new Error(errorMessage); // Error de red (fetch falló)
        }
    }
};


// --- Funciones de Autenticación (sin cambios) ---
export const loginUser = (username, password) => {
    return apiCall('/login', 'POST', { username: username, password: password });
};

export const registerUser = (username, password, nombre, apellido) => {
    const user = {
        username: username,
        password: password,
        nombre: nombre,
        apellido: apellido
    };
    return apiCall('/register', 'POST', user);
};