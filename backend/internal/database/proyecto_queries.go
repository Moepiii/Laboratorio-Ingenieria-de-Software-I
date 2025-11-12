package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"proyecto/internal/models"
)

// --- QUERIES DE PROYECTOS ---

func GetAllProyectos() ([]models.Proyecto, error) {
	rows, err := DB.Query("SELECT id, nombre, fecha_inicio, fecha_cierre, estado, fecha_creacion FROM proyectos ORDER BY id ASC")
	if err != nil {
		log.Printf("Error en GetAllProyectos (Query): %v", err)
		return nil, err
	}
	defer rows.Close()

	var proyectos []models.Proyecto
	for rows.Next() {
		var p models.Proyecto
		if err := rows.Scan(&p.ID, &p.Nombre, &p.FechaInicio, &p.FechaCierre, &p.Estado, &p.FechaCreacion); err != nil {
			log.Printf("Error en GetAllProyectos (Scan): %v", err)
			continue
		}
		proyectos = append(proyectos, p)
	}
	return proyectos, nil
}

func GetProjectByID(id int64) (*models.Proyecto, error) {
	row := DB.QueryRow("SELECT id, nombre, fecha_inicio, fecha_cierre, estado, fecha_creacion FROM proyectos WHERE id = ?", id)
	var p models.Proyecto
	err := row.Scan(&p.ID, &p.Nombre, &p.FechaInicio, &p.FechaCierre, &p.Estado, &p.FechaCreacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Proyecto no encontrado.")
		}
		log.Printf("Error al escanear proyecto (GetProjectByID): %v", err)
		return nil, fmt.Errorf("Error al buscar proyecto: %w", err)
	}
	return &p, nil
}

func CreateProyecto(nombre, fechaInicio, fechaCierre string) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO proyectos (nombre, fecha_inicio, fecha_cierre) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error al preparar inserción (CreateProyecto): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(nombre, fechaInicio, fechaCierre)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: proyectos.nombre") {
			return 0, errors.New("El nombre del proyecto ya existe.")
		}
		return 0, fmt.Errorf("error al ejecutar inserción (CreateProyecto): %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener último ID (CreateProyecto): %w", err)
	}
	return id, nil
}

func UpdateProyecto(id int, nombre, fechaInicio, fechaCierre string) (int64, error) {
	stmt, err := DB.Prepare("UPDATE proyectos SET nombre = ?, fecha_inicio = ?, fecha_cierre = ? WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error al preparar update (UpdateProyecto): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(nombre, fechaInicio, fechaCierre, id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: proyectos.nombre") {
			return 0, errors.New("El nombre del proyecto ya existe.")
		}
		return 0, fmt.Errorf("error al ejecutar update (UpdateProyecto): %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al obtener filas afectadas (UpdateProyecto): %w", err)
	}
	return affected, nil
}

func DeleteProyecto(id int) (int64, error) {
	stmt, err := DB.Prepare("DELETE FROM proyectos WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error al preparar delete (DeleteProyecto): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("error al ejecutar delete (DeleteProyecto): %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al obtener filas afectadas (DeleteProyecto): %w", err)
	}
	return affected, nil
}

func SetProyectoEstado(id int, estado string) (int64, error) {
	stmt, err := DB.Prepare("UPDATE proyectos SET estado = ? WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error al preparar update (SetProyectoEstado): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(estado, id)
	if err != nil {
		return 0, fmt.Errorf("error al ejecutar update (SetProyectoEstado): %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al obtener filas afectadas (SetProyectoEstado): %w", err)
	}
	return affected, nil
}
