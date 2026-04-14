package store

import (
	"testing"

	"github.com/graphsentinel/graphsentinel/pkg/models"
)

func TestMemory_EnqueueGet(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	req := models.AnalyzeRequest{Language: "go", Code: "package main"}

	j, err := s.Enqueue(req)
	if err != nil {
		t.Fatal(err)
	}
	if j.ID == "" || j.Status != models.StatusQueued {
		t.Fatalf("job = %+v", j)
	}
	if j.Request.Language != "go" {
		t.Fatalf("request language = %q", j.Request.Language)
	}

	got, ok := s.Get(j.ID)
	if !ok || got.ID != j.ID {
		t.Fatalf("Get = %v, ok=%v", got, ok)
	}
}

func TestMemory_Get_missing(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	_, ok := s.Get("nope")
	if ok {
		t.Fatal("expected miss")
	}
}

func TestMemory_Snapshot(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	j, err := s.Enqueue(models.AnalyzeRequest{Language: "go", Code: "x"})
	if err != nil {
		t.Fatal(err)
	}
	got, ok := s.Snapshot(j.ID)
	if !ok || got.ID != j.ID || got.Status != models.StatusQueued {
		t.Fatalf("snapshot = %+v", got)
	}
}

func TestMemory_claimComplete(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	j, err := s.Enqueue(models.AnalyzeRequest{Language: "go", Code: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if !s.TryClaim(j.ID) {
		t.Fatal("expected claim")
	}
	got, ok := s.Get(j.ID)
	if !ok || got.Status != models.StatusRunning {
		t.Fatalf("job = %+v", got)
	}
	rep := &models.AnalysisReport{AnalysisID: j.ID, Status: models.StatusCompleted, Language: "go"}
	if err := s.Complete(j.ID, rep); err != nil {
		t.Fatal(err)
	}
	got, ok = s.Get(j.ID)
	if !ok || got.Status != models.StatusCompleted || got.Report == nil {
		t.Fatalf("job = %+v", got)
	}
}

func TestMemory_completeWrongState(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	j, err := s.Enqueue(models.AnalyzeRequest{Language: "go", Code: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Complete(j.ID, &models.AnalysisReport{}); err != ErrJobNotRunning {
		t.Fatalf("err = %v", err)
	}
}

func TestMemory_failWrongState(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	j, err := s.Enqueue(models.AnalyzeRequest{Language: "go", Code: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Fail(j.ID, "x"); err != ErrJobNotRunning {
		t.Fatalf("err = %v", err)
	}
}
