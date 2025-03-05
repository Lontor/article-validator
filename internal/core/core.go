package core

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
