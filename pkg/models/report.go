package models

// ObfuscationSignals is the boolean outcome of each MVP detector.
type ObfuscationSignals struct {
	IdentifierRenaming bool `json:"identifier_renaming"`
	DeadCode           bool `json:"dead_code"`
	ControlFlowChange  bool `json:"control_flow_change"`
}

// AnalysisMetrics holds normalized scores in [0,1] from detectors.
type AnalysisMetrics struct {
	IdentifierEntropyScore float64 `json:"identifier_entropy_score"`
	DeadCodeScore          float64 `json:"dead_code_score"`
	ControlFlowDriftScore  float64 `json:"control_flow_drift_score"`
}

// DetectorOutputs is the structured result of running all detectors before summary generation.
type DetectorOutputs struct {
	IdentifierRenaming IdentifierRenamingOutput
	DeadCode           DeadCodeOutput
	ControlFlow        ControlFlowOutput
}

// IdentifierRenamingOutput is the identifier-focused detector result.
type IdentifierRenamingOutput struct {
	Likely bool
	Score  float64
}

// DeadCodeOutput is the dead-code detector result.
type DeadCodeOutput struct {
	Likely bool
	Score  float64
}

// ControlFlowOutput is the control-flow drift detector result.
type ControlFlowOutput struct {
	Likely bool
	Score  float64
}

// Signals maps detector outputs to API signal booleans.
func (d DetectorOutputs) Signals() ObfuscationSignals {
	return ObfuscationSignals{
		IdentifierRenaming: d.IdentifierRenaming.Likely,
		DeadCode:           d.DeadCode.Likely,
		ControlFlowChange:  d.ControlFlow.Likely,
	}
}

// Metrics maps detector outputs to API metric fields.
func (d DetectorOutputs) Metrics() AnalysisMetrics {
	return AnalysisMetrics{
		IdentifierEntropyScore: d.IdentifierRenaming.Score,
		DeadCodeScore:          d.DeadCode.Score,
		ControlFlowDriftScore:  d.ControlFlow.Score,
	}
}

// AnalysisReport is the machine-readable outcome returned for a completed (or failed) analysis.
type AnalysisReport struct {
	AnalysisID string              `json:"analysis_id"`
	Status     AnalysisStatus      `json:"status"`
	Language   string              `json:"language,omitempty"`
	Signals    *ObfuscationSignals `json:"signals,omitempty"`
	Metrics    *AnalysisMetrics    `json:"metrics,omitempty"`
	Summary    string              `json:"summary,omitempty"`
	Error      string              `json:"error,omitempty"`
}
