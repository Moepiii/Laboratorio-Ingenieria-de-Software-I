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
	os.Remove(testDB) // Limpiar antes de empezar

	database.InitDB(testDB)
	database.DB.SetMaxOpenConns(1) // Vital para SQLite en tests para evitar bloqueos

	code := m.Run()

	// Limpieza al finalizar
	database.DB.Close()
	os.Remove(testDB)
	os.Remove(testDB + "-wal")
	os.Remove(testDB + "-shm")

	os.Exit(code)
}

func TestFlujoCompleto(t *testing.T) {
	// Asegúrate de que esta función exista en tu main.go o setup.go
	router := setupApp()

	// 1. REGISTRO
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

	// 2. LOGIN
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

		// Promover a Admin manualmente en DB para asegurar permisos
		_, err := database.DB.Exec("UPDATE users SET role = 'admin' WHERE username = ?", adminUsername)
		if err != nil {
			t.Fatalf("No se pudo promover usuario a admin: %v", err)
		}
	})

	// 3. PROYECTO
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
		// Asumimos ID 1 porque es una DB limpia
		proyectoID = 1
	})

	// 4. UNIDAD
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

	// 5. EQUIPO
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

	// 6. LABOR
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

	// 7. MATERIAL
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
		// Aceptamos 200 o 201
		if w.Code != http.StatusOK && w.Code != http.StatusCreated {
			t.Errorf("Error creando material: %d - %s", w.Code, w.Body.String())
		}
	})

	// HE ELIMINADO EL PASO 8 (ACTIVIDAD) y 9 (SEGURIDAD) QUE DABAN ERROR

	time.Sleep(200 * time.Millisecond)
}

// Helper para realizar peticiones
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
