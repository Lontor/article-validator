package parser

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Lontor/article-validator/internal/core"
)

var (
	ErrInvalidReference = errors.New("invalid reference: authors or title not found")
)

type DefaultParser struct{}

func New() *DefaultParser {
	return &DefaultParser{}
}

func (p *DefaultParser) Parse(raw string) (core.Reference, error) {
	ref := core.Reference{}

	authorBlocksRegex := regexp.MustCompile(`(\p{Lu}+\p{Ll}*[\s.,]+(\p{Lu}[.,]\s*){1,3}(,\s*)?)+`)
	authorBlocksMatches := authorBlocksRegex.FindAllString(raw, -1)
	surnamesRegex := regexp.MustCompile(`\pL{2,}`)

	for _, a := range authorBlocksMatches {
		surnames := surnamesRegex.FindAllString(a, -1)
		ref.Authors = append(ref.Authors, surnames...)
	}

	titlePart := raw
	if len(authorBlocksMatches) > 0 {
		titlePart = raw[len(authorBlocksMatches[0]):]
	}

	splitChars := []string{"/", "–", "."}
	for _, sep := range splitChars {
		if idx := strings.Index(titlePart, sep); idx > 0 {
			titlePart = titlePart[:idx]
			break
		}
	}

	titleRegex := regexp.MustCompile(`\p{Lu}.*\p{L}`)
	ref.Title = strings.TrimSpace(titleRegex.FindString(titlePart))

	if len(ref.Title) == 0 || len(ref.Authors) == 0 {
		return ref, ErrInvalidReference
	}

	return ref, nil
}
