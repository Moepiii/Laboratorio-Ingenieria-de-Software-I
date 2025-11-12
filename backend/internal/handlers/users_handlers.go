package handlers

import (
	"encoding/json"
	"net/http"

	"proyecto/internal/auth"
	"proyecto/internal/models"
	"proyecto/internal/users" // ⭐️ NUEVO: Importamos el servicio de usuarios
)

// --- 1. EL STRUCT DEL HANDLER ---
type UserHandler struct {
	authSvc auth.AuthService
	userSvc users.UserService
}

// --- 2. EL CONSTRUCTOR DEL HANDLER ---
func NewUserHandler(as auth.AuthService, us users.UserService) *UserHandler {
	return &UserHandler{
		authSvc: as,
		userSvc: us,
	}
}

// --- 3. LOS MÉTODOS (Handlers) ---

func (h *UserHandler) AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AdminActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// ⭐️ ARREGLADO: Usamos el servicio de auth inyectado
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Usamos el servicio de usuario
	users, err := h.userSvc.GetAllUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"users": users})
}

func (h *UserHandler) AdminAddUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// ⭐️ ARREGLADO: Permiso
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Servicio
	// Pasamos solo la parte de 'User', el servicio hace la validación y el hash
	userID, err := h.userSvc.AddUser(req.User)
	if err != nil {
		// El servicio nos devuelve el error (ej. "ya existe" o "campos requeridos")
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]interface{}{"mensaje": "Usuario añadido", "id": userID})
}

func (h *UserHandler) AdminDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// ⭐️ ARREGLADO: Permiso
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Servicio
	affected, err := h.userSvc.DeleteUser(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Usuario no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Usuario borrado."})
}

func (h *UserHandler) AdminUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// ⭐️ ARREGLADO: Permiso
	hasPermission, err := h.authSvc.CheckPermission(req.AdminUsername, "admin")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin.")
		return
	}

	// ⭐️ LÓGICA MOVIDA: Servicio
	affected, err := h.userSvc.UpdateUserRole(req.ID, req.NewRole)
	if err != nil {
		// El servicio nos devuelve error de validación ("Rol debe ser...") o de DB
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Usuario no encontrado.")
		return
	}

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Rol actualizado."})
}

func (h *UserHandler) AdminAssignProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AssignProyectoRequest
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
	affected, err := h.userSvc.AssignProjectToUser(req.UserID, req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if affected == 0 {
		// Nota: El servicio no puede diferenciar "usuario no encontrado" de "proyecto ya asignado".
		// Para esta lógica, está bien.
		respondWithError(w, http.StatusNotFound, "Usuario no encontrado o asignación sin cambios.")
		return
	}

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Proyecto asignado/desasignado."})
}

// --- Handler Usuario ---

func (h *UserHandler) UserProjectDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UserProjectDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// Nota: Este handler no tenía chequeo de permisos.
	// Se podría añadir:
	// 1. Un chequeo de que el usuario que pide es el req.UserID
	// 2. O que el usuario que pide es admin/gerente.
	// Por ahora, lo dejamos como estaba.

	// ⭐️ LÓGICA MOVIDA: Servicio
	details, err := h.userSvc.GetProjectDetailsForUser(req.UserID)
	if err != nil {
		// El servicio nos devuelve "Usuario no encontrado" u otro error
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, details)
}