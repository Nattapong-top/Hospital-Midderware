package infrastructure_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Hospital-Midderware/internal/domain"
	"Hospital-Midderware/internal/infrastructure"
)

func TestHospitalAAPIAdapter_Search_Success(t *testing.T) {
	// 1. Mock Server ของ Hospital A
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/patient/search/1100200300400" {
			t.Errorf("URL Path ไม่ถูกต้อง ได้: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		patient := domain.PatientDTO{
			FirstNameTH: "สมชาย",
			LastNameTH:  "ใจดี",
			NationalID:  "1100200300400",
			PatientHN:   "HN-12345",
			Gender:      "M",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(patient)
		if err != nil {
			return
		}
	}))
	defer mockServer.Close()

	// 2. เรียกใช้ Adapter โดยส่ง Mock Server URL เข้าไป
	adapter := infrastructure.NewHospitalAAPIAdapter(mockServer.URL)
	criteria := domain.SearchCriteria{NationalID: "1100200300400"}

	patient, err := adapter.Search(criteria)

	// 3. Assert Results
	if err != nil {
		t.Fatalf("คาดหวังว่าค้นหาสำเร็จ แต่ได้ error: %v", err)
	}
	if patient == nil {
		t.Fatal("คาดหวังว่าจะเจอข้อมูลผู้ป่วย แต่ได้ nil")
	}
	if patient.FirstNameTH != "สมชาย" {
		t.Errorf("คาดหวังชื่อ 'สมชาย' แต่ได้ %s", patient.FirstNameTH)
	}
}

func TestHospitalAAPIAdapter_Search_NotFound(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	adapter := infrastructure.NewHospitalAAPIAdapter(mockServer.URL)
	criteria := domain.SearchCriteria{NationalID: "9999999999999"}

	patient, err := adapter.Search(criteria)

	if err != nil {
		t.Fatalf("คาดหวัง error = nil เมื่อไม่พบผู้ป่วย แต่ได้: %v", err)
	}
	if patient != nil {
		t.Errorf("คาดหวัง patient = nil เมื่อไม่พบข้อมูล แต่ได้ object กลับมา")
	}
}
