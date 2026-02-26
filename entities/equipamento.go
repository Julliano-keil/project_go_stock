package entities

// Equipamento representa a tabela equipamento.
type Equipamento struct {
	ID               int64  `json:"id"`
	Nome             string `json:"nome"`
	IDSubCategoria   int64  `json:"id_sub_categoria"`
	IDUnidadeEstoque int64  `json:"id_unidade_estoque"`
}
