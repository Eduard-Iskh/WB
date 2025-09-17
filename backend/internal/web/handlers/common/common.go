package common

import (
	"encoding/json"
	"net/http"
)

const (
	errorStatus   = "error"
	successStatus = "success"
)

type ErrorResponseStruct struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type SuccessResponseStruct struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func ErrorResponse(w http.ResponseWriter, err string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponseStruct{Status: errorStatus, Error: err})
}

func SuccessResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponseStruct{Status: successStatus, Data: data})
}
