package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"proyecto/internal/models"
)

// --- QUERIES DE LABORES ---

func GetLaboresByProyectoID(proyectoID int) ([]models.LaborAgronomica, error) {
	rows, err := DB.Query("SELECT id, proyecto_id, codigo_labor, descripcion, estado, fecha_creacion FROM labores_agronomicas WHERE proyecto_id = ? ORDER BY id ASC", proyectoID)
	if err != nil {
		log.Printf("Error en GetLaboresByProyectoID (Query): %v", err)
		return nil, err
	}
	defer rows.Close()

	var labores []models.LaborAgronomica
	for rows.Next() {
		var l models.LaborAgronomica
		if err := rows.Scan(&l.ID, &l.ProyectoID, &l.CodigoLabor, &l.Descripcion, &l.Estado, &l.FechaCreacion); err != nil {
			log.Printf("Error en GetLaboresByProyectoID (Scan): %v", err)
			continue
		}
		labores = append(labores, l)
	}
	return labores, nil
}

func GetLaborByID(id int) (*models.LaborAgronomica, error) {
	row := DB.QueryRow("SELECT id, proyecto_id, codigo_labor, descripcion, estado, fecha_creacion FROM labores_agronomicas WHERE id = ?", id)
	var l models.LaborAgronomica
	err := row.Scan(&l.ID, &l.ProyectoID, &l.CodigoLabor, &l.Descripcion, &l.Estado, &l.FechaCreacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Labor no encontrada.")
		}
		log.Printf("Error al escanear labor (GetLaborByID): %v", err)
		return nil, fmt.Errorf("Error al buscar labor: %w", err)
	}
	return &l, nil
}

func CreateLabor(labor models.LaborAgronomica) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO labores_agronomicas (proyecto_id, codigo_labor, descripcion, estado) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("error al preparar inserción (CreateLabor): %w", err)
	}
	defer stmt.Close()

	// ⭐️ MODIFICADO: Se añade el 'estado' por defecto si no viene
	estado := labor.Estado
	if estado == "" {
		estado = "Activo"
	}

	res, err := stmt.Exec(labor.ProyectoID, labor.CodigoLabor, labor.Descripcion, estado)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: labores_agronomicas.proyecto_id, labores_agronomicas.codigo_labor") {
			return 0, errors.New("El código de labor ya existe para este proyecto.")
		}
		return 0, fmt.Errorf("error al ejecutar inserción (CreateLabor): %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener último ID (CreateLabor): %w", err)
	}
	return id, nil
}

func UpdateLabor(id int, codigo, descripcion, estado string) (int64, error) {
	stmt, err := DB.Prepare("UPDATE labores_agronomicas SET codigo_labor = ?, descripcion = ?, estado = ? WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error al preparar update (UpdateLabor): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(codigo, descripcion, estado, id)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: labores_agronomicas.proyecto_id, labores_agronomicas.codigo_labor") {
			return 0, errors.New("El código de labor ya existe para este proyecto.")
		}
		return 0, fmt.Errorf("error al ejecutar update (UpdateLabor): %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al obtener filas afectadas (UpdateLabor): %w", err)
	}
	return affected, nil
}

func DeleteLabor(id int) (int64, error) {
	stmt, err := DB.Prepare("DELETE FROM labores_agronomicas WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error al preparar delete (DeleteLabor): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("error al ejecutar delete (DeleteLabor): %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error al obtener filas afectadas (DeleteLabor): %w", err)
	}
	return affected, nil
}
