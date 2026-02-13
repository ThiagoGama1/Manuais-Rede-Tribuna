package data

import (
	"manual/internal/models"
)

// aqui ficam as funcoes sql

func InsertManual(manual models.Manual)(int64, error) {
	stmt := "INSERT INTO manuais(titulo, conteudo, secao) VALUES(?, ?, ?)"
	resp, err := DB.Exec(stmt, manual.Titulo, manual.Conteudo, manual.Secao)
	
	if err != nil{
		return 0, err
	}
	id, err := resp.LastInsertId()

	if err != nil{
		return 0, err
	}
	for _, anexo := range manual.Arquivos{
		stmtAnexos := "INSERT INTO anexos(nome, tamanho, caminho, tipo_arquivo, manual_id) VALUES(?, ?, ?, ?, ?)"
		
		_, err = DB.Exec(stmtAnexos, anexo.Nome, anexo.Tamanho_bytes, anexo.Caminho, anexo.Tipo_arquivo, id)
		
	}
	return id, nil
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
	err := DB.QueryRow(query, id).Scan(&result.ID, &result.Titulo, &result.Conteudo, &result.Secao, &result.Arquivos)


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