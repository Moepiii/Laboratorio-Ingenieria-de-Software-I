package handlers

import (
	"encoding/json"
	"net/http"
	"proyecto/internal/auth"
	"proyecto/internal/database"
	"proyecto/internal/logger"
	"proyecto/internal/models"
)

type PlanHandler struct {
	authSvc   auth.AuthService
	loggerSvc logger.LoggerService
}

func NewPlanHandler(as auth.AuthService, ls logger.LoggerService) *PlanHandler {
	return &PlanHandler{authSvc: as, loggerSvc: ls}
}

// CreatePlanHandler guarda un nuevo plan
func (h *PlanHandler) CreatePlanHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	stmt, err := database.DB.Prepare(`
		INSERT INTO planes_accion (proyecto_id, actividad, accion, fecha_inicio, fecha_cierre, horas, responsable, costo_unitario, monto)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(req.ProyectoID, req.Actividad, req.Accion, req.FechaInicio, req.FechaCierre, req.Horas, req.Responsable, req.CostoUnitario, req.Monto)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, _ := res.LastInsertId()
	h.loggerSvc.Log(req.AdminUsername, "admin", "CREACIÓN", "Plan Accion", int(id))

	respondWithJSON(w, http.StatusCreated, models.SimpleResponse{Mensaje: "Plan creado exitosamente"})
}

// GetPlanesHandler obtiene los planes
func (h *PlanHandler) GetPlanesHandler(w http.ResponseWriter, r *http.Request) {
	var req models.GetPlanesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	rows, err := database.DB.Query("SELECT id, proyecto_id, actividad, accion, fecha_inicio, fecha_cierre, horas, responsable, costo_unitario, monto FROM planes_accion WHERE proyecto_id = ?", req.ProyectoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var planes []models.PlanAccion
	for rows.Next() {
		var p models.PlanAccion
		if err := rows.Scan(&p.ID, &p.ProyectoID, &p.Actividad, &p.Accion, &p.FechaInicio, &p.FechaCierre, &p.Horas, &p.Responsable, &p.CostoUnitario, &p.Monto); err != nil {
			continue
		}
		planes = append(planes, p)
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"planes": planes})
}

// ⭐️ NUEVO: UPDATE PLAN
func (h *PlanHandler) UpdatePlanHandler(w http.ResponseWriter, r *http.Request) {
	// Reusamos el struct pero añadimos ID manualmente o creamos uno nuevo.
	// Usaremos un map para flexibilidad o el mismo struct CreatePlanRequest asumiendo que el ID viene aparte o en URL,
	// pero lo mejor es definir un UpdateRequest.
	type UpdatePlanRequest struct {
		ID            int     `json:"id"`
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

	var req UpdatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	stmt, err := database.DB.Prepare(`
		UPDATE planes_accion SET 
			actividad=?, accion=?, fecha_inicio=?, fecha_cierre=?, 
			horas=?, responsable=?, costo_unitario=?, monto=?
		WHERE id=?
	`)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.Actividad, req.Accion, req.FechaInicio, req.FechaCierre, req.Horas, req.Responsable, req.CostoUnitario, req.Monto, req.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.loggerSvc.Log(req.AdminUsername, "admin", "MODIFICACIÓN", "Plan Accion", req.ID)
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Plan actualizado"})
}

// ⭐️ NUEVO: DELETE PLAN
func (h *PlanHandler) DeletePlanHandler(w http.ResponseWriter, r *http.Request) {
	type DeletePlanRequest struct {
		ID            int    `json:"id"`
		AdminUsername string `json:"admin_username"`
	}
	var req DeletePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	stmt, err := database.DB.Prepare("DELETE FROM planes_accion WHERE id=?")
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

	h.loggerSvc.Log(req.AdminUsername, "admin", "ELIMINACIÓN", "Plan Accion", req.ID)
	respondWithJSON(w, http.StatusOK, models.SimpleResponse{Mensaje: "Plan eliminado"})
}
