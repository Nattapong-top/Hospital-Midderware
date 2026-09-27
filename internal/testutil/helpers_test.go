package testutil_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"Hospital-Middleware/internal/domain"
	"Hospital-Middleware/internal/testutil"
)

func TestSetupGinTestContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	c, rec := testutil.SetupGinTestContext(req)
	if c == nil || rec == nil {
		t.Fatal("expected non-nil context and recorder")
	}
}

func TestMockAdapterAndResolver(t *testing.T) {
	patient := &domain.PatientDTO{PatientHN: "HN001"}
	adapter := &testutil.MockAdapter{Patient: patient}

	dto, err := adapter.Search(domain.SearchCriteria{NationalID: "123"})
	if err != nil || dto != patient {
		t.Errorf("MockAdapter unexpected result: %v, %v", dto, err)
	}

	resolver := &testutil.MockResolver{Adapter: adapter}
	resAdapter, err := resolver.Resolve("HOSP_A")
	if err != nil || resAdapter != adapter {
		t.Errorf("MockResolver unexpected result: %v, %v", resAdapter, err)
	}

	errResolver := &testutil.MockResolver{Err: errors.New("resolver error")}
	_, err = errResolver.Resolve("HOSP_A")
	if err == nil {
		t.Error("expected error from MockResolver")
	}
}
