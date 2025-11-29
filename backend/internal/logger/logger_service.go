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

	// DeleteLogs elimina una lista de eventos por sus IDs (Ya lo tenías)
	DeleteLogs(ids []int) error

	// ⭐️ NUEVO: DeleteLogsByRange elimina eventos en un rango de fechas
	DeleteLogsByRange(fechaInicio, fechaFin string) (int64, error)
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

// DeleteLogs: Elimina múltiples logs por sus IDs (Lógica que ya tenías)
func (s *loggerService) DeleteLogs(ids []int) error {
	if len(ids) == 0 {
		return errors.New("no se enviaron IDs para eliminar")
	}

	// Como SQLite no tiene un "DELETE WHERE ID IN (...)" nativo fácil en Go sin armar string manual,
	// iteramos y borramos uno por uno (o asumimos que tienes una función para ello).
	// Aquí usamos una implementación segura iterando:
	for _, id := range ids {
		// Llamamos a la función de DB para borrar individualmente
		// (Asumiendo que ya tienes database.DeleteLog o similar del paso anterior)
		err := database.DeleteLog(id)
		if err != nil {
			log.Printf("Error borrando log ID %d: %v", id, err)
			// Podríamos retornar error o continuar 'best effort'
			return err
		}
	}
	return nil
}

// ⭐️ NUEVO MÉTODO IMPLEMENTADO PARA EL PASO 2 ⭐️
func (s *loggerService) DeleteLogsByRange(fechaInicio, fechaFin string) (int64, error) {
	// Validación básica de negocio
	if fechaInicio == "" || fechaFin == "" {
		return 0, errors.New("las fechas de inicio y fin son requeridas")
	}

	// Llamada a la capa de datos (la función que creamos en el Paso 1)
	return database.DeleteLogsByRange(fechaInicio, fechaFin)
}
