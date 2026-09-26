package domain_test

import (
	"testing"

	"Hospital-Midderware/internal/domain"
)

func TestSearchCriteria_Validate_Success(t *testing.T) {
	criteria := domain.SearchCriteria{
		NationalID: "1100200300400",
	}

	if err := criteria.Validate(); err != nil {
		t.Errorf("คาดหวังว่าผ่านการ Validate แต่เกิด error: %v", err)
	}
}

func TestSearchCriteria_Validate_MissingAllCriteria_ShouldFail(t *testing.T) {
	criteria := domain.SearchCriteria{} // ว่างเปล่าทุก field

	err := criteria.Validate()
	if err == nil {
		t.Error("คาดหวัง error เนื่องจากไม่ได้ระบุเงื่อนไขการค้นหาเลย")
	}
}
