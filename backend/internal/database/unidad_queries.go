package database

import (
	"proyecto/internal/models"
)

// Recibe proyectoID
func GetUnidadesByProyectoID(proyectoID int) ([]models.UnidadMedida, error) {
	rows, err := DB.Query("SELECT id, proyecto_id, nombre, abreviatura, tipo, dimension, fecha_creacion FROM unidades_medida WHERE proyecto_id = ? ORDER BY fecha_creacion DESC", proyectoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unidades []models.UnidadMedida
	for rows.Next() {
		var u models.UnidadMedida
		if err := rows.Scan(&u.ID, &u.ProyectoID, &u.Nombre, &u.Abreviatura, &u.Tipo, &u.Dimension, &u.FechaCreacion); err != nil {
			continue
		}
		unidades = append(unidades, u)
	}
	return unidades, nil
}

func GetUnidadByID(id int) (*models.UnidadMedida, error) {
	row := DB.QueryRow("SELECT id, proyecto_id, nombre, abreviatura, tipo, dimension, fecha_creacion FROM unidades_medida WHERE id = ?", id)
	var u models.UnidadMedida
	err := row.Scan(&u.ID, &u.ProyectoID, &u.Nombre, &u.Abreviatura, &u.Tipo, &u.Dimension, &u.FechaCreacion)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUnidad(u models.UnidadMedida) (int64, error) {
	// proyecto_id
	stmt, err := DB.Prepare("INSERT INTO unidades_medida (proyecto_id, nombre, abreviatura, tipo, dimension) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.ProyectoID, u.Nombre, u.Abreviatura, u.Tipo, u.Dimension)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Update y Delete quedan igual (usan ID)
func UpdateUnidad(id int, nombre, abreviatura, tipo string, dimension float64) (int64, error) {
	stmt, err := DB.Prepare("UPDATE unidades_medida SET nombre = ?, abreviatura = ?, tipo = ?, dimension = ? WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(nombre, abreviatura, tipo, dimension, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func DeleteUnidad(id int) (int64, error) {
	res, err := DB.Exec("DELETE FROM unidades_medida WHERE id = ?", id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
