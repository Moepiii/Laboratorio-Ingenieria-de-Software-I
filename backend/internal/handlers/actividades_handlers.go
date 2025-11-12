package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto/internal/actividades" // ⭐️ NUEVO
	"proyecto/internal/auth"
	"proyecto/internal/models"
)

// --- 1. EL STRUCT DEL HANDLER ---
type ActividadHandler struct {
	authSvc      auth.AuthService
	actividadSvc actividades.ActividadService
}

// --- 2. EL CONSTRUCTOR DEL HANDLER ---
func NewActividadHandler(as auth.AuthService, acs actividades.ActividadService) *ActividadHandler {
	return &ActividadHandler{
		authSvc:      as,
		actividadSvc: acs,
	}
}

// --- 3. LOS MÉTODOS (Handlers) ---

func (h *ActividadHandler) GetDatosProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GetDatosProyectoRequest
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
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Servicio
	// El servicio se encarga de llamar a las 4 funciones de DB
	response, err := h.actividadSvc.GetDatosProyecto(req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *ActividadHandler) CreateActividadHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateActividadRequest
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
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Servicio
	actividades, err := h.actividadSvc.CreateActividad(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// El servicio devuelve la lista actualizada (o un mensaje de error si falló la recarga)
	if err != nil {
		respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: err.Error()})
	} else {
		respondWithJSON(w, http.StatusCreated, map[string]interface{}{"actividades": actividades})
	}
}

func (h *ActividadHandler) UpdateActividadHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateActividadRequest
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
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Servicio
	actividades, err := h.actividadSvc.UpdateActividad(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err != nil {
		respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: err.Error()})
	} else {
		respondWithJSON(w, http.StatusOK, map[string]interface{}{"actividades": actividades})
	}
}

func (h *ActividadHandler) DeleteActividadHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteActividadRequest
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
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Servicio
	affected, err := h.actividadSvc.DeleteActividad(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Actividad no encontrada.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Actividad borrada."})
}