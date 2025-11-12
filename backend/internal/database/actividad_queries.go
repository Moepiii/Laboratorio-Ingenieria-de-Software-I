package database

import (
	"fmt"
	"log"

	"proyecto/internal/models"
)

// --- QUERIES DE ACTIVIDADES ---

func CreateActividad(act models.Actividad) (int64, error) {
	stmt, err := DB.Prepare(`
		INSERT INTO actividades (
			proyecto_id, actividad, labor_agronomica_id, equipo_implemento_id, 
			encargado_id, recurso_humano, costo, observaciones
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`)
	if err != nil {
		return 0, fmt.Errorf("error al preparar inserción (CreateActividad): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		act.ProyectoID, act.Actividad, act.LaborAgronomicaID, act.EquipoImplementoID,
		act.EncargadoID, act.RecursoHumano, act.Costo, act.Observaciones,
	)
	if err != nil {
		return 0, fmt.Errorf("error al ejecutar inserción (CreateActividad): %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener último ID (CreateActividad): %w", err)
	}
	return id, nil
}

// GetActividadesByProyectoID (Función compleja que trae todo)
func GetActividadesByProyectoID(proyectoID int) ([]models.ActividadResponse, error) {
	// ⭐️ MODIFICADO: Se usan LEFT JOINs para que no se rompa si un ID es NULL
	// ⭐️ MODIFICADO: Se usa COALESCE para que los campos NULL devuelvan un string vacío "" en lugar de NULL
	// ⭐️ MODIFICADO: Se usa u.nombre || ' ' || u.apellido para el nombre del encargado
	query := `
		SELECT 
			a.id, a.proyecto_id, a.actividad, a.labor_agronomica_id, a.equipo_implemento_id,
			a.encargado_id, a.recurso_humano, a.costo, a.observaciones, a.fecha_creacion,
			
			COALESCE(l.descripcion, '') AS labor_descripcion,
			COALESCE(e.nombre, '') AS equipo_nombre,
			COALESCE(u.nombre || ' ' || u.apellido, '') AS encargado_nombre

		FROM actividades a
		LEFT JOIN labores_agronomicas l ON a.labor_agronomica_id = l.id
		LEFT JOIN equipos_implementos e ON a.equipo_implemento_id = e.id
		LEFT JOIN users u ON a.encargado_id = u.id
		WHERE a.proyecto_id = ?
		ORDER BY a.id ASC;
	`

	rows, err := DB.Query(query, proyectoID)
	if err != nil {
		log.Printf("Error en GetActividadesByProyectoID (Query): %v", err)
		return nil, err
	}
	defer rows.Close()

	var actividades []models.ActividadResponse
	for rows.Next() {
		var act models.ActividadResponse
		if err := rows.Scan(
			&act.ID, &act.ProyectoID, &act.Actividad, &act.LaborAgronomicaID, &act.EquipoImplementoID,
			&act.EncargadoID, &act.RecursoHumano, &act.Costo, &act.Observaciones, &act.FechaCreacion,
			&act.LaborDescripcion, &act.EquipoNombre, &act.EncargadoNombre,
		); err != nil {
			log.Printf("Error escaneando actividad: %v", err)
			continue
		}
		actividades = append(actividades, act)
	}
	return actividades, nil
}

func UpdateActividad(act models.Actividad) (int64, error) {
	stmt, err := DB.Prepare(`
		UPDATE actividades SET
			actividad = ?, labor_agronomica_id = ?, equipo_implemento_id = ?, 
			encargado_id = ?, recurso_humano = ?, costo = ?, observaciones = ?
		WHERE id = ? AND proyecto_id = ?`)
	if err != nil {
		return 0, fmt.Errorf("error preparando update (UpdateActividad): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		act.Actividad, act.LaborAgronomicaID, act.EquipoImplementoID,
		act.EncargadoID, act.RecursoHumano, act.Costo, act.Observaciones,
		act.ID, act.ProyectoID,
	)
	if err != nil {
		return 0, fmt.Errorf("error ejecutando update (UpdateActividad): %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error obteniendo filas afectadas (UpdateActividad): %w", err)
	}
	return affected, nil
}

func DeleteActividad(id int) (int64, error) {
	stmt, err := DB.Prepare("DELETE FROM actividades WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("error preparando delete (DeleteActividad): %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, fmt.Errorf("error ejecutando delete (DeleteActividad): %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error obteniendo filas afectadas (DeleteActividad): %w", err)
	}
	return affected, nil
}
