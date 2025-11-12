package auth

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"proyecto/internal/database"
	"proyecto/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt" // ⭐️ CORRECCIÓN: Importamos bcrypt
)

// jwtKey se mueve del handler al servicio.
var jwtKey = []byte("mi_llave_secreta_super_segura_12345")

// --- 1. EL CONTRATO (Interface) ---
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

	// Validar longitud de contraseña
	if len(user.Password) < 6 {
		return 0, errors.New("La contraseña debe tener al menos 6 caracteres.")
	}

	// La función 'database.RegisterUser' se encarga de la encriptación
	id, err := database.RegisterUser(user.Username, user.Password, user.Nombre, user.Apellido, user.Cedula)
	if err != nil {
		log.Printf("Error en authService.Register: %v", err)
		// El error de "UNIQUE constraint" ya viene de la base de datos
		return 0, err
	}

	return id, nil
}

// ⭐️ CORRECCIÓN AQUÍ ⭐️
func (s *authService) Login(username, password string) (*models.LoginResponse, error) {
	if username == "" || password == "" {
		return nil, errors.New("Usuario y contraseña son requeridos.")
	}

	user, err := database.GetUserByUsername(username)
	if err != nil {
		// No reveles si el usuario existe o no
		return nil, errors.New("Credenciales inválidas.")
	}

	// Aquí es donde se usa bcrypt, en el servicio
	// Comparamos el hash de la BD (user.HashedPassword) con la contraseña (password)
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		// El hash no coincide
		return nil, errors.New("Credenciales inválidas.")
	}

	// --- Si la contraseña es correcta, genera el token ---
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
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
	role, err := database.GetUserRole(username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("CheckPermission: Usuario no encontrado '%s'", username)
			return false, nil // Usuario no existe, no tiene permiso
		}
		log.Printf("CheckPermission: Error al obtener rol de '%s': %v", username, err)
		return false, err // Otro error de DB
	}

	for _, reqRole := range requiredRoles {
		if strings.EqualFold(role, reqRole) {
			return true, nil // El usuario tiene el rol
		}
	}

	log.Printf("CheckPermission: Acceso denegado. Usuario '%s' (Rol: '%s') no tiene rol requerido (%v)", username, role, requiredRoles)
	return false, nil // No se encontró el rol
}
