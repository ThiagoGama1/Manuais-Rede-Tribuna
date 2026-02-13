package data

import (
	"manual/internal/models"
	"os"
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
	query := "SELECT id, titulo, conteudo, secao FROM manuais"
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
	query := `SELECT id, titulo, conteudo, secao FROM manuais WHERE id = ?`
	var result models.Manual
	err := DB.QueryRow(query, id).Scan(&result.ID, &result.Titulo, &result.Conteudo, &result.Secao)


	if err != nil{
		return models.Manual{}, err
	}
	queryAnexos := `SELECT * FROM anexos WHERE manual_id = ?`

	rows, err := DB.Query(queryAnexos, id) // pega os arquivos

	if err != nil{
		return result, err
	}
	defer rows.Close()

	for rows.Next(){
		var a models.Anexo

		err := rows.Scan(&a.ID, &a.Nome, &a.Tamanho_bytes, &a.Caminho, &a.Tipo_arquivo, &a.Manual_id)

		if err != nil{
			return result, err
		}
		result.Arquivos = append(result.Arquivos, a)
	}

	return result, nil
}

func DeleteManual(id int) error{
	query := `DELETE FROM manuais WHERE id = ?`
	
	manual, err := GetManualByID(id)
	if err != nil{
		return err
	}
	for _, arquivo := range manual.Arquivos{
		caminho := arquivo.Caminho
		os.Remove(caminho)
	}
	queryDeleteAnexos := `DELETE FROM anexos WHERE manual_id = ?`

	_, err = DB.Exec(queryDeleteAnexos, id) // primeiro o filho
	_, err = DB.Exec(query, id) //depois o pai

	
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