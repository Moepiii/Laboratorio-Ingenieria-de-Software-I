package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"proyecto/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// --- QUERIES DE USUARIOS ---

func RegisterUser(username, password, nombre, apellido, cedula string) (int64, error) {
	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("error al hashear password: %w", err)
	}

	// ⭐️ MODIFICADO: Añadido 'nombre', 'apellido' y 'cedula'
	stmt, err := DB.Prepare("INSERT INTO users (username, password, role, nombre, apellido, cedula) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error al preparar inserción: %w", err)
	}
	defer stmt.Close()

	// ⭐️ MODIFICADO: Añadido 'nombre', 'apellido' y 'cedula'
	// Por defecto, el rol es 'user'
	res, err := stmt.Exec(username, string(hashedPassword), "user", nombre, apellido, cedula)
	if err != nil {
		// Manejo de error específico para 'UNIQUE constraint failed: users.username'
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return 0, errors.New("El nombre de usuario ya existe.")
		}
		// Manejo de error específico para 'UNIQUE constraint failed: users.cedula'
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.cedula") {
			return 0, errors.New("La cédula ya está registrada.")
		}
		return 0, fmt.Errorf("error al ejecutar inserción: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener último ID: %w", err)
	}

	return id, nil
}

func GetUserByUsername(username string) (*models.UserDB, error) {
	row := DB.QueryRow("SELECT id, username, password, role, nombre, apellido, cedula, proyecto_id FROM users WHERE username = ?", username)
	var user models.UserDB
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.Role,
		&user.Nombre,   // ⭐️ NUEVO
		&user.Apellido, // ⭐️ NUEVO
		&user.Cedula,   // ⭐️ NUEVO
		&user.ProyectoID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Usuario no encontrado.")
		}
		log.Printf("Error al escanear usuario (GetUserByUsername): %v", err)
		return nil, fmt.Errorf("Error al buscar usuario: %w", err)
	}
	return &user, nil
}

func GetUserRole(username string) (string, error) {
	var role string
	err := DB.QueryRow("SELECT role FROM users WHERE username = ?", username).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("Usuario no encontrado.")
		}
		log.Printf("Error en GetUserRole: %v", err)
		return "", err
	}
	return role, nil
}

func GetAllUsersWithProjectNames() ([]models.UserListResponse, error) {
	// ⭐️ MODIFICADO: Se usa LEFT JOIN para incluir usuarios sin proyecto
	// ⭐️ MODIFICADO: Se seleccionan los campos nuevos (nombre, apellido, cedula)
	rows, err := DB.Query(`
        SELECT u.id, u.username, u.role, u.nombre, u.apellido, u.cedula, u.proyecto_id, p.nombre 
        FROM users u 
        LEFT JOIN proyectos p ON u.proyecto_id = p.id
        ORDER BY u.id ASC
    `)
	if err != nil {
		log.Printf("Error en GetAllUsersWithProjectNames (Query): %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.UserListResponse
	for rows.Next() {
		var user models.UserListResponse
		var proyectoID sql.NullInt64      // Usamos sql.NullInt64 para proyecto_id
		var proyectoNombre sql.NullString // Usamos sql.NullString para p.nombre

		// ⭐️ MODIFICADO: Se escanean los campos nuevos
		if err := rows.Scan(&user.ID, &user.Username, &user.Role, &user.Nombre, &user.Apellido, &user.Cedula, &proyectoID, &proyectoNombre); err != nil {
			log.Printf("Error en GetAllUsersWithProjectNames (Scan): %v", err)
			continue
		}

		// Convertir de sql.NullString a *string
		if proyectoNombre.Valid {
			user.ProyectoNombre = &proyectoNombre.String
		} else {
			user.ProyectoNombre = nil
		}

		// Convertir de sql.NullInt64 a *int
		if proyectoID.Valid {
			idInt := int(proyectoID.Int64)
			user.ProyectoID = &idInt
		} else {
			user.ProyectoID = nil
		}

		users = append(users, user)
	}
	return users, nil
}

func AddUser(user models.User, defaultRole string) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("error al hashear password: %w", err)
	}

	// Asignar rol por defecto si no se especifica
	role := defaultRole
	if role == "" {
		role = "user" // O "encargado" si lo prefieres
	}

	// ⭐️ MODIFICADO: Añadido 'nombre', 'apellido' y 'cedula'
	stmt, err := DB.Prepare("INSERT INTO users (username, password, role, nombre, apellido, cedula) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error al preparar inserción (AddUser): %w", err)
	}
	defer stmt.Close()

	// ⭐️ MODIFICADO: Añadido 'nombre', 'apellido' y 'cedula'
	res, err := stmt.Exec(user.Username, string(hashedPassword), role, user.Nombre, user.Apellido, user.Cedula)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return 0, errors.New("El nombre de usuario ya existe.")
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.cedula") {
			return 0, errors.New("La cédula ya está registrada.")
		}
		return 0, fmt.Errorf("error al ejecutar inserción (AddUser): %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener último ID (AddUser): %w", err)
	}
	return id, nil
}

func DeleteUser(id int) (int64, error) {
	stmt, err := DB.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error al preparar delete (DeleteUser): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("error al ejecutar delete (DeleteUser): %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al obtener filas afectadas (DeleteUser): %w", err)
	}

	return affected, nil
}

func UpdateUserRole(id int, newRole string) (int64, error) {
	stmt, err := DB.Prepare("UPDATE users SET role = ? WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error al preparar update (UpdateUserRole): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(newRole, id)
	if err != nil {
		return 0, fmt.Errorf("error al ejecutar update (UpdateUserRole): %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al obtener filas afectadas (UpdateUserRole): %w", err)
	}

	return affected, nil
}

func AssignProjectToUser(userID int, proyectoID int) (int64, error) {
	var stmt *sql.Stmt
	var err error

	// Si proyectoID es 0, queremos desasignar (poner NULL)
	if proyectoID == 0 {
		stmt, err = DB.Prepare("UPDATE users SET proyecto_id = NULL WHERE id = ?")
		if err != nil {
			return 0, fmt.Errorf("error al preparar update (AssignProjectToUser NULL): %w", err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(userID)
		if err != nil {
			return 0, fmt.Errorf("error al ejecutar update (AssignProjectToUser NULL): %w", err)
		}
		return res.RowsAffected()

	} else {
		// Si proyectoID no es 0, asignamos el proyecto
		stmt, err = DB.Prepare("UPDATE users SET proyecto_id = ? WHERE id = ?")
		if err != nil {
			return 0, fmt.Errorf("error al preparar update (AssignProjectToUser): %w", err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(proyectoID, userID)
		if err != nil {
			return 0, fmt.Errorf("error al ejecutar update (AssignProjectToUser): %w", err)
		}
		return res.RowsAffected()
	}
}

func GetProjectDetailsForUser(userID int) (*models.UserProjectDetailsResponse, error) {
	// 1. Obtener el ID del proyecto del usuario
	var proyectoID sql.NullInt64
	err := DB.QueryRow("SELECT proyecto_id FROM users WHERE id = ?", userID).Scan(&proyectoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Usuario no encontrado.")
		}
		log.Printf("Error obteniendo proyecto_id de usuario %d: %v", userID, err)
		return nil, errors.New("Error al buscar el usuario.")
	}

	// Si no tiene proyecto asignado
	if !proyectoID.Valid {
		return &models.UserProjectDetailsResponse{
			Proyecto: nil,
			Gerentes: nil,
			Miembros: nil,
		}, nil
	}

	// 2. Obtener los detalles del proyecto
	var proyecto models.Proyecto
	projID := proyectoID.Int64
	// ⭐️ AQUI ESTABA UNO DE LOS ERRORES (Falta 'FechaCreacion') ⭐️
	//     Asegúrate de haber añadido FechaCreacion a la struct Proyecto en models.go
	err = DB.QueryRow("SELECT id, nombre, fecha_inicio, fecha_cierre, estado, fecha_creacion FROM proyectos WHERE id = ?", projID).Scan(
		&proyecto.ID, &proyecto.Nombre, &proyecto.FechaInicio, &proyecto.FechaCierre, &proyecto.Estado, &proyecto.FechaCreacion,
	)
	if err != nil {
		log.Printf("Error obteniendo detalles del proyecto %d: %v", projID, err)
		return nil, errors.New("Error al obtener detalles del proyecto.")
	}

	// 3. Obtener los "gerentes" de ese proyecto
	// ⭐️ INICIO DE LA CORRECCIÓN DE TIPO ⭐️
	var gerentes []models.ProjectMember // <- ANTES ERA: models.UserSimple
	rows, err := DB.Query("SELECT id, username, nombre, apellido FROM users WHERE proyecto_id = ? AND role = 'gerente'", projID)
	if err != nil {
		log.Printf("Error obteniendo gerentes del proyecto %d: %v", projID, err)
		return nil, errors.New("Error al obtener gerentes.")
	}
	defer rows.Close()
	for rows.Next() {
		var u models.ProjectMember // <- ANTES ERA: models.UserSimple
		if err := rows.Scan(&u.ID, &u.Username, &u.Nombre, &u.Apellido); err != nil {
			log.Printf("Error escaneando gerente: %v", err)
			continue
		}
		gerentes = append(gerentes, u)
	}

	// 4. Obtener los "miembros" (users) de ese proyecto (excluyendo al usuario actual)
	var miembros []models.ProjectMember // <- ANTES ERA: models.UserSimple
	rowsMiembros, err := DB.Query("SELECT id, username, nombre, apellido FROM users WHERE proyecto_id = ? AND role = 'user' AND id != ?", projID, userID)
	if err != nil {
		log.Printf("Error obteniendo miembros del proyecto %d: %v", projID, err)
		return nil, errors.New("Error al obtener miembros del proyecto.")
	}
	defer rowsMiembros.Close()
	for rowsMiembros.Next() {
		var u models.ProjectMember // <- ANTES ERA: models.UserSimple
		if err := rowsMiembros.Scan(&u.ID, &u.Username, &u.Nombre, &u.Apellido); err != nil {
			log.Printf("Error escaneando miembro: %v", err)
			continue
		}
		miembros = append(miembros, u)
	}
	// ⭐️ FIN DE LA CORRECCIÓN DE TIPO ⭐️

	// 5. Construir la respuesta
	response := &models.UserProjectDetailsResponse{
		Proyecto: &proyecto,
		Gerentes: gerentes,
		Miembros: miembros,
	}

	return response, nil
}

func GetEncargados() ([]models.EncargadoResponse, error) {
	rows, err := DB.Query("SELECT id, nombre, apellido, cedula FROM users WHERE role = 'encargado'")
	if err != nil {
		log.Printf("Error en GetEncargados (Query): %v", err)
		return nil, err
	}
	defer rows.Close()

	var encargados []models.EncargadoResponse
	for rows.Next() {
		var e models.EncargadoResponse
		if err := rows.Scan(&e.ID, &e.Nombre, &e.Apellido, &e.Cedula); err != nil {
			log.Printf("Error en GetEncargados (Scan): %v", err)
			continue
		}
		encargados = append(encargados, e)
	}
	return encargados, nil
}
