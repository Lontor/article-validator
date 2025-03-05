package parser

import "github.com/Lontor/article-validator/internal/core"

type DefaultParser struct{}

func New() *DefaultParser {
	return &DefaultParser{}
}

func (p *DefaultParser) Parse(raw string) (core.Reference, error) {
	return core.Reference{Title: raw}, nil
}
