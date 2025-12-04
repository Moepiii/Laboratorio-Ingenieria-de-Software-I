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

var (
	authToken     string
	proyectoID    int
	laborID       int
	equipoID      int
	adminUsername = "admin_test"
)

func TestMain(m *testing.M) {
	testDB := "./test_integration.db"
	// Limpieza previa por si quedó basura de una ejecución anterior
	os.Remove(testDB)
	os.Remove(testDB + "-wal")
	os.Remove(testDB + "-shm")

	database.InitDB(testDB)
	database.DB.SetMaxOpenConns(1) // Vital para SQLite en tests

	code := m.Run()

	// Limpieza al finalizar
	database.DB.Close()
	os.Remove(testDB)
	os.Remove(testDB + "-wal")
	os.Remove(testDB + "-shm")

	os.Exit(code)
}

func TestFlujoCompleto(t *testing.T) {
	router := setupApp()

	// 1. REGISTRO (Happy Path)
	t.Run("1. Registrar Admin", func(t *testing.T) {
		payload := map[string]string{
			"username": adminUsername,
			"password": "password123",
			"nombre":   "Admin",
			"apellido": "Test",
			"cedula":   "V-123456",
		}
		w := performRequest(router, "POST", "/api/auth/register", payload, "")
		if w.Code != http.StatusCreated {
			t.Errorf("Falló registro. Código: %d, Resp: %s", w.Code, w.Body.String())
		}
	})

	// 2. LOGIN (Happy Path)
	t.Run("2. Login Admin y Obtener Token", func(t *testing.T) {
		payload := map[string]string{
			"username": adminUsername,
			"password": "password123",
		}
		w := performRequest(router, "POST", "/api/auth/login", payload, "")

		if w.Code != http.StatusOK {
			t.Fatalf("Falló login. Código: %d", w.Code)
		}

		var resp models.LoginResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		authToken = resp.Token

		// TRUCO: Promover a Admin manualmente en DB para tener permisos totales
		_, err := database.DB.Exec("UPDATE users SET role = 'admin' WHERE username = ?", adminUsername)
		if err != nil {
			t.Fatalf("No se pudo promover usuario a admin: %v", err)
		}
	})

	// 3. PROYECTO (Happy Path)
	t.Run("3. Crear Proyecto", func(t *testing.T) {
		payload := map[string]interface{}{
			"nombre":         "Proyecto Maíz 2025",
			"fecha_inicio":   "2025-01-01",
			"fecha_cierre":   "2025-12-31",
			"admin_username": adminUsername,
		}
		w := performRequest(router, "POST", "/api/admin/create-proyecto", payload, authToken)

		if w.Code != http.StatusCreated {
			t.Errorf("Error creando proyecto: %d - %s", w.Code, w.Body.String())
		}
		proyectoID = 1 // Asumimos ID 1 en DB limpia
	})

	// 4. UNIDAD (Happy Path)
	t.Run("4. Crear Unidad de Medida", func(t *testing.T) {
		payload := map[string]interface{}{
			"proyecto_id":    proyectoID,
			"nombre":         "Litros",
			"abreviatura":    "Lts",
			"tipo":           "Volumen",
			"dimension":      1,
			"admin_username": adminUsername,
		}
		w := performRequest(router, "POST", "/api/admin/create-unidad", payload, authToken)
		if w.Code != http.StatusCreated {
			t.Errorf("Error creando unidad: %d - %s", w.Code, w.Body.String())
		}
	})

	// 5. EQUIPO (Happy Path)
	t.Run("5. Crear Equipo/Implemento", func(t *testing.T) {
		payload := map[string]interface{}{
			"proyecto_id":    proyectoID,
			"codigo_equipo":  "TR-01",
			"nombre":         "Tractor John Deere",
			"tipo":           "Equipo",
			"estado":         "Operativo",
			"admin_username": adminUsername,
		}
		w := performRequest(router, "POST", "/api/admin/create-equipo", payload, authToken)
		if w.Code != http.StatusCreated {
			t.Errorf("Error creando equipo: %d - %s", w.Code, w.Body.String())
		}
		equipoID = 1
	})

	// 6. LABOR (Happy Path)
	t.Run("6. Crear Labor Agronómica", func(t *testing.T) {
		payload := map[string]interface{}{
			"proyecto_id":    proyectoID,
			"codigo_labor":   "L-01",
			"descripcion":    "Riego por Goteo",
			"admin_username": adminUsername,
		}
		w := performRequest(router, "POST", "/api/admin/create-labor", payload, authToken)
		if w.Code != http.StatusCreated {
			t.Errorf("Error creando labor: %d - %s", w.Code, w.Body.String())
		}
		laborID = 1
	})

	// 7. MATERIAL (Happy Path)
	t.Run("7. Crear Material/Insumo", func(t *testing.T) {
		payload := map[string]interface{}{
			"proyecto_id":    proyectoID,
			"actividad":      "Fertilización",
			"accion":         "Aplicar",
			"categoria":      "Fertilizante",
			"responsable":    "Juan Perez",
			"nombre":         "Urea",
			"unidad":         "Sacos",
			"cantidad":       50,
			"costo_unitario": 25.5,
			"monto":          1275.0,
			"admin_username": adminUsername,
		}
		w := performRequest(router, "POST", "/api/admin/create-material", payload, authToken)
		if w.Code != http.StatusOK && w.Code != http.StatusCreated {
			t.Errorf("Error creando material: %d - %s", w.Code, w.Body.String())
		}
	})

	// 9. SEGURIDAD NEGATIVA (Nuevo Test Agregado)
	t.Run("9. Intento de borrado sin permisos", func(t *testing.T) {
		// A. Registrar un usuario normal (el intruso)
		intruderName := "pepe_intruso"
		regPayload := map[string]string{
			"username": intruderName,
			"password": "password123",
			"nombre":   "Pepe",
			"apellido": "Intruso",
			"cedula":   "V-999999",
		}
		performRequest(router, "POST", "/api/auth/register", regPayload, "")

		// B. Login del intruso para obtener SU token
		loginPayload := map[string]string{
			"username": intruderName,
			"password": "password123",
		}
		wLogin := performRequest(router, "POST", "/api/auth/login", loginPayload, "")

		var resp models.LoginResponse
		json.Unmarshal(wLogin.Body.Bytes(), &resp)
		tokenIntruso := resp.Token

		// C. Intentar borrar el Proyecto (Acción reservada para Admins)
		// NOTA: No le damos update a 'admin' en la DB, así que es un simple mortal.
		delPayload := map[string]interface{}{
			"id":             proyectoID,
			"admin_username": intruderName,
		}

		// Usamos el token del intruso
		w := performRequest(router, "POST", "/api/admin/delete-proyecto", delPayload, tokenIntruso)

		// D. Verificar que el sistema lo rechace
		// Esperamos 403 Forbidden o 401 Unauthorized
		if w.Code != http.StatusForbidden && w.Code != http.StatusUnauthorized {
			t.Errorf("ALERTA DE SEGURIDAD: Usuario sin permisos pudo acceder a ruta protegida. Código: %d", w.Code)
		} else {
			t.Log("✅ Correcto: El sistema bloqueó el acceso no autorizado.")
		}
	})

	time.Sleep(200 * time.Millisecond)
}

// Helper para realizar peticiones HTTP en el test
func performRequest(r http.Handler, method, path string, payload interface{}, token string) *httptest.ResponseRecorder {
	var reqBody []byte
	if payload != nil {
		reqBody, _ = json.Marshal(payload)
	}

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
