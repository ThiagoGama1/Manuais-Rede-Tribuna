package models

type Etapa struct {
	Id                 int     `json:"id"`
	Manual_id          int     `json:"manual_id"`
	Anexos             []Anexo `json:"anexos"`
	Ordem_apresentacao int     `json:"ordem_apresentacao"`
	Conteudo           string  `json:"conteudo"`
}
