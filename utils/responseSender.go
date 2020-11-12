package utils

import (
	"bookList/model"
	"encoding/json"
	"net/http"
)

func SendError(w http.ResponseWriter, statusCode int, err model.ErrorMessage) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}
