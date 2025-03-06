package main

import (
	"fmt"

	"github.com/Lontor/article-validator/internal/apis"
	"github.com/Lontor/article-validator/internal/core"
	"github.com/Lontor/article-validator/internal/parser"
)

func main() {
	p := parser.New()
	client := apis.New("https://api.semanticscholar.org/graph/v1/paper/search/match", 3)
	core := core.New(p, []core.APIClient{client})

	valid, _ := core.Validate("Hanke M. On the shape derivative of polygonal inclusions in the conductivity problem")
	fmt.Printf("Valid: %v\n", valid)
}
