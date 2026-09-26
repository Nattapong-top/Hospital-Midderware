package application_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"Hospital-Midderware/internal/application"
	"Hospital-Midderware/internal/domain"
	"Hospital-Midderware/internal/infrastructure"
)

func TestSearchPatient_MultiHospital_Success(t *testing.T) {
	// 1. Mock Server สำหรับ Hospital A (ตอบกลับเป็น สมชาย)
	serverA := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		patient := domain.PatientDTO{
			FirstNameTH: "สมชาย",
			LastNameTH:  "ใจดี",
			NationalID:  "1111111111111",
			PatientHN:   "HN-HOSP-A",
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(patient)
		if err != nil {
			return
		}
	}))
	defer serverA.Close()

	// 2. Mock Server สำหรับ Hospital B (ตอบกลับเป็น สมหญิง)
	serverB := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		patient := domain.PatientDTO{
			FirstNameTH: "สมหญิง",
			LastNameTH:  "รักดี",
			NationalID:  "2222222222222",
			PatientHN:   "HN-HOSP-B",
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(patient)
		if err != nil {
			return
		}
	}))
	defer serverB.Close()

	// 3. Mock Server สำหรับ Hospital C (ตอบกลับเป็น สมศักดิ์)
	serverC := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		patient := domain.PatientDTO{
			FirstNameTH: "สมศักดิ์",
			LastNameTH:  "ภักดี",
			NationalID:  "3333333333333",
			PatientHN:   "HN-HOSP-C",
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(patient)
		if err != nil {
			return
		}
	}))
	defer serverC.Close()

	// 4. ผูก Adapters กับ Resolver สำหรับทั้ง 3 โรงพยาบาล
	adapterA := infrastructure.NewHospitalAAPIAdapter(serverA.URL)
	adapterB := infrastructure.NewHospitalAAPIAdapter(serverB.URL)
	adapterC := infrastructure.NewHospitalAAPIAdapter(serverC.URL)

	// สร้าง Map / Custom Resolver รองรับ 3 โรงพยาบาล
	multiResolver := &mockMultiHospitalResolver{
		adapters: map[string]domain.ExternalAPIAdapter{
			"HOSP_A": adapterA,
			"HOSP_B": adapterB,
			"HOSP_C": adapterC,
		},
	}

	useCase := application.NewSearchPatient(multiResolver)

	// --- TEST CASE 1: ค้นหาคนไข้ฝั่ง Hospital A ---
	t.Run("Search Hospital A", func(t *testing.T) {
		res, err := useCase.Execute("HOSP_A", domain.SearchCriteria{NationalID: "1111111111111"})
		if err != nil || res.Patient == nil {
			t.Fatalf("คาดหวังว่าเจอคนไข้ Hosp A แต่ได้ error: %v", err)
		}
		if res.Patient.FirstNameTH != "สมชาย" {
			t.Errorf("คาดหวังชื่อ 'สมชาย' แต่ได้: %s", res.Patient.FirstNameTH)
		}
	})

	// --- TEST CASE 2: ค้นหาคนไข้ฝั่ง Hospital B ---
	t.Run("Search Hospital B", func(t *testing.T) {
		res, err := useCase.Execute("HOSP_B", domain.SearchCriteria{NationalID: "2222222222222"})
		if err != nil || res.Patient == nil {
			t.Fatalf("คาดหวังว่าเจอคนไข้ Hosp B แต่ได้ error: %v", err)
		}
		if res.Patient.FirstNameTH != "สมหญิง" {
			t.Errorf("คาดหวังชื่อ 'สมหญิง' แต่ได้: %s", res.Patient.FirstNameTH)
		}
	})

	// --- TEST CASE 3: ค้นหาคนไข้ฝั่ง Hospital C ---
	t.Run("Search Hospital C", func(t *testing.T) {
		res, err := useCase.Execute("HOSP_C", domain.SearchCriteria{NationalID: "3333333333333"})
		if err != nil || res.Patient == nil {
			t.Fatalf("คาดหวังว่าเจอคนไข้ Hosp C แต่ได้ error: %v", err)
		}
		if res.Patient.FirstNameTH != "สมศักดิ์" {
			t.Errorf("คาดหวังชื่อ 'สมศักดิ์' แต่ได้: %s", res.Patient.FirstNameTH)
		}
	})
}

// Mock Resolver สำหรับรองรับหลาย Adapters ใน Test
type mockMultiHospitalResolver struct {
	adapters map[string]domain.ExternalAPIAdapter
}

func (r *mockMultiHospitalResolver) Resolve(hospitalID string) (domain.ExternalAPIAdapter, error) {
	adapter, exists := r.adapters[hospitalID]
	if !exists {
		return nil, fmt.Errorf("hospital id %s not found", hospitalID)
	}
	return adapter, nil
}
