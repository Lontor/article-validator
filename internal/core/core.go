package core

import (
	"fmt"
	"sync"
)

type ValidationResult struct {
	APIName string
	Valid   bool
	Error   error
}

type ValidationResponse struct {
	Results []ValidationResult
	IsValid bool
}

type Reference struct {
	Title   string
	Authors []string
}

type Parser interface {
	Parse(raw string) (Reference, error)
}

type APIClient interface {
	Name() string
	Validate(ref Reference) (bool, error)
}

type Core struct {
	parser  Parser
	clients []APIClient
}

func New(parser Parser, clients []APIClient) *Core {
	return &Core{parser: parser, clients: clients}
}

func (c *Core) Validate(rawReference string) (*ValidationResponse, error) {
	ref, err := c.parser.Parse(rawReference)
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	result := ValidationResponse{Results: make([]ValidationResult, len(c.clients)), IsValid: false}

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(len(c.clients))

	for i, client := range c.clients {
		go func(idx int, cl APIClient) {
			defer wg.Done()

			valid, err := cl.Validate(ref)

			result.Results[idx] = ValidationResult{
				APIName: cl.Name(),
				Valid:   valid,
				Error:   err,
			}

			if valid {
				mu.Lock()
				result.IsValid = true
				mu.Unlock()
			}
		}(i, client)
	}

	wg.Wait()

	return &result, nil
}
