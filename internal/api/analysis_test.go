package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/graphsentinel/graphsentinel/internal/api"
	"github.com/graphsentinel/graphsentinel/internal/reports"
	"github.com/graphsentinel/graphsentinel/internal/store"
	"github.com/graphsentinel/graphsentinel/internal/workers"
	"github.com/graphsentinel/graphsentinel/pkg/models"
)

func TestGetAnalysis_notFound(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(api.NewRouter(store.NewMemory(), nil))
	t.Cleanup(srv.Close)

	res, err := http.Get(srv.URL + "/analysis/nope")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d", res.StatusCode)
	}
}

func TestGetAnalysis_eventuallyCompleted(t *testing.T) {
	t.Parallel()
	s := store.NewMemory()
	pool := workers.NewPool(2, 32, s, func(ctx context.Context, job *models.AnalysisJob) (*models.AnalysisReport, error) {
		return reports.BuildStubReport(job), nil
	})
	pool.Start()
	t.Cleanup(func() {
		pool.Close()
		pool.Wait()
	})

	srv := httptest.NewServer(api.NewRouter(s, pool.Submit))
	t.Cleanup(srv.Close)

	body := `{"language":"c","code":"int main(){return 0;}"}`
	sub, err := http.Post(srv.URL+"/analyze", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Body.Close()
	if sub.StatusCode != http.StatusAccepted {
		t.Fatalf("submit status = %d", sub.StatusCode)
	}
	var subResp models.SubmitAnalysisResponse
	if err := json.NewDecoder(sub.Body).Decode(&subResp); err != nil {
		t.Fatal(err)
	}
	id := subResp.AnalysisID

	deadline := time.Now().Add(2 * time.Second)
	var got models.GetAnalysisResponse
	for time.Now().Before(deadline) {
		res, err := http.Get(srv.URL + "/analysis/" + id)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			res.Body.Close()
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("get status = %d", res.StatusCode)
		}
		if got.Status == models.StatusCompleted {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}

	if got.Status != models.StatusCompleted {
		t.Fatalf("final status = %q", got.Status)
	}
	if got.AnalysisID != id || got.Language != "c" {
		t.Fatalf("response = %+v", got)
	}
	if got.Signals == nil || got.Metrics == nil || got.Summary == "" {
		t.Fatalf("incomplete report: %+v", got)
	}
}
