package http

import (
	"encoding/json"
	"net/http"

	"Hospital-Midderware/internal/application"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hospital string `json:"hospital"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type StaffHandler struct {
	authService *application.AuthService
}

func NewStaffHandler(authService *application.AuthService) *StaffHandler {
	return &StaffHandler{
		authService: authService,
	}
}

func (h *StaffHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request body"})
		return
	}

	token, err := h.authService.Login(req.Username, req.Password, req.Hospital)
	if err != nil {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, LoginResponse{Token: token})
}

func respondJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
