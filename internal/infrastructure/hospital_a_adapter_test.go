package infrastructure_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"Hospital-Middleware/internal/domain"
	"Hospital-Middleware/internal/infrastructure"
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

func TestHospitalAAPIAdapter_Search_ByOtherCriteria(t *testing.T) {
	expectedQuery := url.Values{
		"first_name":    {"สมชาย ใจดี"},
		"middle_name":   {"ก."},
		"last_name":     {"ทดสอบ"},
		"date_of_birth": {"1990-01-02"},
		"phone_number":  {"+66 81 234 5678"},
		"email":         {"somchai+test@example.com"},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("คาดหวัง GET แต่ได้ %s", r.Method)
		}
		if r.URL.Path != "/patient/search" {
			t.Errorf("คาดหวัง path /patient/search แต่ได้ %s", r.URL.Path)
		}
		if got := r.URL.Query(); !equalURLValues(got, expectedQuery) {
			t.Errorf("query parameters ไม่ตรงกัน: got %v, want %v", got, expectedQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(domain.PatientDTO{FirstNameTH: "สมชาย"})
	}))
	defer mockServer.Close()

	adapter := infrastructure.NewHospitalAAPIAdapter(mockServer.URL)
	patient, err := adapter.Search(domain.SearchCriteria{
		FirstName:   "สมชาย ใจดี",
		MiddleName:  "ก.",
		LastName:    "ทดสอบ",
		DateOfBirth: "1990-01-02",
		PhoneNumber: "+66 81 234 5678",
		Email:       "somchai+test@example.com",
	})
	if err != nil {
		t.Fatalf("ค้นหาด้วยเงื่อนไขอื่นไม่สำเร็จ: %v", err)
	}
	if patient == nil || patient.FirstNameTH != "สมชาย" {
		t.Fatalf("ผลลัพธ์ผู้ป่วยไม่ถูกต้อง: %+v", patient)
	}
}

func TestHospitalAAPIAdapter_Search_IDWithAdditionalCriteria(t *testing.T) {
	expectedQuery := url.Values{"email": {"somchai+test@example.com"}}
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/patient/search/1100200300400" {
			t.Errorf("path ไม่ตรงกัน: got %s", r.URL.Path)
		}
		if got := r.URL.Query(); !equalURLValues(got, expectedQuery) {
			t.Errorf("query parameters ไม่ตรงกัน: got %v, want %v", got, expectedQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(domain.PatientDTO{NationalID: "1100200300400"})
	}))
	defer mockServer.Close()

	adapter := infrastructure.NewHospitalAAPIAdapter(mockServer.URL)
	_, err := adapter.Search(domain.SearchCriteria{
		NationalID: "1100200300400",
		Email:      "somchai+test@example.com",
	})
	if err != nil {
		t.Fatalf("ค้นหาด้วย ID และเงื่อนไขเพิ่มเติมไม่สำเร็จ: %v", err)
	}
}

func TestHospitalAAPIAdapter_Search_PassportID(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/patient/search/A12345678" {
			t.Errorf("path ไม่ตรงกัน: got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(domain.PatientDTO{PassportID: "A12345678"})
	}))
	defer mockServer.Close()

	adapter := infrastructure.NewHospitalAAPIAdapter(mockServer.URL)
	_, err := adapter.Search(domain.SearchCriteria{PassportID: "A12345678"})
	if err != nil {
		t.Fatalf("ค้นหาด้วย passport ID ไม่สำเร็จ: %v", err)
	}
}

func TestHospitalAAPIAdapter_Search_EmptyCriteria(t *testing.T) {
	adapter := infrastructure.NewHospitalAAPIAdapter("http://unused.example")
	if _, err := adapter.Search(domain.SearchCriteria{}); err == nil {
		t.Fatal("คาดหวัง error เมื่อไม่มีเงื่อนไขค้นหา")
	}
}

func equalURLValues(got, want url.Values) bool {
	if len(got) != len(want) {
		return false
	}
	for key, wantValues := range want {
		gotValues, ok := got[key]
		if !ok || len(gotValues) != len(wantValues) {
			return false
		}
		for i := range wantValues {
			if gotValues[i] != wantValues[i] {
				return false
			}
		}
	}
	return true
}
