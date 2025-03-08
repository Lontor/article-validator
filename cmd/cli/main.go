package main

import (
	"github.com/Lontor/article-validator/internal/apis"
	"github.com/Lontor/article-validator/internal/cli"
	"github.com/Lontor/article-validator/internal/core"
	"github.com/Lontor/article-validator/internal/parser"
)

func main() {
	p := parser.New()
	client := apis.NewSemanticScholarClient("https://api.semanticscholar.org/graph/v1/paper/search/match", 3)
	core := core.New(p, []core.APIClient{client})

	cli := cli.New(core)
	cli.Run()

}
