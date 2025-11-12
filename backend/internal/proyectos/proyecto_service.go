package proyectos

import (
	"errors"
	// "fmt"
	"log"
	"strings"

	"proyecto/internal/database"
	"proyecto/internal/models"
)

// --- 1. EL CONTRATO (Interface) ---
type ProyectoService interface {
	GetAllProyectos() ([]models.Proyecto, error)
	CreateProyecto(nombre, fechaInicio, fechaCierre string) (*models.Proyecto, error)
	UpdateProyecto(id int, nombre, fechaInicio, fechaCierre string) (*models.Proyecto, error)
	DeleteProyecto(id int) (int64, error)
	SetProyectoEstado(id int, estado string) (int64, error)
}

// --- 2. LA IMPLEMENTACIÓN (Struct) ---
type proyectoService struct {
	// (Dependencias futuras, como repositorios)
}

// --- 3. EL CONSTRUCTOR ---
func NewProyectoService() ProyectoService {
	return &proyectoService{}
}

// --- 4. LOS MÉTODOS (Lógica de Negocio) ---

func (s *proyectoService) GetAllProyectos() ([]models.Proyecto, error) {
	proyectos, err := database.GetAllProyectos()
	if err != nil {
		log.Printf("Error en proyectoService.GetAllProyectos: %v", err)
		return nil, errors.New("Error al obtener proyectos.")
	}
	return proyectos, nil
}

func (s *proyectoService) CreateProyecto(nombre, fechaInicio, fechaCierre string) (*models.Proyecto, error) {
	if nombre == "" || fechaInicio == "" || fechaCierre == "" {
		return nil, errors.New("Nombre, Fecha de Inicio y Fecha de Cierre son requeridos.")
	}

	id, err := database.CreateProyecto(nombre, fechaInicio, fechaCierre)
	if err != nil {
		log.Printf("Error en proyectoService.CreateProyecto: %v", err)
		if strings.Contains(err.Error(), "ya existe") {
			return nil, err
		}
		return nil, errors.New("Error al crear proyecto.")
	}

	// Devolvemos el proyecto recién creado
	proyecto, err := database.GetProjectByID(id) // Asumiendo que GetProjectByID existe
	if err != nil {
		log.Printf("Error al obtener proyecto recién creado (ID: %d): %v", id, err)
		return nil, errors.New("Proyecto creado con éxito, pero no se pudo recuperar.")
	}
	return proyecto, nil
}

func (s *proyectoService) UpdateProyecto(id int, nombre, fechaInicio, fechaCierre string) (*models.Proyecto, error) {
	if id == 0 || nombre == "" || fechaInicio == "" || fechaCierre == "" {
		return nil, errors.New("ID, Nombre, Fecha de Inicio y Fecha de Cierre son requeridos.")
	}

	affected, err := database.UpdateProyecto(id, nombre, fechaInicio, fechaCierre)
	if err != nil {
		log.Printf("Error en proyectoService.UpdateProyecto (ID: %d): %v", id, err)
		if strings.Contains(err.Error(), "ya existe") {
			return nil, err
		}
		return nil, errors.New("Error al actualizar proyecto.")
	}
	if affected == 0 {
		return nil, errors.New("Proyecto no encontrado.")
	}

	// Devolvemos el proyecto actualizado
	proyecto, err := database.GetProjectByID(int64(id))
	if err != nil {
		return nil, errors.New("Proyecto actualizado pero no se pudo recuperar.")
	}
	return proyecto, nil
}

func (s *proyectoService) DeleteProyecto(id int) (int64, error) {
	if id == 0 {
		return 0, errors.New("ID de proyecto requerido.")
	}
	affected, err := database.DeleteProyecto(id)
	if err != nil {
		log.Printf("Error en proyectoService.DeleteProyecto (ID: %d): %v", id, err)
		return 0, errors.New("Error al borrar proyecto.")
	}
	return affected, nil
}

// ⭐️ CORRECCIÓN AQUÍ ⭐️
func (s *proyectoService) SetProyectoEstado(id int, estado string) (int64, error) {
	if id == 0 {
		return 0, errors.New("ID de proyecto requerido.")
	}
	if estado == "" {
		return 0, errors.New("Estado requerido.")
	}

	// El error era un typo. La función se llama 'SetProyectoEstado'
	affected, err := database.SetProyectoEstado(id, estado)
	if err != nil {
		log.Printf("Error en proyectoService.SetProyectoEstado (ID: %d): %v", id, err)
		return 0, errors.New("Error al cambiar estado del proyecto.")
	}
	return affected, nil
}
