package entities

// SubCategoria representa a tabela sub_categoria.
type SubCategoria struct {
	ID          int64  `json:"id"`
	IDCategoria int64  `json:"id_categoria"`
	Nome        string `json:"nome"`
}
