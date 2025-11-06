package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers" // Usas 'handlers' de gorilla para CORS

	"proyecto/internal/database"
	apphandlers "proyecto/internal/handlers" // Tu alias 'apphandlers'
)

func main() {
	database.InitDB("../Base de datos/users.db") // Asegúrate que la ruta a tu DB sea correcta

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

	// Rutas Admin/Gerente (Labores)
	mux.HandleFunc("/api/admin/get-labores", apphandlers.GetLaboresHandler)
	mux.HandleFunc("/api/admin/create-labor", apphandlers.CreateLaborHandler)
	mux.HandleFunc("/api/admin/update-labor", apphandlers.UpdateLaborHandler)
	mux.HandleFunc("/api/admin/delete-labor", apphandlers.DeleteLaborHandler)

	// ⭐️ --- (INICIO) Nuevas Rutas de Equipos --- ⭐️
	mux.HandleFunc("/api/admin/get-equipos", apphandlers.GetEquiposHandler)
	mux.HandleFunc("/api/admin/create-equipo", apphandlers.CreateEquipoHandler)
	mux.HandleFunc("/api/admin/update-equipo", apphandlers.UpdateEquipoHandler)
	mux.HandleFunc("/api/admin/delete-equipo", apphandlers.DeleteEquipoHandler)
	// ⭐️ --- (FIN) Nuevas Rutas de Equipos --- ⭐️

	// Ruta Usuario - Usa 'apphandlers'
	mux.HandleFunc("/api/user/project-details", apphandlers.UserProjectDetailsHandler)

	// --- Iniciar Servidor ---
	log.Println("🚀 Servidor Go escuchando en http://localhost:8080")
	// Aplica el middleware de CORS a tu mux
	log.Fatal(http.ListenAndServe(":8080", corsOptions(mux)))
}
