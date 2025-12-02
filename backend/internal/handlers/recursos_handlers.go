package handlers

import (
	"encoding/json"
	"net/http"
	"proyecto/internal/auth"
	"proyecto/internal/database"
	"proyecto/internal/logger"
	"proyecto/internal/models"
)

type RecursoHandler struct {
	authSvc   auth.AuthService
	loggerSvc logger.LoggerService
}

func NewRecursoHandler(as auth.AuthService, ls logger.LoggerService) *RecursoHandler {
	return &RecursoHandler{authSvc: as, loggerSvc: ls}
}

// CREATE
func (h *RecursoHandler) CreateRecursoHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRecursoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	stmt, err := database.DB.Prepare(`
		INSERT INTO recursos_humanos (proyecto_id, actividad, accion, nombre, cedula, tiempo, cantidad, costo_unitario, monto)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(req.ProyectoID, req.Actividad, req.Accion, req.Nombre, req.Cedula, req.Tiempo, req.Cantidad, req.CostoUnitario, req.Monto)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, _ := res.LastInsertId()
	h.loggerSvc.Log(req.AdminUsername, "admin", "CREACIÓN", "Recurso Humano", int(id))
	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: "Recurso creado"})
}

// GET
func (h *RecursoHandler) GetRecursosHandler(w http.ResponseWriter, r *http.Request) {
	type GetReq struct {
		ProyectoID int `json:"proyecto_id"`
	}
	var req GetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	rows, err := database.DB.Query("SELECT id, proyecto_id, actividad, accion, nombre, cedula, tiempo, cantidad, costo_unitario, monto FROM recursos_humanos WHERE proyecto_id = ?", req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var lista []models.RecursoHumano
	for rows.Next() {
		var r models.RecursoHumano
		if err := rows.Scan(&r.ID, &r.ProyectoID, &r.Actividad, &r.Accion, &r.Nombre, &r.Cedula, &r.Tiempo, &r.Cantidad, &r.CostoUnitario, &r.Monto); err != nil {
			continue
		}
		lista = append(lista, r)
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"recursos": lista})
}

// UPDATE
func (h *RecursoHandler) UpdateRecursoHandler(w http.ResponseWriter, r *http.Request) {
	type UpdateReq struct {
		ID            int     `json:"id"`
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
	var req UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	stmt, err := database.DB.Prepare(`
		UPDATE recursos_humanos SET actividad=?, accion=?, nombre=?, cedula=?, tiempo=?, cantidad=?, costo_unitario=?, monto=? WHERE id=?
	`)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.Actividad, req.Accion, req.Nombre, req.Cedula, req.Tiempo, req.Cantidad, req.CostoUnitario, req.Monto, req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin", "MODIFICACIÓN", "Recurso Humano", req.ID)
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Recurso actualizado"})
}

// DELETE
func (h *RecursoHandler) DeleteRecursoHandler(w http.ResponseWriter, r *http.Request) {
	type DeleteReq struct {
		ID            int    `json:"id"`
		AdminUsername string `json:"admin_username"`
	}
	var req DeleteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	_, err := database.DB.Exec("DELETE FROM recursos_humanos WHERE id=?", req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.loggerSvc.Log(req.AdminUsername, "admin", "ELIMINACIÓN", "Recurso Humano", req.ID)
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Recurso eliminado"})
}
