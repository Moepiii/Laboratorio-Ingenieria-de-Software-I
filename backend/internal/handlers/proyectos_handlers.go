package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"proyecto/internal/auth"
	"proyecto/internal/models"
	"proyecto/internal/proyectos" // ⭐️ NUEVO: Importamos el servicio
)

// --- 1. EL STRUCT DEL HANDLER ---
type ProyectoHandler struct {
	authSvc     auth.AuthService
	proyectoSvc proyectos.ProyectoService
}

// --- 2. EL CONSTRUCTOR DEL HANDLER ---
func NewProyectoHandler(as auth.AuthService, ps proyectos.ProyectoService) *ProyectoHandler {
	return &ProyectoHandler{
		authSvc:     as,
		proyectoSvc: ps,
	}
}

// --- 3. LOS MÉTODOS (Handlers) ---

func (h *ProyectoHandler) AdminGetProyectosHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AdminActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// ⭐️ ARREGLADO: Usamos el servicio de auth
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Usamos el servicio de proyecto
	proyectos, err := h.proyectoSvc.GetAllProyectos()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"proyectos": proyectos})
}

func (h *ProyectoHandler) AdminCreateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// ⭐️ ARREGLADO: Permiso
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Servicio
	proyecto, err := h.proyectoSvc.CreateProyecto(req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		// El servicio nos devuelve el error (ej. "ya existe" o "campos requeridos")
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, proyecto)
}

func (h *ProyectoHandler) AdminUpdateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// ⭐️ ARREGLADO: Permiso
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Servicio
	proyecto, err := h.proyectoSvc.UpdateProyecto(req.ID, req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, proyecto)
}

func (h *ProyectoHandler) AdminDeleteProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// ⭐️ ARREGLADO: Permiso
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Servicio
	affected, err := h.proyectoSvc.DeleteProyecto(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Proyecto no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Proyecto borrado."})
}

func (h *ProyectoHandler) AdminSetProyectoEstadoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SetProyectoEstadoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// ⭐️ ARREGLADO: Permiso
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Servicio
	affected, err := h.proyectoSvc.SetProyectoEstado(req.ID, req.Estado)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Proyecto no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: fmt.Sprintf("Estado actualizado a '%s'.", strings.ToLower(req.Estado))})
}