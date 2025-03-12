package main

import (
	"github.com/Lontor/article-validator/internal/apis"
	"github.com/Lontor/article-validator/internal/cli"
	"github.com/Lontor/article-validator/internal/core"
	"github.com/Lontor/article-validator/internal/parser"
)

func main() {
	p := parser.New()
	semanticScholar := apis.NewSemanticScholarClient(
		"https://api.semanticscholar.org/graph/v1/paper/search/match", 3, nil)
	crossref := apis.NewCrossrefClient(
		"https://api.crossref.org/works", 3, nil)
	core := core.New(p, []core.APIClient{semanticScholar, crossref})

	cli := cli.New(core)
	cli.Run()

}
