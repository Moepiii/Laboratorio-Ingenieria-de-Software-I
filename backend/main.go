package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers" // Este sigue siendo 'handlers' para CORS

	"proyecto/internal/database"
	apphandlers "proyecto/internal/handlers" // <--- ALIAS AÑADIDO 'apphandlers'
)

func main() {
	database.InitDB("../Base de datos/users.db") // Ajusta la ruta si es necesario

	// Configuración de CORS (usa 'handlers' de gorilla)
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	corsOptions := handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)

	// --- Definición de Rutas ---
	mux := http.NewServeMux()

	// Usa 'apphandlers' para llamar a TUS handlers
	mux.HandleFunc("/api/saludo", apphandlers.SaludoHandler)
	mux.HandleFunc("/api/register", apphandlers.RegisterHandler)
	mux.HandleFunc("/api/login", apphandlers.LoginHandler)

	// Rutas Admin/Gerente (Usuarios) - Usa 'apphandlers'
	mux.HandleFunc("/api/admin/users", apphandlers.AdminUsersHandler)
	mux.HandleFunc("/api/admin/add-user", apphandlers.AdminAddUserHandler)
	mux.HandleFunc("/api/admin/delete-user", apphandlers.AdminDeleteUserHandler)
	mux.HandleFunc("/api/admin/update-user", apphandlers.AdminUpdateUserHandler)
	mux.HandleFunc("/api/admin/assign-proyecto", apphandlers.AdminAssignProyectoHandler)

	// Rutas Admin/Gerente (Proyectos) - Usa 'apphandlers'
	mux.HandleFunc("/api/admin/get-proyectos", apphandlers.AdminGetProyectosHandler)
	mux.HandleFunc("/api/admin/create-proyecto", apphandlers.AdminCreateProyectoHandler)
	mux.HandleFunc("/api/admin/delete-proyecto", apphandlers.AdminDeleteProyectoHandler)
	mux.HandleFunc("/api/admin/update-proyecto", apphandlers.AdminUpdateProyectoHandler)
	mux.HandleFunc("/api/admin/set-proyecto-estado", apphandlers.AdminSetProyectoEstadoHandler)

	// Ruta Usuario - Usa 'apphandlers'
	mux.HandleFunc("/api/user/project-details", apphandlers.UserProjectDetailsHandler)

	// Iniciar Servidor
	log.Println("🚀 Servidor Go modularizado en :8080")
	// Aplica CORS (de 'handlers' gorilla) a todas las rutas
	if err := http.ListenAndServe(":8080", corsOptions(mux)); err != nil {
		log.Fatalf("Error al iniciar servidor: %v", err)
	}
}
