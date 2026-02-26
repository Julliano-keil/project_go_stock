package entities

// ErrorStruct representa a estrutura de erro padronizada da API.
type ErrorStruct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorReturn é o envelope de resposta de erro.
type ErrorReturn struct {
	Error ErrorStruct `json:"error"`
}
