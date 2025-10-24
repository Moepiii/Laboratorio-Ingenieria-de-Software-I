package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"

	"proyecto/internal/models" // <--- RUTA CORREGIDA
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
	migrateProyectosTable()
	migrateUsersTable()
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
		id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT NOT NULL UNIQUE, password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user', nombre TEXT NOT NULL DEFAULT '', apellido TEXT NOT NULL DEFAULT ''
	);`
	if _, err := DB.Exec(createTableSQL); err != nil {
		log.Fatalf("Error crear 'users': %v", err)
	}
}

func createProyectosTable() {
	createProyectosTableSQL := `
	CREATE TABLE IF NOT EXISTS proyectos (
		id INTEGER PRIMARY KEY AUTOINCREMENT, nombre TEXT NOT NULL UNIQUE,
		fecha_inicio TEXT NOT NULL DEFAULT '', fecha_cierre TEXT NOT NULL DEFAULT '',
		estado TEXT NOT NULL DEFAULT 'habilitado'
	);`
	if _, err := DB.Exec(createProyectosTableSQL); err != nil {
		log.Fatalf("Error crear 'proyectos': %v", err)
	}
}

func migrateProyectosTable() {
	rows, err := DB.Query("PRAGMA table_info(proyectos)")
	if err != nil {
		log.Fatalf("Error PRAGMA 'proyectos': %v", err)
	}
	defer rows.Close()
	columnExists := map[string]bool{"nombre": false, "fecha_inicio": false, "fecha_cierre": false, "estado": false}
	for rows.Next() {
		var (
			cid        int
			name       string
			ctype      string
			notnull    int
			dflt_value sql.NullString
			pk         int
		)
		err = rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk)
		if err != nil {
			log.Fatalf("Error scan 'proyectos': %v", err)
		}
		if _, ok := columnExists[name]; ok {
			columnExists[name] = true
		}
	}
	rows.Close()
	if !columnExists["nombre"] {
		log.Println("⚠️ Recreando 'proyectos' por falta de 'nombre'...")
		_, err = DB.Exec("DROP TABLE IF EXISTS proyectos;")
		if err != nil {
			log.Fatalf("Error DROP 'proyectos': %v", err)
		}
		createProyectosTable()
		log.Println("✅ Tabla 'proyectos' recreada.")
		return
	}
	if !columnExists["fecha_inicio"] {
		log.Println("⚠️ Migrando 'fecha_inicio' en 'proyectos'...")
		_, err = DB.Exec("ALTER TABLE proyectos ADD COLUMN fecha_inicio TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error migrar 'fecha_inicio': %v", err)
		}
	}
	if !columnExists["fecha_cierre"] {
		log.Println("⚠️ Migrando 'fecha_cierre' en 'proyectos'...")
		_, err = DB.Exec("ALTER TABLE proyectos ADD COLUMN fecha_cierre TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error migrar 'fecha_cierre': %v", err)
		}
	}
	if !columnExists["estado"] {
		log.Println("⚠️ Migrando 'estado' en 'proyectos'...")
		_, err = DB.Exec("ALTER TABLE proyectos ADD COLUMN estado TEXT NOT NULL DEFAULT 'habilitado'")
		if err != nil {
			log.Fatalf("Error migrar 'estado': %v", err)
		}
	}
}

func migrateUsersTable() {
	rows, err := DB.Query("PRAGMA table_info(users)")
	if err != nil {
		log.Fatalf("Error PRAGMA 'users': %v", err)
	}
	defer rows.Close()
	columnExists := map[string]bool{"role": false, "nombre": false, "apellido": false, "proyecto_id": false}
	for rows.Next() {
		var (
			cid        int
			name       string
			ctype      string
			notnull    int
			dflt_value sql.NullString
			pk         int
		)
		err = rows.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk)
		if err != nil {
			log.Fatalf("Error scan 'users': %v", err)
		}
		if _, ok := columnExists[name]; ok {
			columnExists[name] = true
		}
	}
	rows.Close()
	if !columnExists["role"] {
		log.Println("⚠️ Migrando 'role' en 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user'")
		if err != nil {
			log.Fatalf("Error migrar 'role': %v", err)
		}
	}
	if !columnExists["nombre"] {
		log.Println("⚠️ Migrando 'nombre' en 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN nombre TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error migrar 'nombre': %v", err)
		}
	}
	if !columnExists["apellido"] {
		log.Println("⚠️ Migrando 'apellido' en 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN apellido TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Fatalf("Error migrar 'apellido': %v", err)
		}
	}
	if !columnExists["proyecto_id"] {
		log.Println("⚠️ Migrando 'proyecto_id' en 'users'...")
		_, err = DB.Exec("ALTER TABLE users ADD COLUMN proyecto_id INTEGER REFERENCES proyectos(id) ON DELETE SET NULL")
		if err != nil {
			log.Fatalf("Error migrar 'proyecto_id': %v", err)
		}
	}
}

func createDefaultAdmin() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		log.Fatalf("Error check admin: %v", err)
	}
	if count == 0 {
		adminPassword := "admin123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		insertAdminSQL := "INSERT INTO users(username, password, role, nombre, apellido) VALUES(?, ?, ?, ?, ?)"
		_, err = DB.Exec(insertAdminSQL, "admin", string(hashedPassword), "admin", "Admin", "User")
		if err != nil {
			log.Fatalf("Error insert admin: %v", err)
		}
		log.Println("✅ Usuario admin creado.")
	}
}

func GetUserByUsername(username string) (*models.UserDB, error) {
	var userDB models.UserDB
	query := "SELECT id, username, password, role FROM users WHERE username = ?"
	row := DB.QueryRow(query, username)
	err := row.Scan(&userDB.ID, &userDB.Username, &userDB.HashedPassword, &userDB.Role)
	if err != nil {
		return nil, err
	}
	return &userDB, nil
}

func CreateUser(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al hashear password: %w", err)
	}
	query := "INSERT INTO users(username, password, role, nombre, apellido) VALUES(?, ?, ?, ?, ?)"
	_, err = DB.Exec(query, user.Username, string(hashedPassword), "user", user.Nombre, user.Apellido)
	if err != nil {
		return fmt.Errorf("error al insertar usuario: %w", err)
	}
	return nil
}

func GetAllUsersWithProjects() ([]models.UserListResponse, error) {
	query := `SELECT u.id, u.username, u.role, u.nombre, u.apellido, p.id, p.nombre FROM users u LEFT JOIN proyectos p ON u.proyecto_id = p.id ORDER BY u.id ASC`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar usuarios: %w", err)
	}
	defer rows.Close()
	users := []models.UserListResponse{}
	for rows.Next() {
		var u models.UserListResponse
		var pID sql.NullInt64
		var pNombre sql.NullString
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.Nombre, &u.Apellido, &pID, &pNombre); err != nil {
			log.Printf("Error al escanear usuario %s: %v", u.Username, err)
			continue
		}
		if pID.Valid {
			id := int(pID.Int64)
			u.ProyectoID = &id
		}
		if pNombre.Valid {
			u.ProyectoNombre = &pNombre.String
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserRole(username string) (string, error) {
	var role string
	err := DB.QueryRow("SELECT role FROM users WHERE username = ?", username).Scan(&role)
	if err != nil {
		return "", err
	}
	return role, nil
}

func DeleteUserByID(id int) (int64, error) {
	query := "DELETE FROM users WHERE id = ?"
	result, err := DB.Exec(query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func UpdateUserRole(id int, newRole string) (int64, error) {
	query := "UPDATE users SET role = ? WHERE id = ?"
	result, err := DB.Exec(query, newRole, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func AssignProjectToUser(userID int, projectID int) (int64, error) {
	var idToSet interface{}
	if projectID == 0 {
		idToSet = nil
	} else {
		idToSet = projectID
	}
	query := "UPDATE users SET proyecto_id = ? WHERE id = ?"
	result, err := DB.Exec(query, idToSet, userID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func GetAllProjects() ([]models.Proyecto, error) {
	rows, err := DB.Query("SELECT id, nombre, fecha_inicio, fecha_cierre, estado FROM proyectos ORDER BY nombre ASC")
	if err != nil {
		return nil, fmt.Errorf("error al consultar proyectos: %w", err)
	}
	defer rows.Close()
	proyectos := []models.Proyecto{}
	for rows.Next() {
		var p models.Proyecto
		if err := rows.Scan(&p.ID, &p.Nombre, &p.FechaInicio, &p.FechaCierre, &p.Estado); err != nil {
			log.Printf("Error al escanear proyecto: %v", err)
			continue
		}
		proyectos = append(proyectos, p)
	}
	return proyectos, nil
}

func CreateProject(p models.CreateProyectoRequest) error {
	query := "INSERT INTO proyectos(nombre, fecha_inicio, fecha_cierre) VALUES(?, ?, ?)"
	_, err := DB.Exec(query, p.Nombre, p.FechaInicio, p.FechaCierre)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("el nombre del proyecto ya existe")
		}
		return fmt.Errorf("error al crear proyecto: %w", err)
	}
	return nil
}

func DeleteProjectByID(id int) (int64, error) {
	query := "DELETE FROM proyectos WHERE id = ?"
	result, err := DB.Exec(query, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func GetProjectByID(id int64) (*models.Proyecto, error) {
	var p models.Proyecto
	query := "SELECT id, nombre, fecha_inicio, fecha_cierre, estado FROM proyectos WHERE id = ?"
	err := DB.QueryRow(query, id).Scan(&p.ID, &p.Nombre, &p.FechaInicio, &p.FechaCierre, &p.Estado)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func UpdateProject(p models.UpdateProyectoRequest) error {
	query := "UPDATE proyectos SET nombre = ?, fecha_inicio = ?, fecha_cierre = ? WHERE id = ?"
	_, err := DB.Exec(query, p.Nombre, p.FechaInicio, p.FechaCierre, p.ID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("ese nombre de proyecto ya existe")
		}
		return fmt.Errorf("error al actualizar proyecto: %w", err)
	}
	return nil
}

func SetProjectState(id int, newState string) (int64, error) {
	query := "UPDATE proyectos SET estado = ? WHERE id = ?"
	result, err := DB.Exec(query, newState, id)
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
				log.Printf("Error al escanear miembro de proyecto: %v", err)
				continue
			}
			if strings.ToLower(member.Role) == "gerente" {
				response.Gerentes = append(response.Gerentes, member)
			} else {
				response.Miembros = append(response.Miembros, member)
			}
		}
	}
	return response, nil
}
