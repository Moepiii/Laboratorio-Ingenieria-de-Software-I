package database

import (
	"log"
	"strings"

	"proyecto/internal/models"
)

// --- QUERIES DEL LOGGER ---

// InsertLog inserta un nuevo evento en la base de datos
func InsertLog(logEntry models.EventLog) (int64, error) {
	stmt, err := DB.Prepare(`
		INSERT INTO event_logs 
		(timestamp, usuario_username, usuario_rol, accion, entidad, entidad_id) 
		VALUES (strftime('%Y-%m-%d %H:%M:%S', 'now', 'localtime'), ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Printf("Error preparando InsertLog: %v", err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		logEntry.UsuarioUsername,
		logEntry.UsuarioRol,
		logEntry.Accion,
		logEntry.Entidad,
		logEntry.EntidadID,
	)
	if err != nil {
		log.Printf("Error ejecutando InsertLog: %v", err)
		return 0, err
	}

	return res.LastInsertId()
}

// GetLogs recupera los logs con filtros dinámicos
func GetLogs(filtros models.GetLogsRequest) ([]models.EventLogResponse, error) {
	var query strings.Builder
	var args []interface{}

	query.WriteString(`
		SELECT id, timestamp, usuario_username, usuario_rol, accion, entidad, entidad_id 
		FROM event_logs 
		WHERE 1=1
	`)

	// Construcción dinámica de la consulta
	if filtros.UsuarioUsername != "" {
		query.WriteString(" AND usuario_username LIKE ?")
		args = append(args, "%"+filtros.UsuarioUsername+"%")
	}
	if filtros.Accion != "" {
		query.WriteString(" AND accion = ?")
		args = append(args, filtros.Accion)
	}
	if filtros.Entidad != "" {
		query.WriteString(" AND entidad = ?")
		args = append(args, filtros.Entidad)
	}
	if filtros.FechaInicio != "" {
		// Asume formato 'YYYY-MM-DD'
		query.WriteString(" AND date(timestamp) >= date(?)")
		args = append(args, filtros.FechaInicio)
	}
	if filtros.FechaCierre != "" {
		// Asume formato 'YYYY-MM-DD'
		query.WriteString(" AND date(timestamp) <= date(?)")
		args = append(args, filtros.FechaCierre)
	}

	// Ordenar como en tu captura (los más nuevos primero)
	query.WriteString(" ORDER BY timestamp DESC")

	// Límite (opcional, pero buena idea para el rendimiento)
	query.WriteString(" LIMIT 1000")

	// Ejecutar la consulta
	rows, err := DB.Query(query.String(), args...)
	if err != nil {
		log.Printf("Error en GetLogs (Query): %v", err)
		return nil, err
	}
	defer rows.Close()

	var logs []models.EventLogResponse
	for rows.Next() {
		var l models.EventLogResponse
		if err := rows.Scan(
			&l.ID,
			&l.Timestamp,
			&l.UsuarioUsername,
			&l.UsuarioRol,
			&l.Accion,
			&l.Entidad,
			&l.EntidadID,
		); err != nil {
			log.Printf("Error en GetLogs (Scan): %v", err)
			continue
		}
		logs = append(logs, l)
	}

	return logs, nil
}

func DeleteLog(id int) (int64, error) {
	res, err := DB.Exec("DELETE FROM event_logs WHERE id = ?", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}