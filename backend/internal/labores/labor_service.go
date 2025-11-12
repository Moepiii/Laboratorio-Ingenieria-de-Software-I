package labores

import (
	"errors"
	"log"
	"strings"

	"proyecto/internal/database"
	"proyecto/internal/models"
)

// --- 1. EL CONTRATO (Interface) ---
type LaborService interface {
	GetLaboresByProyectoID(proyectoID int) ([]models.LaborAgronomica, error)
	CreateLabor(req models.CreateLaborRequest) (*models.LaborAgronomica, error)
	UpdateLabor(req models.UpdateLaborRequest) (int64, error)
	DeleteLabor(id int) (int64, error)
}

// --- 2. LA IMPLEMENTACIÓN (Struct) ---
type laborService struct {
	// (Dependencias futuras, como repositorios)
}

// --- 3. EL CONSTRUCTOR ---
func NewLaborService() LaborService {
	return &laborService{}
}

// --- 4. LOS MÉTODOS (Lógica de Negocio) ---

func (s *laborService) GetLaboresByProyectoID(proyectoID int) ([]models.LaborAgronomica, error) {
	if proyectoID == 0 {
		return nil, errors.New("ID de proyecto requerido.")
	}
	labores, err := database.GetLaboresByProyectoID(proyectoID)
	if err != nil {
		log.Printf("Error en laborService.GetLaboresByProyectoID: %v", err)
		return nil, errors.New("Error al obtener labores.")
	}
	return labores, nil
}

func (s *laborService) CreateLabor(req models.CreateLaborRequest) (*models.LaborAgronomica, error) {
	if req.ProyectoID == 0 || req.Descripcion == "" || req.CodigoLabor == "" {
		return nil, errors.New("ProyectoID, Código y Descripción son requeridos.")
	}

	estado := req.Estado
	if estado == "" {
		estado = "activa" // Valor por defecto
	}

	labor := models.LaborAgronomica{
		ProyectoID:  req.ProyectoID,
		CodigoLabor: req.CodigoLabor,
		Descripcion: req.Descripcion,
		Estado:      estado,
	}

	laborID, err := database.CreateLabor(labor)
	if err != nil {
		log.Printf("Error en laborService.CreateLabor: %v", err)
		if strings.Contains(err.Error(), "ya existe") {
			return nil, err
		}
		return nil, errors.New("Error al crear labor.")
	}

	// Devolvemos la labor recién creada
	nuevaLabor, err := database.GetLaborByID(int(laborID))
	if err != nil {
		log.Printf("Error al obtener labor recién creada (ID: %d): %v", laborID, err)
		return nil, errors.New("Labor creada con éxito, pero no se pudo recuperar.")
	}

	return nuevaLabor, nil
}

func (s *laborService) UpdateLabor(req models.UpdateLaborRequest) (int64, error) {
	if req.ID == 0 || req.CodigoLabor == "" || req.Descripcion == "" || req.Estado == "" {
		return 0, errors.New("ID, Código, Descripción y Estado son requeridos.")
	}

	affected, err := database.UpdateLabor(req.ID, req.CodigoLabor, req.Descripcion, req.Estado)
	if err != nil {
		log.Printf("Error en laborService.UpdateLabor (ID %d): %v", req.ID, err)
		if strings.Contains(err.Error(), "ya existe") {
			return 0, err
		}
		return 0, errors.New("Error al actualizar la labor.")
	}
	return affected, nil
}

func (s *laborService) DeleteLabor(id int) (int64, error) {
	if id == 0 {
		return 0, errors.New("ID de labor requerido.")
	}
	affected, err := database.DeleteLabor(id)
	if err != nil {
		log.Printf("Error en laborService.DeleteLabor (ID %d): %v", id, err)
		return 0, errors.New("Error al borrar la labor.")
	}
	return affected, nil
}