package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto/internal/auth"
	"proyecto/internal/equipos" // ⭐️ NUEVO
	"proyecto/internal/models"
)

// --- 1. EL STRUCT DEL HANDLER ---
type EquipoHandler struct {
	authSvc   auth.AuthService
	equipoSvc equipos.EquipoService
}

// --- 2. EL CONSTRUCTOR DEL HANDLER ---
func NewEquipoHandler(as auth.AuthService, es equipos.EquipoService) *EquipoHandler {
	return &EquipoHandler{
		authSvc:   as,
		equipoSvc: es,
	}
}

// --- 3. LOS MÉTODOS (Handlers) ---

func (h *EquipoHandler) GetEquiposHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GetEquiposRequest
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
	equipos, err := h.equipoSvc.GetEquiposByProyectoID(req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"equipos": equipos})
}

func (h *EquipoHandler) CreateEquipoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEquipoRequest
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
	nuevoEquipo, err := h.equipoSvc.CreateEquipo(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, nuevoEquipo)
}

func (h *EquipoHandler) UpdateEquipoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateEquipoRequest
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
	affected, err := h.equipoSvc.UpdateEquipo(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Equipo no encontrado.")
		return
	}

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Equipo actualizado."})
}

func (h *EquipoHandler) DeleteEquipoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteEquipoRequest
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
	affected, err := h.equipoSvc.DeleteEquipo(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Equipo no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Equipo borrado."})
}