package entities

// Usuario representa a tabela usuario.
type Usuario struct {
	ID    int64  `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
	Senha string `json:"senha"`
	Salt  string `json:"salt"`
}
