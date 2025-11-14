package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"proyecto/internal/models"
)

// --- QUERIES DE EQUIPOS E IMPLEMENTOS ---

// GetEquiposByProyectoID obtiene todos los equipos de un proyecto
func GetEquiposByProyectoID(proyectoID int) ([]models.EquipoImplemento, error) {
	query := `
        SELECT id, proyecto_id, codigo_equipo, nombre, tipo, estado, fecha_creacion 
        FROM equipos_implementos 
        WHERE proyecto_id = ? 
        ORDER BY fecha_creacion DESC
    `
	rows, err := DB.Query(query, proyectoID)
	if err != nil {
		log.Printf("Error en GetEquiposByProyectoID (Query): %v", err)
		return nil, err
	}
	defer rows.Close()

	var equipos []models.EquipoImplemento
	for rows.Next() {
		var e models.EquipoImplemento
		if err := rows.Scan(&e.ID, &e.ProyectoID, &e.CodigoEquipo, &e.Nombre, &e.Tipo, &e.Estado, &e.FechaCreacion); err != nil {
			log.Printf("Error en GetEquiposByProyectoID (Scan): %v", err)
			continue
		}
		equipos = append(equipos, e)
	}

	return equipos, nil
}

// GetEquipoByID obtiene un equipo específico por su ID
func GetEquipoByID(id int) (*models.EquipoImplemento, error) {
	query := `
        SELECT id, proyecto_id, codigo_equipo, nombre, tipo, estado, fecha_creacion 
        FROM equipos_implementos 
        WHERE id = ?
    `
	row := DB.QueryRow(query, id)
	var e models.EquipoImplemento
	err := row.Scan(&e.ID, &e.ProyectoID, &e.CodigoEquipo, &e.Nombre, &e.Tipo, &e.Estado, &e.FechaCreacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("equipo no encontrado")
		}
		log.Printf("Error en GetEquipoByID: %v", err)
		return nil, err
	}
	return &e, nil
}

// CreateEquipo inserta un nuevo equipo en la DB
func CreateEquipo(equipo models.EquipoImplemento) (int64, error) {
	// Comprobación de unicidad para (proyecto_id, codigo_equipo)
	var exists int
	err := DB.QueryRow("SELECT COUNT(*) FROM equipos_implementos WHERE proyecto_id = ? AND codigo_equipo = ?", equipo.ProyectoID, equipo.CodigoEquipo).Scan(&exists)
	if err != nil {
		log.Printf("Error chequeando unicidad de equipo: %v", err)
		return 0, err
	}
	if exists > 0 {
		return 0, errors.New("el código de equipo ya existe para este proyecto")
	}

	// Inserción
	stmt, err := DB.Prepare(`
        INSERT INTO equipos_implementos 
        (proyecto_id, codigo_equipo, nombre, tipo, estado, fecha_creacion) 
        VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
    `)
	if err != nil {
		log.Printf("Error en CreateEquipo (Prepare): %v", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(equipo.ProyectoID, equipo.CodigoEquipo, equipo.Nombre, equipo.Tipo, equipo.Estado)
	if err != nil {
		log.Printf("Error en CreateEquipo (Exec): %v", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, errors.New("el código de equipo ya existe para este proyecto")
		}
		return 0, err
	}

	return res.LastInsertId()
}

// UpdateEquipo actualiza un equipo existente
func UpdateEquipo(id int, codigoEquipo, nombre, tipo, estado string) (int64, error) {
	stmt, err := DB.Prepare(`
        UPDATE equipos_implementos 
        SET codigo_equipo = ?, nombre = ?, tipo = ?, estado = ?
        WHERE id = ?
    `)
	if err != nil {
		log.Printf("Error en UpdateEquipo (Prepare): %v", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(codigoEquipo, nombre, tipo, estado, id)
	if err != nil {
		log.Printf("Error en UpdateEquipo (Exec): %v", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, errors.New("el código de equipo ya existe para este proyecto")
		}
		return 0, err
	}

	return res.RowsAffected()
}

// DeleteEquipo borra un equipo de la DB
func DeleteEquipo(id int) (int64, error) {
	stmt, err := DB.Prepare("DELETE FROM equipos_implementos WHERE id = ?")
	if err != nil {
		log.Printf("Error en DeleteEquipo (Prepare): %v", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		log.Printf("Error en DeleteEquipo (Exec): %v", err)
		return 0, err
	}

	return res.RowsAffected()
}

// ⭐️ --- INICIO: NUEVA FUNCIÓN AÑADIDA --- ⭐️

// GetNextEquipoCodigo calcula el siguiente código secuencial para un proyecto.
// Trata el 'codigo_equipo' como un número.
func GetNextEquipoCodigo(proyectoID int) (int, error) {
	var nextCodigo int

	// Esta consulta es idéntica a la de Labores, pero
	// apunta a la tabla 'equipos_implementos' y 'codigo_equipo'
	query := `
		SELECT IFNULL(MAX(CAST(codigo_equipo AS INTEGER)), 0) + 1 
		FROM equipos_implementos 
		WHERE proyecto_id = ?;
	`

	err := DB.QueryRow(query, proyectoID).Scan(&nextCodigo)
	if err != nil {
		log.Printf("Error en GetNextEquipoCodigo: %v", err)
		return 0, err
	}

	return nextCodigo, nil
}

// ⭐️ --- FIN: NUEVA FUNCIÓN AÑADIDA --- ⭐️
