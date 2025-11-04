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

// --- HANDLERS BÁSICOS ---

func SaludoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "¡Hola desde el backend de Go!")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	if user.Username == "" || user.Password == "" || user.Nombre == "" || user.Apellido == "" {
		respondWithError(w, http.StatusBadRequest, "Todos los campos (username, password, nombre, apellido) son requeridos.")
		return
	}

	lastID, err := database.RegisterUser(user.Username, user.Password, user.Nombre, user.Apellido)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: fmt.Sprintf("Usuario '%s' (ID: %d) registrado con éxito.", user.Username, lastID)})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	user, err := database.GetUserByUsername(creds.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "Usuario no encontrado.")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error interno del servidor.")
		}
		return
	}

	if !database.CheckPasswordHash(creds.Password, user.HashedPassword) {
		respondWithError(w, http.StatusUnauthorized, "Contraseña incorrecta.")
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al generar el token.")
		return
	}

	dbUser, err := database.GetUserByUsername(creds.Username)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al obtener detalles del usuario.")
		return
	}

	userDetails := UserDetails{
		Username: dbUser.Username,
		Nombre:   dbUser.Nombre,
		Apellido: dbUser.Apellido,
	}

	response := LoginResponse{
		Token:  tokenString,
		User:   userDetails,
		Role:   dbUser.Role,
		UserId: dbUser.ID,
	}

	respondWithJSON(w, http.StatusOK, response)
}

// --- Handlers Admin/Gerente (Usuarios) ---

func AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AdminActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}

	users, err := database.GetAllUsersWithProjectNames()
	if err != nil {
		log.Printf("Error en AdminUsersHandler: %v", err)
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

	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin.")
		return
	}

	if req.User.Username == "" || req.User.Password == "" || req.User.Nombre == "" || req.User.Apellido == "" {
		respondWithError(w, http.StatusBadRequest, "Username, password, nombre y apellido son requeridos.")
		return
	}

	hashedPassword, _ := database.HashPassword(req.User.Password)
	userID, err := database.AddUser(req.User, hashedPassword)
	if err != nil {
		log.Printf("Error en AdminAddUserHandler: %v", err)
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			respondWithError(w, http.StatusBadRequest, "El nombre de usuario ya existe.")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error al añadir usuario.")
		}
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]interface{}{"mensaje": "Usuario añadido", "id": userID})
}

func AdminDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin.")
		return
	}
	if req.ID == 0 {
		respondWithError(w, http.StatusBadRequest, "ID de usuario requerido.")
		return
	}
	affected, err := database.DeleteUser(req.ID)
	if err != nil {
		log.Printf("Error en AdminDeleteUserHandler (ID: %d): %v", req.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al borrar usuario.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Usuario no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Usuario borrado."})
}

func AdminUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin.")
		return
	}

	if req.ID == 0 || req.NewRole == "" {
		respondWithError(w, http.StatusBadRequest, "ID y NewRole son requeridos.")
		return
	}
	if req.NewRole != "admin" && req.NewRole != "gerente" && req.NewRole != "user" {
		respondWithError(w, http.StatusBadRequest, "Rol debe ser 'admin', 'gerente' o 'user'.")
		return
	}

	affected, err := database.UpdateUserRole(req.ID, req.NewRole)
	if err != nil {
		log.Printf("Error en AdminUpdateUserHandler (ID: %d): %v", req.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al actualizar rol.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Usuario no encontrado.")
		return
	}

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Rol actualizado."})
}

func AdminAssignProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AssignProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}
	if req.UserID == 0 {
		respondWithError(w, http.StatusBadRequest, "ID de usuario (user_id) requerido.")
		return
	}

	affected, err := database.AssignProjectToUser(req.UserID, req.ProyectoID)
	if err != nil {
		log.Printf("Error en AdminAssignProyectoHandler (User: %d, Proy: %d): %v", req.UserID, req.ProyectoID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al asignar proyecto.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Usuario no encontrado.")
		return
	}

	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Proyecto asignado/desasignado."})
}

// --- Handlers Admin/Gerente (Proyectos) ---

func AdminGetProyectosHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AdminActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}

	proyectos, err := database.GetAllProyectos()
	if err != nil {
		log.Printf("Error en AdminGetProyectosHandler: %v", err)
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

	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}
	if req.Nombre == "" || req.FechaInicio == "" || req.FechaCierre == "" {
		respondWithError(w, http.StatusBadRequest, "Nombre, fecha_inicio y fecha_cierre son requeridos.")
		return
	}

	proyectoID, err := database.CreateProyecto(req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		log.Printf("Error en AdminCreateProyectoHandler: %v", err)
		if strings.Contains(err.Error(), "ya existe") {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error al crear proyecto.")
		}
		return
	}

	proyecto, err := database.GetProjectByID(proyectoID)
	if err != nil {
		log.Printf("Error al recuperar proyecto recién creado (ID: %d): %v", proyectoID, err)
		respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: "Proyecto creado pero no se pudo recuperar."})
		return
	}
	respondWithJSON(w, http.StatusCreated, proyecto)
}

func AdminUpdateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}

	if req.ID == 0 || req.Nombre == "" || req.FechaInicio == "" || req.FechaCierre == "" {
		respondWithError(w, http.StatusBadRequest, "ID, Nombre, fecha_inicio y fecha_cierre son requeridos.")
		return
	}
	affected, err := database.UpdateProyecto(req.ID, req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		log.Printf("Error en AdminUpdateProyectoHandler (ID: %d): %v", req.ID, err)
		if strings.Contains(err.Error(), "ya existe") {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "Error al actualizar proyecto.")
		}
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Proyecto no encontrado.")
		return
	}

	proyecto, err := database.GetProjectByID(int64(req.ID))
	if err != nil {
		respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Proyecto actualizado pero no se pudo recuperar."})
		return
	}
	respondWithJSON(w, http.StatusOK, proyecto)
}

func AdminDeleteProyectoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}
	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}
	if req.ID == 0 {
		respondWithError(w, http.StatusBadRequest, "ID de proyecto requerido.")
		return
	}

	affected, err := database.DeleteProyecto(req.ID)
	if err != nil {
		log.Printf("Error en AdminDeleteProyectoHandler (ID: %d): %v", req.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al borrar proyecto.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Proyecto no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Proyecto borrado."})
}

func AdminSetProyectoEstadoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SetProyectoEstadoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}
	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
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
			log.Printf("Error en UserProjectDetailsHandler (User: %d): %v", req.UserID, err)
			respondWithError(w, http.StatusInternalServerError, "Error al obtener detalles del proyecto.")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, details)
}

// ⭐️ --- (INICIO) Handlers Labores Agronómicas --- ⭐️

// ⭐️ NUEVO: Handler para OBTENER todas las labores de un proyecto
func GetLaboresHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar
	var req models.GetLaboresRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// 2. Seguridad (Solo Admin o Gerente)
	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado. Se requiere rol de admin o gerente.")
		return
	}

	// 3. Lógica de DB
	if req.ProyectoID == 0 {
		respondWithError(w, http.StatusBadRequest, "ID de proyecto requerido.")
		return
	}

	labores, err := database.GetLaboresByProyectoID(req.ProyectoID)
	if err != nil {
		log.Printf("Error al obtener labores para proyecto %d: %v", req.ProyectoID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al obtener labores.")
		return
	}

	// 4. Responder
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"labores": labores})
}

// ⭐️ NUEVO: Handler para CREAR una nueva labor
func CreateLaborHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar
	var req models.CreateLaborRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// 2. Seguridad
	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}

	// 3. Validar
	if req.ProyectoID == 0 || req.Descripcion == "" {
		respondWithError(w, http.StatusBadRequest, "ProyectoID y Descripción son requeridos.")
		return
	}

	estado := req.Estado
	if estado == "" {
		estado = "activa" // Estado por defecto
	}

	// 4. Lógica de DB
	labor := models.LaborAgronomica{
		ProyectoID:  req.ProyectoID,
		Descripcion: req.Descripcion,
		Estado:      estado,
	}

	laborID, err := database.CreateLabor(labor)
	if err != nil {
		log.Printf("Error al crear labor: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error al crear labor.")
		return
	}

	// 5. Responder
	// Devolvemos la labor completa recién creada
	nuevaLabor, err := database.GetLaborByID(int(laborID))
	if err != nil {
		log.Printf("Error al obtener labor recién creada (ID: %d): %v", laborID, err)
		respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: "Labor creada con éxito, pero no se pudo recuperar."})
		return
	}

	respondWithJSON(w, http.StatusCreated, nuevaLabor)
}

// ⭐️ NUEVO: Handler para ACTUALIZAR una labor
func UpdateLaborHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar
	var req models.UpdateLaborRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// 2. Seguridad
	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}

	// 3. Validar
	if req.ID == 0 || req.Descripcion == "" || req.Estado == "" {
		respondWithError(w, http.StatusBadRequest, "ID, Descripción y Estado son requeridos.")
		return
	}

	// 4. Lógica de DB
	affected, err := database.UpdateLabor(req.ID, req.Descripcion, req.Estado)
	if err != nil {
		log.Printf("Error al actualizar labor ID %d: %v", req.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al actualizar la labor.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Labor no encontrada.")
		return
	}

	// 5. Responder
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Labor actualizada."})
}

// ⭐️ NUEVO: Handler para BORRAR una labor
func DeleteLaborHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar
	var req models.DeleteLaborRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato JSON inválido.")
		return
	}

	// 2. Seguridad
	hasPermission, err := auth.CheckPermission(req.AdminUsername, "admin", "gerente")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error verificando permisos.")
		return
	}
	if !hasPermission {
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}

	// 3. Validar
	if req.ID == 0 {
		respondWithError(w, http.StatusBadRequest, "ID de labor requerido.")
		return
	}

	// 4. Lógica de DB
	affected, err := database.DeleteLabor(req.ID)
	if err != nil {
		log.Printf("Error al borrar labor ID %d: %v", req.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al borrar la labor.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Labor no encontrada.")
		return
	}

	// 5. Responder
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Labor borrada."})
}

// ⭐️ --- (FIN) Handlers Labores Agronómicas --- ⭐️
