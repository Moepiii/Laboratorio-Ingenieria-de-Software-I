package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func setupTestDB() {

	var err error

	DB, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("Error al abrir la base de datos en memoria: %v", err)
	}

	// Crear la tabla 'users'
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user'
	);`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error al crear la tabla 'users' en memoria: %v", err)
	}

	// 2. Crear un usuario administrador por defecto
	adminPassword := "admin123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	insertAdminSQL := "INSERT INTO users(username, password, role) VALUES(?, ?, ?)"
	_, err = DB.Exec(insertAdminSQL, "admin", string(hashedPassword), "admin")
	if err != nil {
		log.Fatalf("Error al insertar admin de prueba: %v", err)
	}
}

// TestMain ejecuta setupTestDB antes de todas las pruebas
func TestMain(m *testing.M) {
	// Configuración antes de todas las pruebas
	setupTestDB()

	// Ejecutar las pruebas
	exitCode := m.Run()

	// Cierre de la base de datos después de todas las pruebas
	if DB != nil {
		DB.Close()
	}

	// Salir con el código de estado
	os.Exit(exitCode)
}

// --- PRUEBAS DE UTILIDADES (Funciones internas) ---

func TestCheckAdminRole(t *testing.T) {
	// Insertar un usuario normal
	userPass := "userpass"
	hashedUserPass, _ := bcrypt.GenerateFromPassword([]byte(userPass), bcrypt.DefaultCost)
	DB.Exec("INSERT INTO users(username, password, role) VALUES(?, ?, ?)", "testuser", string(hashedUserPass), "user")

	tests := []struct {
		username string
		want     bool
		wantErr  bool
	}{
		{"admin", true, false},        // Usuario administrador por defecto
		{"testuser", false, false},    // Usuario normal
		{"nonexistent", false, false}, // Usuario inexistente
	}

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			got, err := checkAdminRole(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkAdminRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkAdminRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- PRUEBAS DE HANDLERS DE AUTENTICACIÓN ---

func TestRegisterHandler_Success(t *testing.T) {
	// Datos de prueba para un nuevo usuario
	user := User{Username: "newuser1", Password: "password123"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler devolvió código de estado incorrecto: got %v, want %v. Body: %s", status, http.StatusCreated, rr.Body.String())
	}

	// Verificar si el usuario fue insertado en la DB
	var role string
	err := DB.QueryRow("SELECT role FROM users WHERE username = ?", user.Username).Scan(&role)
	if err != nil {
		t.Fatalf("Error al verificar usuario en DB: %v", err)
	}
	if role != "user" {
		t.Errorf("El rol del nuevo usuario no es 'user', got: %s", role)
	}
}

func TestRegisterHandler_Conflict(t *testing.T) {
	// Intentar registrar un usuario que ya existe ("admin")
	user := User{Username: "admin", Password: "password123"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(registerHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("Handler devolvió código de estado incorrecto: got %v, want %v", status, http.StatusConflict)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	// El usuario admin ya está en la DB
	user := User{Username: "admin", Password: "admin123"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(loginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Login falló con código: %v, body: %s", status, rr.Body.String())
	}

	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	if response["role"] != "admin" {
		t.Errorf("Rol de respuesta incorrecto: got %v, want admin", response["role"])
	}
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	// Contraseña incorrecta
	user := User{Username: "admin", Password: "wrongpassword"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(loginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler devolvió código de estado incorrecto: got %v, want %v", status, http.StatusUnauthorized)
	}
}

// --- PRUEBAS DE HANDLERS DE ADMINISTRACIÓN (Tus US 2.a y 2.b) ---

// Helper para crear una solicitud de admin
func newAdminRequest(method, url string, adminUsername string, data interface{}) *http.Request {
	var body io.Reader
	if data != nil {
		jsonBody, _ := json.Marshal(data)
		body = bytes.NewBuffer(jsonBody)
	} else {
		// Para GET, solo incluiremos el admin_username en el cuerpo para cumplir con la estructura
		// (aunque para GET lo ideal es usar headers o query params)
		jsonBody, _ := json.Marshal(AdminActionRequest{AdminUsername: adminUsername})
		body = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestAdminUsersHandler_AccessDenied(t *testing.T) {
	// Intentar acceder con un usuario normal ("testuser" creado en TestCheckAdminRole)
	req := newAdminRequest("POST", "/api/admin/users", "testuser", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(adminUsersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler devolvió código de estado incorrecto: got %v, want %v", status, http.StatusForbidden)
	}
}

func TestAdminUsersHandler_Success(t *testing.T) {
	req := newAdminRequest("POST", "/api/admin/users", "admin", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(adminUsersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Handler devolvió código de estado incorrecto: got %v, want %v. Body: %s", status, http.StatusOK, rr.Body.String())
	}

	var response map[string][]UserListResponse
	json.NewDecoder(rr.Body).Decode(&response)

	if len(response["users"]) < 2 { // Debería tener al menos 'admin' y 'testuser'
		t.Errorf("Lista de usuarios incompleta: got %d, want >= 2", len(response["users"]))
	}
}

// Corresponde a la US 2.b: Crear perfiles de Usuario (rol por defecto 'user' aquí)
func TestAdminAddUserHandler_Success(t *testing.T) {
	addUserReq := AddUserRequest{
		User:          User{Username: "usertoadd", Password: "securepass"},
		AdminUsername: "admin",
	}

	req := newAdminRequest("POST", "/api/admin/add-user", "admin", addUserReq)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(adminAddUserHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("Handler devolvió código de estado incorrecto: got %v, want %v. Body: %s", status, http.StatusCreated, rr.Body.String())
	}
}

// Corresponde a la US 2.b: Crear perfiles de Usuario (Fallido por no-admin)
func TestAdminAddUserHandler_AccessDenied(t *testing.T) {
	addUserReq := AddUserRequest{
		User:          User{Username: "usertofail", Password: "securepass"},
		AdminUsername: "testuser", // Usuario normal
	}

	req := newAdminRequest("POST", "/api/admin/add-user", "testuser", addUserReq)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(adminAddUserHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler devolvió código de estado incorrecto: got %v, want %v", status, http.StatusForbidden)
	}
}

// Corresponde a la US 2.a: Autenticar Usuarios y parte de 2.b: Control de Acceso
func TestAdminUpdateUserHandler_Success(t *testing.T) {
	// El ID de 'testuser' debe ser encontrado
	var testUserID int
	err := DB.QueryRow("SELECT id FROM users WHERE username = 'testuser'").Scan(&testUserID)
	if err != nil {
		t.Fatalf("No se pudo obtener el ID de 'testuser': %v", err)
	}

	updateReq := UpdateRoleRequest{
		ID:            testUserID,
		NewRole:       "admin",
		AdminUsername: "admin",
	}

	req := newAdminRequest("PUT", "/api/admin/update-user", "admin", updateReq)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(adminUpdateUserHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Handler devolvió código de estado incorrecto: got %v, want %v. Body: %s", status, http.StatusOK, rr.Body.String())
	}

	// Verificar si el rol fue actualizado en la DB
	var role string
	DB.QueryRow("SELECT role FROM users WHERE id = ?", testUserID).Scan(&role)
	if role != "admin" {
		t.Errorf("El rol no fue actualizado a 'admin', got: %s", role)
	}
}

func TestAdminDeleteUserHandler_Success(t *testing.T) {
	// Insertar un usuario para borrar
	userToDelete := "usertodelete"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	res, _ := DB.Exec("INSERT INTO users(username, password, role) VALUES(?, ?, ?)", userToDelete, string(hashedPass), "user")
	idToDelete, _ := res.LastInsertId()

	deleteReq := DeleteUserRequest{
		ID:            int(idToDelete),
		AdminUsername: "admin",
	}

	req := newAdminRequest("DELETE", "/api/admin/delete-user", "admin", deleteReq)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(adminDeleteUserHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Handler devolvió código de estado incorrecto: got %v, want %v. Body: %s", status, http.StatusOK, rr.Body.String())
	}

	// Verificar que el usuario ya no existe en la DB
	var username string
	err := DB.QueryRow("SELECT username FROM users WHERE id = ?", idToDelete).Scan(&username)
	if err != sql.ErrNoRows {
		t.Errorf("El usuario no fue borrado o se encontró un error inesperado: %v", err)
	}
}

func TestAdminDeleteUserHandler_SelfDeleteForbidden(t *testing.T) {
	// Intentar borrar la cuenta 'admin'
	var adminID int
	DB.QueryRow("SELECT id FROM users WHERE username = 'admin'").Scan(&adminID)

	deleteReq := DeleteUserRequest{
		ID:            adminID,
		AdminUsername: "admin",
	}

	req := newAdminRequest("DELETE", "/api/admin/delete-user", "admin", deleteReq)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(adminDeleteUserHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler devolvió código de estado incorrecto: got %v, want %v", status, http.StatusForbidden)
	}
}
