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

// DB es la conexión global a la base de datos
var DB *sql.DB

// --- ESTRUCTURAS DE DATOS ---

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

// --- ESTRUCTURAS DE PROYECTOS (ACTUALIZADAS) ---

type Proyecto struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	FechaInicio string `json:"fecha_inicio"` // NUEVO
	FechaCierre string `json:"fecha_cierre"` // NUEVO
}

type CreateProyectoRequest struct {
	Nombre        string `json:"nombre"`
	FechaInicio   string `json:"fecha_inicio"` // NUEVO
	FechaCierre   string `json:"fecha_cierre"` // NUEVO
	AdminUsername string `json:"admin_username"`
}

// NUEVA STRUCT
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
	ProyectoID    int    `json:"proyecto_id"` // 0 significará "ninguno"
	AdminUsername string `json:"admin_username"`
}

// --- UTILIDADES DE BASE DE DATOS ---

func initDB() {
	var err error
	DB, err = sql.Open("sqlite", "../Base de datos/users.db")
	if err != nil {
		log.Fatalf("Error al abrir la base de datos: %v", err)
	}

	// Desactivar foreign keys temporalmente para las migraciones
	_, err = DB.Exec("PRAGMA foreign_keys = OFF;")
	if err != nil {
		log.Fatalf("Error al desactivar foreign keys: %v", err)
	}

	// 1. Crear tabla 'users'
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		nombre TEXT NOT NULL DEFAULT '', 
		apellido TEXT NOT NULL DEFAULT ''
	);`
	if _, err = DB.Exec(createTableSQL); err != nil {
		log.Fatalf("Error al crear la tabla 'users': %v", err)
	}

	// 2. Crear tabla 'proyectos' (versión inicial)
	createProyectosTableSQL := `
	CREATE TABLE IF NOT EXISTS proyectos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL UNIQUE
	);`
	if _, err = DB.Exec(createProyectosTableSQL); err != nil {
		log.Fatalf("Error al crear la tabla 'proyectos': %v", err)
	}

	// 3. Migrar 'proyectos' (Añade 'nombre', 'fecha_inicio', 'fecha_cierre' si faltan)
	migrateProyectosTable()

	// 4. Migrar 'users' (Añade 'role', 'nombre', 'apellido', 'proyecto_id' si faltan)
	migrateUsersTable()

	// 5. Crear admin por defecto
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		log.Fatalf("Error al verificar usuarios existentes: %v", err)
	}

	if count == 0 {
		adminPassword := "admin123"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Error al hashear contraseña de admin: %v", err)
		}
		insertAdminSQL := "INSERT INTO users(username, password, role, nombre, apellido) VALUES(?, ?, ?, ?, ?)"
		_, err = DB.Exec(insertAdminSQL, "admin", string(hashedPassword), "admin", "Admin", "User")
		if err != nil {
			log.Fatalf("Error al insertar usuario admin por defecto: %v", err)
		}
		log.Println("✅ Usuario administrador 'admin' (pass: admin123) creado por defecto.")
	}

	// 6. Reactivar las foreign keys
	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalf("Error al reactivar foreign keys: %v", err)
	}

	log.Println("✅ Base de datos SQLite inicializada y tablas 'users' y 'proyectos' listas y migradas.")
}

// ⭐️ --- FUNCIÓN DE MIGRACIÓN 'proyectos' ACTUALIZADA --- ⭐️
// Añade 'fecha_inicio' y 'fecha_cierre' si no existen.
func migrateProyectosTable() {
	rows, err := DB.Query("PRAGMA table_info(proyectos)")
	if err != nil {
		log.Fatalf("Error al leer la información del esquema (PRAGMA) para 'proyectos': %v", err)
	}
	defer rows.Close()

	columnExists := map[string]bool{
		"nombre":       false,
		"fecha_inicio": false,
		"fecha_cierre": false,
	}

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
			log.Fatalf("Error al escanear información de columna de 'proyectos': %v", err)
		}
		if _, ok := columnExists[name]; ok {
			columnExists[name] = true
		}
	}
	rows.Close() // Cerrar antes de ejecutar ALTER o DROP

	// SI LA COLUMNA 'nombre' NO EXISTE, la tabla está corrupta.
	if !columnExists["nombre"] {
		log.Println("⚠️ Tabla 'proyectos' está corrupta (falta 'nombre'). Recreando...")
		_, err = DB.Exec("DROP TABLE IF EXISTS proyectos;")
		if err != nil {
			log.Fatalf("Error crítico al 'DROP' la tabla 'proyectos' corrupta: %v", err)
		}
		createProyectosTableSQL := `
		CREATE TABLE IF NOT EXISTS proyectos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			nombre TEXT NOT NULL UNIQUE,
			fecha_inicio TEXT NOT NULL DEFAULT '',
			fecha_cierre TEXT NOT NULL DEFAULT ''
		);`
		if _, err = DB.Exec(createProyectosTableSQL); err != nil {
			log.Fatalf("Error crítico al RE-CREAR la tabla 'proyectos': %v", err)
		}
		log.Println("✅ Tabla 'proyectos' recreada exitosamente.")
		// Como la recreamos, ya tiene todas las columnas. Salimos de la función.
		return
	}

	// Si 'nombre' existe, revisamos las otras columnas (migración normal)
	if !columnExists["fecha_inicio"] {
		log.Println("⚠️ Migrando: Añadiendo 'fecha_inicio' a 'proyectos'...")
		_, err = DB.Exec("ALTER TABLE proyectos ADD COLUMN fecha_inicio TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error al migrar 'fecha_inicio' en 'proyectos': %v", err)
		}
		log.Println("✅ Columna 'fecha_inicio' agregada exitosamente.")
	}

	if !columnExists["fecha_cierre"] {
		log.Println("⚠️ Migrando: Añadiendo 'fecha_cierre' a 'proyectos'...")
		_, err = DB.Exec("ALTER TABLE proyectos ADD COLUMN fecha_cierre TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error al migrar 'fecha_cierre' en 'proyectos': %v", err)
		}
		log.Println("✅ Columna 'fecha_cierre' agregada exitosamente.")
	}
}

// migrateUsersTable (sin cambios)
func migrateUsersTable() {
	rows, err := DB.Query("PRAGMA table_info(users)")
	if err != nil {
		log.Fatalf("Error al leer la información del esquema (PRAGMA) para 'users': %v", err)
	}
	defer rows.Close()

	columnExists := map[string]bool{
		"role":        false,
		"nombre":      false,
		"apellido":    false,
		"proyecto_id": false,
	}

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
			log.Fatalf("Error al escanear información de columna de 'users': %v", err)
		}
		if _, ok := columnExists[name]; ok {
			columnExists[name] = true
		}
	}
	rows.Close()

	if !columnExists["role"] {
		log.Println("⚠️ Migrando: Añadiendo columna 'role' a 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user'")
		if err != nil {
			log.Fatalf("Error al migrar 'role': %v", err)
		}
	}
	if !columnExists["nombre"] {
		log.Println("⚠️ Migrando: Añadiendo columna 'nombre' a 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN nombre TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error al migrar 'nombre': %v", err)
		}
	}
	if !columnExists["apellido"] {
		log.Println("⚠️ Migrando: Añadiendo columna 'apellido' a 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN apellido TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error al migrar 'apellido': %v", err)
		}
	}
	if !columnExists["proyecto_id"] {
		log.Println("⚠️ Migrando: Añadiendo columna 'proyecto_id' a 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN proyecto_id INTEGER REFERENCES proyectos(id) ON DELETE SET NULL")
		if err != nil {
			log.Fatalf("Error al migrar 'proyecto_id': %v", err)
		}
		log.Println("✅ Columna 'proyecto_id' agregada exitosamente a 'users'.")
	}
}

// checkAdminRole (sin cambios)
func checkAdminRole(username string) (bool, error) {
	var role string
	err := DB.QueryRow("SELECT role FROM users WHERE username = ?", username).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error de DB al verificar rol: %w", err)
	}
	return strings.ToLower(role) == "admin", nil
}

// --- HANDLERS DE AUTENTICACIÓN (Sin cambios) ---

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}
	if len(user.Password) < 6 || user.Nombre == "" || user.Apellido == "" {
		http.Error(w, `{"error": "Todos los campos (usuario, pass > 6, nombre, apellido) son requeridos."}`, http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Error interno al procesar contraseña."}`, http.StatusInternalServerError)
		return
	}
	query := "INSERT INTO users(username, password, role, nombre, apellido) VALUES(?, ?, ?, ?, ?)"
	_, err = DB.Exec(query, user.Username, string(hashedPassword), "user", user.Nombre, user.Apellido)
	if err != nil {
		http.Error(w, `{"error": "El nombre de usuario ya existe. Intenta otro."}`, http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"mensaje": "Usuario registrado exitosamente"}`)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, `{"error": "Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDB.HashedPassword), []byte(userReq.Password))
	if err != nil {
		http.Error(w, `{"error": "Credenciales inválidas"}`, http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, `{"mensaje": "Inicio de sesión exitoso", "usuario": "%s", "id": %d, "role": "%s"}`,
		userDB.Username, userDB.ID, userDB.Role)
}

// --- HANDLERS DE ADMINISTRACIÓN (USUARIOS) ---

// adminUsersHandler (sin cambios)
func adminUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var adminReq AdminActionRequest
	json.NewDecoder(r.Body).Decode(&adminReq)

	isAdmin, err := checkAdminRole(adminReq.AdminUsername)
	if err != nil || !isAdmin {
		http.Error(w, `{"error": "Acceso denegado. Solo administradores."}`, http.StatusForbidden)
		return
	}

	query := `
		SELECT 
			u.id, u.username, u.role, u.nombre, u.apellido,
			p.id, p.nombre 
		FROM users u
		LEFT JOIN proyectos p ON u.proyecto_id = p.id
		ORDER BY u.id ASC
	`
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Error al consultar usuarios con join: %v", err)
		http.Error(w, `{"error": "Error interno al consultar usuarios."}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []UserListResponse{}
	for rows.Next() {
		var u UserListResponse
		var proyectoID sql.NullInt64
		var proyectoNombre sql.NullString
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.Nombre, &u.Apellido, &proyectoID, &proyectoNombre); err != nil {
			log.Printf("Error al escanear usuario: %v", err)
			continue
		}
		if proyectoID.Valid {
			id := int(proyectoID.Int64)
			u.ProyectoID = &id
		}
		if proyectoNombre.Valid {
			u.ProyectoNombre = &proyectoNombre.String
		}
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
}

// adminAddUserHandler (sin cambios)
func adminAddUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}
	isAdmin, err := checkAdminRole(req.AdminUsername)
	if err != nil || !isAdmin {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}
	if len(req.User.Username) == 0 || len(req.User.Password) < 6 || req.User.Nombre == "" || req.User.Apellido == "" {
		http.Error(w, `{"error": "Todos los campos (usuario, pass > 6, nombre, apellido) son requeridos."}`, http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Error interno al procesar contraseña."}`, http.StatusInternalServerError)
		return
	}
	query := "INSERT INTO users(username, password, role, nombre, apellido) VALUES(?, ?, ?, ?, ?)"
	_, err = DB.Exec(query, req.User.Username, string(hashedPassword), "user", req.User.Nombre, req.User.Apellido)
	if err != nil {
		http.Error(w, `{"error": "El nombre de usuario ya existe."}`, http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"mensaje": "Usuario '%s' agregado exitosamente."}`, req.User.Username)
}

// adminDeleteUserHandler (sin cambios)
func adminDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}
	isAdmin, err := checkAdminRole(req.AdminUsername)
	if err != nil || !isAdmin {
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
		http.Error(w, `{"error": "No puedes borrar tu propia cuenta."}`, http.StatusForbidden)
		return
	}
	query := "DELETE FROM users WHERE id = ?"
	_, err = DB.Exec(query, req.ID)
	if err != nil {
		http.Error(w, `{"error": "Error interno al borrar el usuario."}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": fmt.Sprintf("Usuario '%s' borrado exitosamente.", targetUsername),
	})
}

// adminUpdateUserHandler (sin cambios)
func adminUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}
	isAdmin, err := checkAdminRole(req.AdminUsername)
	if err != nil || !isAdmin {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}
	newRole := strings.ToLower(req.NewRole)
	if newRole != "admin" && newRole != "user" {
		http.Error(w, `{"error": "Rol no válido."}`, http.StatusBadRequest)
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
		http.Error(w, `{"error": "Error interno al actualizar el rol."}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": fmt.Sprintf("Rol del usuario ID %d actualizado a %s.", req.ID, newRole),
	})
}

// adminAssignProyectoHandler (sin cambios)
func adminAssignProyectoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req AssignProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}
	isAdmin, err := checkAdminRole(req.AdminUsername)
	if err != nil || !isAdmin {
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
		http.Error(w, `{"error": "Error interno al asignar el proyecto."}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Proyecto asignado/actualizado exitosamente.",
	})
}

// --- HANDLERS DE ADMINISTRACIÓN (PROYECTOS) ---

// adminGetProyectosHandler (ACTUALIZADO)
func adminGetProyectosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var adminReq AdminActionRequest
	json.NewDecoder(r.Body).Decode(&adminReq)

	isAdmin, err := checkAdminRole(adminReq.AdminUsername)
	if err != nil || !isAdmin {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}

	// AHORA SELECCIONA LAS FECHAS
	rows, err := DB.Query("SELECT id, nombre, fecha_inicio, fecha_cierre FROM proyectos ORDER BY nombre ASC")
	if err != nil {
		http.Error(w, `{"error": "Error al consultar proyectos."}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	proyectos := []Proyecto{}
	for rows.Next() {
		var p Proyecto
		// AHORA ESCANEA LAS FECHAS
		if err := rows.Scan(&p.ID, &p.Nombre, &p.FechaInicio, &p.FechaCierre); err != nil {
			log.Printf("Error al escanear proyecto: %v", err)
			continue
		}
		proyectos = append(proyectos, p)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"proyectos": proyectos})
}

// adminCreateProyectoHandler (ACTUALIZADO)
func adminCreateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateProyectoRequest // Struct actualizada
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}

	isAdmin, err := checkAdminRole(req.AdminUsername)
	if err != nil || !isAdmin {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}

	// Validación simple
	if req.Nombre == "" || req.FechaInicio == "" {
		http.Error(w, `{"error": "El nombre y la fecha de inicio son obligatorios."}`, http.StatusBadRequest)
		return
	}

	// AHORA INSERTA LAS FECHAS
	query := "INSERT INTO proyectos(nombre, fecha_inicio, fecha_cierre) VALUES(?, ?, ?)"
	_, err = DB.Exec(query, req.Nombre, req.FechaInicio, req.FechaCierre)
	if err != nil {
		log.Printf("Error al crear proyecto: %v", err)
		http.Error(w, `{"error": "El nombre del proyecto ya existe o hubo un error."}`, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": fmt.Sprintf("Proyecto '%s' creado exitosamente.", req.Nombre),
	})
}

// adminDeleteProyectoHandler (sin cambios)
func adminDeleteProyectoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req DeleteProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}

	isAdmin, err := checkAdminRole(req.AdminUsername)
	if err != nil || !isAdmin {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}

	query := "DELETE FROM proyectos WHERE id = ?"
	result, err := DB.Exec(query, req.ID)
	if err != nil {
		http.Error(w, `{"error": "Error interno al borrar el proyecto."}`, http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error": "El proyecto con ese ID no existe."}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Proyecto borrado exitosamente. Los usuarios asignados ahora están libres.",
	})
}

// ⭐️ --- NUEVO HANDLER PARA MODIFICAR PROYECTOS --- ⭐️
func adminUpdateProyectoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req UpdateProyectoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}

	isAdmin, err := checkAdminRole(req.AdminUsername)
	if err != nil || !isAdmin {
		http.Error(w, `{"error": "Acceso denegado."}`, http.StatusForbidden)
		return
	}

	if req.ID == 0 || req.Nombre == "" || req.FechaInicio == "" {
		http.Error(w, `{"error": "El ID, nombre y fecha de inicio son obligatorios."}`, http.StatusBadRequest)
		return
	}

	query := "UPDATE proyectos SET nombre = ?, fecha_inicio = ?, fecha_cierre = ? WHERE id = ?"
	_, err = DB.Exec(query, req.Nombre, req.FechaInicio, req.FechaCierre, req.ID)
	if err != nil {
		log.Printf("Error al actualizar proyecto: %v", err)
		// Chequea si es por 'UNIQUE constraint'
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, `{"error": "Ese nombre de proyecto ya existe."}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error": "Error interno al actualizar el proyecto."}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": fmt.Sprintf("Proyecto '%s' actualizado exitosamente.", req.Nombre),
	})
}

func main() {
	initDB()
	defer DB.Close()

	// Handlers
	saludoAPI := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"mensaje": "Hola desde Go!"}`)
	})
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

	// NUEVO HANDLER
	adminUpdateProyectoAPI := http.HandlerFunc(adminUpdateProyectoHandler)

	// Configuración de CORS
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

	// NUEVA RUTA
	http.Handle("/api/admin/update-proyecto", corsOptions(adminUpdateProyectoAPI))

	log.Println("🚀 Servidor Go escuchando en :8080 (v5 - Proyectos con fechas y update)")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
