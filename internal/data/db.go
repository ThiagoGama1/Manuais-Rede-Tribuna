package data

import (
	"database/sql"
	"log"

	_ "github.com/glebarez/go-sqlite"
)
var DB *sql.DB
func LoadDatabase(){
	var err error
	DB, err = sql.Open("sqlite", "./manuais.db")
		if err != nil {
			log.Panic("Erro ao conectar no SQLite:", err)
		}
		
		if err = DB.Ping(); err != nil{
			log.Panic("Banco de dados não responde: ", err)
		}

	sqlCreateTable := `
	CREATE TABLE IF NOT EXISTS manuais (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	titulo TEXT,
	secao TEXT
	);`

	sqlCreateTable2 := `
	CREATE TABLE IF NOT EXISTS etapas(
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	manual_id INTEGER,
	ordem_apresentacao INTEGER NOT NULL,
	conteudo TEXT,
	FOREIGN KEY(manual_id) REFERENCES manuais(id)
	);`
	
	sqlCreateTable3 := `
	CREATE TABLE IF NOT EXISTS anexos(
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	nome TEXT,
	tamanho INTEGER,
	caminho TEXT,
	tipo_arquivo TEXT,
	etapa_id INTEGER,
	ordem_apresentacao INTEGER NOT NULL,
	FOREIGN KEY(etapa_id) REFERENCES etapas(id)
	);`


	if _, err := DB.Exec(sqlCreateTable); err != nil{
		log.Panic("Erro ao criar tabela:", err)
	}
	if _, err := DB.Exec(sqlCreateTable2); err != nil{
		log.Panic("Erro ao criar tabela de etapas!:", err)
	}
	if _, err := DB.Exec(sqlCreateTable3); err != nil{
		log.Panic("Erro ao criar tabela de anexos!", err)
	}
	
	log.Println("Banco de dados inicializado com sucesso!")

}
 
