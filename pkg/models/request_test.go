package models

import (
	"strings"
	"testing"
)

func TestAnalyzeRequest_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		req     AnalyzeRequest
		wantErr string
	}{
		{
			name:    "nil body treated as zero value",
			req:     AnalyzeRequest{},
			wantErr: "language is required",
		},
		{
			name: "valid c",
			req: AnalyzeRequest{
				Language: "c",
				Code:     "int main(){return 0;}",
			},
		},
		{
			name: "language normalized uppercase",
			req: AnalyzeRequest{
				Language: "C",
				Code:     "x",
			},
		},
		{
			name: "unsupported language",
			req: AnalyzeRequest{
				Language: "fortran",
				Code:     "program",
			},
			wantErr: `unsupported language "fortran"`,
		},
		{
			name: "empty code",
			req: AnalyzeRequest{
				Language: "go",
				Code:     "   ",
			},
			wantErr: "code is required",
		},
		{
			name: "code too large",
			req: AnalyzeRequest{
				Language: "go",
				Code:     strings.Repeat("a", maxCodeBytes+1),
			},
			wantErr: "code exceeds maximum size",
		},
		{
			name: "whitespace padding cannot bypass size limit",
			req: AnalyzeRequest{
				Language: "go",
				Code:     strings.Repeat(" ", maxCodeBytes) + "x",
			},
			wantErr: "code exceeds maximum size",
		},
		{
			name: "language too long",
			req: AnalyzeRequest{
				Language: strings.Repeat("a", maxLanguageLen+1),
				Code:     "x",
			},
			wantErr: "language exceeds maximum length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.req.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() = %v, want nil", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Validate() error = %v, want substring %q", err, tt.wantErr)
			}
		})
	}
}

func TestAnalyzeRequest_Normalized(t *testing.T) {
	t.Parallel()
	r := AnalyzeRequest{Language: "  Go  ", Code: "  package main "}
	n := r.Normalized()
	if n.Language != "go" {
		t.Fatalf("Language = %q, want go", n.Language)
	}
	if n.Code != "package main" {
		t.Fatalf("Code = %q", n.Code)
	}
	var nilReq *AnalyzeRequest
	if z := nilReq.Normalized(); z.Language != "" || z.Code != "" {
		t.Fatalf("nil Normalized = %+v, want zero", z)
	}
}

func TestDetectorOutputs_SignalsMetrics(t *testing.T) {
	t.Parallel()
	d := DetectorOutputs{
		IdentifierRenaming: IdentifierRenamingOutput{Likely: true, Score: 0.9},
		DeadCode:           DeadCodeOutput{Likely: false, Score: 0.1},
		ControlFlow:        ControlFlowOutput{Likely: true, Score: 0.5},
	}
	s := d.Signals()
	if !s.IdentifierRenaming || s.DeadCode || !s.ControlFlowChange {
		t.Fatalf("Signals = %+v", s)
	}
	m := d.Metrics()
	if m.IdentifierEntropyScore != 0.9 || m.DeadCodeScore != 0.1 || m.ControlFlowDriftScore != 0.5 {
		t.Fatalf("Metrics = %+v", m)
	}
}

func TestAnalysisStatus_Terminal(t *testing.T) {
	t.Parallel()
	if StatusQueued.Terminal() || StatusRunning.Terminal() {
		t.Fatal("non-terminal reported terminal")
	}
	if !StatusCompleted.Terminal() || !StatusFailed.Terminal() {
		t.Fatal("terminal statuses not terminal")
	}
}
