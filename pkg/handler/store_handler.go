package handler

import (
	"context"
	"encoding/json"
	"errors"
	customError "github.com/ashwingopalsamy/backend-services/pkg/errors"
	"net/http"

	"github.com/ashwingopalsamy/backend-services/pkg/store/proto"
)

type StoreHandler struct {
	storeService proto.StoreServiceServer
}

func NewStoreHandler(service proto.StoreServiceServer) *StoreHandler {
	return &StoreHandler{
		storeService: service,
	}
}

func (h *StoreHandler) CreateStoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req proto.CreateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	resp, err := h.storeService.CreateStore(context.Background(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string, errors map[string]string) {
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"status":  "error",
		"message": message,
	}
	if errors != nil {
		response["errors"] = errors
	}
	json.NewEncoder(w).Encode(response)
}

// Handle errors returned by the service layer
func handleServiceError(w http.ResponseWriter, err error) {
	var customErr *customError.CustomError
	if errors.As(err, &customErr) {
		writeErrorResponse(w, customErr.StatusCode, customErr.Message, customErr.Details)
		return
	}

	// Fallback for unexpected errors
	writeErrorResponse(w, http.StatusInternalServerError, "An unexpected error occurred. Please try again later.", nil)
}
