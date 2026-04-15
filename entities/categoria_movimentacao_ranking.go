package entities

// CategoriaMovimentacaoRanking ranking de categorias pelo volume de itens registrados em movimentacao_historico.
type CategoriaMovimentacaoRanking struct {
	Posicao           int    `json:"posicao"`
	ID                int64  `json:"id"`
	Nome              string `json:"nome"`
	TotalMovimentacao int64  `json:"total_movimentacao"`
}
