package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"proyecto/internal/database"
	"proyecto/internal/models" // Importamos los modelos
)

// jwtKey se mueve del handler al servicio.
var jwtKey = []byte("mi_llave_secreta_super_segura_12345")

// --- 1. EL CONTRATO (Interface) ---
// Define QUÉ PUEDE HACER nuestro servicio de autenticación.
type AuthService interface {
	Register(user models.User) (int64, error)
	Login(username, password string) (*models.LoginResponse, error)
	CheckPermission(username string, roles ...string) (bool, error)
}

// --- 2. LA IMPLEMENTACIÓN (Struct) ---
type authService struct {
	// (En el futuro, aquí irán dependencias como un UserRepository)
}

// --- 3. EL CONSTRUCTOR ---
func NewAuthService() AuthService {
	return &authService{}
}

// --- 4. LOS MÉTODOS (Lógica de Negocio) ---

func (s *authService) Register(user models.User) (int64, error) {
	if user.Username == "" || user.Password == "" || user.Nombre == "" || user.Apellido == "" || user.Cedula == "" {
		return 0, errors.New("Todos los campos (username, password, nombre, apellido, cedula) son requeridos.")
	}

	lastID, err := database.RegisterUser(user.Username, user.Password, user.Nombre, user.Apellido, user.Cedula)
	if err != nil {
		log.Printf("Error en authService.Register: %v", err)
		return 0, err
	}
	return lastID, nil
}

func (s *authService) Login(username, password string) (*models.LoginResponse, error) {
	user, err := database.GetUserByUsername(username) // Asumiendo que esta función te devuelve un models.UserDB
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Usuario no encontrado.")
		}
		log.Printf("Error en authService.Login (GetUserByUsername): %v", err)
		return nil, errors.New("Error interno del servidor.")
	}

	if !database.CheckPasswordHash(password, user.HashedPassword) {
		return nil, errors.New("Contraseña incorrecta.")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{ // Usamos el struct de models
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error en authService.Login (SignedString): %v", err)
		return nil, errors.New("Error al generar el token.")
	}

	// Usamos los datos de 'user' que ya obtuvimos
	userDetails := models.UserDetails{
		Username: user.Username,
		Nombre:   user.Nombre,
		Apellido: user.Apellido,
		Cedula:   user.Cedula,
	}

	response := &models.LoginResponse{
		Token:  tokenString,
		User:   userDetails,
		Role:   user.Role,
		UserId: user.ID,
	}

	return response, nil
}

// CheckPermission - Lógica movida de tu auth.go
func (s *authService) CheckPermission(username string, requiredRoles ...string) (bool, error) {
	role, err := database.GetUserRole(username) // Usa la función de database
	if err != nil {
		// Propaga sql.ErrNoRows u otros errores
		return false, fmt.Errorf("error al obtener rol: %w", err) 
	}

	userRoleLower := strings.ToLower(role)
	for _, reqRole := range requiredRoles {
		if userRoleLower == strings.ToLower(reqRole) {
			return true, nil // Permiso concedido
		}
	}
	return false, nil // No tiene el rol
}