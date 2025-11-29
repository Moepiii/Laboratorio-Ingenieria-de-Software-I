package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto/internal/auth"
	"proyecto/internal/logger"
	"proyecto/internal/models"
	"proyecto/internal/proyectos"
)

// --- 1. EL STRUCT DEL HANDLER ---
type ProyectoHandler struct {
	authSvc     auth.AuthService
	proyectoSvc proyectos.ProyectoService
	loggerSvc   logger.LoggerService
}

// --- 2. EL CONSTRUCTOR ---
func NewProyectoHandler(as auth.AuthService, ps proyectos.ProyectoService, ls logger.LoggerService) *ProyectoHandler {
	return &ProyectoHandler{
		authSvc:     as,
		proyectoSvc: ps,
		loggerSvc:   ls,
	}
}

// --- 3. ESTRUCTURAS DE PETICIÓN (DTOs Locales para evitar errores de importación) ---

type CreateProyectoRequest struct {
	Nombre        string `json:"nombre"`
	FechaInicio   string `json:"fecha_inicio"`
	FechaCierre   string `json:"fecha_cierre"`
	AdminUsername string `json:"admin_username"`
}

type UpdateProyectoRequest struct {
	ID            int    `json:"id"`
	Nombre        string `json:"nombre"`
	FechaInicio   string `json:"fecha_inicio"`
	FechaCierre   string `json:"fecha_cierre"`
	AdminUsername string `json:"admin_username"`
}

type DeleteProyectoRequest struct {
	ID            int    `json:"id"`
	AdminUsername string `json:"admin_username"`
}

type SetProyectoEstadoRequest struct {
	ID            int    `json:"id"`
	Estado        string `json:"estado"`
	AdminUsername string `json:"admin_username"`
}

// --- 4. LOS MÉTODOS (Handlers) ---

// AdminGetProyectosHandler: Obtiene la lista de todos los proyectos
func (h *ProyectoHandler) AdminGetProyectosHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AdminActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	// Verificamos permisos (admin o gerente)
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error verificando permisos")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "acceso denegado")
		return
	}

	proyectos, err := h.proyectoSvc.GetAllProyectos()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"proyectos": proyectos})
}

// AdminCreateProyectoHandler: Crea un nuevo proyecto
func (h *ProyectoHandler) AdminCreateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	// Solo Admin o Gerente
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !hasPermission {
		respondWithError(w, http.StatusForbidden, "No autorizado")
		return
	}

	newProj, err := h.proyectoSvc.CreateProyecto(req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Log
	h.loggerSvc.Log(req.AdminUsername, "admin", "CREACIÓN", "Proyectos", newProj.ID)

	respondWithJSON(w, http.StatusCreated, newProj)
}

// AdminUpdateProyectoHandler: Actualiza datos de un proyecto
func (h *ProyectoHandler) AdminUpdateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !hasPermission {
		respondWithError(w, http.StatusForbidden, "No autorizado")
		return
	}

	updatedProj, err := h.proyectoSvc.UpdateProyecto(req.ID, req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Log
	h.loggerSvc.Log(req.AdminUsername, "admin", "MODIFICACIÓN", "Proyectos", req.ID)

	respondWithJSON(w, http.StatusOK, updatedProj)
}

// AdminDeleteProyectoHandler: Elimina un proyecto
func (h *ProyectoHandler) AdminDeleteProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	// Solo Admin (borrar es crítico)
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

	// Log
	h.loggerSvc.Log(req.AdminUsername, "admin", "CAMBIO ESTADO", "Proyectos", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Estado actualizado"})
}
