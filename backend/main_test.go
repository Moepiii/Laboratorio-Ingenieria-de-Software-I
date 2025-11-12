package main

import (
	"log"
	"testing" // La biblioteca de pruebas de Go

	// Importamos todos los servicios que vamos a probar
	"proyecto/internal/actividades"
	"proyecto/internal/auth"
	"proyecto/internal/database"
	"proyecto/internal/equipos"
	"proyecto/internal/labores"
	"proyecto/internal/models"
	"proyecto/internal/proyectos"
	"proyecto/internal/users"
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

// --- 2. Batería de Pruebas de Autenticación (Auth) ---
func TestAuthService_Integration(t *testing.T) {
	// ARRANGE (Global): Creamos el servicio
	s := auth.NewAuthService()

	t.Run("Usuario puede registrarse exitosamente", func(t *testing.T) {
		// ARRANGE (Local): DB limpia
		setupTestDB(t)

		// ⭐️ --- CORRECCIÓN 1 (Línea 67) --- ⭐️
		user := models.User{
			Username: "testuser",
			Password: "password123",
			Nombre:   "Test",
			Apellido: "User",
			Cedula:   "123456",
		}

		// ACT: Ejecutamos la función
		id, err := s.Register(user)

		// ASSERT: Verificamos
		if err != nil {
			t.Fatalf("Register() falló: %v", err)
		}
		if id != 2 { // 1 es el admin, 2 es el nuevo usuario
			t.Fatalf("Se esperaba ID=2, pero se obtuvo ID=%d", id)
		}
	})

	t.Run("Usuario puede iniciar sesión exitosamente", func(t *testing.T) {
		// ARRANGE: DB limpia y un usuario registrado
		setupTestDB(t)

		// ⭐️ --- CORRECCIÓN 2 (Línea 87) --- ⭐️
		s.Register(models.User{
			Username: "testuser",
			Password: "password123",
			Nombre:   "Test",
			Apellido: "User",
			Cedula:   "123456",
		})

		// ACT: Intentamos iniciar sesión
		resp, err := s.Login("testuser", "password123")

		// ASSERT:
		if err != nil {
			t.Fatalf("Login() falló: %v", err)
		}
		if resp.Token == "" {
			t.Fatal("Login() no devolvió un token")
		}
		if resp.User.Username != "testuser" {
			t.Fatal("Login() devolvió un usuario incorrecto")
		}
	})

	t.Run("Usuario no puede iniciar sesión con contraseña incorrecta", func(t *testing.T) {
		// ARRANGE: DB limpia y un usuario registrado
		setupTestDB(t)
		// (No es necesario corregir este, pero lo hacemos por consistencia)
		s.Register(models.User{
			Username: "testuser",
			Password: "password123",
			Nombre:   "Test",
			Apellido: "User",
			Cedula:   "123456",
		})

		// ACT: Intentamos iniciar sesión con pass incorrecta
		resp, err := s.Login("testuser", "MALAPASSWORD")

		// ASSERT:
		if err == nil {
			t.Fatal("Se esperaba un error por credenciales inválidas, pero no falló")
		}
		if resp != nil {
			t.Fatal("No se debía devolver una respuesta en caso de error de login")
		}
	})
}

// --- 3. Batería de Pruebas de Proyectos ---
func TestProyectoService_Integration(t *testing.T) {
	// ARRANGE (Global):
	s := proyectos.NewProyectoService()

	t.Run("Admin puede crear un proyecto", func(t *testing.T) {
		// ARRANGE (Local):
		setupTestDB(t)

		// ACT:
		nombre := "Proyecto de Prueba"
		proyecto, err := s.CreateProyecto(nombre, "2025-01-01", "2025-12-31")

		// ASSERT:
		if err != nil {
			t.Fatalf("CreateProyecto falló: %v", err)
		}
		if proyecto.ID != 1 { // 1 es el primer proyecto
			t.Fatalf("Se esperaba ID=1, pero se obtuvo ID=%d", proyecto.ID)
		}
		if proyecto.Nombre != nombre {
			t.Fatal("El nombre del proyecto no coincide")
		}
	})

	t.Run("Admin puede borrar un proyecto", func(t *testing.T) {
		// ARRANGE:
		setupTestDB(t)
		p, _ := s.CreateProyecto("Proyecto a Borrar", "2025-01-01", "2025-12-31")

		// ACT:
		affected, err := s.DeleteProyecto(p.ID)

		// ASSERT:
		if err != nil {
			t.Fatalf("DeleteProyecto falló: %v", err)
		}
		if affected == 0 {
			t.Fatal("El servicio reportó que 0 filas fueron afectadas")
		}
	})
}

// --- 4. Batería de Pruebas de Usuarios (Admin) ---
func TestUserService_Integration(t *testing.T) {
	// ARRANGE (Global):
	s := users.NewUserService()

	t.Run("Admin puede añadir un usuario (encargado)", func(t *testing.T) {
		// ARRANGE:
		setupTestDB(t)
		// (No es necesario corregir este, pero lo hacemos por consistencia)
		user := models.User{
			Username: "encargado1",
			Password: "password123",
			Nombre:   "Encargado",
			Apellido: "Prueba",
			Cedula:   "78910",
		}

		// ACT:
		id, err := s.AddUser(user)

		// ASSERT:
		if err != nil {
			t.Fatalf("AddUser falló: %v", err)
		}
		if id != 2 { // 1 es el admin
			t.Fatalf("Se esperaba ID=2, pero se obtuvo ID=%d", id)
		}

		// Verificamos que el rol sea 'encargado' (como dice la lógica del servicio)
		var role string
		err = database.DB.QueryRow("SELECT role FROM users WHERE id = ?", id).Scan(&role)
		if err != nil {
			t.Fatalf("Error al consultar la DB: %v", err)
		}
		if role != "encargado" {
			t.Fatalf("Se esperaba rol 'encargado', pero se obtuvo '%s'", role)
		}
	})
}

// --- 5. Batería de Pruebas de Flujo Completo (Labores, Equipos, Actividades) ---
func TestFullFlow_Integration(t *testing.T) {
	// ARRANGE: Preparamos todos los servicios
	setupTestDB(t)
	proyectoSvc := proyectos.NewProyectoService()
	laborSvc := labores.NewLaborService()
	equipoSvc := equipos.NewEquipoService()
	actividadSvc := actividades.NewActividadService()
	userSvc := users.NewUserService() // Para crear un encargado

	// 1. Crear Proyecto
	proyecto, _ := proyectoSvc.CreateProyecto("Proyecto Full Flow", "2025-01-01", "2025-12-31")

	// 2. Crear Encargado (con userSvc)
	// ⭐️ --- CORRECCIÓN 3 (Línea 198) --- ⭐️
	encargadoID, _ := userSvc.AddUser(models.User{
		Username: "encargado_flow",
		Password: "pass",
		Nombre:   "Enc",
		Apellido: "Flow",
		Cedula:   "999",
	})
	encargadoIDInt := int(encargadoID)

	// 3. Crear Labor
	laborReq := models.CreateLaborRequest{ProyectoID: proyecto.ID, CodigoLabor: "L-001", Descripcion: "Arado"}
	labor, _ := laborSvc.CreateLabor(laborReq)

	// 4. Crear Equipo
	equipoReq := models.CreateEquipoRequest{ProyectoID: proyecto.ID, CodigoEquipo: "E-001", Nombre: "Tractor", Tipo: "Equipo"}
	equipo, _ := equipoSvc.CreateEquipo(equipoReq)

	// 5. Crear Actividad (La prueba principal)
	t.Run("Servicio de Actividad puede crear una actividad", func(t *testing.T) {

		actividadReq := models.CreateActividadRequest{
			ProyectoID:         proyecto.ID,
			Actividad:          "Preparación de Suelo",
			LaborAgronomicaID:  &labor.ID,
			EquipoImplementoID: &equipo.ID,
			EncargadoID:        &encargadoIDInt,
			RecursoHumano:      5,
			Costo:              1500.75,
			Observaciones:      "Todo ok",
		}

		// ACT:
		listaActividades, err := actividadSvc.CreateActividad(actividadReq)

		// ASSERT:
		if err != nil {
			t.Fatalf("CreateActividad falló: %v", err)
		}
		if len(listaActividades) != 1 {
			t.Fatal("La lista de actividades no tiene 1 elemento")
		}
		if listaActividades[0].Actividad != "Preparación de Suelo" {
			t.Fatal("El nombre de la actividad no coincide")
		}
		if listaActividades[0].EncargadoNombre.String != "Enc Flow" {
			t.Fatalf("El nombre del encargado no coincide, se obtuvo: %s", listaActividades[0].EncargadoNombre.String)
		}
	})
}
