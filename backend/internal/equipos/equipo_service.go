package equipos

import (
	"errors"
	"log"
	"strconv" // ⭐️ 1. IMPORTAMOS strconv
	"strings"

	"proyecto/internal/database"
	"proyecto/internal/models"
)

// --- 1. EL CONTRATO (Interface) ---
type EquipoService interface {
	GetEquiposByProyectoID(proyectoID int) ([]models.EquipoImplemento, error)
	CreateEquipo(req models.CreateEquipoRequest) (*models.EquipoImplemento, error)
	UpdateEquipo(req models.UpdateEquipoRequest) (int64, error)
	DeleteEquipo(id int) (int64, error)
}

// --- 2. LA IMPLEMENTACIÓN (Struct) ---
type equipoService struct {
	// (Dependencias futuras, como repositorios)
}

// --- 3. EL CONSTRUCTOR ---
func NewEquipoService() EquipoService {
	return &equipoService{}
}

// --- 4. LOS MÉTODOS (Lógica de Negocio) ---

func (s *equipoService) GetEquiposByProyectoID(proyectoID int) ([]models.EquipoImplemento, error) {
	if proyectoID == 0 {
		return nil, errors.New("id de proyecto requerido")
	}
	equipos, err := database.GetEquiposByProyectoID(proyectoID)
	if err != nil {
		log.Printf("Error en equipoService.GetEquiposByProyectoID: %v", err)
		return nil, errors.New("error al obtener equipos")
	}
	return equipos, nil
}

// ⭐️ --- INICIO: FUNCIÓN CreateEquipo MODIFICADA --- ⭐️
func (s *equipoService) CreateEquipo(req models.CreateEquipoRequest) (*models.EquipoImplemento, error) {
	// 1. Validación (ya no se valida CodigoEquipo)
	if req.ProyectoID == 0 || req.Nombre == "" || req.Tipo == "" {
		return nil, errors.New("ProyectoID, Nombre y Tipo son requeridos")
	}

	// 2. LÓGICA NUEVA: Obtener el siguiente código
	nextCodigoInt, err := database.GetNextEquipoCodigo(req.ProyectoID)
	if err != nil {
		log.Printf("Error en equipoService.CreateEquipo (GetNextEquipoCodigo): %v", err)
		return nil, errors.New("error al generar el código de equipo")
	}

	// 3. Convertir el número (ej: 1, 2, 3) a un string ("1", "2", "3")
	nextCodigoStr := strconv.Itoa(nextCodigoInt)

	// 4. Construir el struct EquipoImplemento completo
	// El servicio es ahora responsable de asignar el código.
	equipo := models.EquipoImplemento{
		ProyectoID:   req.ProyectoID,
		CodigoEquipo: nextCodigoStr, // ⬅️ Asignamos el nuevo código
		Nombre:       req.Nombre,
		Tipo:         req.Tipo,
		Estado:       req.Estado,
	}

	// 5. Llamada a la base de datos (esta función no cambia)
	equipoID, err := database.CreateEquipo(equipo)
	if err != nil {
		log.Printf("Error en equipoService.CreateEquipo (CreateEquipo): %v", err)
		if strings.Contains(err.Error(), "ya existe") {
			return nil, err
		}
		return nil, errors.New("error al crear equipo")
	}

	// 6. Devolver el objeto creado
	nuevoEquipo, err := database.GetEquipoByID(int(equipoID))
	if err != nil {
		log.Printf("Error al obtener equipo recién creado (ID: %d): %v", equipoID, err)
		return nil, errors.New("equipo creado con éxito, pero no se pudo recuperar")
	}

	return nuevoEquipo, nil
}

// ⭐️ --- FIN: FUNCIÓN CreateEquipo MODIFICADA --- ⭐️

func (s *equipoService) UpdateEquipo(req models.UpdateEquipoRequest) (int64, error) {
	if req.ID == 0 || req.CodigoEquipo == "" || req.Nombre == "" || req.Tipo == "" || req.Estado == "" {
		return 0, errors.New("ID, Código, Nombre, Tipo y Estado son requeridos")
	}

	affected, err := database.UpdateEquipo(req.ID, req.CodigoEquipo, req.Nombre, req.Tipo, req.Estado)
	if err != nil {
		log.Printf("Error en equipoService.UpdateEquipo (ID %d): %v", req.ID, err)
		if strings.Contains(err.Error(), "ya existe") {
			return 0, err
		}
		return 0, errors.New("error al actualizar el equipo")
	}
	if affected == 0 {
		return 0, errors.New("equipo no encontrado")
	}
	return affected, nil
}

func (s *equipoService) DeleteEquipo(id int) (int64, error) {
	if id == 0 {
		return 0, errors.New("id de equipo requerido")
	}
	affected, err := database.DeleteEquipo(id)
	if err != nil {
		log.Printf("Error en equipoService.DeleteEquipo (ID %d): %v", id, err)
		return 0, errors.New("error al borrar el equipo")
	}
	return affected, nil
}
