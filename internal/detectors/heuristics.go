package detectors

import (
	"math"
	"regexp"
	"strings"
	"unicode"

	"github.com/graphsentinel/graphsentinel/internal/ingestion"
	"github.com/graphsentinel/graphsentinel/pkg/models"
)

// IdentifierRenamingDetector defines the contract for identifier obfuscation detection.
type IdentifierRenamingDetector interface {
	Detect(prepared ingestion.PreparedCode) models.IdentifierRenamingOutput
}

// DeadCodeDetector defines the contract for dead-code pattern detection.
type DeadCodeDetector interface {
	Detect(prepared ingestion.PreparedCode) models.DeadCodeOutput
}

// ControlFlowDetector defines the contract for control-flow drift detection.
type ControlFlowDetector interface {
	Detect(prepared ingestion.PreparedCode) models.ControlFlowOutput
}

// Run executes MVP text-structure heuristics that proxy real detector signals.
func Run(prepared ingestion.PreparedCode) models.DetectorOutputs {
	identifierDetector := HeuristicIdentifierRenamingDetector{}
	deadCodeDetector := HeuristicDeadCodeDetector{}
	controlFlowDetector := HeuristicControlFlowDetector{}
	return models.DetectorOutputs{
		IdentifierRenaming: identifierDetector.Detect(prepared),
		DeadCode:           deadCodeDetector.Detect(prepared),
		ControlFlow:        controlFlowDetector.Detect(prepared),
	}
}

var identRE = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_]*`)

var keywordSet = map[string]struct{}{
	"if": {}, "else": {}, "for": {}, "while": {}, "switch": {}, "case": {}, "break": {}, "continue": {},
	"return": {}, "goto": {}, "try": {}, "catch": {}, "throw": {}, "new": {}, "class": {}, "struct": {},
	"func": {}, "function": {}, "int": {}, "char": {}, "long": {}, "short": {}, "float": {}, "double": {},
	"bool": {}, "void": {}, "static": {}, "const": {}, "public": {}, "private": {}, "protected": {},
	"package": {}, "import": {}, "var": {}, "let": {}, "true": {}, "false": {}, "null": {}, "nil": {},
}

// HeuristicIdentifierRenamingDetector implements MVP identifier entropy heuristics.
type HeuristicIdentifierRenamingDetector struct{}

// Detect scores suspicious short/random/repetitive naming patterns.
func (HeuristicIdentifierRenamingDetector) Detect(p ingestion.PreparedCode) models.IdentifierRenamingOutput {
	ids := extractIdentifiers(p.Raw)
	total := len(ids)
	if total == 0 {
		return models.IdentifierRenamingOutput{}
	}

	shortCount := 0
	repetitiveCount := 0
	randomLikeCount := 0
	for _, id := range ids {
		if len(id) <= 2 {
			shortCount++
		}
		if isRepetitive(id) {
			repetitiveCount++
		}
		if looksRandom(id) {
			randomLikeCount++
		}
	}

	shortRate := float64(shortCount) / float64(total)
	repetitiveRate := float64(repetitiveCount) / float64(total)
	randomLikeRate := float64(randomLikeCount) / float64(total)
	score := clamp01((0.45 * shortRate) + (0.30 * randomLikeRate) + (0.25 * repetitiveRate))

	return models.IdentifierRenamingOutput{
		Likely: score >= 0.20,
		Score:  score,
	}
}

// HeuristicDeadCodeDetector implements MVP dead-code marker heuristics.
type HeuristicDeadCodeDetector struct{}

// Detect scans for classic dead/unreachable constructs and estimates a normalized score.
func (HeuristicDeadCodeDetector) Detect(p ingestion.PreparedCode) models.DeadCodeOutput {
	lower := p.Lower
	strongHits := 0
	weakHits := 0

	strongPatterns := []string{
		"if(false)",
		"if (false)",
		"while(false)",
		"while (false)",
		"if(0)",
		"if (0)",
		"while(0)",
		"while (0)",
		"return;",
		"throw;",
	}
	for _, marker := range strongPatterns {
		strongHits += strings.Count(lower, marker)
	}

	weakPatterns := []string{
		"unreachable",
		"never executed",
		"dead code",
		"dummy branch",
		"if constexpr(false)",
		"if constexpr (false)",
		"if(false){",
		"if (false){",
	}
	for _, marker := range weakPatterns {
		weakHits += strings.Count(lower, marker)
	}

	// Look for clearly unused temp-style declarations often injected by obfuscators.
	for _, marker := range []string{"unused_", "dummy_", "tmp_unused"} {
		weakHits += strings.Count(lower, marker)
	}

	rawScore := (float64(strongHits) * 0.45) + (float64(weakHits) * 0.15)
	score := clamp01(rawScore)
	return models.DeadCodeOutput{
		Likely: score >= 0.35,
		Score:  score,
	}
}

// HeuristicControlFlowDetector implements MVP branch-inflation and nesting heuristics.
type HeuristicControlFlowDetector struct{}

// Detect estimates control-flow drift by combining branch density, nesting pressure, and repeated branching.
func (HeuristicControlFlowDetector) Detect(p ingestion.PreparedCode) models.ControlFlowOutput {
	lower := p.Lower
	branchCount := 0
	for _, kw := range []string{
		"if(", "if (", "else if", "switch", "case ", "goto ", "while(", "while (", "for(", "for (", "catch", "?:",
	} {
		branchCount += strings.Count(lower, kw)
	}
	if branchCount == 0 {
		return models.ControlFlowOutput{}
	}

	lineCount := p.LineCount
	if lineCount <= 0 {
		lineCount = 1
	}
	branchDensity := clamp01(float64(branchCount) / float64(lineCount))

	maxNesting := maxBraceNesting(p.Raw)
	nestingPressure := clamp01(float64(maxNesting-2) / 6.0)

	repeated := repeatedBranchingHits(lower)
	repetitionScore := clamp01(float64(repeated) / 4.0)

	score := clamp01((0.45 * branchDensity) + (0.35 * nestingPressure) + (0.20 * repetitionScore))
	return models.ControlFlowOutput{
		Likely: score >= 0.35,
		Score:  score,
	}
}

func extractIdentifiers(src string) []string {
	raw := identRE.FindAllString(src, -1)
	out := make([]string, 0, len(raw))
	for _, tok := range raw {
		l := strings.ToLower(tok)
		if _, isKeyword := keywordSet[l]; isKeyword {
			continue
		}
		if len(tok) < 2 {
			continue
		}
		out = append(out, tok)
	}
	return out
}

func isRepetitive(id string) bool {
	if len(id) < 4 {
		return false
	}
	uniq := map[rune]struct{}{}
	for _, r := range strings.ToLower(id) {
		uniq[r] = struct{}{}
	}
	ratio := float64(len(uniq)) / float64(len([]rune(id)))
	return ratio <= 0.35
}

func looksRandom(id string) bool {
	if len(id) < 6 {
		return false
	}
	low := strings.ToLower(id)
	containsDigit := false
	containsVowel := false
	for _, r := range low {
		if unicode.IsDigit(r) {
			containsDigit = true
		}
		if strings.ContainsRune("aeiou", r) {
			containsVowel = true
		}
	}

	ent := shannonEntropy(low)
	if ent >= 3.2 && (containsDigit || !containsVowel) {
		return true
	}
	return false
}

func shannonEntropy(s string) float64 {
	if s == "" {
		return 0
	}
	freq := map[rune]float64{}
	for _, r := range s {
		freq[r]++
	}
	total := float64(len([]rune(s)))
	h := 0.0
	for _, c := range freq {
		p := c / total
		h -= p * math.Log2(p)
	}
	return h
}

func maxBraceNesting(src string) int {
	depth := 0
	maxDepth := 0
	for _, r := range src {
		switch r {
		case '{':
			depth++
			if depth > maxDepth {
				maxDepth = depth
			}
		case '}':
			if depth > 0 {
				depth--
			}
		}
	}
	return maxDepth
}

func repeatedBranchingHits(lower string) int {
	hits := 0
	// Repeated marker chains often come from opaque predicate insertion.
	for _, marker := range []string{
		"if (", "if(", "else if", "case ",
	} {
		c := strings.Count(lower, marker)
		if c >= 3 {
			hits++
		}
	}
	// Bonus signal for obvious dummy branching tokens.
	for _, marker := range []string{
		"dummy branch", "opaque predicate",
	} {
		if strings.Contains(lower, marker) {
			hits++
		}
	}
	return hits
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
