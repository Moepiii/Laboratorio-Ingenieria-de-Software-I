package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"proyecto/internal/database"
	"proyecto/internal/models"
	"proyecto/internal/proyectos" // ⭐️ Solo importamos esto para la prueba de borrado
)

// --- 1. El "Setup" de la Prueba ---
func setupTestDB(t *testing.T) {
	database.InitDB("file::memory:")
	t.Cleanup(func() {
		err := database.DB.Close()
		if err != nil {
			log.Printf("Error cerrando la DB de prueba: %v", err)
		}
	})
}

// --- BATERÍA ÚNICA DE PRUEBAS DE HANDLERS HTTP ---
// Al ser la única función Test... no hay riesgo de paralelismo.
func TestHandlers_HTTP_Integration(t *testing.T) {

	// ARRANGE (Global para todos los handlers HTTP)
	setupTestDB(t)       // 1. Prepara la DB en memoria UNA VEZ
	server := setupApp() // 2. "Arma" la app UNA VEZ
	var adminToken string
	var proyectoCreado models.Proyecto // Para compartir entre creación y borrado

	// --- 1. PRUEBA DE REGISTRO PÚBLICO ---
	t.Run("POST /api/auth/register - Éxito 201", func(t *testing.T) {
		newUser := models.User{
			Username: "http_user", Password: "password", Nombre: "HTTP", Apellido: "Test", Cedula: "111",
		}
		body, _ := json.Marshal(newUser)
		req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Fatalf("Código esperado %d, se obtuvo %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
		}
	})

	// --- 2. PRUEBA DE LOGIN ---
	t.Run("POST /api/auth/login - Éxito 200 y devuelve token", func(t *testing.T) {
		loginCreds := map[string]string{"username": "admin", "password": "admin123"}
		body, _ := json.Marshal(loginCreds)
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Código esperado %d, se obtuvo %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}
		var loginResp models.LoginResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &loginResp)
		adminToken = loginResp.Token // ⭐️ GUARDA EL TOKEN
	})

	// --- 3. PRUEBA DE PERMISO DENEGADO ---
	t.Run("POST /api/admin/get-proyectos - Falla 403 (Permiso Denegado)", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"admin_username": "usuario_no_existente"})
		req := httptest.NewRequest("POST", "/api/admin/get-proyectos", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("Código esperado %d, se obtuvo %d. Body: %s", http.StatusForbidden, rr.Code, rr.Body.String())
		}
	})

	// --- 4. PRUEBA DE LECTURA (GET) ---
	t.Run("POST /api/admin/get-proyectos - Éxito 200 (Con Token)", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"admin_username": "admin"})
		req := httptest.NewRequest("POST", "/api/admin/get-proyectos", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("Código esperado %d, se obtuvo %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}
	})

	// --- 5. PRUEBA DE CREACIÓN ---
	t.Run("POST /api/admin/create-proyecto - Éxito 201 (Con Token)", func(t *testing.T) {
		nombreProyecto := "Proyecto Creado por HTTP"
		newProject := models.CreateProyectoRequest{
			Nombre: nombreProyecto, FechaInicio: "2025-01-01", FechaCierre: "2025-12-31", AdminUsername: "admin",
		}
		body, _ := json.Marshal(newProject)
		req := httptest.NewRequest("POST", "/api/admin/create-proyecto", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rr := httptest.NewRecorder()
		server.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("Código esperado %d, se obtuvo %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
		}

		// Guarda el proyecto creado para la prueba de borrado
		if err := json.Unmarshal(rr.Body.Bytes(), &proyectoCreado); err != nil {
			t.Fatal("No se pudo decodificar la respuesta JSON del proyecto creado")
		}

		// Verifica la DB
		var dbNombre string
		err := database.DB.QueryRow("SELECT nombre FROM proyectos WHERE id = ?", proyectoCreado.ID).Scan(&dbNombre)
		if err != nil {
			t.Fatalf("Error al consultar la DB: %v", err)
		}
		if dbNombre != nombreProyecto {
			t.Fatalf("El proyecto no se guardó en la DB")
		}
	})

	// --- 6. PRUEBA DE ELIMINACIÓN ---
	t.Run("POST /api/admin/delete-proyecto - Éxito 200 (Con Token)", func(t *testing.T) {
		// ARRANGE: Usa el proyecto de la prueba anterior
		if proyectoCreado.ID == 0 {
			// Solución de respaldo si la prueba 5 falló, para evitar el pánico
			proyectoSvc := proyectos.NewProyectoService()
			p, _ := proyectoSvc.CreateProyecto("Proyecto Temp para Borrar", "2025-01-01", "2025-12-31")
			if p == nil {
				t.Fatal("Fallo crítico: no se pudo crear el proyecto de respaldo para la prueba de borrado")
			}
			proyectoCreado.ID = p.ID
		}

		deleteReq := models.DeleteProyectoRequest{ID: proyectoCreado.ID, AdminUsername: "admin"}
		body, _ := json.Marshal(deleteReq)
		req := httptest.NewRequest("POST", "/api/admin/delete-proyecto", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rr := httptest.NewRecorder()

		// ACT:
		server.ServeHTTP(rr, req)

		// ASSERT (HTTP):
		if rr.Code != http.StatusOK {
			t.Fatalf("Código esperado %d, se obtuvo %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}
		// ASSERT (DB):
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM proyectos WHERE id = ?", proyectoCreado.ID).Scan(&count)
		if err != nil {
			t.Fatalf("Error al consultar la DB: %v", err)
		}
		if count != 0 {
			t.Fatal("El proyecto no fue borrado de la base de datos")
		}
	})
}
