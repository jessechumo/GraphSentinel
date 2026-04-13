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
