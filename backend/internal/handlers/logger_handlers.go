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
type LoggerHandler struct {
	authSvc   auth.AuthService
	loggerSvc logger.LoggerService
}

// 2. EL CONSTRUCTOR DEL HANDLER
func NewLoggerHandler(as auth.AuthService, ls logger.LoggerService) *LoggerHandler {
	return &LoggerHandler{
		authSvc:   as,
		loggerSvc: ls,
	}
}

//  3. ESTRUCTURAS DE PETICIÓN (DTOs Locales)

// Estructura para recibir el rango de fechas del Frontend
type DeleteLogsRangeRequest struct {
	AdminUsername string `json:"admin_username"`
	FechaInicio   string `json:"fecha_inicio"`
	FechaFin      string `json:"fecha_fin"`
}

// 4. LOS MÉTODOS (Handlers)

func (h *LoggerHandler) GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GetLogsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	// Verificación de permisos
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error verificando permisos")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "acceso denegado. se requiere rol de admin")
		return
	}

	logs, err := h.loggerSvc.GetLogs(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, logs)
}

func (h *LoggerHandler) DeleteLogsHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteLogsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Solo Admin
	perm, err := h.authSvc.CheckPermission(req.AdminUsername, "admin")
	if err != nil || !perm {
		respondWithError(w, http.StatusForbidden, "No autorizado. Solo el administrador puede borrar logs.")
		return
	}

	err = h.loggerSvc.DeleteLogs(req.IDs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Logueamos la acción
	h.loggerSvc.Log(req.AdminUsername, "admin", "ELIMINACIÓN", "Logs", 0)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Logs eliminados correctamente"})
}

func (h *LoggerHandler) DeleteLogsRangeHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteLogsRangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// 1. Validar Permisos
	perm, err := h.authSvc.CheckPermission(req.AdminUsername, "admin")
	if err != nil || !perm {
		respondWithError(w, http.StatusForbidden, "No autorizado. Solo admin puede borrar historial masivo.")
		return
	}

	// 2. Ejecutar borrado
	cantidad, err := h.loggerSvc.DeleteLogsByRange(req.FechaInicio, req.FechaFin)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 3. Auditoría

	logMsg := fmt.Sprintf("ELIMINACIÓN MASIVA (%d eventos entre %s y %s)", cantidad, req.FechaInicio, req.FechaFin)

	h.loggerSvc.Log(req.AdminUsername, "admin", logMsg, "Logs", 0)

	// 4. Responder
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{
		Mensaje: fmt.Sprintf("Se eliminaron %d eventos del historial correctamente.", cantidad),
	})
}
