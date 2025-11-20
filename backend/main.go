package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	"proyecto/internal/actividades"
	"proyecto/internal/auth"
	"proyecto/internal/database"
	"proyecto/internal/equipos"
	apphandlers "proyecto/internal/handlers" // Alias para evitar conflicto de nombres
	"proyecto/internal/labores"
	"proyecto/internal/logger"
	"proyecto/internal/proyectos"
	"proyecto/internal/unidades" // ⭐️ IMPORTANTE: El nuevo módulo
	"proyecto/internal/users"
)

func setupApp() http.Handler {
	// 1. DEFINIR EL ROUTER (Mux) PRIMERO
	// Esto soluciona los errores "undefined: mux"
	mux := http.NewServeMux()

	// 2. INICIALIZAR TODOS LOS SERVICIOS
	// El orden importa: creamos los servicios antes de pasarlos a los handlers.
	
	authService := auth.NewAuthService()
	loggerService := logger.NewLoggerService() // Soluciona "undefined: loggerService"
	
	userService := users.NewUserService()
	proyectoService := proyectos.NewProyectoService()
	laborService := labores.NewLaborService()
	equipoService := equipos.NewEquipoService()
	actividadService := actividades.NewActividadService()
	unidadService := unidades.NewUnidadService() // ⭐️ Nuevo servicio

	// 3. INICIALIZAR LOS HANDLERS (Inyección de Dependencias)
	// Aquí conectamos los servicios y el logger a los controladores.

	authHandler := apphandlers.NewAuthHandler(authService, loggerService)
	userHandler := apphandlers.NewUserHandler(authService, userService, loggerService)
	proyectoHandler := apphandlers.NewProyectoHandler(authService, proyectoService, loggerService)
	laborHandler := apphandlers.NewLaborHandler(authService, laborService, loggerService)
	equipoHandler := apphandlers.NewEquipoHandler(authService, equipoService, loggerService)
	actividadHandler := apphandlers.NewActividadHandler(authService, actividadService, loggerService)
	loggerHandler := apphandlers.NewLoggerHandler(authService, loggerService)
	unidadHandler := apphandlers.NewUnidadHandler(authService, unidadService, loggerService) // ⭐️ Nuevo handler

	// 4. REGISTRAR LAS RUTAS (Endpoints)

	// -- Rutas Públicas --
	mux.HandleFunc("/", apphandlers.SaludoHandler)
	mux.HandleFunc("/api/auth/register", authHandler.RegisterHandler)
	mux.HandleFunc("/api/auth/login", authHandler.LoginHandler)

	// -- Rutas Admin: Gestión de Proyectos --
	mux.HandleFunc("/api/admin/get-proyectos", proyectoHandler.AdminGetProyectosHandler)
	mux.HandleFunc("/api/admin/create-proyecto", proyectoHandler.CreateProyectoHandler)
	mux.HandleFunc("/api/admin/update-proyecto", proyectoHandler.UpdateProyectoHandler)
	mux.HandleFunc("/api/admin/delete-proyecto", proyectoHandler.DeleteProyectoHandler)
	mux.HandleFunc("/api/admin/set-proyecto-estado", proyectoHandler.AdminSetProyectoEstadoHandler)

	// -- Rutas Admin: Gestión de Usuarios --
	mux.HandleFunc("/api/admin/users", userHandler.AdminUsersHandler)
	mux.HandleFunc("/api/admin/add-user", userHandler.AdminAddUserHandler)
	mux.HandleFunc("/api/admin/delete-user", userHandler.AdminDeleteUserHandler)
	mux.HandleFunc("/api/admin/update-user", userHandler.AdminUpdateUserRoleHandler)
	mux.HandleFunc("/api/admin/assign-project", userHandler.AdminAssignProjectToUserHandler)

	// -- Rutas Usuario Normal --
	mux.HandleFunc("/api/user/project-details", userHandler.UserProjectDetailsHandler)

	// -- Rutas Configuración: Labores Agronómicas --
	mux.HandleFunc("/api/admin/get-labores", laborHandler.GetLaboresHandler)
	mux.HandleFunc("/api/admin/create-labor", laborHandler.CreateLaborHandler)
	mux.HandleFunc("/api/admin/update-labor", laborHandler.UpdateLaborHandler)
	mux.HandleFunc("/api/admin/delete-labor", laborHandler.DeleteLaborHandler)

	// -- Rutas Configuración: Equipos e Implementos --
	mux.HandleFunc("/api/admin/get-equipos", equipoHandler.GetEquiposHandler)
	mux.HandleFunc("/api/admin/create-equipo", equipoHandler.CreateEquipoHandler)
	mux.HandleFunc("/api/admin/update-equipo", equipoHandler.UpdateEquipoHandler)
	mux.HandleFunc("/api/admin/delete-equipo", equipoHandler.DeleteEquipoHandler)

	// -- Rutas Configuración: Unidades de Medida (⭐️ NUEVO) --
	mux.HandleFunc("/api/admin/get-unidades", unidadHandler.GetUnidadesHandler)
	mux.HandleFunc("/api/admin/create-unidad", unidadHandler.CreateUnidadHandler)
	mux.HandleFunc("/api/admin/update-unidad", unidadHandler.UpdateUnidadHandler)
	mux.HandleFunc("/api/admin/delete-unidad", unidadHandler.DeleteUnidadHandler)

	// -- Rutas Actividades (Datos del Proyecto) --
	mux.HandleFunc("/api/admin/get-datos-proyecto", actividadHandler.GetDatosProyectoHandler)
	mux.HandleFunc("/api/admin/create-actividad", actividadHandler.CreateActividadHandler)
	mux.HandleFunc("/api/admin/update-actividad", actividadHandler.UpdateActividadHandler)
	mux.HandleFunc("/api/admin/delete-actividad", actividadHandler.DeleteActividadHandler)

	// -- Ruta Logger (Auditoría) --
	mux.HandleFunc("/api/admin/get-logs", loggerHandler.GetLogsHandler)

	// 5. CONFIGURAR MIDDLEWARE CORS
	// Permite que el Frontend (puerto 3000) hable con este Backend (puerto 8080)
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	return corsHandler(mux)
}

func main() {
	// 1. INICIALIZAR LA BASE DE DATOS
	// Crea las tablas si no existen (incluyendo la nueva 'unidades_medida')
	database.InitDB("./users.db")

	// 2. CONFIGURAR LA APLICACIÓN
	handler := setupApp()

	// 3. INICIAR EL SERVIDOR
	log.Println("✅ Servidor corriendo en http://localhost:8080")
	// ListenAndServe bloquea la ejecución, así que va al final
	log.Fatal(http.ListenAndServe(":8080", handler))
}