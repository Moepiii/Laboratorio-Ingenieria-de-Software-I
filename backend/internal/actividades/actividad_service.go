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

func (s *actividadService) GetDatosProyecto(proyectoID int) (*GetDatosProyectoResponse, error) {
	if proyectoID == 0 {
		return nil, errors.New("ID de proyecto requerido.")
	}

	// 1. Obtener Labores
	labores, err := database.GetLaboresByProyectoID(proyectoID)
	if err != nil {
		log.Printf("Error en actividadService.GetDatosProyecto (GetLabores): %v", err)
		return nil, errors.New("Error al obtener labores.")
	}

	// 2. Obtener Equipos
	equipos, err := database.GetEquiposByProyectoID(proyectoID)
	if err != nil {
		log.Printf("Error en actividadService.GetDatosProyecto (GetEquipos): %v", err)
		return nil, errors.New("Error al obtener equipos.")
	}

	// 3. Obtener Encargados
	encargados, err := database.GetEncargadosByProyectoID(proyectoID)
	if err != nil {
		log.Printf("Error en actividadService.GetDatosProyecto (GetEncargados): %v", err)
		return nil, errors.New("Error al obtener encargados.")
	}

	// 4. Obtener Actividades
	actividades, err := database.GetActividadesByProyectoID(proyectoID)
	if err != nil {
		log.Printf("Error en actividadService.GetDatosProyecto (GetActividades): %v", err)
		return nil, errors.New("Error al obtener actividades.")
	}

	// 5. Empaquetar y responder
	response := &GetDatosProyectoResponse{
		Labores:     labores,
		Equipos:     equipos,
		Encargados:  encargados,
		Actividades: actividades,
	}
	return response, nil
}

func (s *actividadService) CreateActividad(req models.CreateActividadRequest) ([]models.ActividadResponse, error) {
	if req.ProyectoID == 0 || req.Actividad == "" {
		return nil, errors.New("ProyectoID y Nombre de Actividad son requeridos.")
	}

	// Convertir punteros de JSON a sql.Null*
	var laborID sql.NullInt64
	if req.LaborAgronomicaID != nil {
		laborID = sql.NullInt64{Int64: int64(*req.LaborAgronomicaID), Valid: true}
	}
	var equipoID sql.NullInt64
	if req.EquipoImplementoID != nil {
		equipoID = sql.NullInt64{Int64: int64(*req.EquipoImplementoID), Valid: true}
	}
	var encargadoID sql.NullInt64
	if req.EncargadoID != nil {
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
		return nil, errors.New("Error al crear actividad.")
	}

	// Devolvemos la lista completa actualizada (como en el handler original)
	actividades, err := database.GetActividadesByProyectoID(req.ProyectoID)
	if err != nil {
		log.Printf("Error recargando actividades post-creación: %v", err)
		// Devolvemos un slice vacío pero sin error fatal, para no romper el frontend
		return []models.ActividadResponse{}, nil
	}

	return actividades, nil
}

func (s *actividadService) UpdateActividad(req models.UpdateActividadRequest) ([]models.ActividadResponse, error) {
	if req.ID == 0 || req.ProyectoID == 0 || req.Actividad == "" {
		return nil, errors.New("ID, ProyectoID y Nombre de Actividad son requeridos.")
	}

	// Convertir punteros de JSON a sql.Null*
	var laborID sql.NullInt64
	if req.LaborAgronomicaID != nil {
		laborID = sql.NullInt64{Int64: int64(*req.LaborAgronomicaID), Valid: true}
	}
	var equipoID sql.NullInt64
	if req.EquipoImplementoID != nil {
		equipoID = sql.NullInt64{Int64: int64(*req.EquipoImplementoID), Valid: true}
	}
	var encargadoID sql.NullInt64
	if req.EncargadoID != nil {
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