package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time" 

	"github.com/golang-jwt/jwt/v5" 
	"golang.org/x/crypto/bcrypt"

	"proyecto/internal/auth"
	"proyecto/internal/database"
	"proyecto/internal/models"
)


var jwtKey = []byte("mi_llave_secreta_super_segura_12345")


type Claims struct {
	UserID int    `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}


type UserDetails struct {
	Username string `json:"username"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

type LoginResponse struct {
	Token  string      `json:"token"`
	User   UserDetails `json:"user"`
	Role   string      `json:"role"`
	UserId int         `json:"userId"`
}

// --- UTILIDADES ---
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, models.SimpleResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// --- HANDLERS ---

func SaludoHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Hola desde Go!"})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}
	if len(user.Password) < 6 || user.Nombre == "" || user.Apellido == "" {
		respondWithError(w, http.StatusBadRequest, "Todos los campos obligatorios.")
		return
	}
	err := database.CreateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "usuario ya existe") {
			respondWithError(w, http.StatusConflict, "El nombre de usuario ya existe.")
		} else {
			log.Printf("Error al crear usuario: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Error interno al registrar usuario.")
		}
		return
	}
	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: "Usuario registrado exitosamente."})
}

// ⭐️ NUEVO: LoginHandler reemplazado por completo
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userReq models.User
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// 1. Busca al usuario en la DB
	// Asumimos que userDB tiene: ID, Username, HashedPassword, Role, Nombre, Apellido
	userDB, err := database.GetUserByUsername(userReq.Username)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Credenciales inválidas")
		return
	}

	// 2. Compara la contraseña
	err = bcrypt.CompareHashAndPassword([]byte(userDB.HashedPassword), []byte(userReq.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Credenciales inválidas")
		return
	}

	// 3. ¡Éxito! Genera el token JWT
	expirationTime := time.Now().Add(24 * time.Hour) // Token válido por 24 horas
	claims := &Claims{
		UserID: userDB.ID,
		Role:   userDB.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error al firmar token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error interno al generar token")
		return
	}

	// 4. Construye la respuesta JSON exacta que React espera
	response := LoginResponse{
		Token: tokenString,
		User: UserDetails{
			Username: userDB.Username,
			Nombre:   userDB.Nombre,
			Apellido: userDB.Apellido,
		},
		Role:   userDB.Role,
		UserId: userDB.ID,
	}

	// 5. Envía la respuesta
	respondWithJSON(w, http.StatusOK, response)
}

// --- Handlers Admin/Gerente (Usuarios) ---

func AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AdminActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	isAllowed, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	users, err := database.GetAllUsersWithProjects()
	if err != nil {
		log.Printf("Error en GetAllUsersWithProjects: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error al obtener usuarios.")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"users": users})
}

func AdminAddUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}
	isAllowed, err := auth.CheckPermission(req.AdminUsername, "admin")
	if err != nil || !isAllowed {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	if len(req.User.Username) == 0 || len(req.User.Password) < 6 || req.User.Nombre == "" || req.User.Apellido == "" {
		respondWithError(w, http.StatusBadRequest, "Campos requeridos.")
		return
	}
	err = database.CreateUser(req.User)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "usuario ya existe") {
			respondWithError(w, http.StatusConflict, "El nombre de usuario ya existe.")
		} else {
			log.Printf("Error al crear usuario (admin): %v", err)
			respondWithError(w, http.StatusInternalServerError, "Error interno al crear usuario.")
		}
		return
	}
	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: fmt.Sprintf("Usuario '%s' agregado.", req.User.Username)})
}

func AdminDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}
	isAllowed, err := auth.CheckPermission(req.AdminUsername, "admin")
	if err != nil || !isAllowed {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	userDB, err := database.GetUserByUsername(req.AdminUsername)
	if err != nil || userDB.ID == req.ID {
		respondWithError(w, http.StatusForbidden, "No puedes borrar tu propia cuenta.")
		return
	}
	affected, err := database.DeleteUserByID(req.ID)
	if err != nil {
		log.Printf("Error al borrar usuario ID %d: %v", req.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al borrar usuario.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Usuario no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Usuario borrado exitosamente."})
}

func AdminUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}
	isAllowed, err := auth.CheckPermission(req.AdminUsername, "admin")
	if err != nil || !isAllowed {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	newRole := strings.ToLower(req.NewRole)
	if newRole != "admin" && newRole != "user" && newRole != "gerente" {
		respondWithError(w, http.StatusBadRequest, "Rol inválido.")
		return
	}
	userDB, err := database.GetUserByUsername(req.AdminUsername)
	if err != nil || userDB.ID == req.ID {
		respondWithError(w, http.StatusForbidden, "No puedes cambiar tu propio rol.")
		return
	}
	affected, err := database.UpdateUserRole(req.ID, newRole)
	if err != nil {
		log.Printf("Error al actualizar rol del usuario ID %d: %v", req.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al actualizar rol.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Usuario no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: fmt.Sprintf("Rol actualizado a %s.", newRole)})
}

func AdminAssignProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AssignProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}
	isAllowed, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	affected, err := database.AssignProjectToUser(req.UserID, req.ProyectoID)
	if err != nil {
		log.Printf("Error al asignar proyecto %d al usuario %d: %v", req.ProyectoID, req.UserID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al asignar proyecto.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Usuario no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Proyecto asignado/actualizado."})
}

// --- Handlers Admin/Gerente (Proyectos) ---

func AdminGetProyectosHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AdminActionRequest
	json.NewDecoder(r.Body).Decode(&req)
	isAllowed, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	proyectos, err := database.GetAllProjects()
	if err != nil {
		log.Printf("Error en GetAllProjects: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error al obtener proyectos.")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"proyectos": proyectos})
}

func AdminCreateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}
	isAllowed, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	if req.Nombre == "" || req.FechaInicio == "" {
		respondWithError(w, http.StatusBadRequest, "Nombre y fecha de inicio requeridos.")
		return
	}
	err = database.CreateProject(req)
	if err != nil {
		if err.Error() == "el nombre del proyecto ya existe" {
			respondWithError(w, http.StatusConflict, err.Error())
		} else {
			log.Printf("Error al crear proyecto: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Error al crear proyecto.")
		}
		return
	}
	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: fmt.Sprintf("Proyecto '%s' creado.", req.Nombre)})
}

func AdminDeleteProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}
	isAllowed, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	affected, err := database.DeleteProjectByID(req.ID)
	if err != nil {
		log.Printf("Error al borrar proyecto ID %d: %v", req.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al borrar proyecto.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Proyecto no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Proyecto borrado."})
}

func AdminUpdateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}
	isAllowed, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	if req.ID == 0 || req.Nombre == "" || req.FechaInicio == "" {
		respondWithError(w, http.StatusBadRequest, "ID, nombre y fecha de inicio requeridos.")
		return
	}
	proyecto, err := database.GetProjectByID(int64(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Proyecto no encontrado.")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error al verificar estado.")
		return
	}
	if proyecto.Estado == "cerrado" {
		respondWithError(w, http.StatusForbidden, "No se puede modificar un proyecto cerrado.")
		return
	}
	err = database.UpdateProject(req)
	if err != nil {
		if err.Error() == "ese nombre de proyecto ya existe" {
			respondWithError(w, http.StatusConflict, err.Error())
		} else {
			log.Printf("Error al actualizar proyecto ID %d: %v", req.ID, err)
			respondWithError(w, http.StatusInternalServerError, "Error al actualizar proyecto.")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: fmt.Sprintf("Proyecto '%s' actualizado.", req.Nombre)})
}

func AdminSetProyectoEstadoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SetProyectoEstadoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}
	isAllowed, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	newState := strings.ToLower(req.Estado)
	if newState != "habilitado" && newState != "cerrado" {
		respondWithError(w, http.StatusBadRequest, "Estado inválido.")
		return
	}
	if req.ID == 0 {
		respondWithError(w, http.StatusBadRequest, "ID de proyecto requerido.")
		return
	}
	affected, err := database.SetProjectState(req.ID, newState)
	if err != nil {
		log.Printf("Error al cambiar estado del proyecto ID %d: %v", req.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al actualizar estado.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Proyecto no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: fmt.Sprintf("Estado actualizado a '%s'.", newState)})
}

// --- Handler Usuario ---

func UserProjectDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UserProjectDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}
	if req.UserID == 0 {
		respondWithError(w, http.StatusBadRequest, "ID de usuario requerido.")
		return
	}
	details, err := database.GetProjectDetailsForUser(req.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Usuario no encontrado.")
		} else {
			log.Printf("Error en GetProjectDetailsForUser ID %d: %v", req.UserID, err)
			respondWithError(w, http.StatusInternalServerError, "Error al obtener detalles del proyecto.")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, details)
}
