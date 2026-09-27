package infrastructure

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"Hospital-Middleware/internal/domain"
)

type HospitalAAPIAdapter struct {
	baseURL    string
	httpClient *http.Client
}

func NewHospitalAAPIAdapter(baseURL string) *HospitalAAPIAdapter {
	if baseURL == "" {
		baseURL = "https://hospital-a.api.co.th"
	}
	return &HospitalAAPIAdapter{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *HospitalAAPIAdapter) Search(criteria domain.SearchCriteria) (*domain.PatientDTO, error) {
	if err := criteria.Validate(); err != nil {
		return nil, fmt.Errorf("invalid search criteria: %w", err)
	}

	targetID := criteria.NationalID
	if targetID == "" {
		targetID = criteria.PassportID
	}

	targetURL := strings.TrimRight(a.baseURL, "/") + "/patient/search"
	if targetID != "" {
		targetURL += "/" + url.PathEscape(targetID)
	}

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถสร้าง HTTP request ได้ครับ: %w", err)
	}

	query := req.URL.Query()
	if criteria.FirstName != "" {
		query.Set("first_name", criteria.FirstName)
	}
	if criteria.MiddleName != "" {
		query.Set("middle_name", criteria.MiddleName)
	}
	if criteria.LastName != "" {
		query.Set("last_name", criteria.LastName)
	}
	if criteria.DateOfBirth != "" {
		query.Set("date_of_birth", criteria.DateOfBirth)
	}
	if criteria.PhoneNumber != "" {
		query.Set("phone_number", criteria.PhoneNumber)
	}
	if criteria.Email != "" {
		query.Set("email", criteria.Email)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("การเชื่อมต่อ external API ล้มเหลว: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 4. จัดการ Response Status Code
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Patient not found (สอดคล้องกับ Rule D05-17)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("โรงพยาบาล A API ส่งค่า Status กลับมาไม่สำเร็จ: %d", resp.StatusCode)
	}

	// 5. Decode Response JSON เข้า PatientDTO
	var patient domain.PatientDTO
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		return nil, fmt.Errorf("ไม่สามารถแปลงข้อมูล JSON ของผู้ป่วยได้: %w", err)
	}

	return &patient, nil
}
