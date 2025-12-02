package models

import (
	"database/sql"

	"github.com/golang-jwt/jwt/v5"
)

// --- ESTRUCTURAS DE DATOS ---

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Cedula   string `json:"cedula"`
}

type UserDB struct {
	ID             int
	Username       string
	HashedPassword string
	Role           string
	Nombre         string
	Apellido       string
	Cedula         string
	ProyectoID     sql.NullInt64
}

type UserListResponse struct {
	ID             int     `json:"id"`
	Username       string  `json:"username"`
	Role           string  `json:"role"`
	Nombre         string  `json:"nombre"`
	Apellido       string  `json:"apellido"`
	Cedula         string  `json:"cedula"`
	ProyectoID     *int    `json:"proyecto_id"`
	ProyectoNombre *string `json:"proyecto_nombre"`
}

type UpdateRoleRequest struct {
	ID            int    `json:"id"`
	NewRole       string `json:"new_role"`
	AdminUsername string `json:"admin_username"`
}

type AdminActionRequest struct {
	AdminUsername string `json:"admin_username"`
}

type AddUserRequest struct {
	User          User   `json:"user"`
	AdminUsername string `json:"admin_username"`
}

type DeleteUserRequest struct {
	ID            int    `json:"id"`
	AdminUsername string `json:"admin_username"`
}

type AssignProjectRequest struct {
	UserID        int    `json:"user_id"`
	ProyectoID    int    `json:"proyecto_id"` // 0 para desasignar
	AdminUsername string `json:"admin_username"`
}

// Para el dashboard del rol 'user'
type UserProjectDetailsRequest struct {
	UserID int `json:"user_id"`
}

type UserProjectDetailsResponse struct {
	Proyecto *Proyecto       `json:"proyecto"` // Puede ser nil si no tiene proyecto
	Gerentes []ProjectMember `json:"gerentes"` // Lista de gerentes
	Miembros []ProjectMember `json:"miembros"` // Lista de compañeros
}

type ProjectMember struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

// --- Proyectos ---
type Proyecto struct {
	ID            int    `json:"id"`
	Nombre        string `json:"nombre"`
	FechaInicio   string `json:"fecha_inicio"`
	FechaCierre   string `json:"fecha_cierre"`
	Estado        string `json:"estado"`
	FechaCreacion string `json:"fecha_creacion"`
}

type CreateProyectoRequest struct {
	Nombre        string `json:"nombre"`
	FechaInicio   string `json:"fecha_inicio"`
	FechaCierre   string `json:"fecha_cierre"`
	AdminUsername string `json:"admin_username"`
}

type UpdateProyectoRequest struct {
	ID            int    `json:"id"`
	Nombre        string `json:"nombre"`
	FechaInicio   string `json:"fecha_inicio"`
	FechaCierre   string `json:"fecha_cierre"`
	AdminUsername string `json:"admin_username"`
}

type DeleteProyectoRequest struct {
	ID            int    `json:"id"`
	AdminUsername string `json:"admin_username"`
}

type SetProyectoEstadoRequest struct {
	ID            int    `json:"id"`
	Estado        string `json:"estado"`
	AdminUsername string `json:"admin_username"`
}

// --- Labores Agronómicas ---
type LaborAgronomica struct {
	ID            int    `json:"id"`
	ProyectoID    int    `json:"proyecto_id"`
	CodigoLabor   string `json:"codigo_labor"`
	Descripcion   string `json:"descripcion"`
	Estado        string `json:"estado"`
	FechaCreacion string `json:"fecha_creacion"`
}

type GetLaboresRequest struct {
	ProyectoID    int    `json:"proyecto_id"`
	AdminUsername string `json:"admin_username"`
}

type CreateLaborRequest struct {
	ProyectoID    int    `json:"proyecto_id"`
	Descripcion   string `json:"descripcion"`
	Estado        string `json:"estado"`
	AdminUsername string `json:"admin_username"`
}

type UpdateLaborRequest struct {
	ID            int    `json:"id"`
	CodigoLabor   string `json:"codigo_labor"`
	Descripcion   string `json:"descripcion"`
	Estado        string `json:"estado"`
	AdminUsername string `json:"admin_username"`
}

type DeleteLaborRequest struct {
	ID            int    `json:"id"`
	AdminUsername string `json:"admin_username"`
}

// --- Equipos e Implementos ---
type EquipoImplemento struct {
	ID            int    `json:"id"`
	ProyectoID    int    `json:"proyecto_id"`
	CodigoEquipo  string `json:"codigo_equipo"`
	Nombre        string `json:"nombre"`
	Tipo          string `json:"tipo"` // "Equipo" o "Implemento"
	Estado        string `json:"estado"`
	FechaCreacion string `json:"fecha_creacion"`
}

type GetEquiposRequest struct {
	ProyectoID    int    `json:"proyecto_id"`
	AdminUsername string `json:"admin_username"`
}

// ⭐️ --- INICIO DEL CAMBIO --- ⭐️
// Este es el struct que hemos modificado
type CreateEquipoRequest struct {
	ProyectoID int `json:"proyecto_id"`
	// CodigoEquipo  string `json:"codigo_equipo"` // CAMPO ELIMINADO
	Nombre        string `json:"nombre"`
	Tipo          string `json:"tipo"`
	Estado        string `json:"estado"`
	AdminUsername string `json:"admin_username"`
}

// ⭐️ --- FIN DEL CAMBIO --- ⭐️

type UpdateEquipoRequest struct {
	ID            int    `json:"id"`
	CodigoEquipo  string `json:"codigo_equipo"`
	Nombre        string `json:"nombre"`
	Tipo          string `json:"tipo"`
	Estado        string `json:"estado"`
	AdminUsername string `json:"admin_username"`
}

type DeleteEquipoRequest struct {
	ID            int    `json:"id"`
	AdminUsername string `json:"admin_username"`
}

// --- Actividades ---
type Actividad struct {
	ID                 int
	ProyectoID         int
	Actividad          string
	LaborAgronomicaID  sql.NullInt64 // Clave foránea opcional
	EquipoImplementoID sql.NullInt64 // Clave foránea opcional
	EncargadoID        sql.NullInt64 // Clave foránea opcional
	RecursoHumano      int
	Costo              float64
	Observaciones      sql.NullString // Texto opcional
	FechaCreacion      string
}

type ActividadResponse struct {
	ID                 int            `json:"id"`
	ProyectoID         int            `json:"proyecto_id"`
	Actividad          string         `json:"actividad"`
	LaborAgronomicaID  sql.NullInt64  `json:"labor_agronomica_id"`
	EquipoImplementoID sql.NullInt64  `json:"equipo_implemento_id"`
	EncargadoID        sql.NullInt64  `json:"encargado_id"`
	RecursoHumano      int            `json:"recurso_humano"`
	Costo              float64        `json:"costo"`
	Observaciones      sql.NullString `json:"observaciones"`
	FechaCreacion      string         `json:"fecha_creacion"`
	LaborDescripcion   sql.NullString `json:"labor_descripcion"` // Nombre/Descripción de la labor
	EquipoNombre       sql.NullString `json:"equipo_nombre"`     // Nombre del equipo
	EncargadoNombre    sql.NullString `json:"encargado_nombre"`  // Nombre del encargado
}

type EncargadoResponse struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Cedula   string `json:"cedula"`
}

type GetDatosProyectoRequest struct {
	ProyectoID    int    `json:"proyecto_id"`
	AdminUsername string `json:"admin_username"`
}

type CreateActividadRequest struct {
	ProyectoID         int     `json:"proyecto_id"`
	Actividad          string  `json:"actividad"`
	LaborAgronomicaID  *int    `json:"labor_agronomica_id"`
	EquipoImplementoID *int    `json:"equipo_implemento_id"`
	EncargadoID        *int    `json:"encargado_id"`
	RecursoHumano      int     `json:"recurso_humano"`
	Costo              float64 `json:"costo"`
	Observaciones      string  `json:"observaciones"`
	AdminUsername      string  `json:"admin_username"`
}

type UpdateActividadRequest struct {
	ID                 int     `json:"id"`
	ProyectoID         int     `json:"proyecto_id"`
	Actividad          string  `json:"actividad"`
	LaborAgronomicaID  *int    `json:"labor_agronomica_id"`
	EquipoImplementoID *int    `json:"equipo_implemento_id"`
	EncargadoID        *int    `json:"encargado_id"`
	RecursoHumano      int     `json:"recurso_humano"`
	Costo              float64 `json:"costo"`
	Observaciones      string  `json:"observaciones"`
	AdminUsername      string  `json:"admin_username"`
}

type DeleteActividadRequest struct {
	ID            int    `json:"id"`
	AdminUsername string `json:"admin_username"`
}

// --- Autenticación ---
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type LoginResponse struct {
	Token  string      `json:"token"`
	User   UserDetails `json:"user"`
	Role   string      `json:"role"`
	UserId int         `json:"userId"`
}

type UserDetails struct {
	Username string `json:"username"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Cedula   string `json:"cedula"`
}

// --- Respuestas Genéricas ---
type SimpleResponse struct {
	Mensaje string `json:"mensaje,omitempty"`
	Error   string `json:"error,omitempty"`
}

// --- Logger ---
type EventLog struct {
	ID              int
	Timestamp       string
	UsuarioUsername string
	UsuarioRol      string
	Accion          string
	Entidad         string
	EntidadID       int
}

type EventLogResponse struct {
	ID              int    `json:"id"`
	Timestamp       string `json:"timestamp"`
	UsuarioUsername string `json:"usuario_username"`
	UsuarioRol      string `json:"usuario_rol"`
	Accion          string `json:"accion"`
	Entidad         string `json:"entidad"`
	EntidadID       int    `json:"entidad_id"`
}

type GetLogsRequest struct {
	AdminUsername   string `json:"admin_username"`
	FechaInicio     string `json:"fecha_inicio"`
	FechaCierre     string `json:"fecha_cierre"`
	UsuarioUsername string `json:"usuario_username"`
	Accion          string `json:"accion"`
	Entidad         string `json:"entidad"`
}

// --- Unidades de Medida ---
type UnidadMedida struct {
	ID            int     `json:"id"`
	ProyectoID    int     `json:"proyecto_id"`
	Nombre        string  `json:"nombre"`
	Abreviatura   string  `json:"abreviatura"`
	Tipo          string  `json:"tipo"` // Peso, Líquido, Longitud
	Dimension     float64 `json:"dimension"`
	FechaCreacion string  `json:"fecha_creacion"`
}

type CreateUnidadRequest struct {
	ProyectoID    int     `json:"proyecto_id"` //
	Nombre        string  `json:"nombre"`
	Abreviatura   string  `json:"abreviatura"`
	Tipo          string  `json:"tipo"`
	Dimension     float64 `json:"dimension"`
	AdminUsername string  `json:"admin_username"`
}

type UpdateUnidadRequest struct {
	ID            int     `json:"id"`
	Nombre        string  `json:"nombre"`
	Abreviatura   string  `json:"abreviatura"`
	Tipo          string  `json:"tipo"`
	Dimension     float64 `json:"dimension"` //
	AdminUsername string  `json:"admin_username"`
}

type DeleteUnidadRequest struct {
	ID            int    `json:"id"`
	AdminUsername string `json:"admin_username"`
}
type GetUnidadesRequest struct {
	ProyectoID    int    `json:"proyecto_id"`
	AdminUsername string `json:"admin_username"`
}

type DeleteLogsRequest struct {
	IDs           []int  `json:"ids"` // Lista de IDs a borrar
	AdminUsername string `json:"admin_username"`
}

// --- models.go (Agrega esto al final) ---

type PlanAccion struct {
	ID            int     `json:"id"`
	ProyectoID    int     `json:"proyecto_id"`
	Actividad     string  `json:"actividad"`
	Accion        string  `json:"accion"`
	FechaInicio   string  `json:"fecha_inicio"`
	FechaCierre   string  `json:"fecha_cierre"`
	Horas         float64 `json:"horas"`
	Responsable   string  `json:"responsable"`
	CostoUnitario float64 `json:"costo_unitario"`
	Monto         float64 `json:"monto"`
}

type CreatePlanRequest struct {
	ProyectoID    int     `json:"proyecto_id"`
	Actividad     string  `json:"actividad"`
	Accion        string  `json:"accion"`
	FechaInicio   string  `json:"fecha_inicio"`
	FechaCierre   string  `json:"fecha_cierre"`
	Horas         float64 `json:"horas"`
	Responsable   string  `json:"responsable"`
	CostoUnitario float64 `json:"costo_unitario"`
	Monto         float64 `json:"monto"`
	AdminUsername string  `json:"admin_username"`
}

type GetPlanesRequest struct {
	ProyectoID    int    `json:"proyecto_id"`
	AdminUsername string `json:"admin_username"`
}

// --- Recurso Humano ---
type RecursoHumano struct {
	ID            int     `json:"id"`
	ProyectoID    int     `json:"proyecto_id"`
	Actividad     string  `json:"actividad"`
	Accion        string  `json:"accion"` // Viene de Labores
	Nombre        string  `json:"nombre"` // Responsable
	Cedula        string  `json:"cedula"`
	Tiempo        float64 `json:"tiempo"`   // Días u Horas
	Cantidad      float64 `json:"cantidad"` // Cantidad de personas
	CostoUnitario float64 `json:"costo_unitario"`
	Monto         float64 `json:"monto"` // Calculado
}

type CreateRecursoRequest struct {
	ProyectoID    int     `json:"proyecto_id"`
	Actividad     string  `json:"actividad"`
	Accion        string  `json:"accion"`
	Nombre        string  `json:"nombre"`
	Cedula        string  `json:"cedula"`
	Tiempo        float64 `json:"tiempo"`
	Cantidad      float64 `json:"cantidad"`
	CostoUnitario float64 `json:"costo_unitario"`
	Monto         float64 `json:"monto"`
	AdminUsername string  `json:"admin_username"`
}

// --- Materiales e Insumos ---
type MaterialInsumo struct {
	ID            int     `json:"id"`
	ProyectoID    int     `json:"proyecto_id"`
	Actividad     string  `json:"actividad"`
	Accion        string  `json:"accion"`    // Viene de Labores
	Categoria     string  `json:"categoria"` // ⭐️ NUEVO (Herbicida, Fertilizante, etc.)
	Nombre        string  `json:"nombre"`    // Nombre del producto
	Unidad        string  `json:"unidad"`    // Lts, Kg, Saco
	Cantidad      float64 `json:"cantidad"`
	CostoUnitario float64 `json:"costo_unitario"`
	Monto         float64 `json:"monto"` // Calculado
}

type CreateMaterialRequest struct {
	ProyectoID    int     `json:"proyecto_id"`
	Actividad     string  `json:"actividad"`
	Accion        string  `json:"accion"`
	Categoria     string  `json:"categoria"`
	Nombre        string  `json:"nombre"`
	Unidad        string  `json:"unidad"`
	Cantidad      float64 `json:"cantidad"`
	CostoUnitario float64 `json:"costo_unitario"`
	Monto         float64 `json:"monto"`
	AdminUsername string  `json:"admin_username"`
}
