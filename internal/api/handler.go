package api

import (
	"github.com/graphsentinel/graphsentinel/internal/store"
)

type handler struct {
	store store.JobStore
}
