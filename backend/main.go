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
	"proyecto/internal/proyectos"
	"proyecto/internal/users"
)

func main() {
	// 1. INICIALIZAR LA BASE DE DATOS
	database.InitDB("../Base de datos/users.db")
	defer database.DB.Close() // ⭐️ ¡CORRECCIÓN!
	log.Println("Base de datos conectada.")

	// 2. "ARMAR" LA APLICACIÓN (Inyección de Dependencias)

	// Creamos todos los servicios
	authService := auth.NewAuthService()
	userService := users.NewUserService()
	proyectoService := proyectos.NewProyectoService()
	laborService := labores.NewLaborService()
	equipoService := equipos.NewEquipoService()
	actividadService := actividades.NewActividadService()

	// Creamos todos los handlers
	authHandler := apphandlers.NewAuthHandler(authService)
	userHandler := apphandlers.NewUserHandler(authService, userService)
	proyectoHandler := apphandlers.NewProyectoHandler(authService, proyectoService)
	laborHandler := apphandlers.NewLaborHandler(authService, laborService)
	equipoHandler := apphandlers.NewEquipoHandler(authService, equipoService)
	actividadHandler := apphandlers.NewActividadHandler(authService, actividadService)

	// 3. CONFIGURACIÓN DE RUTAS Y CORS
	mux := http.NewServeMux()

	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	corsOptions := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)

	// --- Definición de Rutas ---

	// Rutas de autenticación (listas)
	mux.HandleFunc("/api/saludo", apphandlers.SaludoHandler)
	mux.HandleFunc("/api/register", authHandler.RegisterHandler)
	mux.HandleFunc("/api/login", authHandler.LoginHandler)

	// Rutas Admin/Gerente (Usuarios) (listas)
	mux.HandleFunc("/api/admin/users", userHandler.AdminUsersHandler)
	mux.HandleFunc("/api/admin/add-user", userHandler.AdminAddUserHandler)
	mux.HandleFunc("/api/admin/delete-user", userHandler.AdminDeleteUserHandler)
	mux.HandleFunc("/api/admin/update-user", userHandler.AdminUpdateUserHandler)
	mux.HandleFunc("/api/admin/assign-proyecto", userHandler.AdminAssignProyectoHandler)

	// Ruta Usuario (lista)
	mux.HandleFunc("/api/user/proyecto-details", userHandler.UserProjectDetailsHandler)

	// Rutas Admin/Gerente (Proyectos) (listas)
	mux.HandleFunc("/api/admin/get-proyectos", proyectoHandler.AdminGetProyectosHandler)
	mux.HandleFunc("/api/admin/create-proyecto", proyectoHandler.AdminCreateProyectoHandler)
	mux.HandleFunc("/api/admin/update-proyecto", proyectoHandler.AdminUpdateProyectoHandler)
	mux.HandleFunc("/api/admin/delete-proyecto", proyectoHandler.AdminDeleteProyectoHandler)
	mux.HandleFunc("/api/admin/set-proyecto-estado", proyectoHandler.AdminSetProyectoEstadoHandler)

	// Rutas Admin/Gerente (Labores) (listas)
	mux.HandleFunc("/api/admin/get-labores", laborHandler.GetLaboresHandler)
	mux.HandleFunc("/api/admin/create-labor", laborHandler.CreateLaborHandler)
	mux.HandleFunc("/api/admin/update-labor", laborHandler.UpdateLaborHandler)
	mux.HandleFunc("/api/admin/delete-labor", laborHandler.DeleteLaborHandler)

	// Rutas Admin/Gerente (Equipos) (listas)
	mux.HandleFunc("/api/admin/get-equipos", equipoHandler.GetEquiposHandler)
	mux.HandleFunc("/api/admin/create-equipo", equipoHandler.CreateEquipoHandler)
	mux.HandleFunc("/api/admin/update-equipo", equipoHandler.UpdateEquipoHandler)
	mux.HandleFunc("/api/admin/delete-equipo", equipoHandler.DeleteEquipoHandler)

	// Rutas de Actividades (DatosProyecto.js) (listas)
	mux.HandleFunc("/api/admin/get-datos-proyecto", actividadHandler.GetDatosProyectoHandler)
	mux.HandleFunc("/api/admin/create-actividad", actividadHandler.CreateActividadHandler)
	mux.HandleFunc("/api/admin/update-actividad", actividadHandler.UpdateActividadHandler)
	mux.HandleFunc("/api/admin/delete-actividad", actividadHandler.DeleteActividadHandler)

	// 4. INICIAR EL SERVIDOR
	log.Println("Servidor escuchando en http://localhost:8080")
	if err := http.ListenAndServe(":8080", corsOptions(mux)); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}