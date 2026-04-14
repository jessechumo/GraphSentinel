package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/graphsentinel/graphsentinel/internal/api"
	"github.com/graphsentinel/graphsentinel/internal/store"
	"github.com/graphsentinel/graphsentinel/pkg/models"
)

func TestAnalyze_acceptsAndQueues(t *testing.T) {
	t.Parallel()
	s := store.NewMemory()
	srv := httptest.NewServer(api.NewRouter(s, nil))
	t.Cleanup(srv.Close)

	body := `{"language":"c","code":"int main(){return 0;}"}`
	res, err := http.Post(srv.URL+"/analyze", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("status = %d", res.StatusCode)
	}
	var got models.SubmitAnalysisResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Status != models.StatusQueued || got.AnalysisID == "" {
		t.Fatalf("response = %+v", got)
	}
	j, ok := s.Get(got.AnalysisID)
	if !ok || j.Request.Language != "c" {
		t.Fatalf("stored job = %+v ok=%v", j, ok)
	}
}

func TestAnalyze_validationError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(api.NewRouter(store.NewMemory(), nil))
	t.Cleanup(srv.Close)

	body := `{"language":"c","code":""}`
	res, err := http.Post(srv.URL+"/analyze", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d", res.StatusCode)
	}
}

func TestAnalyze_bodyTooLarge(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(api.NewRouter(store.NewMemory(), nil))
	t.Cleanup(srv.Close)

	large := bytes.Repeat([]byte("a"), 600<<10)
	body := `{"language":"go","code":"` + string(large) + `"}`
	res, err := http.Post(srv.URL+"/analyze", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d", res.StatusCode)
	}
}
