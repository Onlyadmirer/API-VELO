package utils

import (
	"encoding/json"
	"net/http"
)

type Responses struct {
	Status  int
	Message string
	Data    any
}

func ResponseSuccess(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	resp := Responses{
		Status:  status,
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(resp)
}

func ResponseError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := Responses{
		Status:  status,
		Message: message,
	}
	json.NewEncoder(w).Encode(resp)
}
