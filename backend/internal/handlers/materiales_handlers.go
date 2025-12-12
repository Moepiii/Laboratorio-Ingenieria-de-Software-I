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
		INSERT INTO materiales_insumos (proyecto_id, actividad, accion, categoria, responsable, nombre, unidad, cantidad, costo_unitario, monto)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.ProyectoID, req.Actividad, req.Accion, req.Categoria, req.Responsable, req.Nombre, req.Unidad, req.Cantidad, req.CostoUnitario, req.Monto)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin", "CREACIÓN", "Material/Insumo", 0)
	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: "Material creado exitosamente"})
}

// READ (GetMateriales)
func (h *MaterialHandler) GetMaterialesHandler(w http.ResponseWriter, r *http.Request) {

	var req models.GetMaterialesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	rows, err := database.DB.Query("SELECT id, actividad, accion, categoria, COALESCE(responsable, ''), nombre, unidad, cantidad, costo_unitario, monto FROM materiales_insumos WHERE proyecto_id = ?", req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var materiales []models.MaterialInsumo
	for rows.Next() {

		var m models.MaterialInsumo
		if err := rows.Scan(&m.ID, &m.Actividad, &m.Accion, &m.Categoria, &m.Responsable, &m.Nombre, &m.Unidad, &m.Cantidad, &m.CostoUnitario, &m.Monto); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		materiales = append(materiales, m)
	}

	if materiales == nil {

		materiales = []models.MaterialInsumo{}
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"materiales": materiales})
}

// UPDATE
func (h *MaterialHandler) UpdateMaterialHandler(w http.ResponseWriter, r *http.Request) {
	type UpdateReq struct {
		models.CreateMaterialRequest
		ID int `json:"id"`
	}
	var updateReq UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	stmt, err := database.DB.Prepare(`
		UPDATE materiales_insumos SET actividad=?, accion=?, categoria=?, responsable=?, nombre=?, unidad=?, cantidad=?, costo_unitario=?, monto=? WHERE id=?
	`)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(updateReq.Actividad, updateReq.Accion, updateReq.Categoria, updateReq.Responsable, updateReq.Nombre, updateReq.Unidad, updateReq.Cantidad, updateReq.CostoUnitario, updateReq.Monto, updateReq.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.loggerSvc.Log(updateReq.AdminUsername, "admin", "MODIFICACIÓN", "Material/Insumo", updateReq.ID)
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

	stmt, err := database.DB.Prepare("DELETE FROM materiales_insumos WHERE id=?")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin", "ELIMINACIÓN", "Material/Insumo", req.ID)
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Material eliminado"})
}
