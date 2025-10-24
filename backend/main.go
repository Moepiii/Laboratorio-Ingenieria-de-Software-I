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
// (User, UserDB, UserListResponse, UpdateRoleRequest, etc. sin cambios)
// ... (Omitidas por brevedad) ...
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
	AdminUsername string `json:"admin_username"` // Quien hace la petición
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
type Proyecto struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	FechaInicio string `json:"fecha_inicio"`
	FechaCierre string `json:"fecha_cierre"`
}
type CreateProyectoRequest struct {
	Nombre        string `json:"nombre"`
	FechaInicio   string `json:"fecha_inicio"`
	FechaCierre   string `json:"fecha_cierre"`
	AdminUsername string `json:"admin_username"`
}
type UpdateProyectoRequest struct {
	ID            int    `json:"id"`
	Nombre        string `json:"nombre"`
	FechaInicio   string `json:"fecha_inicio"`
	FechaCierre   string `json:"fecha_cierre"`
	AdminUsername string `json:"admin_username"`
}
type DeleteProyectoRequest struct {
	ID            int    `json:"id"`
	AdminUsername string `json:"admin_username"`
}
type AssignProyectoRequest struct {
	UserID        int    `json:"user_id"`
	ProyectoID    int    `json:"proyecto_id"`
	AdminUsername string `json:"admin_username"`
}

// --- NUEVAS STRUCTS PARA LA VISTA DE USUARIO ---
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
	Proyecto *Proyecto       `json:"proyecto"` // Puntero para que pueda ser null
	Miembros []ProjectMember `json:"miembros"`
	Gerentes []ProjectMember `json:"gerentes"` // Gerentes asignados al mismo proyecto
}

// --- UTILIDADES DE BASE DE DATOS ---
// (initDB, migraciones, checkPermission sin cambios)
// ... (Omitidas por brevedad) ...
func initDB() {
	var err error
	DB, err = sql.Open("sqlite", "../Base de datos/users.db")
	if err != nil {
		log.Fatalf("Error al abrir la base de datos: %v", err)
	}
	_, err = DB.Exec("PRAGMA foreign_keys = OFF;")
	if err != nil {
		log.Fatalf("Error al desactivar foreign keys: %v", err)
	}
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL UNIQUE, password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user', nombre TEXT NOT NULL DEFAULT '', apellido TEXT NOT NULL DEFAULT ''
	);`
	if _, err = DB.Exec(createTableSQL); err != nil {
		log.Fatalf("Error al crear la tabla 'users': %v", err)
	}
	createProyectosTableSQL := `
	CREATE TABLE IF NOT EXISTS proyectos (
		id INTEGER PRIMARY KEY AUTOINCREMENT, nombre TEXT NOT NULL UNIQUE
	);`
	if _, err = DB.Exec(createProyectosTableSQL); err != nil {
		log.Fatalf("Error al crear la tabla 'proyectos': %v", err)
	}
	migrateProyectosTable()
	migrateUsersTable()
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		log.Fatalf("Error al verificar usuarios existentes: %v", err)
	}
	if count == 0 {
		adminPassword := "admin123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		insertAdminSQL := "INSERT INTO users(username, password, role, nombre, apellido) VALUES(?, ?, ?, ?, ?)"
		_, err = DB.Exec(insertAdminSQL, "admin", string(hashedPassword), "admin", "Admin", "User")
		if err != nil {
			log.Fatalf("Error al insertar usuario admin por defecto: %v", err)
		}
		log.Println("✅ Usuario administrador 'admin' (pass: admin123) creado por defecto.")
	}
	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalf("Error al reactivar foreign keys: %v", err)
	}
	log.Println("✅ Base de datos SQLite inicializada y tablas listas y migradas.")
}
func migrateProyectosTable() { /* ... sin cambios ... */
	rows, err := DB.Query("PRAGMA table_info(proyectos)")
	if err != nil {
		log.Fatalf("Error PRAGMA 'proyectos': %v", err)
	}
	defer rows.Close()
	columnExists := map[string]bool{"nombre": false, "fecha_inicio": false, "fecha_cierre": false}
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
	if !columnExists["nombre"] {
		log.Println("⚠️ Tabla 'proyectos' está corrupta (falta 'nombre'). Recreando...")
		_, err = DB.Exec("DROP TABLE IF EXISTS proyectos;")
		if err != nil {
			log.Fatalf("Error 'DROP' 'proyectos': %v", err)
		}
		createProyectosTableSQL := `
		CREATE TABLE IF NOT EXISTS proyectos (
			id INTEGER PRIMARY KEY AUTOINCREMENT, nombre TEXT NOT NULL UNIQUE,
			fecha_inicio TEXT NOT NULL DEFAULT '', fecha_cierre TEXT NOT NULL DEFAULT ''
		);`
		if _, err = DB.Exec(createProyectosTableSQL); err != nil {
			log.Fatalf("Error RE-CREAR 'proyectos': %v", err)
		}
		log.Println("✅ Tabla 'proyectos' recreada exitosamente.")
		return
	}
	if !columnExists["fecha_inicio"] {
		log.Println("⚠️ Migrando: Añadiendo 'fecha_inicio' a 'proyectos'...")
		_, err = DB.Exec("ALTER TABLE proyectos ADD COLUMN fecha_inicio TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error al migrar 'fecha_inicio': %v", err)
		}
	}
	if !columnExists["fecha_cierre"] {
		log.Println("⚠️ Migrando: Añadiendo 'fecha_cierre' a 'proyectos'...")
		_, err = DB.Exec("ALTER TABLE proyectos ADD COLUMN fecha_cierre TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error al migrar 'fecha_cierre': %v", err)
		}
	}
}
func migrateUsersTable() { /* ... sin cambios ... */
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
		log.Println("⚠️ Migrando: Añadiendo 'role' a 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user'")
		if err != nil {
			log.Fatalf("Error al migrar 'role': %v", err)
		}
	}
	if !columnExists["nombre"] {
		log.Println("⚠️ Migrando: Añadiendo 'nombre' a 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN nombre TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error al migrar 'nombre': %v", err)
		}
	}
	if !columnExists["apellido"] {
		log.Println("⚠️ Migrando: Añadiendo 'apellido' a 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN apellido TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error al migrar 'apellido': %v", err)
		}
	}
	if !columnExists["proyecto_id"] {
		log.Println("⚠️ Migrando: Añadiendo 'proyecto_id' a 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN proyecto_id INTEGER REFERENCES proyectos(id) ON DELETE SET NULL")
		if err != nil {
			log.Fatalf("Error al migrar 'proyecto_id': %v", err)
		}
	}
}
func checkPermission(username string, requiredRoles ...string) (bool, error) { /* ... sin cambios ... */
	var userRole string
	err := DB.QueryRow("SELECT role FROM users WHERE username = ?", username).Scan(&userRole)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error de DB al verificar rol: %w", err)
	}
	userRoleLower := strings.ToLower(userRole)
	for _, reqRole := range requiredRoles {
		if userRoleLower == strings.ToLower(reqRole) {
			return true, nil
		}
	}
	return false, nil
}

// --- HANDLERS DE AUTENTICACIÓN (Sin cambios) ---
// (registerHandler, loginHandler sin cambios)
// ... (Omitidos por brevedad) ...
func registerHandler(w http.ResponseWriter, r *http.Request) { /* ... sin cambios ... */
	w.Header().Set("Content-Type", "application/json")
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}
	if len(user.Password) < 6 || user.Nombre == "" || user.Apellido == "" {
		http.Error(w, `{"error": "Todos los campos obligatorios."}`, http.StatusBadRequest)
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
		http.Error(w, `{"error": "El usuario ya existe."}`, http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"mensaje": "Usuario registrado."}`)
}
func loginHandler(w http.ResponseWriter, r *http.Request) { /* ... sin cambios ... */
	w.Header().Set("Content-Type", "application/json")
	var userReq User
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
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
	fmt.Fprintf(w, `{"mensaje": "Inicio de sesión exitoso", "usuario": "%s", "id": %d, "role": "%s"}`, userDB.Username, userDB.ID, userDB.Role)
}

// --- HANDLERS DE ADMINISTRACIÓN (USUARIOS) ---
// (adminUsersHandler, adminAddUserHandler, adminDeleteUserHandler, adminUpdateUserHandler, adminAssignProyectoHandler sin cambios)
// ... (Omitidos por brevedad) ...
func adminUsersHandler(w http.ResponseWriter, r *http.Request) { /* ... sin cambios ... */
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
func adminAddUserHandler(w http.ResponseWriter, r *http.Request) { /* ... sin cambios ... */
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
func adminDeleteUserHandler(w http.ResponseWriter, r *http.Request) { /* ... sin cambios ... */
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
func adminUpdateUserHandler(w http.ResponseWriter, r *http.Request) { /* ... sin cambios ... */
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
func adminAssignProyectoHandler(w http.ResponseWriter, r *http.Request) { /* ... sin cambios ... */
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

// --- HANDLERS DE ADMINISTRACIÓN (PROYECTOS) ---
// (adminGetProyectosHandler, adminCreateProyectoHandler, adminDeleteProyectoHandler, adminUpdateProyectoHandler sin cambios)
// ... (Omitidos por brevedad) ...
func adminGetProyectosHandler(w http.ResponseWriter, r *http.Request) { /* ... sin cambios ... */
	w.Header().Set("Content-Type", "application/json")
	var adminReq AdminActionRequest
	json.NewDecoder(r.Body).Decode(&adminReq)
	isAllowed, err := checkPermission(adminReq.AdminUsername, "admin", "gerente")
	if err != nil || !isAllowed {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}
	rows, err := DB.Query("SELECT id, nombre, fecha_inicio, fecha_cierre FROM proyectos ORDER BY nombre ASC")
	if err != nil {
		http.Error(w, `{"error": "Error interno."}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	proyectos := []Proyecto{}
	for rows.Next() {
		var p Proyecto
		if err := rows.Scan(&p.ID, &p.Nombre, &p.FechaInicio, &p.FechaCierre); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		proyectos = append(proyectos, p)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"proyectos": proyectos})
}
func adminCreateProyectoHandler(w http.ResponseWriter, r *http.Request) { /* ... sin cambios ... */
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
	query := "INSERT INTO proyectos(nombre, fecha_inicio, fecha_cierre) VALUES(?, ?, ?)"
	_, err = DB.Exec(query, req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		http.Error(w, `{"error": "Proyecto ya existe."}`, http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": fmt.Sprintf("Proyecto '%s' creado.", req.Nombre)})
}
func adminDeleteProyectoHandler(w http.ResponseWriter, r *http.Request) { /* ... sin cambios ... */
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
func adminUpdateProyectoHandler(w http.ResponseWriter, r *http.Request) { /* ... sin cambios ... */
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

// ⭐️ --- NUEVO HANDLER PARA USUARIOS NORMALES --- ⭐️
func userProjectDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req UserProjectDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}

	if req.UserID == 0 {
		http.Error(w, `{"error": "ID de usuario requerido."}`, http.StatusBadRequest)
		return
	}

	// 1. Obtener el proyecto_id del usuario
	var proyectoID sql.NullInt64
	err := DB.QueryRow("SELECT proyecto_id FROM users WHERE id = ?", req.UserID).Scan(&proyectoID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Usuario no encontrado."}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Error al buscar usuario."}`, http.StatusInternalServerError)
		return
	}

	response := UserProjectDetailsResponse{
		Proyecto: nil,
		Miembros: []ProjectMember{},
		Gerentes: []ProjectMember{},
	}

	// 2. Si el usuario tiene un proyecto asignado...
	if proyectoID.Valid {
		pID := proyectoID.Int64

		// 2a. Obtener los detalles del proyecto
		var p Proyecto
		err = DB.QueryRow("SELECT id, nombre, fecha_inicio, fecha_cierre FROM proyectos WHERE id = ?", pID).Scan(&p.ID, &p.Nombre, &p.FechaInicio, &p.FechaCierre)
		if err != nil {
			// Si no encontramos el proyecto, algo está mal (FK debería prevenir esto, pero por si acaso)
			log.Printf("Error: Usuario %d tiene proyecto_id %d inválido: %v", req.UserID, pID, err)
			// Devolvemos la respuesta vacía pero sin error HTTP
		} else {
			response.Proyecto = &p // Asigna el proyecto encontrado
		}

		// 2b. Obtener todos los miembros (usuarios y gerentes) de ese proyecto (excluyendo al propio usuario)
		rows, err := DB.Query(`
            SELECT id, username, nombre, apellido, role
            FROM users
            WHERE proyecto_id = ? AND id != ?
            ORDER BY role, nombre`, pID, req.UserID) // Ordena para poner gerentes primero

		if err != nil {
			log.Printf("Error al buscar miembros del proyecto %d: %v", pID, err)
			http.Error(w, `{"error": "Error al buscar miembros del proyecto."}`, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var member ProjectMember
			if err := rows.Scan(&member.ID, &member.Username, &member.Nombre, &member.Apellido, &member.Role); err != nil {
				log.Printf("Error al escanear miembro de proyecto: %v", err)
				continue
			}
			// Separa en listas diferentes según el rol
			if strings.ToLower(member.Role) == "gerente" {
				response.Gerentes = append(response.Gerentes, member)
			} else { // Asume que cualquier otro rol es un miembro normal
				response.Miembros = append(response.Miembros, member)
			}
		}
	}

	// 3. Devolver la respuesta (puede tener proyecto=nil si no está asignado)
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
	// ⭐️ NUEVO HANDLER
	userProjectDetailsAPI := http.HandlerFunc(userProjectDetailsHandler)

	// Configuración de CORS (sin cambios)
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	corsOptions := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)

	// Rutas
	http.Handle("/api/saludo", corsOptions(saludoAPI))
	http.Handle("/api/register", corsOptions(registerAPI))
	http.Handle("/api/login", corsOptions(loginAPI))
	// Rutas Admin/Gerente
	http.Handle("/api/admin/users", corsOptions(adminUsersAPI))
	http.Handle("/api/admin/update-user", corsOptions(adminUpdateUserAPI))
	http.Handle("/api/admin/add-user", corsOptions(adminAddUserAPI))
	http.Handle("/api/admin/delete-user", corsOptions(adminDeleteUserAPI))
	http.Handle("/api/admin/get-proyectos", corsOptions(adminGetProyectosAPI))
	http.Handle("/api/admin/create-proyecto", corsOptions(adminCreateProyectoAPI))
	http.Handle("/api/admin/delete-proyecto", corsOptions(adminDeleteProyectoAPI))
	http.Handle("/api/admin/assign-proyecto", corsOptions(adminAssignProyectoAPI))
	http.Handle("/api/admin/update-proyecto", corsOptions(adminUpdateProyectoAPI))
	// ⭐️ NUEVA RUTA PARA USUARIOS
	http.Handle("/api/user/project-details", corsOptions(userProjectDetailsAPI))

	log.Println("🚀 Servidor Go escuchando en :8080 (v7 - Vista de Usuario)")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
