package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"

	"proyecto/internal/models"
)

var DB *sql.DB

func InitDB(dbPath string) {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Error al abrir DB: %v", err)
	}
	_, err = DB.Exec("PRAGMA foreign_keys = OFF;")
	if err != nil {
		log.Fatalf("Error PRAGMA OFF: %v", err)
	}
	createUsersTable()
	createProyectosTable()
	createLaboresTable()
	createEquiposTable()
	migrateProyectosTable()
	migrateUsersTable()
	migrateAppTables() // ⭐️ NUEVO: Llamada a la nueva función de migración
	createDefaultAdmin()
	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalf("Error PRAGMA ON: %v", err)
	}
	log.Println("✅ DB inicializada y tablas listas y migradas.")
}

func createUsersTable() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		nombre TEXT NOT NULL DEFAULT '',
		apellido TEXT NOT NULL DEFAULT '',
		cedula TEXT NOT NULL UNIQUE DEFAULT '',
		proyecto_id INTEGER,
		FOREIGN KEY (proyecto_id) REFERENCES proyectos(id) ON DELETE SET NULL
	);`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error al crear tabla users: %v", err)
	}
}

func createProyectosTable() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS proyectos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT NOT NULL UNIQUE,
		fecha_inicio TEXT NOT NULL DEFAULT '',
		fecha_cierre TEXT NOT NULL DEFAULT '',
		estado TEXT NOT NULL DEFAULT 'habilitado'
	);`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error al crear tabla proyectos: %v", err)
	}
}

func createLaboresTable() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS labores_agronomicas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		proyecto_id INTEGER NOT NULL,
		codigo_labor TEXT NOT NULL UNIQUE DEFAULT '', -- ⭐️ NUEVO
		descripcion TEXT NOT NULL,
		estado TEXT NOT NULL DEFAULT 'activa',
		fecha_creacion TEXT NOT NULL,
		FOREIGN KEY (proyecto_id) REFERENCES proyectos(id) ON DELETE CASCADE
	);`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error al crear tabla labores_agronomicas: %v", err)
	}
}

func createEquiposTable() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS equipos_implementos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		proyecto_id INTEGER NOT NULL,
		codigo_equipo TEXT NOT NULL UNIQUE DEFAULT '', -- ⭐️ NUEVO
		nombre TEXT NOT NULL,
		tipo TEXT NOT NULL DEFAULT 'implemento',
		estado TEXT NOT NULL DEFAULT 'disponible',
		fecha_creacion TEXT NOT NULL,
		FOREIGN KEY (proyecto_id) REFERENCES proyectos(id) ON DELETE CASCADE
	);`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error al crear tabla equipos_implementos: %v", err)
	}
}

// ⭐️ NUEVO: Función para migrar las tablas de Labores y Equipos
func migrateAppTables() {
	// 1. Migrar labores_agronomicas
	rowsLabores, err := DB.Query("PRAGMA table_info(labores_agronomicas)")
	if err != nil {
		log.Fatalf("Error PRAGMA table_info(labores_agronomicas): %v", err)
	}
	defer rowsLabores.Close()

	columnsLabores := make(map[string]bool)
	for rowsLabores.Next() {
		var (
			cid        int
			name       string
			dtype      string
			notnull    int
			dflt_value sql.NullString
			pk         int
		)
		rowsLabores.Scan(&cid, &name, &dtype, &notnull, &dflt_value, &pk)
		columnsLabores[name] = true
	}

	if !columnsLabores["codigo_labor"] {
		_, err := DB.Exec("ALTER TABLE labores_agronomicas ADD COLUMN codigo_labor TEXT NOT NULL UNIQUE DEFAULT ''")
		if err != nil {
			log.Printf("Advertencia al migrar labores (add codigo_labor): %v. Puede que ya exista.", err)
		} else {
			log.Println("Migración: Columna 'codigo_labor' añadida a 'labores_agronomicas'.")
		}
	}

	// 2. Migrar equipos_implementos
	rowsEquipos, err := DB.Query("PRAGMA table_info(equipos_implementos)")
	if err != nil {
		log.Fatalf("Error PRAGMA table_info(equipos_implementos): %v", err)
	}
	defer rowsEquipos.Close()

	columnsEquipos := make(map[string]bool)
	for rowsEquipos.Next() {
		var (
			cid        int
			name       string
			dtype      string
			notnull    int
			dflt_value sql.NullString
			pk         int
		)
		rowsEquipos.Scan(&cid, &name, &dtype, &notnull, &dflt_value, &pk)
		columnsEquipos[name] = true
	}

	if !columnsEquipos["codigo_equipo"] {
		_, err := DB.Exec("ALTER TABLE equipos_implementos ADD COLUMN codigo_equipo TEXT NOT NULL UNIQUE DEFAULT ''")
		if err != nil {
			log.Printf("Advertencia al migrar equipos (add codigo_equipo): %v. Puede que ya exista.", err)
		} else {
			log.Println("Migración: Columna 'codigo_equipo' añadida a 'equipos_implementos'.")
		}
	}
}

func migrateProyectosTable() {
	rows, err := DB.Query("PRAGMA table_info(proyectos)")
	if err != nil {
		log.Fatalf("Error PRAGMA table_info(proyectos): %v", err)
	}
	defer rows.Close()

	columns := make(map[string]bool)
	for rows.Next() {
		var (
			cid        int
			name       string
			dtype      string
			notnull    int
			dflt_value sql.NullString
			pk         int
		)
		if err := rows.Scan(&cid, &name, &dtype, &notnull, &dflt_value, &pk); err != nil {
			log.Fatalf("Error escaneando table_info(proyectos): %v", err)
		}
		columns[name] = true
	}

	if !columns["fecha_inicio"] {
		_, err := DB.Exec("ALTER TABLE proyectos ADD COLUMN fecha_inicio TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error al migrar proyectos (add fecha_inicio): %v", err)
		}
		log.Println("Migración: Columna 'fecha_inicio' añadida a 'proyectos'.")
	}
	if !columns["fecha_cierre"] {
		_, err := DB.Exec("ALTER TABLE proyectos ADD COLUMN fecha_cierre TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error al migrar proyectos (add fecha_cierre): %v", err)
		}
		log.Println("Migración: Columna 'fecha_cierre' añadida a 'proyectos'.")
	}
	if !columns["estado"] {
		_, err := DB.Exec("ALTER TABLE proyectos ADD COLUMN estado TEXT NOT NULL DEFAULT 'habilitado'")
		if err != nil {
			log.Fatalf("Error al migrar proyectos (add estado): %v", err)
		}
		log.Println("Migración: Columna 'estado' añadida a 'proyectos'.")
	}
}

func migrateUsersTable() {
	rows, err := DB.Query("PRAGMA table_info(users)")
	if err != nil {
		log.Fatalf("Error PRAGMA table_info(users): %v", err)
	}
	defer rows.Close()

	columns := make(map[string]bool)
	for rows.Next() {
		var (
			cid        int
			name       string
			dtype      string
			notnull    int
			dflt_value sql.NullString
			pk         int
		)
		if err := rows.Scan(&cid, &name, &dtype, &notnull, &dflt_value, &pk); err != nil {
			log.Fatalf("Error escaneando table_info(users): %v", err)
		}
		columns[name] = true
	}

	if !columns["nombre"] {
		_, err := DB.Exec("ALTER TABLE users ADD COLUMN nombre TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error al migrar users (add nombre): %v", err)
		}
		log.Println("Migración: Columna 'nombre' añadida a 'users'.")
	}
	if !columns["apellido"] {
		_, err := DB.Exec("ALTER TABLE users ADD COLUMN apellido TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error al migrar users (add apellido): %v", err)
		}
		log.Println("Migración: Columna 'apellido' añadida a 'users'.")
	}
	if !columns["proyecto_id"] {
		_, err := DB.Exec("ALTER TABLE users ADD COLUMN proyecto_id INTEGER REFERENCES proyectos(id) ON DELETE SET NULL")
		if err != nil {
			log.Fatalf("Error al migrar users (add proyecto_id): %v", err)
		}
		log.Println("Migración: Columna 'proyecto_id' añadida a 'users'.")
	}
	if !columns["cedula"] {
		_, err := DB.Exec("ALTER TABLE users ADD COLUMN cedula TEXT NOT NULL UNIQUE DEFAULT ''")
		if err != nil {
			log.Printf("Advertencia al migrar users (add cedula): %v. Puede que ya exista.", err)
		} else {
			log.Println("Migración: Columna 'cedula' añadida a 'users'.")
		}
	}
}

func createDefaultAdmin() {
	var exists int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&exists)
	if err != nil {
		log.Fatalf("Error al verificar admin: %v", err)
	}

	if exists == 0 {
		hashedPassword, _ := HashPassword("admin123")
		_, err := DB.Exec("INSERT INTO users (username, password, role, nombre, apellido, cedula) VALUES (?, ?, ?, ?, ?, ?)",
			"admin", hashedPassword, "admin", "Administrador", "Del Sistema", "00000001")
		if err != nil {
			log.Fatalf("Error al crear admin por defecto: %v", err)
		}
		log.Println("Usuario 'admin' por defecto creado con contraseña 'admin123'.")
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserRole(username string) (string, error) {
	var role string
	err := DB.QueryRow("SELECT role FROM users WHERE username = ?", username).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("usuario '%s' no encontrado", username)
		}
		return "", fmt.Errorf("error al obtener rol: %w", err)
	}
	return role, nil
}

func GetUserByUsername(username string) (models.UserDB, error) {
	var user models.UserDB
	err := DB.QueryRow("SELECT id, username, password, role, nombre, apellido, cedula, proyecto_id FROM users WHERE username = ?", username).Scan(
		&user.ID, &user.Username, &user.HashedPassword, &user.Role, &user.Nombre, &user.Apellido, &user.Cedula, &user.ProyectoID)
	return user, err
}

func RegisterUser(username, password, nombre, apellido, cedula string) (int64, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("error al hashear password: %w", err)
	}
	stmt, err := DB.Prepare("INSERT INTO users (username, password, nombre, apellido, cedula) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error preparando statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(username, hashedPassword, nombre, apellido, cedula)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return 0, fmt.Errorf("el nombre de usuario '%s' ya existe", username)
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.cedula") {
			return 0, fmt.Errorf("la cédula '%s' ya está registrada", cedula)
		}
		return 0, fmt.Errorf("error ejecutando insert: %w", err)
	}
	return res.LastInsertId()
}

func GetAllUsersWithProjectNames() ([]models.UserListResponse, error) {
	query := `
	SELECT u.id, u.username, u.role, u.nombre, u.apellido, u.cedula, u.proyecto_id, p.nombre
	FROM users u
	LEFT JOIN proyectos p ON u.proyecto_id = p.id
	ORDER BY u.id
	`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error en query GetAllUsers: %w", err)
	}
	defer rows.Close()

	var users []models.UserListResponse
	for rows.Next() {
		var user models.UserListResponse
		var proyectoID sql.NullInt64
		var proyectoNombre sql.NullString

		if err := rows.Scan(&user.ID, &user.Username, &user.Role, &user.Nombre, &user.Apellido, &user.Cedula, &proyectoID, &proyectoNombre); err != nil {
			log.Printf("Error escaneando usuario: %v", err)
			continue
		}
		if proyectoID.Valid {
			id := int(proyectoID.Int64)
			user.ProyectoID = &id
		}
		if proyectoNombre.Valid {
			user.ProyectoNombre = &proyectoNombre.String
		}

		users = append(users, user)
	}
	return users, nil
}

func AddUser(user models.User, hashedPassword string) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO users (username, password, nombre, apellido, cedula, role) VALUES (?, ?, ?, ?, ?, 'user')")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.Username, hashedPassword, user.Nombre, user.Apellido, user.Cedula)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return 0, fmt.Errorf("el nombre de usuario '%s' ya existe", user.Username)
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.cedula") {
			return 0, fmt.Errorf("la cédula '%s' ya está registrada", user.Cedula)
		}
		return 0, err
	}
	return res.LastInsertId()
}

func DeleteUser(id int) (int64, error) {
	stmt, err := DB.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func UpdateUserRole(id int, newRole string) (int64, error) {
	stmt, err := DB.Prepare("UPDATE users SET role = ? WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(newRole, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func AssignProjectToUser(userID int, proyectoID int) (int64, error) {
	stmt, err := DB.Prepare("UPDATE users SET proyecto_id = ? WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error preparando update: %w", err)
	}
	defer stmt.Close()

	var res sql.Result
	if proyectoID == 0 {
		res, err = stmt.Exec(sql.NullInt64{}, userID)
	} else {
		res, err = stmt.Exec(proyectoID, userID)
	}

	if err != nil {
		return 0, fmt.Errorf("error ejecutando update: %w", err)
	}
	return res.RowsAffected()
}

func GetProjectByID(id int64) (*models.Proyecto, error) {
	var p models.Proyecto
	err := DB.QueryRow("SELECT id, nombre, fecha_inicio, fecha_cierre, estado FROM proyectos WHERE id = ?", id).Scan(&p.ID, &p.Nombre, &p.FechaInicio, &p.FechaCierre, &p.Estado)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func CreateProyecto(nombre, fechaInicio, fechaCierre string) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO proyectos (nombre, fecha_inicio, fecha_cierre, estado) VALUES (?, ?, ?, 'habilitado')")
	if err != nil {
		return 0, fmt.Errorf("error preparando statement: %w", err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(nombre, fechaInicio, fechaCierre)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: proyectos.nombre") {
			return 0, fmt.Errorf("el nombre de proyecto '%s' ya existe", nombre)
		}
		return 0, fmt.Errorf("error ejecutando insert: %w", err)
	}
	return res.LastInsertId()
}

func GetAllProyectos() ([]models.Proyecto, error) {
	rows, err := DB.Query("SELECT id, nombre, fecha_inicio, fecha_cierre, estado FROM proyectos ORDER BY id DESC")
	if err != nil {
		return nil, fmt.Errorf("error en query GetAllProyectos: %w", err)
	}
	defer rows.Close()
	var proyectos []models.Proyecto
	for rows.Next() {
		var p models.Proyecto
		if err := rows.Scan(&p.ID, &p.Nombre, &p.FechaInicio, &p.FechaCierre, &p.Estado); err != nil {
			log.Printf("Error escaneando proyecto: %v", err)
			continue
		}
		proyectos = append(proyectos, p)
	}
	return proyectos, nil
}

func UpdateProyecto(id int, nombre, fechaInicio, fechaCierre string) (int64, error) {
	stmt, err := DB.Prepare("UPDATE proyectos SET nombre = ?, fecha_inicio = ?, fecha_cierre = ? WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(nombre, fechaInicio, fechaCierre, id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: proyectos.nombre") {
			return 0, fmt.Errorf("el nombre de proyecto '%s' ya existe", nombre)
		}
		return 0, err
	}
	return result.RowsAffected()
}

func DeleteProyecto(id int) (int64, error) {
	stmt, err := DB.Prepare("DELETE FROM proyectos WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func SetProjectState(id int, estado string) (int64, error) {
	stmt, err := DB.Prepare("UPDATE proyectos SET estado = ? WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(estado, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func GetProjectDetailsForUser(userID int) (*models.UserProjectDetailsResponse, error) {
	var proyectoID sql.NullInt64
	err := DB.QueryRow("SELECT proyecto_id FROM users WHERE id = ?", userID).Scan(&proyectoID)
	if err != nil {
		return nil, fmt.Errorf("error al buscar usuario: %w", err)
	}
	response := &models.UserProjectDetailsResponse{Proyecto: nil, Miembros: []models.ProjectMember{}, Gerentes: []models.ProjectMember{}}
	if proyectoID.Valid {
		pID := proyectoID.Int64
		proyecto, err := GetProjectByID(pID)
		if err != nil {
			log.Printf("Advertencia: Usuario %d tiene proyecto_id %d inválido: %v", userID, pID, err)
		} else {
			response.Proyecto = proyecto
		}
		rows, err := DB.Query(`SELECT id, username, nombre, apellido, role FROM users WHERE proyecto_id = ? AND id != ? ORDER BY role, nombre`, pID, userID)
		if err != nil {
			return nil, fmt.Errorf("error al buscar miembros del proyecto %d: %w", pID, err)
		}
		defer rows.Close()
		for rows.Next() {
			var member models.ProjectMember
			if err := rows.Scan(&member.ID, &member.Username, &member.Nombre, &member.Apellido, &member.Role); err != nil {
				log.Printf("Error escaneando miembro: %v", err)
				continue
			}
			if member.Role == "gerente" {
				response.Gerentes = append(response.Gerentes, member)
			} else {
				response.Miembros = append(response.Miembros, member)
			}
		}
	}
	return response, nil
}

// --- Funciones CRUD para Labores Agronómicas ---

// ⭐️ MODIFICADO: Acepta 'codigo_labor'
func CreateLabor(labor models.LaborAgronomica) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO labores_agronomicas (proyecto_id, codigo_labor, descripcion, estado, fecha_creacion) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error preparando statement: %w", err)
	}
	defer stmt.Close()
	fechaCreacion := time.Now().Format(time.RFC3339)
	res, err := stmt.Exec(labor.ProyectoID, labor.CodigoLabor, labor.Descripcion, labor.Estado, fechaCreacion)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: labores_agronomicas.codigo_labor") {
			return 0, fmt.Errorf("el código de labor '%s' ya existe", labor.CodigoLabor)
		}
		return 0, fmt.Errorf("error ejecutando insert: %w", err)
	}
	return res.LastInsertId()
}

// ⭐️ MODIFICADO: Devuelve 'codigo_labor'
func GetLaboresByProyectoID(proyectoID int) ([]models.LaborAgronomica, error) {
	rows, err := DB.Query("SELECT id, proyecto_id, codigo_labor, descripcion, estado, fecha_creacion FROM labores_agronomicas WHERE proyecto_id = ? ORDER BY fecha_creacion DESC", proyectoID)
	if err != nil {
		return nil, fmt.Errorf("error en query GetLaboresByProyectoID: %w", err)
	}
	defer rows.Close()

	var labores []models.LaborAgronomica
	for rows.Next() {
		var labor models.LaborAgronomica
		if err := rows.Scan(&labor.ID, &labor.ProyectoID, &labor.CodigoLabor, &labor.Descripcion, &labor.Estado, &labor.FechaCreacion); err != nil {
			log.Printf("Error escaneando labor: %v", err)
			continue
		}
		labores = append(labores, labor)
	}
	return labores, nil
}

// ⭐️ MODIFICADO: Devuelve 'codigo_labor'
func GetLaborByID(id int) (*models.LaborAgronomica, error) {
	var labor models.LaborAgronomica
	err := DB.QueryRow("SELECT id, proyecto_id, codigo_labor, descripcion, estado, fecha_creacion FROM labores_agronomicas WHERE id = ?", id).Scan(
		&labor.ID, &labor.ProyectoID, &labor.CodigoLabor, &labor.Descripcion, &labor.Estado, &labor.FechaCreacion)
	if err != nil {
		return nil, err
	}
	return &labor, nil
}

// ⭐️ MODIFICADO: Actualiza 'codigo_labor'
func UpdateLabor(id int, codigo_labor string, descripcion string, estado string) (int64, error) {
	stmt, err := DB.Prepare("UPDATE labores_agronomicas SET codigo_labor = ?, descripcion = ?, estado = ? WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(codigo_labor, descripcion, estado, id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: labores_agronomicas.codigo_labor") {
			return 0, fmt.Errorf("el código de labor '%s' ya existe", codigo_labor)
		}
		return 0, err
	}
	return result.RowsAffected()
}

func DeleteLabor(id int) (int64, error) {
	stmt, err := DB.Prepare("DELETE FROM labores_agronomicas WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// --- Funciones CRUD para Equipos e Implementos ---

// ⭐️ MODIFICADO: Acepta 'codigo_equipo'
func CreateEquipo(equipo models.EquipoImplemento) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO equipos_implementos (proyecto_id, codigo_equipo, nombre, tipo, estado, fecha_creacion) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error preparando statement: %w", err)
	}
	defer stmt.Close()
	fechaCreacion := time.Now().Format(time.RFC3339)
	res, err := stmt.Exec(equipo.ProyectoID, equipo.CodigoEquipo, equipo.Nombre, equipo.Tipo, equipo.Estado, fechaCreacion)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: equipos_implementos.codigo_equipo") {
			return 0, fmt.Errorf("el código de equipo '%s' ya existe", equipo.CodigoEquipo)
		}
		return 0, fmt.Errorf("error ejecutando insert: %w", err)
	}
	return res.LastInsertId()
}

// ⭐️ MODIFICADO: Devuelve 'codigo_equipo'
func GetEquiposByProyectoID(proyectoID int) ([]models.EquipoImplemento, error) {
	rows, err := DB.Query("SELECT id, proyecto_id, codigo_equipo, nombre, tipo, estado, fecha_creacion FROM equipos_implementos WHERE proyecto_id = ? ORDER BY fecha_creacion DESC", proyectoID)
	if err != nil {
		return nil, fmt.Errorf("error en query GetEquiposByProyectoID: %w", err)
	}
	defer rows.Close()

	var equipos []models.EquipoImplemento
	for rows.Next() {
		var equipo models.EquipoImplemento
		if err := rows.Scan(&equipo.ID, &equipo.ProyectoID, &equipo.CodigoEquipo, &equipo.Nombre, &equipo.Tipo, &equipo.Estado, &equipo.FechaCreacion); err != nil {
			log.Printf("Error escaneando equipo: %v", err)
			continue
		}
		equipos = append(equipos, equipo)
	}
	return equipos, nil
}

// ⭐️ MODIFICADO: Devuelve 'codigo_equipo'
func GetEquipoByID(id int) (*models.EquipoImplemento, error) {
	var equipo models.EquipoImplemento
	err := DB.QueryRow("SELECT id, proyecto_id, codigo_equipo, nombre, tipo, estado, fecha_creacion FROM equipos_implementos WHERE id = ?", id).Scan(
		&equipo.ID, &equipo.ProyectoID, &equipo.CodigoEquipo, &equipo.Nombre, &equipo.Tipo, &equipo.Estado, &equipo.FechaCreacion)
	if err != nil {
		return nil, err
	}
	return &equipo, nil
}

// ⭐️ MODIFICADO: Actualiza 'codigo_equipo'
func UpdateEquipo(id int, codigo_equipo string, nombre string, tipo string, estado string) (int64, error) {
	stmt, err := DB.Prepare("UPDATE equipos_implementos SET codigo_equipo = ?, nombre = ?, tipo = ?, estado = ? WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(codigo_equipo, nombre, tipo, estado, id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: equipos_implementos.codigo_equipo") {
			return 0, fmt.Errorf("el código de equipo '%s' ya existe", codigo_equipo)
		}
		return 0, err
	}
	return result.RowsAffected()
}

func DeleteEquipo(id int) (int64, error) {
	stmt, err := DB.Prepare("DELETE FROM equipos_implementos WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
