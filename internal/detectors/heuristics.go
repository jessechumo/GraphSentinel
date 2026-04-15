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

// Run executes MVP text-structure heuristics that proxy real detector signals.
func Run(prepared ingestion.PreparedCode) models.DetectorOutputs {
	identifierDetector := HeuristicIdentifierRenamingDetector{}
	return models.DetectorOutputs{
		IdentifierRenaming: identifierDetector.Detect(prepared),
		DeadCode:           models.DeadCodeOutput{},
		ControlFlow:        models.ControlFlowOutput{},
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

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
