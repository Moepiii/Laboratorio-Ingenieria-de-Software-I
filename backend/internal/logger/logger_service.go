package logger

import (
	"errors"
	"log"
	"proyecto/internal/database"
	"proyecto/internal/models"
)

// --- 1. EL CONTRATO (Interface) ---
type LoggerService interface {
	// Log escribe un evento en la base de datos de forma asíncrona
	Log(usuarioUsername string, usuarioRol string, accion string, entidad string, entidadID int)

	// GetLogs obtiene los eventos con filtros
	GetLogs(filtros models.GetLogsRequest) ([]models.EventLogResponse, error)

	// DeleteLogs elimina una lista de eventos por sus IDs (⭐️ NUEVO)
	DeleteLogs(ids []int) error
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

// Log: Registra un evento en segundo plano
func (s *loggerService) Log(usuarioUsername string, usuarioRol string, accion string, entidad string, entidadID int) {

	logEntry := models.EventLog{
		UsuarioUsername: usuarioUsername,
		UsuarioRol:      usuarioRol,
		Accion:          accion,
		Entidad:         entidad,
		EntidadID:       entidadID,
	}

	// Usamos una goroutine (go func()) para que el log no detenga la operación principal
	go func() {
		_, err := database.InsertLog(logEntry)
		if err != nil {
			// Si falla el log, solo lo mostramos en consola, no rompemos el flujo del usuario
			log.Printf("ERROR CRÍTICO: No se pudo guardar el evento de log en DB: %v", err)
		}
	}()
}

// GetLogs: Obtiene logs filtrados
func (s *loggerService) GetLogs(filtros models.GetLogsRequest) ([]models.EventLogResponse, error) {
	return database.GetLogs(filtros)
}

// DeleteLogs: Elimina múltiples logs (⭐️ NUEVO)
func (s *loggerService) DeleteLogs(ids []int) error {
	if len(ids) == 0 {
		return errors.New("no se enviaron IDs para eliminar")
	}

	// Recorremos la lista de IDs y borramos uno por uno
	// Esta es una estrategia segura y simple para SQLite
	for _, id := range ids {
		// Llamamos a la función que creamos en logger_queries.go
		_, err := database.DeleteLog(id)
		if err != nil {
			log.Printf("Error borrando log ID %d: %v", id, err)
			// Si falla uno, detenemos el proceso y retornamos error
			return err
		}
	}

	return nil
}