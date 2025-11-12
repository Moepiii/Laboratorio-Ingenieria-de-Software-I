package auth

import (
	"errors"
	"log"
	"strings" // ⭐️ 1. IMPORTAMOS "strings"
	"time"

	"proyecto/internal/database"
	"proyecto/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
		return 0, errors.New("todos los campos (username, password, nombre, apellido, cedula) son requeridos")
	}

	// Validar longitud de contraseña
	if len(user.Password) < 6 {
		return 0, errors.New("la contraseña debe tener al menos 6 caracteres")
	}

	id, err := database.RegisterUser(user.Username, user.Password, user.Nombre, user.Apellido, user.Cedula)
	if err != nil {
		log.Printf("Error en authService.Register: %v", err)
		return 0, err
	}

	return id, nil
}

func (s *authService) Login(username, password string) (*models.LoginResponse, error) {
	if username == "" || password == "" {
		return nil, errors.New("usuario y contraseña son requeridos")
	}

	user, err := database.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, errors.New("credenciales inválidas")
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
		return nil, errors.New("error al generar el token")
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
// ⭐️ --- INICIO DE LA CORRECCIÓN --- ⭐️
func (s *authService) CheckPermission(username string, requiredRoles ...string) (bool, error) {
	role, err := database.GetUserRole(username)
	if err != nil {
		// Si el error es "Usuario no encontrado", no es un error 500.
		// Es un simple "no tiene permiso".
		// Comparamos contra el string de error que definimos en user_queries.go
		if strings.Contains(err.Error(), "Usuario no encontrado") {
			log.Printf("CheckPermission: Usuario no encontrado '%s'", username)
			return false, nil // ⬅️ Devolvemos (false, nil)
		}

		// (Tu log original de sql.ErrNoRows ya no era necesario porque
		// GetUserRole devuelve un error personalizado)

		// Otro error (ej. DB desconectada) SÍ es un 500.
		log.Printf("CheckPermission: Error al obtener rol de '%s': %v", username, err)
		return false, err // ⬅️ Devolvemos (false, err)
	}

	// (El resto de la función es idéntica)
	for _, reqRole := range requiredRoles {
		if strings.EqualFold(role, reqRole) {
			return true, nil // El usuario tiene el rol
		}
	}

	log.Printf("CheckPermission: Acceso denegado. Usuario '%s' (Rol: '%s') no tiene rol requerido (%v)", username, role, requiredRoles)
	return false, nil // No se encontró el rol
}

// ⭐️ --- FIN DE LA CORRECCIÓN --- ⭐️
