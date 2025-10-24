package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

// --- ESTRUCTURAS DE DATOS ---
// (User, UserDB, UserListResponse, etc. sin cambios)
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}
type UserDB struct {
	ID             int
	Username       string
	HashedPassword string
	Role           string
	Nombre         string
	Apellido       string
	ProyectoID     sql.NullInt64
}
type UserListResponse struct {
	ID             int     `json:"id"`
	Username       string  `json:"username"`
	Role           string  `json:"role"`
	Nombre         string  `json:"nombre"`
	Apellido       string  `json:"apellido"`
	ProyectoID     *int    `json:"proyecto_id"`
	ProyectoNombre *string `json:"proyecto_nombre"`
}
type UpdateRoleRequest struct {
	ID            int    `json:"id"`
	NewRole       string `json:"new_role"`
	AdminUsername string `json:"admin_username"`
}
type AdminActionRequest struct {
	AdminUsername string `json:"admin_username"`
}
type AddUserRequest struct {
	User
	AdminUsername string `json:"admin_username"`
}
type DeleteUserRequest struct {
	ID            int    `json:"id"`
	AdminUsername string `json:"admin_username"`
}
type AssignProyectoRequest struct {
	UserID        int    `json:"user_id"`
	ProyectoID    int    `json:"proyecto_id"`
	AdminUsername string `json:"admin_username"`
}

// Estructuras de Proyectos (ACTUALIZADAS)
type Proyecto struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	FechaInicio string `json:"fecha_inicio"`
	FechaCierre string `json:"fecha_cierre"`
	Estado      string `json:"estado"` // NUEVO
}
type CreateProyectoRequest struct { // (Sin cambios, el estado es default)
	Nombre        string `json:"nombre"`
	FechaInicio   string `json:"fecha_inicio"`
	FechaCierre   string `json:"fecha_cierre"`
	AdminUsername string `json:"admin_username"`
}
type UpdateProyectoRequest struct { // (Sin cambios, no se actualiza estado aquí)
	ID            int    `json:"id"`
	Nombre        string `json:"nombre"`
	FechaInicio   string `json:"fecha_inicio"`
	FechaCierre   string `json:"fecha_cierre"`
	AdminUsername string `json:"admin_username"`
}
type DeleteProyectoRequest struct { // (Sin cambios)
	ID            int    `json:"id"`
	AdminUsername string `json:"admin_username"`
}

// NUEVA STRUCT para cambiar estado
type SetProyectoEstadoRequest struct {
	ID            int    `json:"id"`
	Estado        string `json:"estado"` // 'habilitado' o 'cerrado'
	AdminUsername string `json:"admin_username"`
}

// (UserProjectDetailsRequest, ProjectMember, UserProjectDetailsResponse sin cambios)
type UserProjectDetailsRequest struct {
	UserID int `json:"user_id"`
}
type ProjectMember struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Role     string `json:"role"`
}
type UserProjectDetailsResponse struct {
	Proyecto *Proyecto       `json:"proyecto"`
	Miembros []ProjectMember `json:"miembros"`
	Gerentes []ProjectMember `json:"gerentes"`
}

// --- UTILIDADES DE BASE DE DATOS ---

func initDB() {
	var err error
	DB, err = sql.Open("sqlite", "../Base de datos/users.db")
	if err != nil {
		log.Fatalf("Error al abrir DB: %v", err)
	}
	_, err = DB.Exec("PRAGMA foreign_keys = OFF;")
	if err != nil {
		log.Fatalf("Error PRAGMA OFF: %v", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL UNIQUE, password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user', nombre TEXT NOT NULL DEFAULT '', apellido TEXT NOT NULL DEFAULT ''
	);`
	if _, err = DB.Exec(createTableSQL); err != nil {
		log.Fatalf("Error crear 'users': %v", err)
	}

	// Crear 'proyectos' - AHORA CON ESTADO
	createProyectosTableSQL := `
	CREATE TABLE IF NOT EXISTS proyectos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL UNIQUE,
		fecha_inicio TEXT NOT NULL DEFAULT '',
		fecha_cierre TEXT NOT NULL DEFAULT '',
		estado TEXT NOT NULL DEFAULT 'habilitado' -- NUEVO
	);`
	if _, err = DB.Exec(createProyectosTableSQL); err != nil {
		log.Fatalf("Error crear 'proyectos': %v", err)
	}

	// Migraciones (migrateProyectosTable ahora incluye 'estado')
	migrateProyectosTable()
	migrateUsersTable()

	// Crear admin (sin cambios)
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		log.Fatalf("Error check admin: %v", err)
	}
	if count == 0 {
		adminPassword := "admin123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		insertAdminSQL := "INSERT INTO users(username, password, role, nombre, apellido) VALUES(?, ?, ?, ?, ?)"
		_, err = DB.Exec(insertAdminSQL, "admin", string(hashedPassword), "admin", "Admin", "User")
		if err != nil {
			log.Fatalf("Error insert admin: %v", err)
		}
		log.Println("✅ Usuario admin creado.")
	}

	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalf("Error PRAGMA ON: %v", err)
	}
	log.Println("✅ DB inicializada.")
}

// migrateProyectosTable (ACTUALIZADO para incluir 'estado')
func migrateProyectosTable() {
	rows, err := DB.Query("PRAGMA table_info(proyectos)")
	if err != nil {
		log.Fatalf("Error PRAGMA 'proyectos': %v", err)
	}
	defer rows.Close()
	columnExists := map[string]bool{"nombre": false, "fecha_inicio": false, "fecha_cierre": false, "estado": false}
	for rows.Next() {
		var (
			cid        int
			name       string
			ctype      string
			notnull    int
			dflt_value sql.NullString
			pk         int
		)
		err = rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk)
		if err != nil {
			log.Fatalf("Error scan 'proyectos': %v", err)
		}
		if _, ok := columnExists[name]; ok {
			columnExists[name] = true
		}
	}
	rows.Close()

	if !columnExists["nombre"] { // Si falta 'nombre', recreamos todo
		log.Println("⚠️ Recreando 'proyectos' por falta de 'nombre'...")
		_, err = DB.Exec("DROP TABLE IF EXISTS proyectos;")
		if err != nil {
			log.Fatalf("Error DROP 'proyectos': %v", err)
		}
		createProyectosTableSQL := `
		CREATE TABLE proyectos (
			id INTEGER PRIMARY KEY AUTOINCREMENT, nombre TEXT NOT NULL UNIQUE,
			fecha_inicio TEXT NOT NULL DEFAULT '', fecha_cierre TEXT NOT NULL DEFAULT '',
			estado TEXT NOT NULL DEFAULT 'habilitado'
		);`
		if _, err = DB.Exec(createProyectosTableSQL); err != nil {
			log.Fatalf("Error RE-CREAR 'proyectos': %v", err)
		}
		log.Println("✅ Tabla 'proyectos' recreada.")
		return
	}
	// Migraciones individuales si 'nombre' existe
	if !columnExists["fecha_inicio"] {
		log.Println("⚠️ Migrando 'fecha_inicio' en 'proyectos'...")
		_, err = DB.Exec("ALTER TABLE proyectos ADD COLUMN fecha_inicio TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error migrar 'fecha_inicio': %v", err)
		}
	}
	if !columnExists["fecha_cierre"] {
		log.Println("⚠️ Migrando 'fecha_cierre' en 'proyectos'...")
		_, err = DB.Exec("ALTER TABLE proyectos ADD COLUMN fecha_cierre TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error migrar 'fecha_cierre': %v", err)
		}
	}
	if !columnExists["estado"] {
		log.Println("⚠️ Migrando 'estado' en 'proyectos'...")
		_, err = DB.Exec("ALTER TABLE proyectos ADD COLUMN estado TEXT NOT NULL DEFAULT 'habilitado'")
		if err != nil {
			log.Fatalf("Error migrar 'estado': %v", err)
		}
	}
}

// migrateUsersTable (sin cambios)
func migrateUsersTable() {
	rows, err := DB.Query("PRAGMA table_info(users)")
	if err != nil {
		log.Fatalf("Error PRAGMA 'users': %v", err)
	}
	defer rows.Close()
	columnExists := map[string]bool{"role": false, "nombre": false, "apellido": false, "proyecto_id": false}
	for rows.Next() {
		var (
			cid        int
			name       string
			ctype      string
			notnull    int
			dflt_value sql.NullString
			pk         int
		)
		err = rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk)
		if err != nil {
			log.Fatalf("Error scan 'users': %v", err)
		}
		if _, ok := columnExists[name]; ok {
			columnExists[name] = true
		}
	}
	rows.Close()
	if !columnExists["role"] {
		log.Println("⚠️ Migrando 'role' en 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user'")
		if err != nil {
			log.Fatalf("Error migrar 'role': %v", err)
		}
	}
	if !columnExists["nombre"] {
		log.Println("⚠️ Migrando 'nombre' en 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN nombre TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error migrar 'nombre': %v", err)
		}
	}
	if !columnExists["apellido"] {
		log.Println("⚠️ Migrando 'apellido' en 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN apellido TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error migrar 'apellido': %v", err)
		}
	}
	if !columnExists["proyecto_id"] {
		log.Println("⚠️ Migrando 'proyecto_id' en 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN proyecto_id INTEGER REFERENCES proyectos(id) ON DELETE SET NULL")
		if err != nil {
			log.Fatalf("Error migrar 'proyecto_id': %v", err)
		}
	}
}

// checkPermission (sin cambios)
func checkPermission(username string, requiredRoles ...string) (bool, error) {
	var userRole string
	err := DB.QueryRow("SELECT role FROM users WHERE username = ?", username).Scan(&userRole)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("DB error check role: %w", err)
	}
	userRoleLower := strings.ToLower(userRole)
	for _, reqRole := range requiredRoles {
		if userRoleLower == strings.ToLower(reqRole) {
			return true, nil
		}
	}
	return false, nil
}

// --- HANDLERS AUTENTICACIÓN (Sin cambios) ---
func registerHandler(w http.ResponseWriter, r *http.Request) { /* ... */
	w.Header().Set("Content-Type", "application/json")
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "JSON inválido."}`, http.StatusBadRequest)
		return
	}
	if len(user.Password) < 6 || user.Nombre == "" || user.Apellido == "" {
		http.Error(w, `{"error": "Campos requeridos."}`, http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Error interno."}`, http.StatusInternalServerError)
		return
	}
	query := "INSERT INTO users(username, password, role, nombre, apellido) VALUES(?, ?, ?, ?, ?)"
	_, err = DB.Exec(query, user.Username, string(hashedPassword), "user", user.Nombre, user.Apellido)
	if err != nil {
		http.Error(w, `{"error": "Usuario ya existe."}`, http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"mensaje": "Usuario registrado."}`)
}
func loginHandler(w http.ResponseWriter, r *http.Request) { /* ... */
	w.Header().Set("Content-Type", "application/json")
	var userReq User
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, `{"error": "JSON inválido."}`, http.StatusBadRequest)
		return
	}
	var userDB UserDB
	query := "SELECT id, username, password, role FROM users WHERE username = ?"
	row := DB.QueryRow(query, userReq.Username)
	err := row.Scan(&userDB.ID, &userDB.Username, &userDB.HashedPassword, &userDB.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Credenciales inválidas"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"error": "Error interno."}`, http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDB.HashedPassword), []byte(userReq.Password))
	if err != nil {
		http.Error(w, `{"error": "Credenciales inválidas"}`, http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, `{"mensaje": "Login exitoso", "usuario": "%s", "id": %d, "role": "%s"}`, userDB.Username, userDB.ID, userDB.Role)
}

// --- HANDLERS ADMIN (USUARIOS) (Sin cambios) ---
func adminUsersHandler(w http.ResponseWriter, r *http.Request) { /* ... */
	w.Header().Set("Content-Type", "application/json")
	var adminReq AdminActionRequest
	json.NewDecoder(r.Body).Decode(&adminReq)
	isAllowed, err := checkPermission(adminReq.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}
	query := `SELECT u.id, u.username, u.role, u.nombre, u.apellido, p.id, p.nombre FROM users u LEFT JOIN proyectos p ON u.proyecto_id = p.id ORDER BY u.id ASC`
	rows, err := DB.Query(query)
	if err != nil {
		http.Error(w, `{"error": "Error interno."}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	users := []UserListResponse{}
	for rows.Next() {
		var u UserListResponse
		var pID sql.NullInt64
		var pNombre sql.NullString
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.Nombre, &u.Apellido, &pID, &pNombre); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		if pID.Valid {
			id := int(pID.Int64)
			u.ProyectoID = &id
		}
		if pNombre.Valid {
			u.ProyectoNombre = &pNombre.String
		}
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
}
func adminAddUserHandler(w http.ResponseWriter, r *http.Request) { /* ... */
	w.Header().Set("Content-Type", "application/json")
	var req AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "JSON inválido."}`, http.StatusBadRequest)
		return
	}
	isAllowed, err := checkPermission(req.AdminUsername, "admin")
	if err != nil || !isAllowed {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}
	if len(req.User.Username) == 0 || len(req.User.Password) < 6 || req.User.Nombre == "" || req.User.Apellido == "" {
		http.Error(w, `{"error": "Campos requeridos."}`, http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Error interno."}`, http.StatusInternalServerError)
		return
	}
	query := "INSERT INTO users(username, password, role, nombre, apellido) VALUES(?, ?, ?, ?, ?)"
	_, err = DB.Exec(query, req.User.Username, string(hashedPassword), "user", req.User.Nombre, req.User.Apellido)
	if err != nil {
		http.Error(w, `{"error": "Usuario ya existe."}`, http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"mensaje": "Usuario '%s' agregado."}`, req.User.Username)
}
func adminDeleteUserHandler(w http.ResponseWriter, r *http.Request) { /* ... */
	w.Header().Set("Content-Type", "application/json")
	var req DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "JSON inválido."}`, http.StatusBadRequest)
		return
	}
	isAllowed, err := checkPermission(req.AdminUsername, "admin")
	if err != nil || !isAllowed {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}
	var targetUsername string
	err = DB.QueryRow("SELECT username FROM users WHERE id = ?", req.ID).Scan(&targetUsername)
	if err != nil {
		http.Error(w, `{"error": "Usuario no encontrado."}`, http.StatusNotFound)
		return
	}
	if targetUsername == req.AdminUsername {
		http.Error(w, `{"error": "No puedes borrarte a ti mismo."}`, http.StatusForbidden)
		return
	}
	query := "DELETE FROM users WHERE id = ?"
	_, err = DB.Exec(query, req.ID)
	if err != nil {
		http.Error(w, `{"error": "Error interno."}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"mensaje": fmt.Sprintf("Usuario '%s' borrado.", targetUsername)})
}
func adminUpdateUserHandler(w http.ResponseWriter, r *http.Request) { /* ... */
	w.Header().Set("Content-Type", "application/json")
	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "JSON inválido."}`, http.StatusBadRequest)
		return
	}
	isAllowed, err := checkPermission(req.AdminUsername, "admin")
	if err != nil || !isAllowed {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}
	newRole := strings.ToLower(req.NewRole)
	if newRole != "admin" && newRole != "user" && newRole != "gerente" {
		http.Error(w, `{"error": "Rol inválido."}`, http.StatusBadRequest)
		return
	}
	var targetUsername string
	err = DB.QueryRow("SELECT username FROM users WHERE id = ?", req.ID).Scan(&targetUsername)
	if err != nil {
		http.Error(w, `{"error": "Usuario no encontrado."}`, http.StatusNotFound)
		return
	}
	if targetUsername == req.AdminUsername {
		http.Error(w, `{"error": "No puedes cambiar tu propio rol."}`, http.StatusForbidden)
		return
	}
	query := "UPDATE users SET role = ? WHERE id = ?"
	_, err = DB.Exec(query, newRole, req.ID)
	if err != nil {
		http.Error(w, `{"error": "Error interno."}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"mensaje": fmt.Sprintf("Rol de %s actualizado a %s.", targetUsername, newRole)})
}
func adminAssignProyectoHandler(w http.ResponseWriter, r *http.Request) { /* ... */
	w.Header().Set("Content-Type", "application/json")
	var req AssignProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "JSON inválido."}`, http.StatusBadRequest)
		return
	}
	isAllowed, err := checkPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}
	var idToSet interface{}
	if req.ProyectoID == 0 {
		idToSet = nil
	} else {
		idToSet = req.ProyectoID
	}
	query := "UPDATE users SET proyecto_id = ? WHERE id = ?"
	_, err = DB.Exec(query, idToSet, req.UserID)
	if err != nil {
		http.Error(w, `{"error": "Error interno."}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Proyecto asignado."})
}

// --- HANDLERS ADMIN (PROYECTOS) ---

// adminGetProyectosHandler (ACTUALIZADO para devolver estado)
func adminGetProyectosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var adminReq AdminActionRequest
	json.NewDecoder(r.Body).Decode(&adminReq)
	isAllowed, err := checkPermission(adminReq.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}

	// Selecciona el estado
	rows, err := DB.Query("SELECT id, nombre, fecha_inicio, fecha_cierre, estado FROM proyectos ORDER BY nombre ASC")
	if err != nil {
		http.Error(w, `{"error": "Error interno."}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	proyectos := []Proyecto{}
	for rows.Next() {
		var p Proyecto
		// Escanea el estado
		if err := rows.Scan(&p.ID, &p.Nombre, &p.FechaInicio, &p.FechaCierre, &p.Estado); err != nil {
			log.Printf("Scan error proyecto: %v", err)
			continue
		}
		proyectos = append(proyectos, p)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"proyectos": proyectos})
}

// adminCreateProyectoHandler (ACTUALIZADO para insertar estado default)
func adminCreateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "JSON inválido."}`, http.StatusBadRequest)
		return
	}
	isAllowed, err := checkPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}
	if req.Nombre == "" || req.FechaInicio == "" {
		http.Error(w, `{"error": "Campos requeridos."}`, http.StatusBadRequest)
		return
	}

	// El estado se inserta con el DEFAULT 'habilitado'
	query := "INSERT INTO proyectos(nombre, fecha_inicio, fecha_cierre) VALUES(?, ?, ?)"
	_, err = DB.Exec(query, req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		http.Error(w, `{"error": "Proyecto ya existe."}`, http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": fmt.Sprintf("Proyecto '%s' creado.", req.Nombre)})
}

// adminDeleteProyectoHandler (Sin cambios)
func adminDeleteProyectoHandler(w http.ResponseWriter, r *http.Request) { /* ... */
	w.Header().Set("Content-Type", "application/json")
	var req DeleteProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "JSON inválido."}`, http.StatusBadRequest)
		return
	}
	isAllowed, err := checkPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}
	query := "DELETE FROM proyectos WHERE id = ?"
	result, err := DB.Exec(query, req.ID)
	if err != nil {
		http.Error(w, `{"error": "Error interno."}`, http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error": "Proyecto no existe."}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Proyecto borrado."})
}

// adminUpdateProyectoHandler (ACTUALIZADO para prevenir edición si está cerrado)
func adminUpdateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req UpdateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "JSON inválido."}`, http.StatusBadRequest)
		return
	}
	isAllowed, err := checkPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}
	if req.ID == 0 || req.Nombre == "" || req.FechaInicio == "" {
		http.Error(w, `{"error": "Campos requeridos."}`, http.StatusBadRequest)
		return
	}

	// ⭐️ VERIFICACIÓN DE ESTADO ANTES DE ACTUALIZAR
	var currentState string
	err = DB.QueryRow("SELECT estado FROM proyectos WHERE id = ?", req.ID).Scan(&currentState)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Proyecto no encontrado."}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Error al verificar estado."}`, http.StatusInternalServerError)
		return
	}
	if strings.ToLower(currentState) == "cerrado" {
		http.Error(w, `{"error": "No se puede modificar un proyecto cerrado."}`, http.StatusForbidden)
		return
	}

	query := "UPDATE proyectos SET nombre = ?, fecha_inicio = ?, fecha_cierre = ? WHERE id = ?"
	_, err = DB.Exec(query, req.Nombre, req.FechaInicio, req.FechaCierre, req.ID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, `{"error": "Nombre de proyecto ya existe."}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error": "Error interno."}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"mensaje": fmt.Sprintf("Proyecto '%s' actualizado.", req.Nombre)})
}

// ⭐️ --- NUEVO HANDLER PARA CAMBIAR ESTADO DEL PROYECTO --- ⭐️
func adminSetProyectoEstadoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req SetProyectoEstadoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "JSON inválido."}`, http.StatusBadRequest)
		return
	}

	// Permiso: Admin y Gerente
	isAllowed, err := checkPermission(req.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}

	newState := strings.ToLower(req.Estado)
	if newState != "habilitado" && newState != "cerrado" {
		http.Error(w, `{"error": "Estado inválido. Debe ser 'habilitado' o 'cerrado'."}`, http.StatusBadRequest)
		return
	}
	if req.ID == 0 {
		http.Error(w, `{"error": "ID de proyecto requerido."}`, http.StatusBadRequest)
		return
	}

	query := "UPDATE proyectos SET estado = ? WHERE id = ?"
	result, err := DB.Exec(query, newState, req.ID)
	if err != nil {
		http.Error(w, `{"error": "Error interno al actualizar estado."}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error": "Proyecto no encontrado."}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"mensaje": fmt.Sprintf("Estado del proyecto actualizado a '%s'.", newState)})
}

// --- HANDLER USUARIO (Sin cambios) ---
func userProjectDetailsHandler(w http.ResponseWriter, r *http.Request) { /* ... */
	w.Header().Set("Content-Type", "application/json")
	var req UserProjectDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "JSON inválido."}`, http.StatusBadRequest)
		return
	}
	if req.UserID == 0 {
		http.Error(w, `{"error": "ID Usuario requerido."}`, http.StatusBadRequest)
		return
	}
	var proyectoID sql.NullInt64
	err := DB.QueryRow("SELECT proyecto_id FROM users WHERE id = ?", req.UserID).Scan(&proyectoID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Usuario no encontrado."}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Error buscar usuario."}`, http.StatusInternalServerError)
		return
	}
	response := UserProjectDetailsResponse{Proyecto: nil, Miembros: []ProjectMember{}, Gerentes: []ProjectMember{}}
	if proyectoID.Valid {
		pID := proyectoID.Int64
		var p Proyecto
		err = DB.QueryRow("SELECT id, nombre, fecha_inicio, fecha_cierre, estado FROM proyectos WHERE id = ?", pID).Scan(&p.ID, &p.Nombre, &p.FechaInicio, &p.FechaCierre, &p.Estado)
		if err != nil {
			log.Printf("Error: Usuario %d tiene proyecto_id %d inválido: %v", req.UserID, pID, err)
		} else {
			response.Proyecto = &p
		}
		rows, err := DB.Query(`SELECT id, username, nombre, apellido, role FROM users WHERE proyecto_id = ? AND id != ? ORDER BY role, nombre`, pID, req.UserID)
		if err != nil {
			log.Printf("Error buscar miembros proyecto %d: %v", pID, err)
			http.Error(w, `{"error": "Error buscar miembros."}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var member ProjectMember
			if err := rows.Scan(&member.ID, &member.Username, &member.Nombre, &member.Apellido, &member.Role); err != nil {
				log.Printf("Error scan miembro: %v", err)
				continue
			}
			if strings.ToLower(member.Role) == "gerente" {
				response.Gerentes = append(response.Gerentes, member)
			} else {
				response.Miembros = append(response.Miembros, member)
			}
		}
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	initDB()
	defer DB.Close()
	// Handlers
	saludoAPI := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { /* ... */ })
	registerAPI := http.HandlerFunc(registerHandler)
	loginAPI := http.HandlerFunc(loginHandler)
	adminUsersAPI := http.HandlerFunc(adminUsersHandler)
	adminUpdateUserAPI := http.HandlerFunc(adminUpdateUserHandler)
	adminAddUserAPI := http.HandlerFunc(adminAddUserHandler)
	adminDeleteUserAPI := http.HandlerFunc(adminDeleteUserHandler)
	adminGetProyectosAPI := http.HandlerFunc(adminGetProyectosHandler)
	adminCreateProyectoAPI := http.HandlerFunc(adminCreateProyectoHandler)
	adminDeleteProyectoAPI := http.HandlerFunc(adminDeleteProyectoHandler)
	adminAssignProyectoAPI := http.HandlerFunc(adminAssignProyectoHandler)
	adminUpdateProyectoAPI := http.HandlerFunc(adminUpdateProyectoHandler)
	userProjectDetailsAPI := http.HandlerFunc(userProjectDetailsHandler)
	// NUEVO HANDLER DE ESTADO
	adminSetProyectoEstadoAPI := http.HandlerFunc(adminSetProyectoEstadoHandler)

	// CORS (sin cambios)
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	corsOptions := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)

	// Rutas
	http.Handle("/api/saludo", corsOptions(saludoAPI))
	http.Handle("/api/register", corsOptions(registerAPI))
	http.Handle("/api/login", corsOptions(loginAPI))
	http.Handle("/api/admin/users", corsOptions(adminUsersAPI))
	http.Handle("/api/admin/update-user", corsOptions(adminUpdateUserAPI))
	http.Handle("/api/admin/add-user", corsOptions(adminAddUserAPI))
	http.Handle("/api/admin/delete-user", corsOptions(adminDeleteUserAPI))
	http.Handle("/api/admin/get-proyectos", corsOptions(adminGetProyectosAPI))
	http.Handle("/api/admin/create-proyecto", corsOptions(adminCreateProyectoAPI))
	http.Handle("/api/admin/delete-proyecto", corsOptions(adminDeleteProyectoAPI))
	http.Handle("/api/admin/assign-proyecto", corsOptions(adminAssignProyectoAPI))
	http.Handle("/api/admin/update-proyecto", corsOptions(adminUpdateProyectoAPI))
	http.Handle("/api/user/project-details", corsOptions(userProjectDetailsAPI))
	// NUEVA RUTA DE ESTADO
	http.Handle("/api/admin/set-proyecto-estado", corsOptions(adminSetProyectoEstadoAPI))

	log.Println("🚀 Servidor Go v8 (Estado Proyecto) en :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error server: %v", err)
	}
}
