package handlers

import (
	"encoding/json"
	"net/http"
	"proyecto/internal/auth"
	"proyecto/internal/database"
	"proyecto/internal/logger"
	"proyecto/internal/models"
)

type MaterialHandler struct {
	authSvc   auth.AuthService
	loggerSvc logger.LoggerService
}

func NewMaterialHandler(as auth.AuthService, ls logger.LoggerService) *MaterialHandler {
	return &MaterialHandler{authSvc: as, loggerSvc: ls}
}

// CREATE
func (h *MaterialHandler) CreateMaterialHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMaterialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	stmt, err := database.DB.Prepare(`
		INSERT INTO materiales_insumos (proyecto_id, actividad, accion, categoria, nombre, unidad, cantidad, costo_unitario, monto)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(req.ProyectoID, req.Actividad, req.Accion, req.Categoria, req.Nombre, req.Unidad, req.Cantidad, req.CostoUnitario, req.Monto)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, _ := res.LastInsertId()
	h.loggerSvc.Log(req.AdminUsername, "admin", "CREACIÓN", "Material/Insumo", int(id))
	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: "Material creado"})
}

// GET
func (h *MaterialHandler) GetMaterialesHandler(w http.ResponseWriter, r *http.Request) {
	type GetReq struct {
		ProyectoID int `json:"proyecto_id"`
	}
	var req GetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	rows, err := database.DB.Query("SELECT id, proyecto_id, actividad, accion, categoria, nombre, unidad, cantidad, costo_unitario, monto FROM materiales_insumos WHERE proyecto_id = ?", req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var lista []models.MaterialInsumo
	for rows.Next() {
		var m models.MaterialInsumo
		if err := rows.Scan(&m.ID, &m.ProyectoID, &m.Actividad, &m.Accion, &m.Categoria, &m.Nombre, &m.Unidad, &m.Cantidad, &m.CostoUnitario, &m.Monto); err != nil {
			continue
		}
		lista = append(lista, m)
	}
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"materiales": lista})
}

// UPDATE
func (h *MaterialHandler) UpdateMaterialHandler(w http.ResponseWriter, r *http.Request) {
	type UpdateReq struct {
		ID            int     `json:"id"`
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
	var req UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	stmt, err := database.DB.Prepare(`
		UPDATE materiales_insumos SET actividad=?, accion=?, categoria=?, nombre=?, unidad=?, cantidad=?, costo_unitario=?, monto=? WHERE id=?
	`)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.Actividad, req.Accion, req.Categoria, req.Nombre, req.Unidad, req.Cantidad, req.CostoUnitario, req.Monto, req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin", "MODIFICACIÓN", "Material/Insumo", req.ID)
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Material actualizado"})
}

// DELETE
func (h *MaterialHandler) DeleteMaterialHandler(w http.ResponseWriter, r *http.Request) {
	type DeleteReq struct {
		ID            int    `json:"id"`
		AdminUsername string `json:"admin_username"`
	}
	var req DeleteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	_, err := database.DB.Exec("DELETE FROM materiales_insumos WHERE id=?", req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.loggerSvc.Log(req.AdminUsername, "admin", "ELIMINACIÓN", "Material/Insumo", req.ID)
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Material eliminado"})
}
