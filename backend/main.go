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
	"proyecto/internal/unidades"
	"proyecto/internal/users"
)

func setupApp() http.Handler {
	// 1. DEFINIR EL ROUTER (Mux) PRIMERO
	mux := http.NewServeMux()

	// 2. INICIALIZAR TODOS LOS SERVICIOS
	// El orden importa: creamos los servicios antes de pasarlos a los handlers.

	authService := auth.NewAuthService()
	loggerService := logger.NewLoggerService()

	userService := users.NewUserService()
	proyectoService := proyectos.NewProyectoService()
	laborService := labores.NewLaborService()
	equipoService := equipos.NewEquipoService()
	actividadService := actividades.NewActividadService()
	unidadService := unidades.NewUnidadService()

	// 3. INICIALIZAR HANDLERS (Controladores)
	// Inyectamos los servicios necesarios en cada Handler

	// AuthHandler necesita AuthService y LoggerService
	authHandler := apphandlers.NewAuthHandler(authService, loggerService)

	// UserHandler necesita AuthService, UserService y LoggerService
	userHandler := apphandlers.NewUserHandler(authService, userService, loggerService)

	// ProyectoHandler necesita AuthService, ProyectoService y LoggerService
	proyectoHandler := apphandlers.NewProyectoHandler(authService, proyectoService, loggerService)

	// LaborHandler necesita Auth, LaborService y Logger
	laborHandler := apphandlers.NewLaborHandler(authService, laborService, loggerService)

	// EquipoHandler necesita Auth, EquipoService y Logger
	equipoHandler := apphandlers.NewEquipoHandler(authService, equipoService, loggerService)

	// ActividadHandler necesita Auth, ActividadService y Logger
	actividadHandler := apphandlers.NewActividadHandler(authService, actividadService, loggerService)

	// UnidadHandler necesita Auth, UnidadService y Logger
	unidadHandler := apphandlers.NewUnidadHandler(authService, unidadService, loggerService)

	// LoggerHandler necesita Auth y LoggerService
	loggerHandler := apphandlers.NewLoggerHandler(authService, loggerService)

	// 4. DEFINIR RUTAS (Endpoints)

	// -- Rutas Públicas --
	mux.HandleFunc("/api/saludo", apphandlers.SaludoHandler)
	mux.HandleFunc("/api/auth/register", authHandler.RegisterHandler)
	mux.HandleFunc("/api/auth/login", authHandler.LoginHandler)

	// -- Rutas Usuarios --
	mux.HandleFunc("/api/admin/users", userHandler.AdminUsersHandler)
	mux.HandleFunc("/api/admin/add-user", userHandler.AdminAddUserHandler)
	mux.HandleFunc("/api/admin/delete-user", userHandler.AdminDeleteUserHandler)
	mux.HandleFunc("/api/admin/update-user", userHandler.AdminUpdateUserRoleHandler)         // Cambiar Rol
	mux.HandleFunc("/api/admin/assign-project", userHandler.AdminAssignProjectToUserHandler) // Asignar Proyecto
	mux.HandleFunc("/api/user/project-details", userHandler.UserProjectDetailsHandler)       // Dashboard Usuario

	// -- Rutas Proyectos --
	mux.HandleFunc("/api/admin/get-proyectos", proyectoHandler.AdminGetProyectosHandler)
	mux.HandleFunc("/api/admin/create-proyecto", proyectoHandler.AdminCreateProyectoHandler)
	mux.HandleFunc("/api/admin/update-proyecto", proyectoHandler.AdminUpdateProyectoHandler)
	mux.HandleFunc("/api/admin/delete-proyecto", proyectoHandler.AdminDeleteProyectoHandler)
	mux.HandleFunc("/api/admin/set-proyecto-estado", proyectoHandler.AdminSetProyectoEstadoHandler)

	// -- Rutas Labores Agronómicas --
	mux.HandleFunc("/api/admin/get-labores", laborHandler.GetLaboresHandler)
	mux.HandleFunc("/api/admin/create-labor", laborHandler.CreateLaborHandler)
	mux.HandleFunc("/api/admin/update-labor", laborHandler.UpdateLaborHandler)
	mux.HandleFunc("/api/admin/delete-labor", laborHandler.DeleteLaborHandler)

	// -- Rutas Equipos e Implementos --
	mux.HandleFunc("/api/admin/get-equipos", equipoHandler.GetEquiposHandler)
	mux.HandleFunc("/api/admin/create-equipo", equipoHandler.CreateEquipoHandler)
	mux.HandleFunc("/api/admin/update-equipo", equipoHandler.UpdateEquipoHandler)
	mux.HandleFunc("/api/admin/delete-equipo", equipoHandler.DeleteEquipoHandler)

	// -- Rutas Unidades de Medida --
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
	mux.HandleFunc("/api/admin/delete-logs", loggerHandler.DeleteLogsHandler)
	// ⭐️ NUEVA RUTA AGREGADA PARA BORRADO POR RANGO ⭐️
	mux.HandleFunc("/api/admin/delete-logs-range", loggerHandler.DeleteLogsRangeHandler)

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
	// Crea las tablas si no existen
	database.InitDB("./users.db")

	// Aseguramos que se cierre al terminar
	defer database.DB.Close()

	// 2. CONFIGURAR EL SERVIDOR
	router := setupApp()

	log.Println("🚀 Servidor corriendo en http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
