package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"proyecto/internal/auth"
	"proyecto/internal/logger" // ⭐️ 1. IMPORTAMOS EL LOGGER
	"proyecto/internal/models"
	"proyecto/internal/proyectos"
)

// --- 1. EL STRUCT DEL HANDLER ---\
type ProyectoHandler struct {
	authSvc     auth.AuthService
	proyectoSvc proyectos.ProyectoService
	loggerSvc   logger.LoggerService // ⭐️ 2. AÑADIMOS EL SERVICIO DE LOGGER
}

// --- 2. EL CONSTRUCTOR DEL HANDLER ---
func NewProyectoHandler(as auth.AuthService, ps proyectos.ProyectoService, ls logger.LoggerService) *ProyectoHandler {
	return &ProyectoHandler{
		authSvc:     as,
		proyectoSvc: ps,
		loggerSvc:   ls, // ⭐️ 3. INYECTAMOS EL SERVICIO
	}
}

// --- 3. LOS MÉTODOS (Handlers) ---

func (h *ProyectoHandler) AdminGetProyectosHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AdminActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error verificando permisos")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "acceso denegado. se requiere rol de admin o gerente")
		return
	}

	proyectos, err := h.proyectoSvc.GetAllProyectos()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string][]models.Proyecto{"proyectos": proyectos})
}

func (h *ProyectoHandler) AdminCreateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error verificando permisos")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "acceso denegado. se requiere rol de admin o gerente")
		return
	}

	proyecto, err := h.proyectoSvc.CreateProyecto(req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		if strings.Contains(err.Error(), "ya existe") {
			respondWithError(w, http.StatusConflict, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	// ⭐️ 4. REGISTRAMOS EL EVENTO
	if proyecto != nil {
		h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "CREACIÓN", "Proyectos", proyecto.ID)
	}

	respondWithJSON(w, http.StatusCreated, proyecto)
}

func (h *ProyectoHandler) AdminUpdateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error verificando permisos")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "acceso denegado. se requiere rol de admin o gerente")
		return
	}

	proyecto, err := h.proyectoSvc.UpdateProyecto(req.ID, req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		if strings.Contains(err.Error(), "ya existe") {
			respondWithError(w, http.StatusConflict, err.Error())
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	// ⭐️ 4. REGISTRAMOS EL EVENTO
	if proyecto != nil {
		h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "MODIFICACIÓN", "Proyectos", proyecto.ID)
	}

	respondWithJSON(w, http.StatusOK, proyecto)
}

func (h *ProyectoHandler) AdminDeleteProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error verificando permisos")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "acceso denegado. se requiere rol de admin")
		return
	}

	affected, err := h.proyectoSvc.DeleteProyecto(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "proyecto no encontrado")
		return
	}

	// ⭐️ 4. REGISTRAMOS EL EVENTO
	h.loggerSvc.Log(req.AdminUsername, "admin", "ELIMINACIÓN", "Proyectos", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Proyecto borrado."})
}

func (h *ProyectoHandler) AdminSetProyectoEstadoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SetProyectoEstadoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error verificando permisos")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "acceso denegado. se requiere rol de admin o gerente")
		return
	}

	affected, err := h.proyectoSvc.SetProyectoEstado(req.ID, req.Estado)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "proyecto no encontrado")
		return
	}

	// ⭐️ 4. REGISTRAMOS EL EVENTO
	logMsg := fmt.Sprintf("CAMBIO DE ESTADO (a %s)", req.Estado)
	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", logMsg, "Proyectos", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Estado del proyecto actualizado."})
}
