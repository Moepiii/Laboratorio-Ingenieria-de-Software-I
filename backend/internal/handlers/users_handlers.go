package handlers

import (
	"encoding/json"
	"fmt" // ⭐️ 1. IMPORTAMOS FMT (para formatear strings)
	"net/http"

	"proyecto/internal/auth"
	"proyecto/internal/logger" // ⭐️ 1. IMPORTAMOS EL LOGGER
	"proyecto/internal/models"
	"proyecto/internal/users"
)

// --- 1. EL STRUCT DEL HANDLER ---\
type UserHandler struct {
	authSvc   auth.AuthService
	userSvc   users.UserService
	loggerSvc logger.LoggerService // ⭐️ 2. AÑADIMOS EL SERVICIO DE LOGGER
}

// --- 2. EL CONSTRUCTOR DEL HANDLER ---\
func NewUserHandler(as auth.AuthService, us users.UserService, ls logger.LoggerService) *UserHandler {
	return &UserHandler{
		authSvc:   as,
		userSvc:   us,
		loggerSvc: ls, // ⭐️ 3. INYECTAMOS EL SERVICIO
	}
}

// --- 3. LOS MÉTODOS (Handlers) ---

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
		respondWithError(w, http.StatusForbidden, "acceso denegado. se requiere rol de admin o gerente")
		return
	}

	users, err := h.userSvc.GetAllUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string][]models.UserListResponse{"users": users})
}

func (h *UserHandler) AdminAddUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AddUserRequest
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

	lastID, err := h.userSvc.AddUser(req.User)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// ⭐️ 4. REGISTRAMOS EL EVENTO
	// El rol por defecto ("encargado") se asigna en el servicio
	h.loggerSvc.Log(req.AdminUsername, "admin", "CREACIÓN (encargado)", "Usuarios", int(lastID))

	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: fmt.Sprintf("Usuario '%s' (ID: %d) agregado con éxito.", req.User.Username, lastID)})
}

func (h *UserHandler) AdminDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteUserRequest
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

	affected, err := h.userSvc.DeleteUser(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "usuario no encontrado")
		return
	}

	// ⭐️ 4. REGISTRAMOS EL EVENTO
	h.loggerSvc.Log(req.AdminUsername, "admin", "ELIMINACIÓN", "Usuarios", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Usuario borrado."})
}

func (h *UserHandler) AdminUpdateRoleHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateRoleRequest
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

	affected, err := h.userSvc.UpdateUserRole(req.ID, req.NewRole)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "usuario no encontrado")
		return
	}

	// ⭐️ 4. REGISTRAMOS EL EVENTO
	logMsg := fmt.Sprintf("CAMBIO DE ROL (a %s)", req.NewRole)
	h.loggerSvc.Log(req.AdminUsername, "admin", logMsg, "Usuarios", req.ID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Rol de usuario actualizado."})
}

func (h *UserHandler) AdminAssignProjectToUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AssignProjectRequest
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

	affected, err := h.userSvc.AssignProjectToUser(req.UserID, req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "usuario no encontrado o asignación sin cambios")
		return
	}

	// ⭐️ 4. REGISTRAMOS EL EVENTO
	logMsg := fmt.Sprintf("ASIGNACIÓN DE PROYECTO (ProyectoID: %d)", req.ProyectoID)
	if req.ProyectoID == 0 {
		logMsg = "DESASIGNACIÓN DE PROYECTO"
	}
	h.loggerSvc.Log(req.AdminUsername, "admin/gerente", logMsg, "Usuarios", req.UserID)

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Proyecto asignado/desasignado."})
}

// --- Handler Usuario ---

func (h *UserHandler) UserProjectDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UserProjectDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	// (Este handler es solo de lectura, no necesita logging de evento)
	details, err := h.userSvc.GetProjectDetailsForUser(req.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, details)
}
