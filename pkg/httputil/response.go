package httputil

import (
	"context"
	"encoding/json"
	"ktrain/pkg/storage"
	"log"
	"net/http"
	"strconv"
)

type baseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func respondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(statusCode)
	_, _ = w.Write(data)
}

func RespondSuccess(w http.ResponseWriter, httpStatusCode int, payload interface{}, message string) {
	base := &baseResponse{
		Success: true,
		Message: message,
		Data:    payload,
	}
	respondJSON(w, httpStatusCode, base)
}

func RespondSuccessWithData(w http.ResponseWriter, httpStatusCode int, payload interface{}) {
	RespondSuccess(w, httpStatusCode, payload, "")
}

func RespondSuccessWithMessage(w http.ResponseWriter, httpStatusCode int, message string) {
	RespondSuccess(w, httpStatusCode, nil, message)
}

func RespondError(w http.ResponseWriter, httpStatusCode int, message string) {
	base := &baseResponse{
		Success: false,
		Message: message,
		Data:    nil,
	}
	respondJSON(w, httpStatusCode, base)
}

func DisconnectDb(mongDB *storage.MongoDBManager, ctx context.Context) {
	if err := mongDB.Disconnect(ctx); err != nil {
		log.Fatalf("Error when disconect database, err: %v", err)
		return
	}
}
