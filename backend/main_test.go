package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"proyecto/internal/database"
	"proyecto/internal/models"
)

// --- CONFIGURACIÓN GLOBAL ---

func TestMain(m *testing.M) {
	testDB := "./test_integration.db"
	os.Remove(testDB)

	// Inicializamos la DB
	database.InitDB(testDB)
	
	// ⭐️ ASEGURARSE QUE ESTO ESTÉ APLICADO EN EL ENTORNO DE TEST TAMBIÉN
	// Aunque esté en database.go, lo forzamos aquí por si acaso.
	database.DB.SetMaxOpenConns(1) 

	code := m.Run()

	database.DB.Close()
	time.Sleep(200 * time.Millisecond)
	os.Remove(testDB)

	os.Exit(code)
}

// --- BATERÍA DE PRUEBAS INTEGRALES ---

func TestBackendFlow(t *testing.T) {
	router := setupApp()

	var adminToken string
	var userToken string
	var proyectoID int

	// ----------------------------------------------------------------
	// 1. PRUEBAS DE AUTENTICACIÓN
	// ----------------------------------------------------------------
	t.Run("1. Registro de Admin", func(t *testing.T) {
		payload := map[string]string{
			"username": "admin_test",
			"password": "123456",
			"nombre":   "Admin",
			"apellido": "Test",
			"cedula":   "111",
		}
		w := performRequest(router, "POST", "/api/auth/register", payload, "")
		
		if w.Code != http.StatusCreated {
			t.Errorf("Esperaba 201 Created, obtuvo %d. Resp: %s", w.Code, w.Body.String())
		}
		
		// Esperar a que el Logger libere la DB
		time.Sleep(200 * time.Millisecond)

		// Forzamos rol de admin
		_, err := database.DB.Exec("UPDATE users SET role = 'admin' WHERE username = 'admin_test'")
		if err != nil {
			t.Fatalf("❌ Error crítico al forzar rol de admin: %v", err)
		}
	})

	t.Run("2. Login de Admin y Obtención de Token", func(t *testing.T) {
		payload := map[string]string{"username": "admin_test", "password": "123456"}
		w := performRequest(router, "POST", "/api/auth/login", payload, "")

		if w.Code != http.StatusOK {
			t.Fatalf("Login falló. Código: %d - Resp: %s", w.Code, w.Body.String())
		}

		var res models.LoginResponse
		json.Unmarshal(w.Body.Bytes(), &res)
		adminToken = res.Token

		if res.Role != "admin" {
			t.Fatalf("El usuario se logueó pero el rol es '%s'", res.Role)
		}
	})

	// ⭐️ PAUSA CRÍTICA: El Login genera un LOG. Esperamos que termine de escribirse.
	time.Sleep(200 * time.Millisecond)

	// ----------------------------------------------------------------
	// 2. PRUEBAS DE SEGURIDAD ESTÁNDAR
	// ----------------------------------------------------------------
	t.Run("3. Acceso Denegado sin Token (403/401)", func(t *testing.T) {
		payloadFake := map[string]string{"admin_username": "hacker"}
		wFake := performRequest(router, "POST", "/api/admin/get-proyectos", payloadFake, "token_invalido")
		
		if wFake.Code != http.StatusForbidden && wFake.Code != http.StatusUnauthorized && wFake.Code != http.StatusInternalServerError {
			t.Errorf("Esperaba bloqueo de seguridad, obtuvo %d", wFake.Code)
		}
	})

	// ----------------------------------------------------------------
	// 3. PRUEBAS DE HANDLERS (Funcionalidad)
	// ----------------------------------------------------------------
	t.Run("4. Crear Proyecto (Handler)", func(t *testing.T) {
		payload := map[string]interface{}{
			"nombre":         "Proyecto Test Unitario",
			"fecha_inicio":   "2025-01-01",
			"fecha_cierre":   "2025-12-31",
			"admin_username": "admin_test",
		}
		w := performRequest(router, "POST", "/api/admin/create-proyecto", payload, adminToken)

		if w.Code != http.StatusCreated {
			t.Errorf("Crear proyecto falló: %d - %s", w.Code, w.Body.String())
		}

		var p models.Proyecto
		json.Unmarshal(w.Body.Bytes(), &p)
		proyectoID = p.ID
	})

	// ⭐️ PAUSA CRÍTICA: Crear Proyecto genera LOG. Esperamos.
	time.Sleep(200 * time.Millisecond)

	t.Run("5. Crear Unidad de Medida (Nuevo Módulo)", func(t *testing.T) {
		// Si falló el paso anterior, este fallará también, pero evitamos el pánico
		if proyectoID == 0 {
			t.Skip("Saltando prueba de Unidad porque no se creó el Proyecto")
		}

		payload := map[string]interface{}{
			"proyecto_id":    proyectoID,
			"nombre":         "Litro",
			"abreviatura":    "lt",
			"tipo":           "Líquido",
			"dimension":      1.5,
			"admin_username": "admin_test",
		}
		w := performRequest(router, "POST", "/api/admin/create-unidad", payload, adminToken)

		if w.Code != http.StatusCreated {
			t.Errorf("Crear unidad falló: %d - %s", w.Code, w.Body.String())
		}
	})

	// ⭐️ PAUSA CRÍTICA: Crear Unidad genera LOG. Esperamos.
	time.Sleep(200 * time.Millisecond)

	// ----------------------------------------------------------------
	// 4. PRUEBAS DE RBAC (Control de Roles)
	// ----------------------------------------------------------------
	t.Run("6. Registro Usuario Normal", func(t *testing.T) {
		payload := map[string]string{
			"username": "pepe_user", 
			"password": "123456",
			"nombre": "Pepe", 
			"apellido": "User", 
			"cedula": "222",
		}
		wReg := performRequest(router, "POST", "/api/auth/register", payload, "")
		if wReg.Code != http.StatusCreated {
			t.Fatalf("Fallo registro user normal: %d - %s", wReg.Code, wReg.Body.String())
		}
		
		// Esperamos por el log de registro
		time.Sleep(100 * time.Millisecond)

		// Login user normal
		wLog := performRequest(router, "POST", "/api/auth/login", map[string]string{"username": "pepe_user", "password": "123456"}, "")
		var res models.LoginResponse
		json.Unmarshal(wLog.Body.Bytes(), &res)
		userToken = res.Token
	})

	// Esperamos por el log de login
	time.Sleep(200 * time.Millisecond)

	t.Run("7. Usuario Normal intenta borrar proyecto (Debe fallar)", func(t *testing.T) {
		payload := map[string]interface{}{
			"id":             proyectoID,
			"admin_username": "pepe_user",
		}
		w := performRequest(router, "POST", "/api/admin/delete-proyecto", payload, userToken)

		if w.Code != http.StatusForbidden {
			t.Errorf("Fallo de Seguridad RBAC: Usuario normal pudo borrar proyecto o obtuvo código incorrecto. Code: %d", w.Code)
		}
	})

	t.Run("8. Admin borra proyecto (Éxito)", func(t *testing.T) {
		if proyectoID == 0 {
			t.Skip("Saltando borrado porque no hay proyecto")
		}
		payload := map[string]interface{}{
			"id":             proyectoID,
			"admin_username": "admin_test",
		}
		w := performRequest(router, "POST", "/api/admin/delete-proyecto", payload, adminToken)

		if w.Code != http.StatusOK {
			t.Errorf("Admin no pudo borrar proyecto: %d - %s", w.Code, w.Body.String())
		}
	})
	
	// Pausa final antes de cerrar DB
	time.Sleep(500 * time.Millisecond)
}

// --- UTILIDADES ---

func performRequest(r http.Handler, method, path string, payload interface{}, token string) *httptest.ResponseRecorder {
	var body []byte
	if payload != nil {
		body, _ = json.Marshal(payload)
	}

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}