package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite" // Driver SQLite

	"proyecto/internal/auth"
	"proyecto/internal/database"             // Usaremos database.DB
	apphandlers "proyecto/internal/handlers" 
	"proyecto/internal/models"
)

// setupTestDB crea una base de datos en memoria para las pruebas
func setupTestDB() {
	var err error
	// Usa database.DB globalmente para las pruebas
	database.DB, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatalf("Error al abrir la base de datos en memoria: %v", err)
	}

	_, err = database.DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalf("Error al habilitar foreign keys: %v", err)
	}

	// Crear tabla 'users' con todas las columnas
	createUsersSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL UNIQUE, password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user', nombre TEXT NOT NULL DEFAULT '', apellido TEXT NOT NULL DEFAULT '',
		proyecto_id INTEGER REFERENCES proyectos(id) ON DELETE SET NULL
	);`
	if _, err = database.DB.Exec(createUsersSQL); err != nil {
		log.Fatalf("Error al crear tabla 'users' en memoria: %v", err)
	}

	// Crear tabla 'proyectos' con todas las columnas
	createProyectosSQL := `
	CREATE TABLE IF NOT EXISTS proyectos (
		id INTEGER PRIMARY KEY AUTOINCREMENT, nombre TEXT NOT NULL UNIQUE,
		fecha_inicio TEXT NOT NULL DEFAULT '', fecha_cierre TEXT NOT NULL DEFAULT '',
		estado TEXT NOT NULL DEFAULT 'habilitado'
	);`
	if _, err = database.DB.Exec(createProyectosSQL); err != nil {
		log.Fatalf("Error al crear tabla 'proyectos' en memoria: %v", err)
	}

	// Crear usuario admin de prueba
	adminPassword := "admin123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	insertAdminSQL := "INSERT INTO users(username, password, role, nombre, apellido) VALUES(?, ?, ?, ?, ?)"
	_, err = database.DB.Exec(insertAdminSQL, "admin", string(hashedPassword), "admin", "Admin", "User")
	if err != nil {
		log.Fatalf("Error al insertar admin de prueba: %v", err)
	}
}

// TestMain ejecuta setupTestDB antes de todas las pruebas
func TestMain(m *testing.M) {
	setupTestDB()
	exitCode := m.Run()
	if database.DB != nil {
		database.DB.Close()
	}
	os.Exit(exitCode)
}

// --- PRUEBAS DE UTILIDADES (Funciones internas) ---

// Renombrado de TestCheckAdminRole a TestCheckPermission
func TestCheckPermission(t *testing.T) {
	// Insertar usuarios de prueba adicionales
	userPass := "userpass"
	hashedUserPass, _ := bcrypt.GenerateFromPassword([]byte(userPass), bcrypt.DefaultCost)
	database.DB.Exec("INSERT INTO users(username, password, role) VALUES(?, ?, ?)", "testuser", string(hashedUserPass), "user")
	gerentePass := "gerentepass"
	hashedGerentePass, _ := bcrypt.GenerateFromPassword([]byte(gerentePass), bcrypt.DefaultCost)
	database.DB.Exec("INSERT INTO users(username, password, role) VALUES(?, ?, ?)", "testgerente", string(hashedGerentePass), "gerente")

	tests := []struct {
		username      string
		requiredRoles []string
		want          bool
		wantErr       bool
	}{
		{"admin", []string{"admin"}, true, false},
		{"admin", []string{"admin", "gerente"}, true, false},
		{"admin", []string{"user"}, false, false},
		{"testuser", []string{"user"}, true, false},
		{"testuser", []string{"admin"}, false, false},
		{"testuser", []string{"admin", "gerente"}, false, false},
		{"testgerente", []string{"gerente"}, true, false},
		{"testgerente", []string{"admin", "gerente"}, true, false},
		{"testgerente", []string{"admin"}, false, false},
		{"nonexistent", []string{"admin"}, false, true}, // Esperamos error porque no existe
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_req_%s", tt.username, tt.requiredRoles), func(t *testing.T) {
			got, err := auth.CheckPermission(tt.username, tt.requiredRoles...) // Llama a la función de auth
			if (err != nil) != tt.wantErr {
				// Si el error es sql.ErrNoRows y esperábamos error, está bien
				if !(err == sql.ErrNoRows && tt.wantErr) {
					t.Errorf("CheckPermission() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("CheckPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- PRUEBAS DE HANDLERS DE AUTENTICACIÓN ---

func TestRegisterHandler_Success(t *testing.T) {
	// Usa models.User
	user := models.User{Username: "newuser1", Password: "password123", Nombre: "New", Apellido: "User"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	// Usa apphandlers.RegisterHandler
	handler := http.HandlerFunc(apphandlers.RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler devolvió código de estado incorrecto: got %v, want %v. Body: %s", status, http.StatusCreated, rr.Body.String())
	}

	var role string
	err := database.DB.QueryRow("SELECT role FROM users WHERE username = ?", user.Username).Scan(&role)
	if err != nil {
		t.Fatalf("Error al verificar usuario en DB: %v", err)
	}
	if role != "user" {
		t.Errorf("El rol del nuevo usuario no es 'user', got: %s", role)
	}
}

func TestRegisterHandler_Conflict(t *testing.T) {
	// Usa models.User
	user := models.User{Username: "admin", Password: "password123", Nombre: "Admin", Apellido: "UserTest"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	// Usa apphandlers.RegisterHandler
	handler := http.HandlerFunc(apphandlers.RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("Handler devolvió código de estado incorrecto: got %v, want %v", status, http.StatusConflict)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	// Usa models.User
	user := models.User{Username: "admin", Password: "admin123"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	// Usa apphandlers.LoginHandler
	handler := http.HandlerFunc(apphandlers.LoginHandler)
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
	// Usa models.User
	user := models.User{Username: "admin", Password: "wrongpassword"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))

	rr := httptest.NewRecorder()
	// Usa apphandlers.LoginHandler
	handler := http.HandlerFunc(apphandlers.LoginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler devolvió código de estado incorrecto: got %v, want %v", status, http.StatusUnauthorized)
	}
}

// --- PRUEBAS DE HANDLERS DE ADMINISTRACIÓN ---

func newAdminRequest(method, url string, adminUsername string, data interface{}) *http.Request {
	var body io.Reader
	if data != nil {
		jsonBody, _ := json.Marshal(data)
		body = bytes.NewBuffer(jsonBody)
	} else {
		// Envía el username en el cuerpo incluso para GET (como lo hacen los handlers)
		jsonBody, _ := json.Marshal(models.AdminActionRequest{AdminUsername: adminUsername})
		body = bytes.NewBuffer(jsonBody)
	}
	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestAdminUsersHandler_AccessDenied_User(t *testing.T) {
	// Usuario normal "testuser" no debería poder acceder
	req := newAdminRequest("POST", "/api/admin/users", "testuser", nil)
	rr := httptest.NewRecorder()
	// Usa apphandlers.AdminUsersHandler
	handler := http.HandlerFunc(apphandlers.AdminUsersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler (user) devolvió código incorrecto: got %v, want %v", status, http.StatusForbidden)
	}
}
func TestAdminUsersHandler_Success_Admin(t *testing.T) {
	req := newAdminRequest("POST", "/api/admin/users", "admin", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apphandlers.AdminUsersHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Handler (admin) devolvió código incorrecto: got %v, want %v. Body: %s", status, http.StatusOK, rr.Body.String())
	}
	var response map[string][]models.UserListResponse
	json.NewDecoder(rr.Body).Decode(&response)
	if len(response["users"]) < 3 {
		t.Errorf("Lista de usuarios incompleta: got %d, want >= 3 (admin, testuser, testgerente)", len(response["users"]))
	}
}

func TestAdminUsersHandler_Success_Gerente(t *testing.T) {
	req := newAdminRequest("POST", "/api/admin/users", "testgerente", nil) // Prueba con gerente
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apphandlers.AdminUsersHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Handler (gerente) devolvió código incorrecto: got %v, want %v. Body: %s", status, http.StatusOK, rr.Body.String())
	}
	var response map[string][]models.UserListResponse
	json.NewDecoder(rr.Body).Decode(&response)
	if len(response["users"]) < 3 {
		t.Errorf("Lista de usuarios incompleta: got %d, want >= 3", len(response["users"]))
	}
}

func TestAdminAddUserHandler_Success(t *testing.T) {
	// Usa models.AddUserRequest y models.User
	addUserReq := models.AddUserRequest{
		User:          models.User{Username: "usertoadd", Password: "securepass", Nombre: "ToAdd", Apellido: "User"},
		AdminUsername: "admin",
	}
	req := newAdminRequest("POST", "/api/admin/add-user", "admin", addUserReq)
	rr := httptest.NewRecorder()
	// Usa apphandlers.AdminAddUserHandler
	handler := http.HandlerFunc(apphandlers.AdminAddUserHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("Handler devolvió código incorrecto: got %v, want %v. Body: %s", status, http.StatusCreated, rr.Body.String())
	}
}

func TestAdminAddUserHandler_AccessDenied_Gerente(t *testing.T) {
	addUserReq := models.AddUserRequest{
		User:          models.User{Username: "usertofailGerente", Password: "securepass", Nombre: "Fail", Apellido: "Gerente"},
		AdminUsername: "testgerente", // Gerente intenta añadir
	}
	req := newAdminRequest("POST", "/api/admin/add-user", "testgerente", addUserReq)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apphandlers.AdminAddUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler (gerente) devolvió código incorrecto: got %v, want %v", status, http.StatusForbidden)
	}
}

func TestAdminUpdateUserHandler_Success(t *testing.T) {
	var testUserID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = 'testuser'").Scan(&testUserID)
	if err != nil {
		t.Fatalf("No se pudo obtener ID de 'testuser': %v", err)
	}

	// Usa models.UpdateRoleRequest
	updateReq := models.UpdateRoleRequest{
		ID:            testUserID,
		NewRole:       "gerente", // Cambiamos a gerente
		AdminUsername: "admin",
	}
	req := newAdminRequest("POST", "/api/admin/update-user", "admin", updateReq) // POST ahora
	rr := httptest.NewRecorder()
	// Usa apphandlers.AdminUpdateUserHandler
	handler := http.HandlerFunc(apphandlers.AdminUpdateUserHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Handler devolvió código incorrecto: got %v, want %v. Body: %s", status, http.StatusOK, rr.Body.String())
	}
	var role string
	database.DB.QueryRow("SELECT role FROM users WHERE id = ?", testUserID).Scan(&role)
	if role != "gerente" {
		t.Errorf("El rol no fue actualizado a 'gerente', got: %s", role)
	}
}

func TestAdminUpdateUserHandler_AccessDenied_Gerente(t *testing.T) {
	var testUserID int
	database.DB.QueryRow("SELECT id FROM users WHERE username = 'testuser'").Scan(&testUserID)
	updateReq := models.UpdateRoleRequest{ID: testUserID, NewRole: "admin", AdminUsername: "testgerente"} // Gerente intenta cambiar rol
	req := newAdminRequest("POST", "/api/admin/update-user", "testgerente", updateReq)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apphandlers.AdminUpdateUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler (gerente) devolvió código incorrecto: got %v, want %v", status, http.StatusForbidden)
	}
}

func TestAdminDeleteUserHandler_Success(t *testing.T) {
	userToDelete := "usertodelete"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	res, _ := database.DB.Exec("INSERT INTO users(username, password, role) VALUES(?, ?, ?)", userToDelete, string(hashedPass), "user")
	idToDelete, _ := res.LastInsertId()

	// Usa models.DeleteUserRequest
	deleteReq := models.DeleteUserRequest{
		ID:            int(idToDelete),
		AdminUsername: "admin",
	}
	req := newAdminRequest("POST", "/api/admin/delete-user", "admin", deleteReq) // POST ahora
	rr := httptest.NewRecorder()
	// Usa apphandlers.AdminDeleteUserHandler
	handler := http.HandlerFunc(apphandlers.AdminDeleteUserHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Handler devolvió código incorrecto: got %v, want %v. Body: %s", status, http.StatusOK, rr.Body.String())
	}
	var username string
	err := database.DB.QueryRow("SELECT username FROM users WHERE id = ?", idToDelete).Scan(&username)
	if err != sql.ErrNoRows {
		t.Errorf("El usuario no fue borrado o error inesperado: %v", err)
	}
}

func TestAdminDeleteUserHandler_AccessDenied_Gerente(t *testing.T) {
	var testUserID int
	database.DB.QueryRow("SELECT id FROM users WHERE username = 'testuser'").Scan(&testUserID) // Usamos testuser ID
	deleteReq := models.DeleteUserRequest{ID: testUserID, AdminUsername: "testgerente"}        // Gerente intenta borrar
	req := newAdminRequest("POST", "/api/admin/delete-user", "testgerente", deleteReq)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apphandlers.AdminDeleteUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler (gerente) devolvió código incorrecto: got %v, want %v", status, http.StatusForbidden)
	}
}

func TestAdminDeleteUserHandler_SelfDeleteForbidden(t *testing.T) {
	var adminID int
	database.DB.QueryRow("SELECT id FROM users WHERE username = 'admin'").Scan(&adminID)
	// Usa models.DeleteUserRequest
	deleteReq := models.DeleteUserRequest{ID: adminID, AdminUsername: "admin"}
	req := newAdminRequest("POST", "/api/admin/delete-user", "admin", deleteReq) // POST ahora
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apphandlers.AdminDeleteUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler devolvió código incorrecto: got %v, want %v", status, http.StatusForbidden)
	}
}

// --- PRUEBAS DE HANDLERS DE PROYECTOS ---

func TestAdminCreateProyectoHandler_Success(t *testing.T) {

	createReq := models.CreateProyectoRequest{
		Nombre:        "Proyecto Secreto Alfa",
		FechaInicio:   "2025-01-01",
		FechaCierre:   "2025-12-31",
		AdminUsername: "admin", // Solo un 'admin' puede crear proyectos
	}

	// 2. Crear la petición
	req := newAdminRequest("POST", "/api/admin/create-proyecto", "admin", createReq)
	rr := httptest.NewRecorder()

	// 3. Ejecutar el handler
	handler := http.HandlerFunc(apphandlers.AdminCreateProyectoHandler)
	handler.ServeHTTP(rr, req)


	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("Handler devolvió código incorrecto: got %v, want %v. Body: %s",
			status, http.StatusCreated, rr.Body.String())
	}


	var nombre string
	// Usamos createReq.Nombre para verificar
	err := database.DB.QueryRow("SELECT nombre FROM proyectos WHERE nombre = ?", createReq.Nombre).Scan(&nombre)
	if err != nil {
		t.Fatalf("Error al verificar el proyecto en la DB: %v", err)
	}

	if nombre != createReq.Nombre {
		t.Errorf("El nombre del proyecto en la DB es incorrecto: got %s, want %s", nombre, createReq.Nombre)
	}
}


