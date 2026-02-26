package entities

import "time"

// Movimentacao representa a tabela movimentacao.
type Movimentacao struct {
	ID               int64     `json:"id"`
	TipoMovimentacao string    `json:"tipo_movimentacao"`
	Data             time.Time `json:"data"`
	IDUser           int64     `json:"id_user"`
}
