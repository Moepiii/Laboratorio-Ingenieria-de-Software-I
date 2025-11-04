package models

import "database/sql"

// --- ESTRUCTURAS DE DATOS ---

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

type UserDB struct {
	ID             int
	Username       string
	HashedPassword string
	Role           string
	Nombre         string
	Apellido       string
	ProyectoID     sql.NullInt64
}

type UserListResponse struct {
	ID             int     `json:"id"`
	Username       string  `json:"username"`
	Role           string  `json:"role"`
	Nombre         string  `json:"nombre"`
	Apellido       string  `json:"apellido"`
	ProyectoID     *int    `json:"proyecto_id"`
	ProyectoNombre *string `json:"proyecto_nombre"`
}

type UpdateRoleRequest struct {
	ID            int    `json:"id"`
	NewRole       string `json:"new_role"`
	AdminUsername string `json:"admin_username"` // Quien hace la petición
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

// NUEVA STRUCT para cambiar estado
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

// ⭐️ --- (INICIO) Estructuras para Labores Agronómicas --- ⭐️

// ⭐️ NUEVO: Estructura de la tabla (la que database.go usa)
type LaborAgronomica struct {
	ID            int    `json:"id"`
	ProyectoID    int    `json:"proyecto_id"`
	Descripcion   string `json:"descripcion"`
	Estado        string `json:"estado"`
	FechaCreacion string `json:"fecha_creacion"`
}

// ⭐️ NUEVO: Para la petición GET del frontend
type GetLaboresRequest struct {
	ProyectoID    int    `json:"proyecto_id"`
	AdminUsername string `json:"admin_username"`
}

// ⭐️ NUEVO: Para la petición CREATE del frontend
type CreateLaborRequest struct {
	ProyectoID    int    `json:"proyecto_id"`
	Descripcion   string `json:"descripcion"`
	Estado        string `json:"estado"` // Opcional, el estado por defecto es 'activa'
	AdminUsername string `json:"admin_username"`
}

// ⭐️ NUEVO: Para la petición UPDATE del frontend
type UpdateLaborRequest struct {
	ID            int    `json:"id"` // ID de la labor
	Descripcion   string `json:"descripcion"`
	Estado        string `json:"estado"`
	AdminUsername string `json:"admin_username"`
}

// ⭐️ NUEVO: Para la petición DELETE del frontend
type DeleteLaborRequest struct {
	ID            int    `json:"id"` // ID de la labor
	AdminUsername string `json:"admin_username"`
}

// ⭐️ --- (FIN) Estructuras para Labores Agronómicas --- ⭐️
