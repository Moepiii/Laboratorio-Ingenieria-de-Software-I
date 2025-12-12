package handlers

import (
	"encoding/json"
	"net/http"
	"proyecto/internal/auth"
	"proyecto/internal/logger"
	"proyecto/internal/models"
	"proyecto/internal/unidades"
)

type UnidadHandler struct {
	authSvc   auth.AuthService
	unidadSvc unidades.UnidadService
	loggerSvc logger.LoggerService
}

func NewUnidadHandler(as auth.AuthService, us unidades.UnidadService, ls logger.LoggerService) *UnidadHandler {
	return &UnidadHandler{authSvc: as, unidadSvc: us, loggerSvc: ls}
}

func (h *UnidadHandler) GetUnidadesHandler(w http.ResponseWriter, r *http.Request) {

	var req models.GetUnidadesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Validar permisos
	perm, _ := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if !perm {
		respondWithError(w, http.StatusForbidden, "No autorizado")
		return
	}

	unidades, err := h.unidadSvc.GetUnidadesByProyectoID(req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, unidades)
}

// CreateUnidadHandler
func (h *UnidadHandler) CreateUnidadHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUnidadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	perm, _ := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if !perm {
		respondWithError(w, http.StatusForbidden, "No autorizado")
		return
	}
	nueva, err := h.unidadSvc.CreateUnidad(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.loggerSvc.Log(req.AdminUsername, "admin", "CREACIÓN", "Unidades Medida", nueva.ID)
	respondWithJSON(w, http.StatusCreated, nueva)
}

// Update y Delete Handler
func (h *UnidadHandler) UpdateUnidadHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateUnidadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	perm, _ := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if !perm {
		respondWithError(w, http.StatusForbidden, "No autorizado")
		return
	}
	_, err := h.unidadSvc.UpdateUnidad(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.loggerSvc.Log(req.AdminUsername, "admin", "MODIFICACIÓN", "Unidades Medida", req.ID)
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Actualizado"})
}
func (h *UnidadHandler) DeleteUnidadHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteUnidadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	perm, _ := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if !perm {
		respondWithError(w, http.StatusForbidden, "No autorizado")
		return
	}
	_, err := h.unidadSvc.DeleteUnidad(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.loggerSvc.Log(req.AdminUsername, "admin", "ELIMINACIÓN", "Unidades Medida", req.ID)
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Eliminado"})
}
