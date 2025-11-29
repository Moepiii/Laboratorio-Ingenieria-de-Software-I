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

	query.WriteString("SELECT id, timestamp, usuario_username, usuario_rol, accion, entidad, entidad_id FROM event_logs WHERE 1=1")

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
		query.WriteString(" AND date(timestamp) >= date(?)")
		args = append(args, filtros.FechaInicio)
	}
	if filtros.FechaCierre != "" {
		query.WriteString(" AND date(timestamp) <= date(?)")
		args = append(args, filtros.FechaCierre)
	}

	query.WriteString(" ORDER BY timestamp DESC")
	query.WriteString(" LIMIT 1000")

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

// ⭐️ ESTA ES LA FUNCIÓN QUE FALTABA ⭐️
// DeleteLog elimina un log específico por ID
func DeleteLog(id int) error {
	stmt, err := DB.Prepare("DELETE FROM event_logs WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// DeleteLogsByRange elimina logs dentro de un rango de fechas (inclusivo).
func DeleteLogsByRange(fechaInicio, fechaFin string) (int64, error) {
	query := `DELETE FROM event_logs WHERE date(timestamp) >= date(?) AND date(timestamp) <= date(?)`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(fechaInicio, fechaFin)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
