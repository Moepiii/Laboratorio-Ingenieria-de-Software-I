package unidades

import (
	"errors"
	"log"
	"proyecto/internal/database"
	"proyecto/internal/models"
)

type UnidadService interface {
	GetUnidadesByProyectoID(proyectoID int) ([]models.UnidadMedida, error)
	CreateUnidad(req models.CreateUnidadRequest) (*models.UnidadMedida, error)
	UpdateUnidad(req models.UpdateUnidadRequest) (int64, error)
	DeleteUnidad(id int) (int64, error)
}

type unidadService struct{}

func NewUnidadService() UnidadService {
	return &unidadService{}
}

// Acepta ID de proyecto
func (s *unidadService) GetUnidadesByProyectoID(proyectoID int) ([]models.UnidadMedida, error) {
	if proyectoID == 0 {
		return nil, errors.New("ID de proyecto requerido")
	}
	return database.GetUnidadesByProyectoID(proyectoID)
}

func (s *unidadService) CreateUnidad(req models.CreateUnidadRequest) (*models.UnidadMedida, error) {
	if req.ProyectoID == 0 {
		return nil, errors.New("ID de proyecto requerido")
	}
	if req.Nombre == "" || req.Abreviatura == "" || req.Tipo == "" {
		return nil, errors.New("nombre, abreviatura y tipo son requeridos")
	}

	id, err := database.CreateUnidad(models.UnidadMedida{
		ProyectoID:  req.ProyectoID, // Guardamos el ID
		Nombre:      req.Nombre,
		Abreviatura: req.Abreviatura,
		Tipo:        req.Tipo,
		Dimension:   req.Dimension,
	})
	if err != nil {
		log.Printf("Error creando unidad: %v", err)
		return nil, errors.New("error al crear unidad")
	}

	return database.GetUnidadByID(int(id))
}

func (s *unidadService) UpdateUnidad(req models.UpdateUnidadRequest) (int64, error) {
	if req.ID == 0 || req.Nombre == "" || req.Abreviatura == "" {
		return 0, errors.New("datos incompletos")
	}
	return database.UpdateUnidad(req.ID, req.Nombre, req.Abreviatura, req.Tipo, req.Dimension)
}
func (s *unidadService) DeleteUnidad(id int) (int64, error) { return database.DeleteUnidad(id) }
