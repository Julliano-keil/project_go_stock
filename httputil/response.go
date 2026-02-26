package httputil

import (
	"encoding/json"
	"net/http"

	"lince/entities"
)

// WriteError envia resposta de erro padronizada.
func WriteError(w http.ResponseWriter, err entities.ErrorStruct) {
	var errorCode int
	switch err.Code {
	case 1:
		errorCode = http.StatusInternalServerError
	case 2:
		errorCode = http.StatusForbidden
	case 3:
		errorCode = http.StatusUnauthorized
	case 4, 5:
		errorCode = http.StatusForbidden
	default:
		errorCode = http.StatusInternalServerError
	}

	wsError := entities.ErrorReturn{Error: err}
	jsonError, _ := json.Marshal(wsError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	_, _ = w.Write(jsonError)
}
