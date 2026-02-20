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

		for _, anexo := range etapa.Anexos{
			anexo.Etapa_id = int(idEtapa)
			InsertAnexo(anexo)
		}
	
	}
	return id, nil
}
func GetManuais()[]models.Manual{
	query := "SELECT id, titulo, secao FROM manuais"
	var list []models.Manual
	queryResult, err := DB.Query(query)
	
	if err != nil{
		return list
	}
	defer queryResult.Close()

	for queryResult.Next(){
		var m models.Manual
		queryResult.Scan(&m.ID, &m.Titulo, &m.Secao)
		list = append(list, m)
	}
	
	return list
}

func GetManualByID(id int) (models.Manual, error){
	query := `SELECT id, titulo, secao FROM manuais WHERE id = ?`
	var result models.Manual
	err := DB.QueryRow(query, id).Scan(&result.ID, &result.Titulo, &result.Secao)


	if err != nil{
		return models.Manual{}, err
	}
	queryEtapas := `SELECT id, manual_id, ordem_apresentacao, conteudo FROM etapas WHERE manual_id = ? ORDER BY ordem_apresentacao ASC`


	rows, err := DB.Query(queryEtapas, id)

	if err != nil{
		return result, err
	}
	defer rows.Close()

	for rows.Next(){
		var a models.Etapa

		err := rows.Scan(&a.Id, &a.Manual_id, &a.Ordem_apresentacao, &a.Conteudo)
		queryAnexos := `SELECT id, nome, tamanho, caminho, tipo_arquivo, etapa_id, ordem_apresentacao FROM anexos WHERE etapa_id = ?`

		rowsAnexos, err := DB.Query(queryAnexos, a.Id)
		if err != nil{
			return models.Manual{}, err
		}
		defer rowsAnexos.Close()

		for rowsAnexos.Next(){
			var novoAnexo models.Anexo
			rowsAnexos.Scan(&novoAnexo.ID, &novoAnexo.Nome, &novoAnexo.Tamanho_bytes, &novoAnexo.Caminho, &novoAnexo.Tipo_arquivo, &novoAnexo.Etapa_id, novoAnexo.OrdemApresentacao)
			a.Anexos = append(a.Anexos, novoAnexo)
		}
		
		result.Etapas = append(result.Etapas, a)
	}

	

	return result, nil
}

func DeleteManual(id int) error{
	manual, err := GetManualByID(id)
	if err != nil{
		return err
	}
	query := `DELETE FROM manuais WHERE id = ?`
	queryEtapas := `DELETE FROM etapas WHERE manual_id = ?`
	queryAnexos := `DELETE FROM anexos WHERE etapa_id = ?`
	
	for _, etapa := range manual.Etapas{
		for _, anexo := range etapa.Anexos{
			os.Remove(anexo.Caminho)
		}
		_, err = DB.Exec(queryAnexos, etapa.Id)
		if err != nil{
			return err
		}

	}
	_, err = DB.Exec(queryEtapas, manual.ID)
	if err != nil{
		return err
	}

	_, err = DB.Exec(query, manual.ID)
	
	return err
}

func UpdateManual(m models.Manual) error{
	query := `UPDATE manuais 
			  SET titulo = ?,
			  secao = ?
			  WHERE id = ?`

	_, err := DB.Exec(query, m.Titulo, m.Secao, m.ID)

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
	queryMaiorNumero := `SELECT COALESCE(MAX(ordem_apresentacao), 0) FROM anexos WHERE etapa_id = ?`

	err := DB.QueryRow(queryMaiorNumero, anexo.Etapa_id).Scan(&maiorOrdem)
	if err != nil{
		return err
	}
	var proxOrdem = maiorOrdem + 1


	query := "INSERT INTO anexos(nome, tamanho, caminho, tipo_arquivo, etapa_id, ordem_apresentacao) VALUES(?, ?, ?, ?, ?, ?)"
    _, err = DB.Exec(query, anexo.Nome, anexo.Tamanho_bytes, anexo.Caminho, anexo.Tipo_arquivo, anexo.Etapa_id, proxOrdem)
    return err
}
func ReordenarAnexoData(idAnexo int, direcao string) error{
	var ordemAtual int
	var etapaId int

	queryAnexoAtual := `SELECT ordem_apresentacao, etapa_id FROM anexos WHERE id = ?`
	err := DB.QueryRow(queryAnexoAtual, idAnexo).Scan(&ordemAtual, &etapaId)
	if err != nil{
		return err
	}
	if direcao == "subir"{
		var ordemProx int
		var proxAnexoID int

		queryProxAnexo := `SELECT ordem_apresentacao, id FROM anexos WHERE etapa_id = ? AND ordem_apresentacao < ?
		ORDER BY ordem_apresentacao DESC LIMIT 1`
		
		err = DB.QueryRow(queryProxAnexo, etapaId, ordemAtual).Scan(&ordemProx, &proxAnexoID)

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

		queryAnexoAnterior := `SELECT ordem_apresentacao, id FROM anexos WHERE etapa_id = ? AND ordem_apresentacao > ?
		ORDER BY ordem_apresentacao ASC LIMIT 1`

		err = DB.QueryRow(queryAnexoAnterior, etapaId, ordemAtual).Scan(&ordemAnterior, &anexoAnteriorID)
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