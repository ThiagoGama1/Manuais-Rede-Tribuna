package data

import (
	"manual/internal/models"
)

// aqui ficam as funcoes sql

func InsertManual(manual models.Manual){
	stmt := "INSERT INTO manuais(titulo, conteudo, secao) VALUES(?, ?, ?)"
	_, err := DB.Exec(stmt, manual.Titulo, manual.Conteudo, manual.Secao)
	
	if err != nil{
		return
	}
}
func GetManuais()[]models.Manual{
	query := "SELECT * FROM manuais"
	var list []models.Manual
	queryResult, err := DB.Query(query)
	
	if err != nil{
		return list
	}
	defer queryResult.Close()

	for queryResult.Next(){
		var m models.Manual
		queryResult.Scan(&m.ID, &m.Titulo, &m.Conteudo, &m.Secao)
		list = append(list, m)
	}
	
	return list
}

func GetManualByID(id int) (models.Manual, error){
	query := `SELECT * FROM manuais WHERE id = ?`
	var result models.Manual
	err := DB.QueryRow(query, id).Scan(&result.ID, &result.Titulo, &result.Conteudo, &result.Secao)


	if err == nil{
		return result, nil
	}
	return models.Manual{}, err
}

func DeleteManual(id int) error{
	query := `DELETE FROM manuais WHERE id = ?`

	_, err := DB.Exec(query, id)

	
	return err
}

func UpdateManual(m models.Manual) error{
	query := `UPDATE manuais 
			  SET titulo = ?,
			  conteudo = ?,
			  secao = ?
			  WHERE id = ?`

	_, err := DB.Exec(query, m.Titulo, m.Conteudo, m.Secao, m.ID)

	if err != nil {
		return err
	}
	return nil
}