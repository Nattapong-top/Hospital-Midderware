package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"Hospital-Midderware/internal/domain"
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
	// 1. เลือก ID ที่จะเอาไปยิงใส่ URL Path (National ID หรือ Passport ID)
	targetID := criteria.NationalID
	if targetID == "" {
		targetID = criteria.PassportID
	}

	if targetID == "" {
		return nil, errors.New("hospital A search requires national_id or passport_id")
	}

	// 2. สร้าง Request URL
	url := fmt.Sprintf("%s/patient/search/%s", a.baseURL, targetID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 3. ยิง HTTP Request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("external API call failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 4. จัดการ Response Status Code
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Patient not found (สอดคล้องกับ Rule D05-17)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hospital A API returned status: %d", resp.StatusCode)
	}

	// 5. Decode Response JSON เข้า PatientDTO
	var patient domain.PatientDTO
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		return nil, fmt.Errorf("failed to decode patient data: %w", err)
	}

	return &patient, nil
}
