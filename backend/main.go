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
}

type UserDB struct {
	ID             int
	Username       string
	HashedPassword string
	Role           string
}

type UserListResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
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

// --- UTILIDADES DE BASE DE DATOS ---

// initDB inicializa la conexión a la base de datos SQLite y crea la tabla 'users'.
func initDB() {
	var err error
	// ⭐️ ¡CAMBIO AQUÍ!
	// Usamos `../` para subir al directorio padre (Laboratorio-Ingenieria-de-Software-I)
	// y luego acceder a "Base de datos/users.db"
	DB, err = sql.Open("sqlite", "../Base de datos/users.db")
	if err != nil {
		log.Fatalf("Error al abrir la base de datos: %v", err)
	}

	// 1. SQL para crear la tabla con el campo 'role'
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user' -- ROL AÑADIDO con valor por defecto 'user'
	);`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error al crear la tabla 'users': %v", err)
	}

	// ⭐️ INICIO LÓGICA DE MIGRACIÓN ROBUSTA
	rows, err := DB.Query("PRAGMA table_info(users)")
	if err != nil {
		log.Fatalf("Error al leer la información del esquema (PRAGMA): %v", err)
	}
	defer rows.Close()

	columnExists := false
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
			log.Fatalf("Error al escanear información de columna: %v", err)
		}
		if name == "role" {
			columnExists = true
			break
		}
	}

	if !columnExists {
		log.Println("⚠️ Columna 'role' ausente. Ejecutando migración...")
		_, migrateErr := DB.Exec("ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user'")
		if migrateErr != nil {
			log.Fatalf("Error crítico al migrar la tabla: %v", migrateErr)
		}
		log.Println("✅ Columna 'role' agregada exitosamente.")
	}
	// ⭐️ FIN LÓGICA DE MIGRACIÓN

	// 2. Crear un usuario administrador por defecto si no existe
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

		insertAdminSQL := "INSERT INTO users(username, password, role) VALUES(?, ?, ?)"
		_, err = DB.Exec(insertAdminSQL, "admin", string(hashedPassword), "admin")
		if err != nil {
			log.Fatalf("Error al insertar usuario admin por defecto: %v", err)
		}
		log.Println("✅ Usuario administrador 'admin' (pass: admin123) creado por defecto.")
	}

	log.Println("✅ Base de datos SQLite inicializada y tabla 'users' creada.")
}

// checkAdminRole verifica si el usuario dado tiene el rol de administrador.
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

// --- HANDLERS DE AUTENTICACIÓN ---

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}

	if len(user.Password) < 6 {
		http.Error(w, `{"error": "La contraseña debe tener al menos 6 caracteres."}`, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error al hashear: %v", err)
		http.Error(w, `{"error": "Error interno al procesar contraseña."}`, http.StatusInternalServerError)
		return
	}

	// Insertamos con el rol por defecto 'user'
	query := "INSERT INTO users(username, password, role) VALUES(?, ?, ?)"
	_, err = DB.Exec(query, user.Username, string(hashedPassword), "user")
	if err != nil {
		// Detectar error de usuario duplicado (UNIQUE constraint)
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
		log.Printf("Error al buscar usuario en DB: %v", err)
		http.Error(w, `{"error": "Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDB.HashedPassword), []byte(userReq.Password))

	if err != nil {
		// Las contraseñas no coinciden
		http.Error(w, `{"error": "Credenciales inválidas"}`, http.StatusUnauthorized)
		return
	}

	// ¡Autenticación exitosa! Devolvemos el rol.
	fmt.Fprintf(w, `{"mensaje": "Inicio de sesión exitoso", "usuario": "%s", "id": %d, "role": "%s"}`,
		userDB.Username, userDB.ID, userDB.Role)
}

// --- HANDLERS DE ADMINISTRACIÓN ---

func adminUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var adminReq AdminActionRequest
	json.NewDecoder(r.Body).Decode(&adminReq)

	if adminReq.AdminUsername == "" {
		log.Println("Advertencia: admin_username ausente en /api/admin/users.")
	}

	isAdmin, err := checkAdminRole(adminReq.AdminUsername)
	if err != nil || !isAdmin {
		http.Error(w, `{"error": "Acceso denegado. Solo administradores pueden ver esta lista."}`, http.StatusForbidden)
		return
	}

	rows, err := DB.Query("SELECT id, username, role FROM users ORDER BY id ASC")
	if err != nil {
		log.Printf("Error al consultar usuarios: %v", err)
		http.Error(w, `{"error": "Error interno al consultar usuarios."}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []UserListResponse{}
	for rows.Next() {
		var u UserListResponse
		if err := rows.Scan(&u.ID, &u.Username, &u.Role); err != nil {
			log.Printf("Error al escanear usuario: %v", err)
			continue
		}
		users = append(users, u)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
}

func adminAddUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}

	isAdmin, err := checkAdminRole(req.AdminUsername)
	if err != nil || !isAdmin {
		http.Error(w, `{"error": "Acceso denegado. Solo administradores pueden agregar usuarios."}`, http.StatusForbidden)
		return
	}

	if len(req.User.Username) == 0 || len(req.User.Password) < 6 {
		http.Error(w, `{"error": "Nombre de usuario o contraseña (mínimo 6 caracteres) no válidos."}`, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error al hashear: %v", err)
		http.Error(w, `{"error": "Error interno al procesar contraseña."}`, http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users(username, password, role) VALUES(?, ?, ?)"
	_, err = DB.Exec(query, req.User.Username, string(hashedPassword), "user")
	if err != nil {
		http.Error(w, `{"error": "El nombre de usuario ya existe. Intenta otro."}`, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"mensaje": "Usuario '%s' agregado exitosamente con rol 'user'."}`, req.User.Username)
}

func adminDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}

	isAdmin, err := checkAdminRole(req.AdminUsername)
	if err != nil || !isAdmin {
		http.Error(w, `{"error": "Acceso denegado. Solo administradores pueden borrar usuarios."}`, http.StatusForbidden)
		return
	}

	var targetUsername string
	err = DB.QueryRow("SELECT username FROM users WHERE id = ?", req.ID).Scan(&targetUsername)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "Usuario a borrar no encontrado."}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "Error al buscar usuario."}`, http.StatusInternalServerError)
		return
	}

	if targetUsername == req.AdminUsername {
		http.Error(w, `{"error": "No puedes borrar tu propia cuenta de administrador por seguridad."}`, http.StatusForbidden)
		return
	}

	query := "DELETE FROM users WHERE id = ?"
	result, err := DB.Exec(query, req.ID)
	if err != nil {
		log.Printf("Error al borrar usuario: %v", err)
		http.Error(w, `{"error": "Error interno al borrar el usuario."}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error": "El usuario con ese ID no existe."}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": fmt.Sprintf("Usuario con ID %d ('%s') borrado exitosamente.", req.ID, targetUsername),
	})
}

func adminUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Formato JSON inválido."}`, http.StatusBadRequest)
		return
	}

	isAdmin, err := checkAdminRole(req.AdminUsername)
	if err != nil || !isAdmin {
		http.Error(w, `{"error": "Acceso denegado. Solo administradores pueden modificar roles."}`, http.StatusForbidden)
		return
	}

	newRole := strings.ToLower(req.NewRole)
	if newRole != "admin" && newRole != "user" {
		http.Error(w, `{"error": "Rol no válido. Debe ser 'admin' o 'user'."}`, http.StatusBadRequest)
		return
	}

	var targetUsername string
	err = DB.QueryRow("SELECT username FROM users WHERE id = ?", req.ID).Scan(&targetUsername)
	if err != nil {
		http.Error(w, `{"error": "Usuario a modificar no encontrado."}`, http.StatusNotFound)
		return
	}

	if targetUsername == req.AdminUsername {
		http.Error(w, `{"error": "No puedes cambiar tu propio rol por seguridad."}`, http.StatusForbidden)
		return
	}

	query := "UPDATE users SET role = ? WHERE id = ?"
	result, err := DB.Exec(query, newRole, req.ID)
	if err != nil {
		log.Printf("Error al actualizar rol: %v", err)
		http.Error(w, `{"error": "Error interno al actualizar el rol."}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error": "El usuario con ese ID no existe."}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": fmt.Sprintf("Rol del usuario ID %d actualizado a %s.", req.ID, newRole),
	})
}

func main() {
	// ⭐️ 1. Inicializar la base de datos
	initDB()
	defer DB.Close()

	// Definir handlers para las rutas
	saludoAPI := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"mensaje": "Hola desde Go, tu backend está corriendo y conectado a SQLite!"}`)
	})
	registerAPI := http.HandlerFunc(registerHandler)
	loginAPI := http.HandlerFunc(loginHandler)

	// Handlers de Admin
	adminUsersAPI := http.HandlerFunc(adminUsersHandler)
	adminUpdateUserAPI := http.HandlerFunc(adminUpdateUserHandler)
	adminAddUserAPI := http.HandlerFunc(adminAddUserHandler)
	adminDeleteUserAPI := http.HandlerFunc(adminDeleteUserHandler)

	// Configuración de CORS
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	corsOptions := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)

	// Asignar los handlers con CORS a sus rutas
	http.Handle("/api/saludo", corsOptions(saludoAPI))
	http.Handle("/api/register", corsOptions(registerAPI))
	http.Handle("/api/login", corsOptions(loginAPI))

	// Rutas de Administración
	http.Handle("/api/admin/users", corsOptions(adminUsersAPI))
	http.Handle("/api/admin/update-user", corsOptions(adminUpdateUserAPI))
	http.Handle("/api/admin/add-user", corsOptions(adminAddUserAPI))
	http.Handle("/api/admin/delete-user", corsOptions(adminDeleteUserAPI))

	// Iniciar el servidor en el puerto 8080
	log.Println("🚀 Servidor Go escuchando en :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
