package models

type Manual struct {
	ID         int    `json:"id"`
	Titulo     string `json:"titulo"`
	Conteudo   string `json:"conteudo"`
	Secao      string `json:"secao"`
	Arquivos []Anexo `json:"anexos"`
}