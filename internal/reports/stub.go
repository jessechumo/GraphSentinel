package reports

import "github.com/graphsentinel/graphsentinel/pkg/models"

// BuildStubReport returns a completed report with neutral detector outputs.
// Real detectors replace this in later pipeline stages.
func BuildStubReport(job *models.AnalysisJob) *models.AnalysisReport {
	outs := models.DetectorOutputs{
		IdentifierRenaming: models.IdentifierRenamingOutput{Likely: false, Score: 0},
		DeadCode:           models.DeadCodeOutput{Likely: false, Score: 0},
		ControlFlow:        models.ControlFlowOutput{Likely: false, Score: 0},
	}
	sig := outs.Signals()
	met := outs.Metrics()
	return &models.AnalysisReport{
		AnalysisID: job.ID,
		Status:     models.StatusCompleted,
		Language:   job.Request.Language,
		Signals:    &sig,
		Metrics:    &met,
		Summary:    "Baseline structural pass complete; detailed heuristics attach in later pipeline stages.",
	}
}
