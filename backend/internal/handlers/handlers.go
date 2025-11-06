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
	Cedula   string `json:"cedula"`
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

	if user.Username == "" || user.Password == "" || user.Nombre == "" || user.Apellido == "" || user.Cedula == "" {
		respondWithError(w, http.StatusBadRequest, "Todos los campos (username, password, nombre, apellido, cedula) son requeridos.")
		return
	}

	lastID, err := database.RegisterUser(user.Username, user.Password, user.Nombre, user.Apellido, user.Cedula)
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
		Cedula:   dbUser.Cedula,
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

	if req.User.Username == "" || req.User.Password == "" || req.User.Nombre == "" || req.User.Apellido == "" || req.User.Cedula == "" {
		respondWithError(w, http.StatusBadRequest, "Username, password, nombre, apellido y cedula son requeridos.")
		return
	}

	hashedPassword, _ := database.HashPassword(req.User.Password)
	userID, err := database.AddUser(req.User, hashedPassword)
	if err != nil {
		log.Printf("Error en AdminAddUserHandler: %v", err)
		if strings.Contains(err.Error(), "ya existe") || strings.Contains(err.Error(), "ya está registrada") {
			respondWithError(w, http.StatusBadRequest, err.Error())
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

	// ⭐️ MODIFICADO: Añadido "encargado" a la validación
	if req.NewRole != "admin" && req.NewRole != "gerente" && req.NewRole != "user" && req.NewRole != "encargado" {
		respondWithError(w, http.StatusBadRequest, "Rol debe ser 'admin', 'gerente', 'encargado' o 'user'.")
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

// --- Handlers Labores Agronómicas ---

func GetLaboresHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GetLaboresRequest
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
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"labores": labores})
}

func CreateLaborHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLaborRequest
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
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	if req.ProyectoID == 0 || req.Descripcion == "" {
		respondWithError(w, http.StatusBadRequest, "ProyectoID y Descripción son requeridos.")
		return
	}
	estado := req.Estado
	if estado == "" {
		estado = "activa"
	}
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
	nuevaLabor, err := database.GetLaborByID(int(laborID))
	if err != nil {
		log.Printf("Error al obtener labor recién creada (ID: %d): %v", laborID, err)
		respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: "Labor creada con éxito, pero no se pudo recuperar."})
		return
	}
	respondWithJSON(w, http.StatusCreated, nuevaLabor)
}

func UpdateLaborHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateLaborRequest
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
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	if req.ID == 0 || req.Descripcion == "" || req.Estado == "" {
		respondWithError(w, http.StatusBadRequest, "ID, Descripción y Estado son requeridos.")
		return
	}
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
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Labor actualizada."})
}

func DeleteLaborHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteLaborRequest
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
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	if req.ID == 0 {
		respondWithError(w, http.StatusBadRequest, "ID de labor requerido.")
		return
	}
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
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Labor borrada."})
}

// --- Handlers Equipos e Implementos ---

func GetEquiposHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GetEquiposRequest
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
	if req.ProyectoID == 0 {
		respondWithError(w, http.StatusBadRequest, "ID de proyecto requerido.")
		return
	}
	equipos, err := database.GetEquiposByProyectoID(req.ProyectoID)
	if err != nil {
		log.Printf("Error al obtener equipos para proyecto %d: %v", req.ProyectoID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al obtener equipos.")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"equipos": equipos})
}

func CreateEquipoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEquipoRequest
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
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	if req.ProyectoID == 0 || req.Nombre == "" {
		respondWithError(w, http.StatusBadRequest, "ProyectoID y Nombre son requeridos.")
		return
	}
	tipo := req.Tipo
	if tipo == "" {
		tipo = "implemento"
	}
	estado := req.Estado
	if estado == "" {
		estado = "disponible"
	}
	equipo := models.EquipoImplemento{
		ProyectoID: req.ProyectoID,
		Nombre:     req.Nombre,
		Tipo:       tipo,
		Estado:     estado,
	}
	equipoID, err := database.CreateEquipo(equipo)
	if err != nil {
		log.Printf("Error al crear equipo: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error al crear equipo.")
		return
	}
	nuevoEquipo, err := database.GetEquipoByID(int(equipoID))
	if err != nil {
		log.Printf("Error al obtener equipo recién creado (ID: %d): %v", equipoID, err)
		respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: "Equipo creado con éxito, pero no se pudo recuperar."})
		return
	}
	respondWithJSON(w, http.StatusCreated, nuevoEquipo)
}

func UpdateEquipoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateEquipoRequest
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
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	if req.ID == 0 || req.Nombre == "" || req.Tipo == "" || req.Estado == "" {
		respondWithError(w, http.StatusBadRequest, "ID, Nombre, Tipo y Estado son requeridos.")
		return
	}
	affected, err := database.UpdateEquipo(req.ID, req.Nombre, req.Tipo, req.Estado)
	if err != nil {
		log.Printf("Error al actualizar equipo ID %d: %v", req.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al actualizar el equipo.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Equipo no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Equipo actualizado."})
}

func DeleteEquipoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteEquipoRequest
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
		respondWithError(w, http.StatusForbidden, "Acceso denegado.")
		return
	}
	if req.ID == 0 {
		respondWithError(w, http.StatusBadRequest, "ID de equipo requerido.")
		return
	}
	affected, err := database.DeleteEquipo(req.ID)
	if err != nil {
		log.Printf("Error al borrar equipo ID %d: %v", req.ID, err)
		respondWithError(w, http.StatusInternalServerError, "Error al borrar el equipo.")
		return
	}
	if affected == 0 {
		respondWithError(w, http.StatusNotFound, "Equipo no encontrado.")
		return
	}
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Equipo borrado."})
}
