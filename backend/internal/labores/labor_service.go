package labores

import (
	"errors"
	"log"
	"strconv" // ⭐️ 1. IMPORTAMOS strconv
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
		return nil, errors.New("id de proyecto requerido")
	}
	labores, err := database.GetLaboresByProyectoID(proyectoID)
	if err != nil {
		log.Printf("Error en laborService.GetLaboresByProyectoID: %v", err)
		return nil, errors.New("error al obtener labores")
	}
	return labores, nil
}

// ⭐️ --- INICIO: FUNCIÓN CreateLabor MODIFICADA --- ⭐️
func (s *laborService) CreateLabor(req models.CreateLaborRequest) (*models.LaborAgronomica, error) {
	// 1. Validación (ya no se valida CodigoLabor)
	if req.ProyectoID == 0 || req.Descripcion == "" {
		return nil, errors.New("ProyectoID y Descripcion son requeridos")
	}

	// 2. LÓGICA NUEVA: Obtener el siguiente código
	nextCodigoInt, err := database.GetNextLaborCodigo(req.ProyectoID)
	if err != nil {
		log.Printf("Error en laborService.CreateLabor (GetNextLaborCodigo): %v", err)
		return nil, errors.New("error al generar el código de labor")
	}

	// 3. Convertir el número (ej: 1, 2, 3) a un string ("1", "2", "3")
	nextCodigoStr := strconv.Itoa(nextCodigoInt)

	// 4. Construir el struct LaborAgronomica completo
	// El servicio es ahora responsable de asignar el código.
	labor := models.LaborAgronomica{
		ProyectoID:  req.ProyectoID,
		CodigoLabor: nextCodigoStr, // ⬅️ Asignamos el nuevo código
		Descripcion: req.Descripcion,
		Estado:      req.Estado,
	}

	// 5. Llamada a la base de datos (esta función no cambia)
	laborID, err := database.CreateLabor(labor)
	if err != nil {
		log.Printf("Error en laborService.CreateLabor (CreateLabor): %v", err)
		// Este error es menos probable ahora, pero lo mantenemos por si acaso
		if strings.Contains(err.Error(), "ya existe") {
			return nil, errors.New("el código de labor ya existe")
		}
		return nil, errors.New("error al crear la labor")
	}

	// 6. Devolver el objeto creado
	nuevaLabor, err := database.GetLaborByID(int(laborID))
	if err != nil {
		log.Printf("Error al obtener labor recién creada (ID: %d): %v", laborID, err)
		return nil, errors.New("labor creada con éxito, pero no se pudo recuperar")
	}

	return nuevaLabor, nil
}

// ⭐️ --- FIN: FUNCIÓN CreateLabor MODIFICADA --- ⭐️

func (s *laborService) UpdateLabor(req models.UpdateLaborRequest) (int64, error) {
	if req.ID == 0 || req.CodigoLabor == "" || req.Descripcion == "" || req.Estado == "" {
		return 0, errors.New("ID, Código, Descripcion y Estado son requeridos")
	}

	affected, err := database.UpdateLabor(req.ID, req.CodigoLabor, req.Descripcion, req.Estado)
	if err != nil {
		log.Printf("Error en laborService.UpdateLabor (ID %d): %v", req.ID, err)
		if strings.Contains(err.Error(), "ya existe") {
			return 0, errors.New("el código de labor ya existe para este proyecto")
		}
		return 0, errors.New("error al actualizar la labor")
	}
	if affected == 0 {
		return 0, errors.New("labor no encontrada")
	}
	return affected, nil
}

func (s *laborService) DeleteLabor(id int) (int64, error) {
	if id == 0 {
		return 0, errors.New("id de labor requerido")
	}
	affected, err := database.DeleteLabor(id)
	if err != nil {
		log.Printf("Error en laborService.DeleteLabor (ID %d): %v", id, err)
		return 0, errors.New("error al borrar la labor")
	}
	return affected, nil
}
