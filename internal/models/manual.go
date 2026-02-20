package models

type Manual struct {
	ID     int     `json:"id"`
	Titulo string  `json:"titulo"`
	Secao  string  `json:"secao"`
	Etapas []Etapa `json:"etapas"` //agora manuais contem etapas, e etapas contem anexos
}