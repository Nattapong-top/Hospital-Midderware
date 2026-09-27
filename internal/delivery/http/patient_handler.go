package http

import (
	"net/http"

	"Hospital-Midderware/internal/application"
	"Hospital-Midderware/internal/domain"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	searchPatientUseCase *application.SearchPatient
}

func NewPatientHandler(searchPatientUseCase *application.SearchPatient) *PatientHandler {
	return &PatientHandler{
		searchPatientUseCase: searchPatientUseCase,
	}
}

func (h *PatientHandler) SearchPatient(c *gin.Context) {
	// 1. ดึง hospital_id จาก Context ที่ AuthMiddleware สกัดมาจาก JWT Token
	hospitalID := c.GetString("hospital")
	if hospitalID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: missing hospital context"})
		return
	}

	// 2. รับค่า Query Parameters ทั้ง 8 Fields ตามโจทย์ Agnos
	criteria := domain.SearchCriteria{
		NationalID:  c.Query("national_id"),
		PassportID:  c.Query("passport_id"),
		FirstName:   c.Query("first_name"),
		MiddleName:  c.Query("middle_name"),
		LastName:    c.Query("last_name"),
		DateOfBirth: c.Query("date_of_birth"),
		PhoneNumber: c.Query("phone_number"),
		Email:       c.Query("email"),
	}

	// 3. เรียก Use Case ให้ทำงาน
	result, err := h.searchPatientUseCase.Execute(hospitalID, criteria)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. ส่งผลลัพธ์กลับ 200 OK
	c.JSON(http.StatusOK, result)
}
