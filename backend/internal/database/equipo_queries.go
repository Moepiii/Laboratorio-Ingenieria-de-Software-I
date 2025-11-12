package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"proyecto/internal/models"
)

// --- QUERIES DE EQUIPOS E IMPLEMENTOS ---

func GetEquiposByProyectoID(proyectoID int) ([]models.EquipoImplemento, error) {
	rows, err := DB.Query("SELECT id, proyecto_id, codigo_equipo, nombre, tipo, estado, fecha_creacion FROM equipos_implementos WHERE proyecto_id = ? ORDER BY id ASC", proyectoID)
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

func GetEquipoByID(id int) (*models.EquipoImplemento, error) {
	row := DB.QueryRow("SELECT id, proyecto_id, codigo_equipo, nombre, tipo, estado, fecha_creacion FROM equipos_implementos WHERE id = ?", id)
	var e models.EquipoImplemento
	err := row.Scan(&e.ID, &e.ProyectoID, &e.CodigoEquipo, &e.Nombre, &e.Tipo, &e.Estado, &e.FechaCreacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Equipo no encontrado.")
		}
		log.Printf("Error al escanear equipo (GetEquipoByID): %v", err)
		return nil, fmt.Errorf("Error al buscar equipo: %w", err)
	}
	return &e, nil
}

func CreateEquipo(equipo models.EquipoImplemento) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO equipos_implementos (proyecto_id, codigo_equipo, nombre, tipo, estado) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error al preparar inserción (CreateEquipo): %w", err)
	}
	defer stmt.Close()

	// ⭐️ MODIFICADO: Se añade el 'estado' por defecto si no viene
	estado := equipo.Estado
	if estado == "" {
		estado = "Activo"
	}

	res, err := stmt.Exec(equipo.ProyectoID, equipo.CodigoEquipo, equipo.Nombre, equipo.Tipo, estado)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: equipos_implementos.proyecto_id, equipos_implementos.codigo_equipo") {
			return 0, errors.New("El código de equipo ya existe para este proyecto.")
		}
		return 0, fmt.Errorf("error al ejecutar inserción (CreateEquipo): %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener último ID (CreateEquipo): %w", err)
	}
	return id, nil
}

func UpdateEquipo(id int, codigo, nombre, tipo, estado string) (int64, error) {
	stmt, err := DB.Prepare("UPDATE equipos_implementos SET codigo_equipo = ?, nombre = ?, tipo = ?, estado = ? WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error al preparar update (UpdateEquipo): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(codigo, nombre, tipo, estado, id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: equipos_implementos.proyecto_id, equipos_implementos.codigo_equipo") {
			return 0, errors.New("El código de equipo ya existe para este proyecto.")
		}
		return 0, fmt.Errorf("error al ejecutar update (UpdateEquipo): %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al obtener filas afectadas (UpdateEquipo): %w", err)
	}
	return affected, nil
}

func DeleteEquipo(id int) (int64, error) {
	stmt, err := DB.Prepare("DELETE FROM equipos_implementos WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error al preparar delete (DeleteEquipo): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("error al ejecutar delete (DeleteEquipo): %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al obtener filas afectadas (DeleteEquipo): %w", err)
	}
	return affected, nil
}
