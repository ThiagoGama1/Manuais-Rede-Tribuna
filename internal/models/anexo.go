package models

type Anexo struct {
	ID                int    `json:"id"`
	Nome              string `json:"nome"`
	Tamanho_bytes     int    `json:"tamanho"`
	Caminho           string `json:"caminho"`
	Tipo_arquivo      string `json:"tipo_arquivo"`
	Etapa_id          int    `json:"etapa_id"`
	OrdemApresentacao int    `json:"ordem_apresentacao"`
}