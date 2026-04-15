package detectors

import (
	"math"
	"strings"

	"github.com/graphsentinel/graphsentinel/internal/ingestion"
	"github.com/graphsentinel/graphsentinel/pkg/models"
)

// Run executes MVP text-structure heuristics that proxy real detector signals.
func Run(prepared ingestion.PreparedCode) models.DetectorOutputs {
	return models.DetectorOutputs{
		IdentifierRenaming: detectIdentifierRenaming(prepared),
		DeadCode:           detectDeadCode(prepared),
		ControlFlow:        detectControlFlow(prepared),
	}
}

func detectIdentifierRenaming(p ingestion.PreparedCode) models.IdentifierRenamingOutput {
	shortTokens := 0
	totalTokens := 0
	for _, tok := range strings.Fields(p.Compact) {
		totalTokens++
		if len(tok) <= 2 {
			shortTokens++
		}
	}
	if totalTokens == 0 {
		return models.IdentifierRenamingOutput{}
	}
	score := clamp01(float64(shortTokens) / float64(totalTokens))
	return models.IdentifierRenamingOutput{
		Likely: score >= 0.30,
		Score:  score,
	}
}

func detectDeadCode(p ingestion.PreparedCode) models.DeadCodeOutput {
	lower := p.Lower
	hits := 0
	for _, marker := range []string{
		"if (false)",
		"if(false)",
		"while(false)",
		"for (;;)",
		"/* dead */",
		"unreachable",
	} {
		if strings.Contains(lower, marker) {
			hits++
		}
	}
	if hits == 0 {
		return models.DeadCodeOutput{}
	}
	score := clamp01(float64(hits) / 3.0)
	return models.DeadCodeOutput{
		Likely: score >= 0.34,
		Score:  score,
	}
}

func detectControlFlow(p ingestion.PreparedCode) models.ControlFlowOutput {
	lower := p.Lower
	branchTerms := 0
	for _, kw := range []string{
		"goto ",
		"switch",
		"continue",
		"break",
		"try",
		"catch",
	} {
		branchTerms += strings.Count(lower, kw)
	}
	lineFactor := 0.0
	if p.LineCount > 0 {
		lineFactor = float64(branchTerms) / float64(p.LineCount)
	}
	shapeFactor := 0.0
	if p.NewlineRate > 0 {
		shapeFactor = math.Min(1.0, p.NewlineRate*40.0)
	}
	score := clamp01((lineFactor * 0.7) + (shapeFactor * 0.3))
	return models.ControlFlowOutput{
		Likely: score >= 0.35,
		Score:  score,
	}
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
