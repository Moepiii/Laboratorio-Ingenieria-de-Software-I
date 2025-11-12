package users

import (
	"database/sql"
	"errors"
	"log"

	// "strings" // No es necesario si no se usa

	"proyecto/internal/database"
	"proyecto/internal/models"
	// 'bcrypt' no es necesario aquí
)

// --- 1. EL CONTRATO (Interface) ---
type UserService interface {
	GetAllUsers() ([]models.UserListResponse, error)
	AddUser(user models.User) (int64, error)
	DeleteUser(id int) (int64, error)
	UpdateUserRole(id int, newRole string) (int64, error)
	AssignProjectToUser(userID int, proyectoID int) (int64, error)
	GetProjectDetailsForUser(userID int) (*models.UserProjectDetailsResponse, error)
}

// --- 2. LA IMPLEMENTACIÓN (Struct) ---
type userService struct {
	// (En el futuro, aquí irán dependencias de repositorios)
}

// --- 3. EL CONSTRUCTOR ---
func NewUserService() UserService {
	return &userService{}
}

// --- 4. LOS MÉTODOS (Lógica de Negocio) ---

func (s *userService) GetAllUsers() ([]models.UserListResponse, error) {
	users, err := database.GetAllUsersWithProjectNames()
	if err != nil {
		log.Printf("Error en userService.GetAllUsers: %v", err)
		// ⭐️ ST1005: minúscula y sin punto
		return nil, errors.New("error al obtener usuarios")
	}
	return users, nil
}

func (s *userService) AddUser(user models.User) (int64, error) {
	// La única lógica del servicio es validar.
	if user.Username == "" || user.Password == "" || user.Nombre == "" || user.Apellido == "" || user.Cedula == "" {
		// ⭐️ ST1005: minúscula y sin punto
		return 0, errors.New("todos los campos (username, password, nombre, apellido, cedula) son requeridos")
	}

	// La función 'database.AddUser' (en user_queries.go) se encarga de la encriptación
	// Asignamos "encargado" como rol por defecto desde este servicio
	id, err := database.AddUser(user, "encargado")
	if err != nil {
		log.Printf("Error en userService.AddUser: %v", err)
		// El error de "UNIQUE constraint" ya viene de la base de datos
		return 0, err
	}

	return id, nil
}

func (s *userService) DeleteUser(id int) (int64, error) {
	if id == 0 {
		// ⭐️ ST1005: minúscula y sin punto
		return 0, errors.New("id de usuario requerido")
	}

	// (Opcional: Verificar que no se borre el admin)
	// if id == 1 {
	// 	return 0, errors.New("no se puede borrar al usuario admin")
	// }

	affected, err := database.DeleteUser(id)
	if err != nil {
		log.Printf("Error en userService.DeleteUser (ID: %d): %v", id, err)
		// ⭐️ ST1005: minúscula y sin punto
		return 0, errors.New("error al borrar usuario")
	}
	return affected, nil
}

func (s *userService) UpdateUserRole(id int, newRole string) (int64, error) {
	if id == 0 || newRole == "" {
		// ⭐️ ST1005: minúscula y sin punto
		return 0, errors.New("id y newRole son requeridos")
	}
	if newRole != "admin" && newRole != "gerente" && newRole != "user" && newRole != "encargado" {
		// ⭐️ ST1005: minúscula y sin punto
		return 0, errors.New("rol debe ser 'admin', 'gerente', 'encargado' o 'user'")
	}

	affected, err := database.UpdateUserRole(id, newRole)
	if err != nil {
		log.Printf("Error en userService.UpdateUserRole (ID: %d): %v", id, err)
		// ⭐️ ST1005: minúscula y sin punto
		return 0, errors.New("error al actualizar rol")
	}
	return affected, nil
}

func (s *userService) AssignProjectToUser(userID int, proyectoID int) (int64, error) {
	if userID == 0 {
		// ⭐️ ST1005: minúscula y sin punto
		return 0, errors.New("id de usuario (user_id) requerido")
	}
	// (proyectoID 0 es válido para desasignar)

	affected, err := database.AssignProjectToUser(userID, proyectoID)
	if err != nil {
		log.Printf("Error en userService.AssignProjectToUser (User: %d, Proy: %d): %v", userID, proyectoID, err)
		// ⭐️ ST1005: minúscula y sin punto
		return 0, errors.New("error al asignar proyecto")
	}
	return affected, nil
}

func (s *userService) GetProjectDetailsForUser(userID int) (*models.UserProjectDetailsResponse, error) {
	if userID == 0 {
		// ⭐️ ST1005: minúscula y sin punto
		return nil, errors.New("id de usuario requerido")
	}

	details, err := database.GetProjectDetailsForUser(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// ⭐️ ST1005: minúscula y sin punto
			return nil, errors.New("usuario no encontrado")
		}
		log.Printf("Error en userService.GetProjectDetailsForUser (User: %d): %v", userID, err)
		// ⭐️ ST1005: minúscula y sin punto
		return nil, errors.New("error al obtener detalles del proyecto")
	}

	// Si el usuario no tiene proyecto, 'details' no será nil,
	// pero 'details.Proyecto' sí lo será, lo cual es correcto.
	return details, nil
}
