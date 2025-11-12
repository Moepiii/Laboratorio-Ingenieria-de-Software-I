package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto/internal/auth"
	"proyecto/internal/logger" // ⭐️ Importamos el nuevo servicio
	"proyecto/internal/models"
)

// --- 1. EL STRUCT DEL HANDLER ---\
type LoggerHandler struct {
	authSvc   auth.AuthService
	loggerSvc logger.LoggerService
}

// --- 2. EL CONSTRUCTOR DEL HANDLER ---
func NewLoggerHandler(as auth.AuthService, ls logger.LoggerService) *LoggerHandler {
	return &LoggerHandler{
		authSvc:   as,
		loggerSvc: ls,
	}
}

// --- 3. LOS MÉTODOS (Handlers) ---

func (h *LoggerHandler) GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GetLogsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	// ⭐️ Verificación de permisos: Solo los 'admin' pueden ver los logs
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error verificando permisos")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "acceso denegado. se requiere rol de admin")
		return
	}

	// Llamamos al servicio de logger
	logs, err := h.loggerSvc.GetLogs(req)
	if err != nil {
		// El servicio ya logueó el error, solo respondemos
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Respondemos con los logs
	respondWithJSON(w, http.StatusOK, logs)
}
