package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto/internal/actividades"
	"proyecto/internal/auth"
	"proyecto/internal/logger"
	"proyecto/internal/models"
)

// 1. EL STRUCT DEL HANDLER
type ActividadHandler struct {
	authSvc      auth.AuthService
	actividadSvc actividades.ActividadService
	loggerSvc    logger.LoggerService
}

// 2. EL CONSTRUCTOR DEL HANDLER
func NewActividadHandler(as auth.AuthService, acs actividades.ActividadService, ls logger.LoggerService) *ActividadHandler {
	return &ActividadHandler{
		authSvc:      as,
		actividadSvc: acs,
		loggerSvc:    ls,
	}
}

//  3. LOS MÉTODOS (Handlers)

func (h *ActividadHandler) GetDatosProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GetDatosProyectoRequest
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
		respondWithError(w, http.StatusForbidden, "acceso denegado")
		return
	}

	datos, err := h.actividadSvc.GetDatosProyecto(req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, datos)
}

func (h *ActividadHandler) CreateActividadHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateActividadRequest
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
		respondWithError(w, http.StatusForbidden, "acceso denegado")
		return
	}

	actividades, err := h.actividadSvc.CreateActividad(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "CREACIÓN (Actividad)", "Proyectos", req.ProyectoID)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"actividades": actividades})
}

func (h *ActividadHandler) UpdateActividadHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateActividadRequest
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
		respondWithError(w, http.StatusForbidden, "acceso denegado")
		return
	}

	actividades, err := h.actividadSvc.UpdateActividad(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "MODIFICACIÓN", "Actividades", req.ID)

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"actividades": actividades})
}

func (h *ActividadHandler) DeleteActividadHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteActividadRequest
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
		respondWithError(w, http.StatusForbidden, "acceso denegado")
		return
	}

	affected, err := h.actividadSvc.DeleteActividad(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "actividad no encontrada")
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "ELIMINACIÓN", "Actividades", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Actividad borrada."})
}
