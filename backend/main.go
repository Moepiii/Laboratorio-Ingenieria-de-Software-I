package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"

	"proyecto/internal/actividades"
	"proyecto/internal/auth"
	"proyecto/internal/database"
	"proyecto/internal/equipos"
	apphandlers "proyecto/internal/handlers"
	"proyecto/internal/labores"
	"proyecto/internal/logger"
	"proyecto/internal/proyectos"
	"proyecto/internal/users"
)

// ⭐️ setupApp() es la lógica que extrajimos de main().
// Ahora es reutilizable y podemos llamarla desde main_test.go
func setupApp() http.Handler {
	// 1. "ARMAR" LA APLICACIÓN (Inyección de Dependencias)

	// Creamos todos los servicios
	authService := auth.NewAuthService()
	userService := users.NewUserService()
	proyectoService := proyectos.NewProyectoService()
	laborService := labores.NewLaborService()
	equipoService := equipos.NewEquipoService()
	actividadService := actividades.NewActividadService()
	loggerService := logger.NewLoggerService()

	// Creamos todos los handlers, inyectando el loggerService
	authHandler := apphandlers.NewAuthHandler(authService, loggerService)
	userHandler := apphandlers.NewUserHandler(authService, userService, loggerService)
	proyectoHandler := apphandlers.NewProyectoHandler(authService, proyectoService, loggerService)
	laborHandler := apphandlers.NewLaborHandler(authService, laborService, loggerService)
	equipoHandler := apphandlers.NewEquipoHandler(authService, equipoService, loggerService)
	actividadHandler := apphandlers.NewActividadHandler(authService, actividadService, loggerService)
	loggerHandler := apphandlers.NewLoggerHandler(authService, loggerService)

	// 2. CONFIGURAR EL ROUTER
	mux := http.NewServeMux()

	// (Aquí van todas tus rutas, las he omitido por brevedad)
	// Rutas de Saludo
	mux.HandleFunc("/api/saludo", apphandlers.SaludoHandler)
	// Rutas de Autenticación
	mux.HandleFunc("/api/auth/register", authHandler.RegisterHandler)
	mux.HandleFunc("/api/auth/login", authHandler.LoginHandler)
	// Rutas de Admin (Usuarios)
	mux.HandleFunc("/api/admin/users", userHandler.AdminUsersHandler)
	mux.HandleFunc("/api/admin/add-user", userHandler.AdminAddUserHandler)
	mux.HandleFunc("/api/admin/delete-user", userHandler.AdminDeleteUserHandler)
	mux.HandleFunc("/api/admin/update-user", userHandler.AdminUpdateRoleHandler)
	mux.HandleFunc("/api/admin/assign-project", userHandler.AdminAssignProjectToUserHandler)
	// Rutas de Usuario (Dashboard)
	mux.HandleFunc("/api/user/project-details", userHandler.UserProjectDetailsHandler)
	// Rutas Admin/Gerente (Proyectos)
	mux.HandleFunc("/api/admin/get-proyectos", proyectoHandler.AdminGetProyectosHandler)
	mux.HandleFunc("/api/admin/create-proyecto", proyectoHandler.AdminCreateProyectoHandler)
	mux.HandleFunc("/api/admin/update-proyecto", proyectoHandler.AdminUpdateProyectoHandler)
	mux.HandleFunc("/api/admin/delete-proyecto", proyectoHandler.AdminDeleteProyectoHandler)
	mux.HandleFunc("/api/admin/set-proyecto-estado", proyectoHandler.AdminSetProyectoEstadoHandler)
	// Rutas Admin/Gerente (Labores)
	mux.HandleFunc("/api/admin/get-labores", laborHandler.GetLaboresHandler)
	mux.HandleFunc("/api/admin/create-labor", laborHandler.CreateLaborHandler)
	mux.HandleFunc("/api/admin/update-labor", laborHandler.UpdateLaborHandler)
	mux.HandleFunc("/api/admin/delete-labor", laborHandler.DeleteLaborHandler)
	// Rutas Admin/Gerente (Equipos)
	mux.HandleFunc("/api/admin/get-equipos", equipoHandler.GetEquiposHandler)
	mux.HandleFunc("/api/admin/create-equipo", equipoHandler.CreateEquipoHandler)
	mux.HandleFunc("/api/admin/update-equipo", equipoHandler.UpdateEquipoHandler)
	mux.HandleFunc("/api/admin/delete-equipo", equipoHandler.DeleteEquipoHandler)
	// Rutas de Actividades (DatosProyecto.js)
	mux.HandleFunc("/api/admin/get-datos-proyecto", actividadHandler.GetDatosProyectoHandler)
	mux.HandleFunc("/api/admin/create-actividad", actividadHandler.CreateActividadHandler)
	mux.HandleFunc("/api/admin/update-actividad", actividadHandler.UpdateActividadHandler)
	mux.HandleFunc("/api/admin/delete-actividad", actividadHandler.DeleteActividadHandler)
	// Ruta del Logger
	mux.HandleFunc("/api/admin/get-logs", loggerHandler.GetLogsHandler)

	// 3. CONFIGURAR MIDDLEWARES (CORS)
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	return corsHandler(mux)
}

func main() {
	// 1. INICIALIZAR LA BASE DE DATOS REAL
	database.InitDB("./users.db")
	defer database.DB.Close()
	log.Println("Base de datos conectada.")

	// 2. OBTENER EL ROUTER ARMADO
	app := setupApp()

	// 3. INICIAR EL SERVIDOR
	port := ":8080"
	log.Printf("Servidor escuchando en http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, app))
}
