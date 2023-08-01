package models

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string      `json:"Message"`
	Data    interface{} `json:"Data,omitempty"`
}

func HandleResponse(w http.ResponseWriter, code int, message string, data interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Message: message,
		Data:    data,
	})
}
