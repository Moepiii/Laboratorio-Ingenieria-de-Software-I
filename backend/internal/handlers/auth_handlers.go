package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"proyecto/internal/auth"   // Importamos la INTERFAZ del servicio
	"proyecto/internal/models" // Importamos los modelos
)

// --- 1. EL STRUCT DEL HANDLER ---
type AuthHandler struct {
	authSvc auth.AuthService
}

// --- 2. EL CONSTRUCTOR DEL HANDLER ---
func NewAuthHandler(as auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: as,
	}
}

// --- 3. LOS MÉTODOS (Handlers) ---

// SaludoHandler no necesita dependencias, puede quedar como estaba
func SaludoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "¡Hola desde el backend de Go!")
}

// RegisterHandler ahora es un método de AuthHandler
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// Dejamos que el servicio haga el trabajo
	lastID, err := h.authSvc.Register(user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: fmt.Sprintf("Usuario '%s' (ID: %d) registrado con éxito.", user.Username, lastID)})
}

// LoginHandler ahora es un método de AuthHandler
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// Dejamos que el servicio haga el trabajo
	loginResponse, err := h.authSvc.Login(creds.Username, creds.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, loginResponse)
}