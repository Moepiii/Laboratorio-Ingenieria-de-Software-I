package auth

import (
	"fmt"
	"strings"

	"proyecto/internal/database" // <--- RUTA CORREGIDA
)

// CheckPermission verifica si un usuario tiene al menos uno de los roles requeridos
func CheckPermission(username string, requiredRoles ...string) (bool, error) {
	role, err := database.GetUserRole(username) // Usa la función de database
	if err != nil {
		return false, fmt.Errorf("error al obtener rol: %w", err) // Propaga sql.ErrNoRows u otros errores
	}

	userRoleLower := strings.ToLower(role)
	for _, reqRole := range requiredRoles {
		if userRoleLower == strings.ToLower(reqRole) {
			return true, nil // Permiso concedido
		}
	}
	return false, nil // No tiene el rol
}
