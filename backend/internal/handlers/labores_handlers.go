package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto/internal/auth"
	"proyecto/internal/labores"
	"proyecto/internal/logger"
	"proyecto/internal/models"
)

// 1. EL STRUCT DEL HANDLER
type LaborHandler struct {
	authSvc   auth.AuthService
	laborSvc  labores.LaborService
	loggerSvc logger.LoggerService
}

// 2. EL CONSTRUCTOR DEL HANDLER
func NewLaborHandler(as auth.AuthService, ls labores.LaborService, logs logger.LoggerService) *LaborHandler {
	return &LaborHandler{
		authSvc:   as,
		laborSvc:  ls,
		loggerSvc: logs,
	}
}

//  3. LOS MÉTODOS (Handlers)

func (h *LaborHandler) GetLaboresHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GetLaboresRequest
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

	labores, err := h.laborSvc.GetLaboresByProyectoID(req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string][]models.LaborAgronomica{"labores": labores})
}

func (h *LaborHandler) CreateLaborHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLaborRequest
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

	nuevaLabor, err := h.laborSvc.CreateLabor(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if nuevaLabor != nil {
		h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "CREACIÓN", "Labores", nuevaLabor.ID)
	}

	respondWithJSON(w, http.StatusCreated, nuevaLabor)
}

func (h *LaborHandler) UpdateLaborHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateLaborRequest
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

	affected, err := h.laborSvc.UpdateLabor(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "labor no encontrada")
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "MODIFICACIÓN", "Labores", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Labor actualizada."})
}

func (h *LaborHandler) DeleteLaborHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteLaborRequest
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

	affected, err := h.laborSvc.DeleteLabor(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "labor no encontrada")
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "ELIMINACIÓN", "Labores", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Labor borrada."})
}
