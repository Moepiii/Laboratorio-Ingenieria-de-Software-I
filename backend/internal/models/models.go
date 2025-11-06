package models

import "database/sql"

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
	User
	AdminUsername string `json:"admin_username"`
}

type DeleteUserRequest struct {
	ID            int    `json:"id"`
	AdminUsername string `json:"admin_username"`
}

// --- ESTRUCTURAS DE RESPUESTA ---

type SimpleResponse struct {
	Mensaje string `json:"mensaje,omitempty"`
	Error   string `json:"error,omitempty"`
}

// --- ESTRUCTURAS DE PROYECTOS ---

type Proyecto struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`
	FechaInicio string `json:"fecha_inicio"`
	FechaCierre string `json:"fecha_cierre"`
	Estado      string `json:"estado"`
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

type AssignProyectoRequest struct {
	UserID        int    `json:"user_id"`
	ProyectoID    int    `json:"proyecto_id"`
	AdminUsername string `json:"admin_username"`
}

// Para la vista de usuario
type UserProjectDetailsRequest struct {
	UserID int `json:"user_id"`
}

type ProjectMember struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Role     string `json:"role"`
}

type UserProjectDetailsResponse struct {
	Proyecto *Proyecto       `json:"proyecto"`
	Miembros []ProjectMember `json:"miembros"`
	Gerentes []ProjectMember `json:"gerentes"`
}

// --- Estructuras para Labores Agronómicas ---

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
	CodigoLabor   string `json:"codigo_labor"`
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

// --- Estructuras para Equipos e Implementos ---

type EquipoImplemento struct {
	ID            int    `json:"id"`
	ProyectoID    int    `json:"proyecto_id"`
	CodigoEquipo  string `json:"codigo_equipo"`
	Nombre        string `json:"nombre"`
	Tipo          string `json:"tipo"`
	Estado        string `json:"estado"`
	FechaCreacion string `json:"fecha_creacion"`
}

type GetEquiposRequest struct {
	ProyectoID    int    `json:"proyecto_id"`
	AdminUsername string `json:"admin_username"`
}

type CreateEquipoRequest struct {
	ProyectoID    int    `json:"proyecto_id"`
	CodigoEquipo  string `json:"codigo_equipo"`
	Nombre        string `json:"nombre"`
	Tipo          string `json:"tipo"`
	Estado        string `json:"estado"`
	AdminUsername string `json:"admin_username"`
}

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

// ⭐️ --- (INICIO) Estructuras para Actividades (DatosProyecto.js) --- ⭐️

// ⭐️ NUEVO: Estructura de la tabla `actividades`
// Usamos sql.Null* para las claves foráneas que pueden ser nulas
type Actividad struct {
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
}

// ⭐️ NUEVO: Estructura para la respuesta de la lista de actividades (con JOINs)
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
	LaborDescripcion   sql.NullString `json:"labor_descripcion"`
	EquipoNombre       sql.NullString `json:"equipo_nombre"`
	EncargadoNombre    sql.NullString `json:"encargado_nombre"`
}

// ⭐️ NUEVO: Estructura para la lista de Encargados (rol 'encargado')
type EncargadoResponse struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

// ⭐️ NUEVO: Para la petición GET de datos de la página (Labores, Equipos, Encargados, Actividades)
type GetDatosProyectoRequest struct {
	ProyectoID    int    `json:"proyecto_id"`
	AdminUsername string `json:"admin_username"`
}

// ⭐️ NUEVO: Para la petición de CREAR una actividad
// El frontend enviará 0 o null para los IDs opcionales
type CreateActividadRequest struct {
	ProyectoID         int     `json:"proyecto_id"`
	Actividad          string  `json:"actividad"`
	LaborAgronomicaID  *int    `json:"labor_agronomica_id"` // Usamos punteros para manejar 'null' desde JSON
	EquipoImplementoID *int    `json:"equipo_implemento_id"`
	EncargadoID        *int    `json:"encargado_id"`
	RecursoHumano      int     `json:"recurso_humano"`
	Costo              float64 `json:"costo"`
	Observaciones      string  `json:"observaciones"` // El frontend enviará string vacío
	AdminUsername      string  `json:"admin_username"`
}

// ⭐️ NUEVO: Para la petición de ACTUALIZAR una actividad
type UpdateActividadRequest struct {
	ID                 int     `json:"id"`
	ProyectoID         int     `json:"proyecto_id"` // Para seguridad
	Actividad          string  `json:"actividad"`
	LaborAgronomicaID  *int    `json:"labor_agronomica_id"`
	EquipoImplementoID *int    `json:"equipo_implemento_id"`
	EncargadoID        *int    `json:"encargado_id"`
	RecursoHumano      int     `json:"recurso_humano"`
	Costo              float64 `json:"costo"`
	Observaciones      string  `json:"observaciones"`
	AdminUsername      string  `json:"admin_username"`
}

// ⭐️ NUEVO: Para la petición de BORRAR una actividad
type DeleteActividadRequest struct {
	ID            int    `json:"id"`
	AdminUsername string `json:"admin_username"`
}

// ⭐️ --- (FIN) Estructuras para Actividades --- ⭐️
