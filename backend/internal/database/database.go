package database

import (
	"database/sql"
	"log"

	// 'time' se mantiene por la creación de tablas
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dbPath string) {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Error al abrir DB: %v", err)
	}

	// Habilita el modo WAL (Write-Ahead Logging)
	_, err = DB.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		log.Fatalf("Error al habilitar modo WAL: %v", err)
	}
	log.Println("✅ Modo WAL habilitado en SQLite.")

	_, err = DB.Exec("PRAGMA foreign_keys = OFF;")
	if err != nil {
		log.Fatalf("Error PRAGMA OFF: %v", err)
	}
	createUsersTable()
	createProyectosTable()
	createLaboresTable()
	createEquiposTable()
	createActividadesTable()

	// ⭐️ 1. LLAMADA A LA NUEVA FUNCIÓN
	createEventLogsTable()

	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalf("Error PRAGMA ON: %v", err)
	}
}

// --- CREACIÓN DE TABLAS ---

func createUsersTable() {
	_, err := DB.Exec(`
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        role TEXT NOT NULL DEFAULT 'user',
        nombre TEXT NOT NULL,
        apellido TEXT NOT NULL,
        cedula TEXT NOT NULL UNIQUE,
        proyecto_id INTEGER,
        FOREIGN KEY (proyecto_id) REFERENCES proyectos(id) ON DELETE SET NULL
    );
    `)
	if err != nil {
		log.Fatalf("Error al crear tabla users: %v", err)
	}

	// Crear usuario admin si no existe
	row := DB.QueryRow("SELECT id FROM users WHERE username = 'admin'")
	var id int
	if err := row.Scan(&id); err == sql.ErrNoRows {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Error al hashear password de admin: %v", err)
		}
		_, err = DB.Exec("INSERT INTO users (username, password, role, nombre, apellido, cedula) VALUES (?, ?, ?, ?, ?, ?)",
			"admin", string(hashedPassword), "admin", "Administrador", "Del Sistema", "000000")
		if err != nil {
			log.Fatalf("Error al crear usuario admin: %v", err)
		}
		log.Println("Usuario 'admin' (pass: 'admin123') creado.")
	}
}

func createProyectosTable() {
	_, err := DB.Exec(`
    CREATE TABLE IF NOT EXISTS proyectos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        nombre TEXT NOT NULL UNIQUE,
        fecha_inicio TEXT NOT NULL,
        fecha_cierre TEXT NOT NULL,
        estado TEXT NOT NULL DEFAULT 'Activo',
        fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `)
	if err != nil {
		log.Fatalf("Error al crear tabla proyectos: %v", err)
	}
}

func createLaboresTable() {
	_, err := DB.Exec(`
    CREATE TABLE IF NOT EXISTS labores_agronomicas (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        proyecto_id INTEGER NOT NULL,
        codigo_labor TEXT NOT NULL,
        descripcion TEXT NOT NULL,
        estado TEXT NOT NULL DEFAULT 'Activo',
        fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (proyecto_id) REFERENCES proyectos(id) ON DELETE CASCADE,
        UNIQUE(proyecto_id, codigo_labor)
    );
    `)
	if err != nil {
		log.Fatalf("Error al crear tabla labores_agronomicas: %v", err)
	}
}

func createEquiposTable() {
	_, err := DB.Exec(`
    CREATE TABLE IF NOT EXISTS equipos_implementos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        proyecto_id INTEGER NOT NULL,
        codigo_equipo TEXT NOT NULL,
        nombre TEXT NOT NULL,
        tipo TEXT NOT NULL CHECK (tipo IN ('Equipo', 'Implemento')),
        estado TEXT NOT NULL DEFAULT 'Activo',
        fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (proyecto_id) REFERENCES proyectos(id) ON DELETE CASCADE,
        UNIQUE(proyecto_id, codigo_equipo)
    );
    `)
	if err != nil {
		log.Fatalf("Error al crear tabla equipos_implementos: %v", err)
	}
}

func createActividadesTable() {
	_, err := DB.Exec(`
    CREATE TABLE IF NOT EXISTS actividades (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        proyecto_id INTEGER NOT NULL,
        actividad TEXT NOT NULL,
        labor_agronomica_id INTEGER,
        equipo_implemento_id INTEGER,
        encargado_id INTEGER,
        recurso_humano INTEGER NOT NULL,
        costo REAL NOT NULL,
        observaciones TEXT,
        fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (proyecto_id) REFERENCES proyectos(id) ON DELETE CASCADE,
        FOREIGN KEY (labor_agronomica_id) REFERENCES labores_agronomicas(id) ON DELETE SET NULL,
        FOREIGN KEY (equipo_implemento_id) REFERENCES equipos_implementos(id) ON DELETE SET NULL,
        FOREIGN KEY (encargado_id) REFERENCES users(id) ON DELETE SET NULL
    );
    `)
	if err != nil {
		log.Fatalf("Error al crear tabla actividades: %v", err)
	}
}

// ⭐️ 2. FUNCIÓN DE LA NUEVA TABLA AÑADIDA (AL FINAL DEL ARCHIVO)
func createEventLogsTable() {
	_, err := DB.Exec(`
    CREATE TABLE IF NOT EXISTS event_logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        usuario_username TEXT,
        usuario_rol TEXT,
        accion TEXT,
        entidad TEXT,
        entidad_id INTEGER
    );
    `)
	if err != nil {
		log.Fatalf("Error al crear tabla event_logs: %v", err)
	}
}
