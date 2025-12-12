package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto/internal/auth"
	"proyecto/internal/equipos"
	"proyecto/internal/logger"
	"proyecto/internal/models"
)

// 1. EL STRUCT DEL HANDLER
type EquipoHandler struct {
	authSvc   auth.AuthService
	equipoSvc equipos.EquipoService
	loggerSvc logger.LoggerService
}

// 2. EL CONSTRUCTOR DEL HANDLER
func NewEquipoHandler(as auth.AuthService, es equipos.EquipoService, logs logger.LoggerService) *EquipoHandler {
	return &EquipoHandler{
		authSvc:   as,
		equipoSvc: es,
		loggerSvc: logs, // ⭐️ 3. INYECTAMOS EL SERVICIO
	}
}

// 3. LOS MÉTODOS (Handlers)

func (h *EquipoHandler) GetEquiposHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GetEquiposRequest
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

	equipos, err := h.equipoSvc.GetEquiposByProyectoID(req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string][]models.EquipoImplemento{"equipos": equipos})
}

func (h *EquipoHandler) CreateEquipoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEquipoRequest
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

	nuevoEquipo, err := h.equipoSvc.CreateEquipo(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if nuevoEquipo != nil {
		h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "CREACIÓN", "Equipos/Implementos", nuevoEquipo.ID)
	}

	respondWithJSON(w, http.StatusCreated, nuevoEquipo)
}

func (h *EquipoHandler) UpdateEquipoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateEquipoRequest
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

	affected, err := h.equipoSvc.UpdateEquipo(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "equipo no encontrado")
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "MODIFICACIÓN", "Equipos/Implementos", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Equipo actualizado."})
}

func (h *EquipoHandler) DeleteEquipoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteEquipoRequest
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

	affected, err := h.equipoSvc.DeleteEquipo(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "equipo no encontrado")
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "ELIMINACIÓN", "Equipos/Implementos", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Equipo borrado."})
}
