package actividades

import (
	"database/sql"
	"errors"
	"log"

	"proyecto/internal/database"
	"proyecto/internal/models"
)

// --- 1. EL CONTRATO (Interface) ---

// GetDatosProyectoResponse es un struct para agrupar la respuesta
type GetDatosProyectoResponse struct {
	Labores     []models.LaborAgronomica   `json:"labores"`
	Equipos     []models.EquipoImplemento  `json:"equipos"`
	Encargados  []models.EncargadoResponse `json:"encargados"`
	Actividades []models.ActividadResponse `json:"actividades"`
}

type ActividadService interface {
	GetDatosProyecto(proyectoID int) (*GetDatosProyectoResponse, error)
	CreateActividad(req models.CreateActividadRequest) ([]models.ActividadResponse, error)
	UpdateActividad(req models.UpdateActividadRequest) ([]models.ActividadResponse, error)
	DeleteActividad(id int) (int64, error)
}

// --- 2. LA IMPLEMENTACIÓN (Struct) ---
type actividadService struct {
	// (Dependencias futuras)
}

// --- 3. EL CONSTRUCTOR ---
func NewActividadService() ActividadService {
	return &actividadService{}
}

// --- 4. LOS MÉTODOS (Lógica de Negocio) ---

// ⭐️ CORRECCIÓN AQUÍ ⭐️
func (s *actividadService) GetDatosProyecto(proyectoID int) (*GetDatosProyectoResponse, error) {
	labores, err := database.GetLaboresByProyectoID(proyectoID)
	if err != nil {
		log.Printf("Error en GetDatosProyecto (GetLabores): %v", err)
		return nil, errors.New("Error al obtener labores.")
	}

	equipos, err := database.GetEquiposByProyectoID(proyectoID)
	if err != nil {
		log.Printf("Error en GetDatosProyecto (GetEquipos): %v", err)
		return nil, errors.New("Error al obtener equipos.")
	}

	// El error estaba aquí. La función es 'GetEncargados' (sin ID de proyecto)
	encargados, err := database.GetEncargados()
	if err != nil {
		log.Printf("Error en GetDatosProyecto (GetEncargados): %v", err)
		return nil, errors.New("Error al obtener encargados.")
	}

	actividades, err := database.GetActividadesByProyectoID(proyectoID)
	if err != nil {
		log.Printf("Error en GetDatosProyecto (GetActividades): %v", err)
		return nil, errors.New("Error al obtener actividades.")
	}

	return &GetDatosProyectoResponse{
		Labores:     labores,
		Equipos:     equipos,
		Encargados:  encargados,
		Actividades: actividades,
	}, nil
}

func (s *actividadService) CreateActividad(req models.CreateActividadRequest) ([]models.ActividadResponse, error) {
	if req.ProyectoID == 0 || req.Actividad == "" || req.Costo == 0 || req.RecursoHumano == 0 {
		return nil, errors.New("ProyectoID, Actividad, RecursoHumano y Costo son requeridos.")
	}

	// Manejo de valores opcionales (IDs)
	var laborID, equipoID, encargadoID sql.NullInt64
	if req.LaborAgronomicaID != nil && *req.LaborAgronomicaID != 0 {
		laborID = sql.NullInt64{Int64: int64(*req.LaborAgronomicaID), Valid: true}
	}
	if req.EquipoImplementoID != nil && *req.EquipoImplementoID != 0 {
		equipoID = sql.NullInt64{Int64: int64(*req.EquipoImplementoID), Valid: true}
	}
	if req.EncargadoID != nil && *req.EncargadoID != 0 {
		encargadoID = sql.NullInt64{Int64: int64(*req.EncargadoID), Valid: true}
	}
	var observaciones sql.NullString
	if req.Observaciones != "" {
		observaciones = sql.NullString{String: req.Observaciones, Valid: true}
	}

	actividad := models.Actividad{
		ProyectoID:         req.ProyectoID,
		Actividad:          req.Actividad,
		LaborAgronomicaID:  laborID,
		EquipoImplementoID: equipoID,
		EncargadoID:        encargadoID,
		RecursoHumano:      req.RecursoHumano,
		Costo:              req.Costo,
		Observaciones:      observaciones,
	}

	_, err := database.CreateActividad(actividad)
	if err != nil {
		log.Printf("Error en actividadService.CreateActividad: %v", err)
		return nil, errors.New("Error al crear la actividad.")
	}

	// Devolvemos la lista actualizada
	actividades, err := database.GetActividadesByProyectoID(req.ProyectoID)
	if err != nil {
		log.Printf("Error recargando actividades post-creación: %v", err)
		return []models.ActividadResponse{}, nil
	}

	return actividades, nil
}

func (s *actividadService) UpdateActividad(req models.UpdateActividadRequest) ([]models.ActividadResponse, error) {
	if req.ID == 0 || req.ProyectoID == 0 || req.Actividad == "" || req.Costo == 0 || req.RecursoHumano == 0 {
		return nil, errors.New("ID, ProyectoID, Actividad, RecursoHumano y Costo son requeridos.")
	}

	// Manejo de valores opcionales
	var laborID, equipoID, encargadoID sql.NullInt64
	if req.LaborAgronomicaID != nil && *req.LaborAgronomicaID != 0 {
		laborID = sql.NullInt64{Int64: int64(*req.LaborAgronomicaID), Valid: true}
	}
	if req.EquipoImplementoID != nil && *req.EquipoImplementoID != 0 {
		equipoID = sql.NullInt64{Int64: int64(*req.EquipoImplementoID), Valid: true}
	}
	if req.EncargadoID != nil && *req.EncargadoID != 0 {
		encargadoID = sql.NullInt64{Int64: int64(*req.EncargadoID), Valid: true}
	}
	var observaciones sql.NullString
	if req.Observaciones != "" {
		observaciones = sql.NullString{String: req.Observaciones, Valid: true}
	}

	actividad := models.Actividad{
		ID:                 req.ID,
		ProyectoID:         req.ProyectoID,
		Actividad:          req.Actividad,
		LaborAgronomicaID:  laborID,
		EquipoImplementoID: equipoID,
		EncargadoID:        encargadoID,
		RecursoHumano:      req.RecursoHumano,
		Costo:              req.Costo,
		Observaciones:      observaciones,
	}

	affected, err := database.UpdateActividad(actividad)
	if err != nil {
		log.Printf("Error en actividadService.UpdateActividad (ID %d): %v", req.ID, err)
		return nil, errors.New("Error al actualizar la actividad.")
	}
	if affected == 0 {
		return nil, errors.New("Actividad no encontrada.")
	}

	// Devolvemos la lista actualizada
	actividades, err := database.GetActividadesByProyectoID(req.ProyectoID)
	if err != nil {
		log.Printf("Error recargando actividades post-update: %v", err)
		return []models.ActividadResponse{}, nil
	}

	return actividades, nil
}

func (s *actividadService) DeleteActividad(id int) (int64, error) {
	if id == 0 {
		return 0, errors.New("ID de actividad requerido.")
	}
	affected, err := database.DeleteActividad(id)
	if err != nil {
		log.Printf("Error en actividadService.DeleteActividad (ID %d): %v", id, err)
		return 0, errors.New("Error al borrar la actividad.")
	}
	return affected, nil
}
