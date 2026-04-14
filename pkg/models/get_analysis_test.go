package models

import (
	"testing"
	"time"
)

func TestNewGetAnalysisResponse_completed(t *testing.T) {
	t.Parallel()
	j := &AnalysisJob{
		ID:        "abc",
		Status:    StatusCompleted,
		Request:   AnalyzeRequest{Language: "c", Code: "x"},
		CreatedAt: time.Now().UTC(),
		Report: &AnalysisReport{
			AnalysisID: "abc",
			Status:     StatusCompleted,
			Language:   "c",
			Summary:    "done",
		},
	}
	j.Report.Signals = &ObfuscationSignals{}
	j.Report.Metrics = &AnalysisMetrics{}

	got := NewGetAnalysisResponse(j)
	if got.Status != StatusCompleted || got.Summary != "done" {
		t.Fatalf("%+v", got)
	}
}

func TestNewGetAnalysisResponse_failed(t *testing.T) {
	t.Parallel()
	j := &AnalysisJob{
		ID:      "abc",
		Status:  StatusFailed,
		Request: AnalyzeRequest{Language: "go", Code: "x"},
		Error:   "boom",
	}
	got := NewGetAnalysisResponse(j)
	if got.Error != "boom" || got.Status != StatusFailed {
		t.Fatalf("%+v", got)
	}
}
