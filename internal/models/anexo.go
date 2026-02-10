package models

type Anexo struct {
	ID            int    `json:"id"`
	Nome          string `json:"nome"`
	Tamanho_bytes int    `json:"tamanho"`
	Caminho       string `json:"caminho"`
	Tipo_arquivo  string `json:"tipo_arquivo"`
	manual_id     int    `json:"idmanual"`
}