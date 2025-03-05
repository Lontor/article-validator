package core

import "fmt"

type Reference struct {
	Title   string
	Authors []string
}

type Parser interface {
	Parse(raw string) (Reference, error)
}

type APIClient interface {
	Validate(ref Reference) (bool, error)
}

type Core struct {
	parser  Parser
	clients []APIClient
}

func New(parser Parser, clients []APIClient) *Core {
	return &Core{parser: parser, clients: clients}
}

func (c *Core) Validate(rawReference string) (bool, error) {
	ref, err := c.parser.Parse(rawReference)
	if err != nil {
		return false, fmt.Errorf("parsing failed: %w", err)
	}

	for _, client := range c.clients {
		valid, _ := client.Validate(ref)
		if valid {
			return true, nil
		}
	}
	return false, nil
}
