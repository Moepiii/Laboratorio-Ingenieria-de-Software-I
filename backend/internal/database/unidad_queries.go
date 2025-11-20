package database

import (
	"proyecto/internal/models"
)

// --- QUERIES DE UNIDADES DE MEDIDA ---

func GetAllUnidades() ([]models.UnidadMedida, error) {
	// ⭐️ AÑADIDO: dimension en el SELECT
	rows, err := DB.Query("SELECT id, nombre, abreviatura, tipo, dimension, fecha_creacion FROM unidades_medida ORDER BY fecha_creacion DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unidades []models.UnidadMedida
	for rows.Next() {
		var u models.UnidadMedida
		// ⭐️ AÑADIDO: &u.Dimension en el Scan
		if err := rows.Scan(&u.ID, &u.Nombre, &u.Abreviatura, &u.Tipo, &u.Dimension, &u.FechaCreacion); err != nil {
			continue
		}
		unidades = append(unidades, u)
	}
	return unidades, nil
}

func GetUnidadByID(id int) (*models.UnidadMedida, error) {
	// ⭐️ AÑADIDO: dimension en el SELECT
	row := DB.QueryRow("SELECT id, nombre, abreviatura, tipo, dimension, fecha_creacion FROM unidades_medida WHERE id = ?", id)
	var u models.UnidadMedida
	// ⭐️ AÑADIDO: &u.Dimension en el Scan
	err := row.Scan(&u.ID, &u.Nombre, &u.Abreviatura, &u.Tipo, &u.Dimension, &u.FechaCreacion)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUnidad(u models.UnidadMedida) (int64, error) {
	// ⭐️ AÑADIDO: dimension en el INSERT
	stmt, err := DB.Prepare("INSERT INTO unidades_medida (nombre, abreviatura, tipo, dimension) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(u.Nombre, u.Abreviatura, u.Tipo, u.Dimension)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func UpdateUnidad(id int, nombre, abreviatura, tipo string, dimension float64) (int64, error) {
	// ⭐️ AÑADIDO: dimension = ? en el UPDATE
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