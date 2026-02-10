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
	id INT NOT NULL PRIMARY KEY, 
	titulo TEXT,
	conteudo TEXT,
	secao TEXT,
	arquivo TEXT
	);`

	sqlCreateTable2 := `
	CREATE TABLE IF NOT EXISTS anexos(
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	nome TEXT,
	tamanho INT,
	caminho TEXT,
	tipo_arquivo TEXT,
	manual_id INT,
	FOREIGN KEY(manual_id) REFERENCES manuais(id)
	);`

	if _, err := DB.Exec(sqlCreateTable); err != nil{
		log.Panic("Erro ao criar tabela:", err)
	}
	if _, err := DB.Exec(sqlCreateTable2); err != nil{
		log.Panic("Erro ao criar tabela de anexos:", err)
	}
	
	log.Println("Banco de dados inicializado com sucesso!")

}
 
