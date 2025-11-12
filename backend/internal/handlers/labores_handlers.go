package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto/internal/auth"
	"proyecto/internal/labores" // ⭐️ NUEVO
	"proyecto/internal/models"
)

// --- 1. EL STRUCT DEL HANDLER ---
type LaborHandler struct {
	authSvc  auth.AuthService
	laborSvc labores.LaborService
}

// --- 2. EL CONSTRUCTOR DEL HANDLER ---
func NewLaborHandler(as auth.AuthService, ls labores.LaborService) *LaborHandler {
	return &LaborHandler{
		authSvc:  as,
		laborSvc: ls,
	}
}

// --- 3. LOS MÉTODOS (Handlers) ---

func (h *LaborHandler) GetLaboresHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GetLaboresRequest
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
	labores, err := h.laborSvc.GetLaboresByProyectoID(req.ProyectoID)
	if err != nil {
		// El servicio ya validó el ID y manejó errores de DB
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"labores": labores})
}

func (h *LaborHandler) CreateLaborHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLaborRequest
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
	// Pasamos el request completo, el servicio se encarga de la lógica
	nuevaLabor, err := h.laborSvc.CreateLabor(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, nuevaLabor)
}

func (h *LaborHandler) UpdateLaborHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateLaborRequest
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
	affected, err := h.laborSvc.UpdateLabor(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Labor no encontrada.")
		return
	}

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Labor actualizada."})
}

func (h *LaborHandler) DeleteLaborHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteLaborRequest
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
	affected, err := h.laborSvc.DeleteLabor(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Labor no encontrada.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Labor borrada."})
}