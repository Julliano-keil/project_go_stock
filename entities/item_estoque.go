package entities

// ItemEstoque representa a tabela item_estoque.
type ItemEstoque struct {
	ID               int64  `json:"id"`
	IDEquipamento    int64  `json:"id_equipamento"`
	IDUnidadeEstoque int64  `json:"id_unidade_estoque"`
	StatusCode       string `json:"status_code"`
	Codigo           string `json:"codigo"`
}
