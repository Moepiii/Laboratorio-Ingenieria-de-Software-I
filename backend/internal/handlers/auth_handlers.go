package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"proyecto/internal/auth"
	"proyecto/internal/logger"
	"proyecto/internal/models"
)

// 1. EL STRUCT DEL HANDLER
type AuthHandler struct {
	authSvc   auth.AuthService
	loggerSvc logger.LoggerService
}

// 2. EL CONSTRUCTOR DEL HANDLER
func NewAuthHandler(as auth.AuthService, ls logger.LoggerService) *AuthHandler {
	return &AuthHandler{
		authSvc:   as,
		loggerSvc: ls, // ⭐️ 3. INYECTAMOS EL SERVICIO
	}
}

// 3. LOS MÉTODOS (Handlers)

// SaludoHandler
func SaludoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "¡Hola desde el backend de Go!")
}

// RegisterHandler
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	lastID, err := h.authSvc.Register(user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.loggerSvc.Log(user.Username, "user", "REGISTRO", "Usuarios", int(lastID))

	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: fmt.Sprintf("Usuario '%s' (ID: %d) registrado con éxito.", user.Username, lastID)})
}

// LoginHandler
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	loginResponse, err := h.authSvc.Login(creds.Username, creds.Password)
	if err != nil {

		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	h.loggerSvc.Log(loginResponse.User.Username, loginResponse.Role, "INICIO DE SESIÓN", "Auth", loginResponse.UserId)

	respondWithJSON(w, http.StatusOK, loginResponse)
}
