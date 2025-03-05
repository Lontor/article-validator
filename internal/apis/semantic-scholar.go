package apis

import "github.com/Lontor/article-validator/internal/core"

type SemanticScholarClient struct{}

func New() *SemanticScholarClient {
	return &SemanticScholarClient{}
}

func (c *SemanticScholarClient) Validate(ref core.Reference) (bool, error) {
	return true, nil
}
