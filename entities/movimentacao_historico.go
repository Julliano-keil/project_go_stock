package entities

// MovimentacaoHistorico representa a tabela movimentacao_historico.
type MovimentacaoHistorico struct {
	ID              int64 `json:"id"`
	IDMovimentacao  int64 `json:"id_movimentacao"`
	IDItemEstoque   int64 `json:"id_item_estoque"`
}
