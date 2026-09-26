package http

import (
	"Hospital-Midderware/internal/application"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
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
	authService  *application.AuthService
	staffService *application.StaffService
}

type CreateStaffRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hospital string `json:"hospital"`
}

type CreateStaffResponse struct {
	Username string `json:"username"`
	Hospital string `json:"hospital"`
}

func NewStaffHandler(authService *application.AuthService, staffService *application.StaffService) *StaffHandler {
	return &StaffHandler{
		authService:  authService,
		staffService: staffService,
	}
}

func (h *StaffHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	token, err := h.authService.Login(req.Username, req.Password, req.Hospital)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *StaffHandler) CreateStaff(c *gin.Context) {
	var req CreateStaffRequest
	// 1. Bind JSON จาก Request Body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return
	}

	// 2. เรียก Business Logic ใน Application Layer
	err := h.staffService.CreateStaff(c.Request.Context(), application.CreateStaffRequest(req))
	if err != nil {
		// สามารถเช็ก Error Type คืนค่า 400 Bad Request หรือ 409 Conflict ตาม Logic ได้
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Response Success HTTP 201 Created
	c.JSON(http.StatusCreated, CreateStaffResponse{
		Username: req.Username,
		Hospital: req.Hospital,
	})
}

func respondJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
