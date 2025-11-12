package logger

import (
	"errors"
	"log"

	"proyecto/internal/database"
	"proyecto/internal/models"
)

// --- 1. EL CONTRATO (Interface) ---
type LoggerService interface {
	// Log escribe un evento en la base de datos
	Log(usuarioUsername string, usuarioRol string, accion string, entidad string, entidadID int)

	// GetLogs obtiene los eventos con filtros
	GetLogs(filtros models.GetLogsRequest) ([]models.EventLogResponse, error)
}

// --- 2. LA IMPLEMENTACIÓN (Struct) ---
type loggerService struct {
	// (Dependencias futuras)
}

// --- 3. EL CONSTRUCTOR ---
func NewLoggerService() LoggerService {
	return &loggerService{}
}

// --- 4. LOS MÉTODOS (Lógica de Negocio) ---

func (s *loggerService) Log(usuarioUsername string, usuarioRol string, accion string, entidad string, entidadID int) {

	logEntry := models.EventLog{
		UsuarioUsername: usuarioUsername,
		UsuarioRol:      usuarioRol,
		Accion:          accion,
		Entidad:         entidad,
		EntidadID:       entidadID,
	}

	// ⭐️ Usamos una goroutine (go func()) aquí.
	// Esto hace que el registro del log se ejecute en segundo plano.
	// Así, si el guardado del log falla, NO detiene la acción principal del usuario
	// (como crear un proyecto).
	go func() {
		_, err := database.InsertLog(logEntry)
		if err != nil {
			// Si falla el log, solo lo registramos en la consola del servidor
			log.Printf("ERROR: No se pudo guardar el evento de log: %v", err)
		}
	}()
}

func (s *loggerService) GetLogs(filtros models.GetLogsRequest) ([]models.EventLogResponse, error) {
	// ⭐️ ST1005: Corregimos los mensajes de error
	if filtros.AdminUsername == "" {
		return nil, errors.New("nombre de usuario admin requerido")
	}

	logs, err := database.GetLogs(filtros)
	if err != nil {
		log.Printf("Error en loggerService.GetLogs: %v", err)
		return nil, errors.New("error al obtener logs")
	}

	return logs, nil
}
