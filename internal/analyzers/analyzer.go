package analyzers

import (
	"context"
	"strings"

	"github.com/graphsentinel/graphsentinel/internal/detectors"
	"github.com/graphsentinel/graphsentinel/internal/ingestion"
	"github.com/graphsentinel/graphsentinel/pkg/models"
)

// Analyze runs normalized ingestion + detector heuristics and returns a report.
func Analyze(_ context.Context, job *models.AnalysisJob) (*models.AnalysisReport, error) {
	prepared := ingestion.Prepare(job.Request.Code)
	outs := detectors.Run(prepared)
	signals := outs.Signals()
	metrics := outs.Metrics()

	return &models.AnalysisReport{
		AnalysisID: job.ID,
		Status:     models.StatusCompleted,
		Language:   job.Request.Language,
		Signals:    &signals,
		Metrics:    &metrics,
		Summary:    summarize(outs),
	}, nil
}

func summarize(outs models.DetectorOutputs) string {
	flags := make([]string, 0, 3)
	if outs.IdentifierRenaming.Likely {
		flags = append(flags, "identifier renaming")
	}
	if outs.DeadCode.Likely {
		flags = append(flags, "dead code")
	}
	if outs.ControlFlow.Likely {
		flags = append(flags, "control-flow drift")
	}
	if len(flags) == 0 {
		return "No strong semantics-preserving transformation signals detected."
	}
	if len(flags) == 1 {
		return "Likely obfuscation via " + flags[0] + "."
	}
	return "Likely obfuscation via " + strings.Join(flags, ", ") + "."
}
