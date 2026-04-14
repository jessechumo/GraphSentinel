package models

// GetAnalysisResponse is the JSON body for GET /analysis/{id}.
type GetAnalysisResponse struct {
	AnalysisID string              `json:"analysis_id"`
	Status     AnalysisStatus      `json:"status"`
	Language   string              `json:"language,omitempty"`
	Signals    *ObfuscationSignals `json:"signals,omitempty"`
	Metrics    *AnalysisMetrics    `json:"metrics,omitempty"`
	Summary    string              `json:"summary,omitempty"`
	Error      string              `json:"error,omitempty"`
}

// NewGetAnalysisResponse maps a stored job to the public analysis API shape.
func NewGetAnalysisResponse(j *AnalysisJob) GetAnalysisResponse {
	if j == nil {
		return GetAnalysisResponse{}
	}
	out := GetAnalysisResponse{
		AnalysisID: j.ID,
		Status:     j.Status,
		Language:   j.Request.Language,
	}
	switch j.Status {
	case StatusFailed:
		out.Error = j.Error
	case StatusCompleted:
		if j.Report != nil {
			if j.Report.Language != "" {
				out.Language = j.Report.Language
			}
			out.Signals = j.Report.Signals
			out.Metrics = j.Report.Metrics
			out.Summary = j.Report.Summary
			if j.Report.Error != "" {
				out.Error = j.Report.Error
			}
		}
	}
	return out
}
