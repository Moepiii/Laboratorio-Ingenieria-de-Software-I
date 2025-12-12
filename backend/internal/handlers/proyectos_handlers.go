package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto/internal/auth"
	"proyecto/internal/logger"
	"proyecto/internal/models"
	"proyecto/internal/proyectos"
)

// 1. EL STRUCT DEL HANDLER
type ProyectoHandler struct {
	authSvc     auth.AuthService
	proyectoSvc proyectos.ProyectoService
	loggerSvc   logger.LoggerService
}

// 2. EL CONSTRUCTOR
func NewProyectoHandler(as auth.AuthService, ps proyectos.ProyectoService, ls logger.LoggerService) *ProyectoHandler {
	return &ProyectoHandler{
		authSvc:     as,
		proyectoSvc: ps,
		loggerSvc:   ls,
	}
}

//  3. LOS MÉTODOS (Handlers)

// GetProyectosHandler: Obtiene la lista de proyectos
func (h *ProyectoHandler) GetProyectosHandler(w http.ResponseWriter, r *http.Request) {

	var req models.AdminActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	// Validamos que sea admin o gerente
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !hasPermission {
		respondWithError(w, http.StatusForbidden, "No autorizado")
		return
	}

	proyectos, err := h.proyectoSvc.GetAllProyectos()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"proyectos": proyectos})
}

// CreateProyectoHandler: Crea un nuevo proyecto
func (h *ProyectoHandler) CreateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !hasPermission {
		respondWithError(w, http.StatusForbidden, "No autorizado")
		return
	}

	nuevoProyecto, err := h.proyectoSvc.CreateProyecto(req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Log
	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "CREACIÓN", "Proyectos", nuevoProyecto.ID)

	respondWithJSON(w, http.StatusCreated, nuevoProyecto)
}

// UpdateProyectoHandler: Actualiza un proyecto
func (h *ProyectoHandler) UpdateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !hasPermission {
		respondWithError(w, http.StatusForbidden, "No autorizado")
		return
	}

	proyectoActualizado, err := h.proyectoSvc.UpdateProyecto(req.ID, req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Log
	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "MODIFICACIÓN", "Proyectos", req.ID)

	respondWithJSON(w, http.StatusOK, proyectoActualizado)
}

// DeleteProyectoHandler: Elimina un proyecto
func (h *ProyectoHandler) DeleteProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	// Solo Admin puede borrar proyectos
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin")
	if err != nil || !hasPermission {
		respondWithError(w, http.StatusForbidden, "No autorizado (Solo Admin)")
		return
	}

	_, err = h.proyectoSvc.DeleteProyecto(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Log
	h.loggerSvc.Log(req.AdminUsername, "admin", "ELIMINACIÓN", "Proyectos", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Proyecto eliminado"})
}

// AdminSetProyectoEstadoHandler: Cambia estado (Activo/Cerrado)
func (h *ProyectoHandler) AdminSetProyectoEstadoHandler(w http.ResponseWriter, r *http.Request) {
	// Definimos struct local para esta petición específica
	type SetProyectoEstadoRequest struct {
		ID            int    `json:"id"`
		Estado        string `json:"estado"`
		AdminUsername string `json:"admin_username"`
	}

	var req SetProyectoEstadoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !hasPermission {
		respondWithError(w, http.StatusForbidden, "No autorizado")
		return
	}

	_, err = h.proyectoSvc.SetProyectoEstado(req.ID, req.Estado)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "CAMBIO ESTADO", "Proyectos", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Estado actualizado"})
}
