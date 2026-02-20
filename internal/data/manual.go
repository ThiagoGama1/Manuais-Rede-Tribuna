package data

import (
	"manual/internal/models"
	"os"
)


func InsertManual(manual models.Manual)(int64, error) {
	stmt := "INSERT INTO manuais(titulo, secao) VALUES(?, ?)"
	resp, err := DB.Exec(stmt, manual.Titulo, manual.Secao)
	
	if err != nil{
		return 0, err
	}
	id, err := resp.LastInsertId()

	if err != nil{
		return 0, err
	}
	for i, etapa := range manual.Etapas{
		etapa.Manual_id = int(id)
		stmtEtapa := `INSERT INTO etapas(manual_id, ordem_apresentacao, conteudo) VALUES(?, ?, ?)`

		respEtapa, err := DB.Exec(stmtEtapa, etapa.Manual_id,i+1, etapa.Conteudo)
		if err != nil{
			return 0, err
		}
		var idEtapa int64

		idEtapa, err = respEtapa.LastInsertId()

		InsertAnexo(etapa.Anexos[idEtapa])
	
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
	queryAnexos := `SELECT id, nome, tamanho, caminho, tipo_arquivo, manual_id, ordem_apresentacao 
	FROM anexos WHERE manual_id = ? ORDER BY ordem_apresentacao ASC`

	rows, err := DB.Query(queryAnexos, id)

	if err != nil{
		return result, err
	}
	defer rows.Close()

	for rows.Next(){
		var a models.Anexo

		err := rows.Scan(&a.ID, &a.Nome, &a.Tamanho_bytes, &a.Caminho, &a.Tipo_arquivo, &a.Manual_id, &a.OrdemApresentacao)

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

	_, err = DB.Exec(queryDeleteAnexos, id)
	_, err = DB.Exec(query, id)

	
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
func DeleteAnexo(idAnexo int) error{
	query := `SELECT caminho FROM anexos WHERE id = ?`
	var caminho string
	err := DB.QueryRow(query, idAnexo).Scan(&caminho)
	if err != nil{
		return err
	}
	err = os.Remove(caminho)
	if err != nil{
		return err
	}
	query2 := `DELETE FROM anexos WHERE id = ?`

	_, err = DB.Exec(query2, idAnexo)

	return err
}
func InsertAnexo(anexo models.Anexo) error{
	//preciso descobrir o maior numero da ordem
	var maiorOrdem int
	queryMaiorNumero := `SELECT COALESCE(MAX(ordem_apresentacao), 0) FROM anexos WHERE manual_id = ?`

	err := DB.QueryRow(queryMaiorNumero, anexo.Manual_id).Scan(&maiorOrdem)
	if err != nil{
		return err
	}
	var proxOrdem = maiorOrdem + 1


	query := "INSERT INTO anexos(nome, tamanho, caminho, tipo_arquivo, manual_id, ordem_apresentacao) VALUES(?, ?, ?, ?, ?, ?)"
    _, err = DB.Exec(query, anexo.Nome, anexo.Tamanho_bytes, anexo.Caminho, anexo.Tipo_arquivo, anexo.Manual_id, proxOrdem)
    return err
}
func ReordenarAnexoData(idAnexo int, direcao string) error{
	var ordemAtual int
	var manualId int

	queryAnexoAtual := `SELECT ordem_apresentacao, manual_id FROM anexos WHERE id = ?`
	err := DB.QueryRow(queryAnexoAtual, idAnexo).Scan(&ordemAtual, &manualId)
	if err != nil{
		return err
	}
	if direcao == "subir"{
		var ordemProx int
		var proxAnexoID int

		queryProxAnexo := `SELECT ordem_apresentacao, id FROM anexos WHERE manual_id = ? AND ordem_apresentacao < ?
		ORDER BY ordem_apresentacao DESC LIMIT 1`
		
		err = DB.QueryRow(queryProxAnexo, manualId, ordemAtual).Scan(&ordemProx, &proxAnexoID)

		queryUpdateAnexoAtual := `UPDATE anexos
		SET ordem_apresentacao = ?
		WHERE id = ?`

		_, err = DB.Exec(queryUpdateAnexoAtual, ordemProx, idAnexo)

		if err != nil{
			return err
		}
		queryUpdateProxAnexo := `UPDATE anexos
		SET ordem_apresentacao = ?
		WHERE id = ?`

		_, err = DB.Exec(queryUpdateProxAnexo, ordemAtual, proxAnexoID)

		if err != nil{
			return err
		}
		return nil
	}

	if direcao == "descer"{
		var ordemAnterior int
		var anexoAnteriorID int

		queryAnexoAnterior := `SELECT ordem_apresentacao, id FROM anexos WHERE manual_id = ? AND ordem_apresentacao > ?
		ORDER BY ordem_apresentacao ASC LIMIT 1`

		err = DB.QueryRow(queryAnexoAnterior, manualId, ordemAtual).Scan(&ordemAnterior, &anexoAnteriorID)
		if err != nil{
			return err
		}
		queryUpdateAnexoAtual := `UPDATE anexos
		SET ordem_apresentacao = ?
		WHERE id = ?`

		_, err = DB.Exec(queryUpdateAnexoAtual, ordemAnterior, idAnexo)

		if err != nil{
			return err
		}
		queryUpdateAnexoAnterior := `UPDATE anexos
		SET ordem_apresentacao = ?
		WHERE id = ?`

		_, err = DB.Exec(queryUpdateAnexoAnterior, ordemAtual, anexoAnteriorID)

		if err != nil{
			return err
		}
		
		
	}

	return nil
} 