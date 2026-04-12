package models

import (
	"errors"
	"fmt"
	"strings"
)

const (
	maxLanguageLen = 32
	maxCodeBytes   = 256 << 10 // 256 KiB
)

// AnalyzeRequest is the payload for submitting source code for analysis.
type AnalyzeRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// NormalizeLanguage returns a canonical lowercase language tag for storage and reporting.
func NormalizeLanguage(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Validate checks required fields and size limits. Language is validated against a small allowlist.
func (r *AnalyzeRequest) Validate() error {
	if r == nil {
		return errors.New("request is required")
	}
	lang := NormalizeLanguage(r.Language)
	if lang == "" {
		return errors.New("language is required")
	}
	if len(lang) > maxLanguageLen {
		return fmt.Errorf("language exceeds maximum length of %d", maxLanguageLen)
	}
	if !isAllowedLanguage(lang) {
		return fmt.Errorf("unsupported language %q", lang)
	}
	if len(r.Code) > maxCodeBytes {
		return fmt.Errorf("code exceeds maximum size of %d bytes", maxCodeBytes)
	}
	code := strings.TrimSpace(r.Code)
	if code == "" {
		return errors.New("code is required")
	}
	return nil
}

// Normalized returns a copy with trimmed code and normalized language for persistence.
func (r *AnalyzeRequest) Normalized() AnalyzeRequest {
	if r == nil {
		return AnalyzeRequest{}
	}
	return AnalyzeRequest{
		Language: NormalizeLanguage(r.Language),
		Code:     strings.TrimSpace(r.Code),
	}
}

func isAllowedLanguage(lang string) bool {
	_, ok := supportedLanguages[lang]
	return ok
}

// supportedLanguages is the MVP allowlist; extend as parsers are added.
var supportedLanguages = map[string]struct{}{
	"c":          {},
	"cpp":        {},
	"c++":        {},
	"java":       {},
	"go":         {},
	"python":     {},
	"javascript": {},
	"js":         {},
	"typescript": {},
	"ts":         {},
	"rust":       {},
	"csharp":     {},
	"cs":         {},
	"unknown":    {}, // explicit escape hatch for experiments
}
