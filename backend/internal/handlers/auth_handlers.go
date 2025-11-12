package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"proyecto/internal/auth"
	"proyecto/internal/logger" // ⭐️ 1. IMPORTAMOS EL LOGGER
	"proyecto/internal/models"
)

// --- 1. EL STRUCT DEL HANDLER ---\
type AuthHandler struct {
	authSvc   auth.AuthService
	loggerSvc logger.LoggerService // ⭐️ 2. AÑADIMOS EL SERVICIO DE LOGGER
}

// --- 2. EL CONSTRUCTOR DEL HANDLER ---\
func NewAuthHandler(as auth.AuthService, ls logger.LoggerService) *AuthHandler {
	return &AuthHandler{
		authSvc:   as,
		loggerSvc: ls, // ⭐️ 3. INYECTAMOS EL SERVICIO
	}
}

// --- 3. LOS MÉTODOS (Handlers) ---\

// SaludoHandler no necesita dependencias, puede quedar como estaba
func SaludoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "¡Hola desde el backend de Go!")
}

// RegisterHandler ahora es un método de AuthHandler
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	// Dejamos que el servicio haga el trabajo
	lastID, err := h.authSvc.Register(user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// ⭐️ 4. REGISTRAMOS EL EVENTO
	// (El rol por defecto 'user' se asigna en el servicio de auth)
	h.loggerSvc.Log(user.Username, "user", "REGISTRO", "Usuarios", int(lastID))

	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: fmt.Sprintf("Usuario '%s' (ID: %d) registrado con éxito.", user.Username, lastID)})
}

// LoginHandler ahora es un método de AuthHandler
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	// Dejamos que el servicio haga el trabajo
	loginResponse, err := h.authSvc.Login(creds.Username, creds.Password)
	if err != nil {
		// Aquí NO logueamos el fallo en la bitácora (para evitar spam)
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// ⭐️ 4. REGISTRAMOS EL EVENTO (SOLO SI ES EXITOSO)
	h.loggerSvc.Log(loginResponse.User.Username, loginResponse.Role, "INICIO DE SESIÓN", "Auth", loginResponse.UserId)

	respondWithJSON(w, http.StatusOK, loginResponse)
}
