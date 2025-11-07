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
	"proyecto/internal/database" // Usaremos database.DB
	apphandlers "proyecto/internal/handlers"
	"proyecto/internal/models"
)

// TestMain inicializa la DB en memoria y crea usuarios de prueba
func TestMain(m *testing.M) {
	// 1. Inicializa la DB en memoria
	// Usará las funciones (createUsersTable, createProyectosTable, etc.)
	// de TU ARCHIVO ORIGINAL database.go
	database.InitDB(":memory:")

	// 2. CREAMOS LOS USUARIOS DE PRUEBA (ya que el InitDB original no lo hace)
	// (Usamos bcrypt.DefaultCost, que es el correcto)
	adminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	userPass, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	gerentePass, _ := bcrypt.GenerateFromPassword([]byte("gerente123"), bcrypt.DefaultCost)

	// (Añadimos todos los campos que espera la tabla original: nombre, apellido, cedula)
	_, err := database.DB.Exec("INSERT INTO users(username, hashed_password, role, nombre, apellido, cedula) VALUES(?, ?, ?, ?, ?, ?)", "admin", string(adminPass), "admin", "Admin", "Test", "111111")
	if err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}

	_, err = database.DB.Exec("INSERT INTO users(username, hashed_password, role, nombre, apellido, cedula) VALUES(?, ?, ?, ?, ?, ?)", "testuser", string(userPass), "user", "User", "Test", "222222")
	if err != nil {
		log.Fatalf("Failed to seed testuser: %v", err)
	}

	_, err = database.DB.Exec("INSERT INTO users(username, hashed_password, role, nombre, apellido, cedula) VALUES(?, ?, ?, ?, ?, ?)", "testgerente", string(gerentePass), "gerente", "Gerente", "Test", "333333")
	if err != nil {
		log.Fatalf("Failed to seed testgerente: %v", err)
	}

	// 3. Ejecuta las pruebas
	exitCode := m.Run()
	if database.DB != nil {
		database.DB.Close()
	}
	os.Exit(exitCode)
}

// --- PRUEBAS DE UTILIDADES (Funciones internas) ---

func TestCheckPermission(t *testing.T) {
	// (Los usuarios ya fueron creados en TestMain)
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
	// (Tu models.User original sí tiene Cedula)
	user := models.User{Username: "newuser1", Password: "password123", Nombre: "New", Apellido: "User", Cedula: "444444"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
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
	// (Tu RegisterUser original SÍ asigna 'user' por defecto)
	if role != "user" {
		t.Errorf("El rol del nuevo usuario no es 'user', got: %s", role)
	}
}

func TestRegisterHandler_Conflict(t *testing.T) {
	// (Prueba usando el 'admin' que creamos en TestMain)
	user := models.User{Username: "admin", Password: "password123", Nombre: "Admin", Apellido: "UserTest", Cedula: "111111"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apphandlers.RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("Handler devolvió código de estado incorrecto: got %v, want %v", status, http.StatusConflict)
	}
}

func TestLoginHandler_Success(t *testing.T) {
	// (El admin (admin/admin123) fue creado en TestMain)
	user := models.User{Username: "admin", Password: "admin123"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
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
	user := models.User{Username: "admin", Password: "wrongpassword"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apphandlers.LoginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler devolvió código de estado incorrecto: got %v, want %v", status, http.StatusUnauthorized)
	}
}

// --- PRUEBAS DE HANDLERS DE ADMINISTRACIÓN ---

// newAdminRequest (Helper)
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
	req := newAdminRequest("POST", "/api/admin/users", "testuser", nil)
	rr := httptest.NewRecorder()
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
}

func TestAdminAddUserHandler_Success(t *testing.T) {
	addUserReq := models.AddUserRequest{
		User:          models.User{Username: "usertoadd", Password: "securepass", Nombre: "ToAdd", Apellido: "User", Cedula: "555555"},
		AdminUsername: "admin",
	}
	req := newAdminRequest("POST", "/api/admin/add-user", "admin", addUserReq)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apphandlers.AdminAddUserHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("Handler devolvió código incorrecto: got %v, want %v. Body: %s", status, http.StatusCreated, rr.Body.String())
	}
}

func TestAdminAddUserHandler_AccessDenied_Gerente(t *testing.T) {
	addUserReq := models.AddUserRequest{
		User:          models.User{Username: "usertofailGerente", Password: "securepass", Nombre: "Fail", Apellido: "Gerente", Cedula: "666666"},
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

func TestAdminUpdateUserRoleHandler_Success(t *testing.T) {
	var testUserID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = 'testuser'").Scan(&testUserID)
	if err != nil {
		t.Fatalf("No se pudo obtener ID de 'testuser': %v", err)
	}

	updateReq := models.UpdateRoleRequest{
		ID:            testUserID,
		NewRole:       "gerente", // Cambiamos a gerente
		AdminUsername: "admin",
	}
	req := newAdminRequest("POST", "/api/admin/update-user", "admin", updateReq)
	rr := httptest.NewRecorder()

	// --- INICIO CORRECCIÓN 1 ---
	// (El handler original se llamaba AdminUpdateUserHandler)
	handler := http.HandlerFunc(apphandlers.AdminUpdateUserHandler)
	// --- FIN CORRECCIÓN 1 ---

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

func TestAdminUpdateUserRoleHandler_AccessDenied_Gerente(t *testing.T) {
	var testUserID int
	database.DB.QueryRow("SELECT id FROM users WHERE username = 'testuser'").Scan(&testUserID)
	updateReq := models.UpdateRoleRequest{ID: testUserID, NewRole: "admin", AdminUsername: "testgerente"} // Gerente intenta cambiar rol
	req := newAdminRequest("POST", "/api/admin/update-user", "testgerente", updateReq)
	rr := httptest.NewRecorder()

	// --- INICIO CORRECCIÓN 2 ---
	// (El handler original se llamaba AdminUpdateUserHandler)
	handler := http.HandlerFunc(apphandlers.AdminUpdateUserHandler)
	// --- FIN CORRECCIÓN 2 ---

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler (gerente) devolvió código incorrecto: got %v, want %v", status, http.StatusForbidden)
	}
}

func TestAdminDeleteUserHandler_Success(t *testing.T) {
	userToDelete := "usertodelete"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	res, _ := database.DB.Exec("INSERT INTO users(username, hashed_password, role, nombre, apellido, cedula) VALUES(?, ?, ?, ?, ?, ?)", userToDelete, string(hashedPass), "user", "ToDelete", "User", "777777")
	idToDelete, _ := res.LastInsertId()

	deleteReq := models.DeleteUserRequest{
		ID:            int(idToDelete),
		AdminUsername: "admin",
	}
	req := newAdminRequest("POST", "/api/admin/delete-user", "admin", deleteReq)
	rr := httptest.NewRecorder()
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
	database.DB.QueryRow("SELECT id FROM users WHERE username = 'testuser'").Scan(&testUserID)
	deleteReq := models.DeleteUserRequest{ID: testUserID, AdminUsername: "testgerente"}
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
	deleteReq := models.DeleteUserRequest{ID: adminID, AdminUsername: "admin"}
	req := newAdminRequest("POST", "/api/admin/delete-user", "admin", deleteReq)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apphandlers.AdminDeleteUserHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("Handler devolvió código incorrecto: got %v, want %v", status, http.StatusForbidden)
	}
}

// --- PRUEBAS DE HANDLERS DE PROYECTOS ---

func TestAdminCreateProyectoHandler_Success(t *testing.T) {

	// --- INICIO CORRECCIÓN 3 ---
	// (Tu 'models.go' original usa una estructura "anidada"
	// y tiene 'FechaCierre' (string) y 'Estado' (string)
	// y NO tiene 'GerenteID')
	proyectoData := models.Proyecto{
		Nombre:      "Proyecto Secreto Alfa",
		FechaInicio: "2025-01-01",
		FechaCierre: "2025-12-31", // Fix: FechaFin -> FechaCierre
		Estado:      "Iniciado",   // Fix: 1 -> "Iniciado"
		// GerenteID no existe en el struct original
	}

	// (El CreateProjectRequest original es anidado)
	createReq := models.CreateProjectRequest{
		Proyecto:      proyectoData,
		AdminUsername: "admin",
	}
	// --- FIN CORRECCIÓN 3 ---

	req := newAdminRequest("POST", "/api/admin/create-proyecto", "admin", createReq)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(apphandlers.AdminCreateProyectoHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Fatalf("Handler devolvió código incorrecto: got %v, want %v. Body: %s",
			status, http.StatusCreated, rr.Body.String())
	}

	var nombre string
	err := database.DB.QueryRow("SELECT nombre FROM proyectos WHERE nombre = ?", proyectoData.Nombre).Scan(&nombre)
	if err != nil {
		t.Fatalf("Error al verificar el proyecto en la DB: %v", err)
	}
	if nombre != proyectoData.Nombre {
		t.Errorf("El nombre del proyecto en la DB es incorrecto: got %s, want %s", nombre, proyectoData.Nombre)
	}
}
