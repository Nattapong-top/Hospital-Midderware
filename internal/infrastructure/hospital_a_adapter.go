package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	targetID := criteria.NationalID
	if targetID == "" {
		targetID = criteria.PassportID
	}

	if targetID == "" {
		return nil, errors.New("การค้นหาผู้ป่วยของโรงพยาบาล A จำเป็นต้องระบุเลขประจำตัวประชาชน (national_id) หรือหนังสือเดินทาง (passport_id) ครับ")
	}

	// 2. สร้าง Request URL ตาม Spec /patient/search/{id}
	url := fmt.Sprintf("%s/patient/search/%s", a.baseURL, targetID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ไม่สามารถสร้าง HTTP request ได้ครับ: %w", err)
	}

	// 3. ยิง HTTP Request
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
