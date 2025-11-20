package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"proyecto/internal/auth"
	"proyecto/internal/logger"
	"proyecto/internal/models"
	"proyecto/internal/users"
)

type UserHandler struct {
	authSvc   auth.AuthService
	userSvc   users.UserService
	loggerSvc logger.LoggerService
}

func NewUserHandler(as auth.AuthService, us users.UserService, ls logger.LoggerService) *UserHandler {
	return &UserHandler{
		authSvc:   as,
		userSvc:   us,
		loggerSvc: ls,
	}
}

// --- AdminUsersHandler: Listar usuarios ---
func (h *UserHandler) AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusForbidden, "acceso denegado")
		return
	}

	usersList, err := h.userSvc.GetAllUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"users": usersList})
}

// --- AdminAddUserHandler: Crear usuario (desde admin) ---
func (h *UserHandler) AdminAddUserHandler(w http.ResponseWriter, r *http.Request) {
	// Definimos una estructura auxiliar para recibir el JSON complejo { user: {...}, admin_username: "..." }
	type AdminAddUserRequest struct {
		User          models.User `json:"user"`
		AdminUsername string      `json:"admin_username"`
	}

	var req AdminAddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !hasPermission {
		respondWithError(w, http.StatusForbidden, "No autorizado")
		return
	}

	lastID, err := h.userSvc.AddUser(req.User)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", "CREACIÓN", "Usuarios", int(lastID))

	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: "Usuario creado exitosamente"})
}

// --- AdminDeleteUserHandler: Borrar usuario ---
func (h *UserHandler) AdminDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin")
	if err != nil || !hasPermission {
		respondWithError(w, http.StatusForbidden, "No autorizado (Solo Admin)")
		return
	}

	_, err = h.userSvc.DeleteUser(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin", "ELIMINACIÓN", "Usuarios", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Usuario eliminado"})
}

// --- AdminUpdateUserRoleHandler: Actualizar rol de usuario ---
func (h *UserHandler) AdminUpdateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin")
	if err != nil || !hasPermission {
		respondWithError(w, http.StatusForbidden, "No autorizado (Solo Admin)")
		return
	}

	_, err = h.userSvc.UpdateUserRole(req.ID, req.NewRole)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin", "CAMBIO ROL", "Usuarios", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Rol actualizado"})
}

// --- AdminAssignProjectToUserHandler: Asignar usuario a proyecto ---
func (h *UserHandler) AdminAssignProjectToUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AssignProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !hasPermission {
		respondWithError(w, http.StatusForbidden, "No autorizado")
		return
	}

	_, err = h.userSvc.AssignProjectToUser(req.UserID, req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logMsg := fmt.Sprintf("ASIGNACIÓN PROYECTO %d", req.ProyectoID)
	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", logMsg, "Usuarios", req.UserID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Asignación actualizada"})
}

// --- UserProjectDetailsHandler: Dashboard de usuario ---
func (h *UserHandler) UserProjectDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UserProjectDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	response, err := h.userSvc.GetProjectDetailsForUser(req.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}