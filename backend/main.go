package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	"proyecto/internal/actividades"
	"proyecto/internal/auth"
	"proyecto/internal/database"
	"proyecto/internal/equipos"
	apphandlers "proyecto/internal/handlers" // Usamos alias para evitar conflictos
	"proyecto/internal/labores"
	"proyecto/internal/logger"
	"proyecto/internal/proyectos"
	"proyecto/internal/unidades"
	"proyecto/internal/users"
)

func setupApp() http.Handler {
	// 1. DEFINIR EL ROUTER (Mux)
	mux := http.NewServeMux()

	// 2. INICIALIZAR TODOS LOS SERVICIOS
	authService := auth.NewAuthService()
	loggerService := logger.NewLoggerService()

	userService := users.NewUserService()
	proyectoService := proyectos.NewProyectoService()
	laborService := labores.NewLaborService()
	equipoService := equipos.NewEquipoService()
	actividadService := actividades.NewActividadService()
	unidadService := unidades.NewUnidadService()
	// Nota: No creamos un planService separado porque usamos la lógica directa en el handler por ahora,
	// pero si el proyecto crece, deberías crear "proyecto/internal/planes".

	// 3. INICIALIZAR HANDLERS (Controladores)
	// Inyectamos los servicios necesarios en cada Handler
	authHandler := apphandlers.NewAuthHandler(authService, loggerService)
	userHandler := apphandlers.NewUserHandler(authService, userService, loggerService)
	proyectoHandler := apphandlers.NewProyectoHandler(authService, proyectoService, loggerService)
	laborHandler := apphandlers.NewLaborHandler(authService, laborService, loggerService)
	equipoHandler := apphandlers.NewEquipoHandler(authService, equipoService, loggerService)
	unidadHandler := apphandlers.NewUnidadHandler(authService, unidadService, loggerService)
	actividadHandler := apphandlers.NewActividadHandler(authService, actividadService, loggerService)
	loggerHandler := apphandlers.NewLoggerHandler(authService, loggerService)

	// ⭐️ NUEVO: Handler de Planes de Acción
	planHandler := apphandlers.NewPlanHandler(authService, loggerService)

	// 4. REGISTRAR RUTAS
	mux.HandleFunc("/", apphandlers.SaludoHandler)

	// -- Rutas de Autenticación --
	mux.HandleFunc("/api/auth/register", authHandler.RegisterHandler)
	mux.HandleFunc("/api/auth/login", authHandler.LoginHandler)

	// -- Rutas de Usuarios --
	mux.HandleFunc("/api/admin/users", userHandler.AdminUsersHandler)
	mux.HandleFunc("/api/admin/add-user", userHandler.AdminAddUserHandler)
	mux.HandleFunc("/api/admin/delete-user", userHandler.AdminDeleteUserHandler)
	mux.HandleFunc("/api/admin/update-user", userHandler.AdminUpdateUserRoleHandler)
	mux.HandleFunc("/api/admin/assign-project", userHandler.AdminAssignProjectToUserHandler)
	mux.HandleFunc("/api/user/project-details", userHandler.UserProjectDetailsHandler)

	// -- Rutas de Proyectos --
	mux.HandleFunc("/api/admin/get-proyectos", proyectoHandler.GetProyectosHandler)
	mux.HandleFunc("/api/admin/create-proyecto", proyectoHandler.CreateProyectoHandler)
	mux.HandleFunc("/api/admin/update-proyecto", proyectoHandler.UpdateProyectoHandler)
	mux.HandleFunc("/api/admin/delete-proyecto", proyectoHandler.DeleteProyectoHandler)
	mux.HandleFunc("/api/admin/set-proyecto-estado", proyectoHandler.AdminSetProyectoEstadoHandler)

	// -- Rutas de Labores Agronómicas --
	mux.HandleFunc("/api/admin/get-labores", laborHandler.GetLaboresHandler)
	mux.HandleFunc("/api/admin/create-labor", laborHandler.CreateLaborHandler)
	mux.HandleFunc("/api/admin/update-labor", laborHandler.UpdateLaborHandler)
	mux.HandleFunc("/api/admin/delete-labor", laborHandler.DeleteLaborHandler)

	// -- Rutas de Equipos e Implementos --
	mux.HandleFunc("/api/admin/get-equipos", equipoHandler.GetEquiposHandler)
	mux.HandleFunc("/api/admin/create-equipo", equipoHandler.CreateEquipoHandler)
	mux.HandleFunc("/api/admin/update-equipo", equipoHandler.UpdateEquipoHandler)
	mux.HandleFunc("/api/admin/delete-equipo", equipoHandler.DeleteEquipoHandler)

	// -- Rutas de Unidades de Medida --
	mux.HandleFunc("/api/admin/get-unidades", unidadHandler.GetUnidadesHandler)
	mux.HandleFunc("/api/admin/create-unidad", unidadHandler.CreateUnidadHandler)
	mux.HandleFunc("/api/admin/update-unidad", unidadHandler.UpdateUnidadHandler)
	mux.HandleFunc("/api/admin/delete-unidad", unidadHandler.DeleteUnidadHandler)

	// -- Rutas de Actividades (Datos del Proyecto) --
	mux.HandleFunc("/api/admin/get-datos-proyecto", actividadHandler.GetDatosProyectoHandler)
	mux.HandleFunc("/api/admin/create-actividad", actividadHandler.CreateActividadHandler)
	mux.HandleFunc("/api/admin/update-actividad", actividadHandler.UpdateActividadHandler)
	mux.HandleFunc("/api/admin/delete-actividad", actividadHandler.DeleteActividadHandler)

	// -- Rutas de Planes de Acción (⭐️ NUEVO) --
	mux.HandleFunc("/api/admin/create-plan", planHandler.CreatePlanHandler)
	mux.HandleFunc("/api/admin/get-planes", planHandler.GetPlanesHandler)

	// -- Ruta Logger (Auditoría) --
	mux.HandleFunc("/api/admin/get-logs", loggerHandler.GetLogsHandler)
	mux.HandleFunc("/api/admin/delete-logs", loggerHandler.DeleteLogsHandler)
	mux.HandleFunc("/api/admin/delete-logs-range", loggerHandler.DeleteLogsRangeHandler)

	// 5. CONFIGURAR MIDDLEWARE CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	return corsHandler(mux)
}

func main() {
	// 1. INICIALIZAR LA BASE DE DATOS
	database.InitDB("./users.db")

	// 2. INICIAR SERVIDOR
	server := setupApp()
	log.Println("🚀 Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
