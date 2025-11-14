package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"proyecto/internal/models"
)

// --- QUERIES DE LABORES AGRONÓMICAS ---

// GetLaboresByProyectoID obtiene todas las labores de un proyecto
func GetLaboresByProyectoID(proyectoID int) ([]models.LaborAgronomica, error) {
	query := `
        SELECT id, proyecto_id, codigo_labor, descripcion, estado, fecha_creacion 
        FROM labores_agronomicas 
        WHERE proyecto_id = ? 
        ORDER BY fecha_creacion DESC
    `
	rows, err := DB.Query(query, proyectoID)
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

// GetLaborByID obtiene una labor específica por su ID
func GetLaborByID(id int) (*models.LaborAgronomica, error) {
	query := `
        SELECT id, proyecto_id, codigo_labor, descripcion, estado, fecha_creacion 
        FROM labores_agronomicas 
        WHERE id = ?
    `
	row := DB.QueryRow(query, id)
	var l models.LaborAgronomica
	err := row.Scan(&l.ID, &l.ProyectoID, &l.CodigoLabor, &l.Descripcion, &l.Estado, &l.FechaCreacion)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("labor no encontrada")
		}
		log.Printf("Error en GetLaborByID: %v", err)
		return nil, err
	}
	return &l, nil
}

// CreateLabor inserta una nueva labor en la DB
func CreateLabor(labor models.LaborAgronomica) (int64, error) {
	// Comprobación de unicidad para (proyecto_id, codigo_labor)
	var exists int
	err := DB.QueryRow("SELECT COUNT(*) FROM labores_agronomicas WHERE proyecto_id = ? AND codigo_labor = ?", labor.ProyectoID, labor.CodigoLabor).Scan(&exists)
	if err != nil {
		log.Printf("Error chequeando unicidad de labor: %v", err)
		return 0, err
	}
	if exists > 0 {
		return 0, errors.New("el código de labor ya existe para este proyecto")
	}

	// Inserción
	stmt, err := DB.Prepare(`
        INSERT INTO labores_agronomicas 
        (proyecto_id, codigo_labor, descripcion, estado, fecha_creacion) 
        VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
    `)
	if err != nil {
		log.Printf("Error en CreateLabor (Prepare): %v", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(labor.ProyectoID, labor.CodigoLabor, labor.Descripcion, labor.Estado)
	if err != nil {
		log.Printf("Error en CreateLabor (Exec): %v", err)
		// Verificamos si es un error de unicidad (aunque ya lo chequeamos, es una buena práctica)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, errors.New("el código de labor ya existe para este proyecto")
		}
		return 0, err
	}

	return res.LastInsertId()
}

// UpdateLabor actualiza una labor existente
func UpdateLabor(id int, codigoLabor, descripcion, estado string) (int64, error) {
	// (Aquí deberíamos chequear la unicidad del código si cambia, omitido por brevedad)

	stmt, err := DB.Prepare(`
        UPDATE labores_agronomicas 
        SET codigo_labor = ?, descripcion = ?, estado = ?
        WHERE id = ?
    `)
	if err != nil {
		log.Printf("Error en UpdateLabor (Prepare): %v", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(codigoLabor, descripcion, estado, id)
	if err != nil {
		log.Printf("Error en UpdateLabor (Exec): %v", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, errors.New("el código de labor ya existe para este proyecto")
		}
		return 0, err
	}

	return res.RowsAffected()
}

// DeleteLabor borra una labor de la DB
func DeleteLabor(id int) (int64, error) {
	stmt, err := DB.Prepare("DELETE FROM labores_agronomicas WHERE id = ?")
	if err != nil {
		log.Printf("Error en DeleteLabor (Prepare): %v", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		log.Printf("Error en DeleteLabor (Exec): %v", err)
		return 0, err
	}

	return res.RowsAffected()
}

// ⭐️ --- INICIO: NUEVA FUNCIÓN AÑADIDA --- ⭐️

// GetNextLaborCodigo calcula el siguiente código secuencial para un proyecto.
// Trata el 'codigo_labor' como un número.
func GetNextLaborCodigo(proyectoID int) (int, error) {
	var nextCodigo int

	// Esta consulta:
	// 1. Busca en labores_agronomicas para un proyecto_id.
	// 2. Convierte el codigo_labor (que es TEXT) a INTEGER.
	// 3. Encuentra el MÁXIMO valor.
	// 4. Si no hay labores (IFNULL), devuelve 0.
	// 5. Le suma 1.
	query := `
		SELECT IFNULL(MAX(CAST(codigo_labor AS INTEGER)), 0) + 1 
		FROM labores_agronomicas 
		WHERE proyecto_id = ?;
	`

	err := DB.QueryRow(query, proyectoID).Scan(&nextCodigo)
	if err != nil {
		log.Printf("Error en GetNextLaborCodigo: %v", err)
		return 0, err
	}

	return nextCodigo, nil
}

// ⭐️ --- FIN: NUEVA FUNCIÓN AÑADIDA --- ⭐️
