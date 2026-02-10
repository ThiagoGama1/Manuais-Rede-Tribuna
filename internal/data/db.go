package data

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)
var DB *sql.DB
func LoadDatabase(){
	var err error
	DB, err = sql.Open("sqlite3", "./manuais.db")
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

	if _, err := DB.Exec(sqlCreateTable); err != nil{
		log.Panic("Erro ao criar tabela:", err)
	}
	
	log.Println("Banco de dados inicializado com sucesso!")

}
 
