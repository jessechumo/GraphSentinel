package analyzers

import (
	"context"
	"strings"
	"testing"

	"github.com/graphsentinel/graphsentinel/pkg/models"
)

func TestAnalyze_returnsCompletedReport(t *testing.T) {
	t.Parallel()

	job := &models.AnalysisJob{
		ID: "job-1",
		Request: models.AnalyzeRequest{
			Language: "c",
			Code: `int main() {
  if (false) { return 2; }
  goto done;
done:
  return 0;
}`,
		},
	}

	rep, err := Analyze(context.Background(), job)
	if err != nil {
		t.Fatal(err)
	}
	if rep == nil || rep.Status != models.StatusCompleted {
		t.Fatalf("report = %+v", rep)
	}
	if rep.AnalysisID != "job-1" || rep.Language != "c" {
		t.Fatalf("report identity mismatch: %+v", rep)
	}
	if rep.Signals == nil || rep.Metrics == nil {
		t.Fatalf("missing report fields: %+v", rep)
	}
	if rep.Summary == "" {
		t.Fatal("expected summary")
	}
}

func TestAnalyze_nilJob(t *testing.T) {
	t.Parallel()
	rep, err := Analyze(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil job")
	}
	if rep != nil {
		t.Fatalf("report = %+v, want nil", rep)
	}
}

func TestSummarize(t *testing.T) {
	t.Parallel()

	none := summarize(models.DetectorOutputs{})
	if none == "" {
		t.Fatal("empty fallback summary")
	}

	some := summarize(models.DetectorOutputs{
		IdentifierRenaming: models.IdentifierRenamingOutput{Likely: true, Score: 0.8},
		DeadCode:           models.DeadCodeOutput{Likely: true, Score: 0.9},
		ControlFlow:        models.ControlFlowOutput{Likely: true, Score: 0.7},
	})
	if some == none {
		t.Fatalf("summary did not change: %q", some)
	}
	if !strings.Contains(some, "Likely obfuscation via") {
		t.Fatalf("unexpected summary format: %q", some)
	}
}
