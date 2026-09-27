package http_test

import (
	"testing"

	"Hospital-Middleware/internal/domain"
)

func TestSearchCriteria_Validate_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		criteria    domain.SearchCriteria
		expectError bool
	}{
		{
			name:        "Valid with NationalID only",
			criteria:    domain.SearchCriteria{NationalID: "1100200300400"},
			expectError: false,
		},
		{
			name:        "Valid with PassportID only",
			criteria:    domain.SearchCriteria{PassportID: "A12345678"},
			expectError: false,
		},
		{
			name:        "Valid with FirstName only",
			criteria:    domain.SearchCriteria{FirstName: "Somchai"},
			expectError: false,
		},
		{
			name:        "Valid with MiddleName only",
			criteria:    domain.SearchCriteria{MiddleName: "Ananda"},
			expectError: false,
		},
		{
			name:        "Valid with LastName only",
			criteria:    domain.SearchCriteria{LastName: "Jaidee"},
			expectError: false,
		},
		{
			name:        "Valid with DateOfBirth only",
			criteria:    domain.SearchCriteria{DateOfBirth: "1990-01-01"},
			expectError: false,
		},
		{
			name:        "Valid with PhoneNumber only",
			criteria:    domain.SearchCriteria{PhoneNumber: "0812345678"},
			expectError: false,
		},
		{
			name:        "Valid with Email only",
			criteria:    domain.SearchCriteria{Email: "somchai@example.com"},
			expectError: false,
		},
		{
			name: "Valid with Multiple Criteria",
			criteria: domain.SearchCriteria{
				FirstName: "Somchai",
				LastName:  "Jaidee",
				Email:     "somchai@example.com",
			},
			expectError: false,
		},
		{
			name:        "Invalid with All Empty Fields",
			criteria:    domain.SearchCriteria{}, // ว่างเปล่าทุก Field
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.criteria.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("[%s] คาดหวังว่าจะเกิด Error เนื่องจากไม่มีการระบุเงื่อนไขการค้นหา แต่ได้ nil", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("[%s] คาดหวังว่าผ่านการ Validate แต่เกิด Error: %v", tt.name, err)
				}
			}
		})
	}
}
