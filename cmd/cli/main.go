package main

import (
	"github.com/Lontor/article-validator/cmd/cli/app"
	"github.com/Lontor/article-validator/internal/apis"
	"github.com/Lontor/article-validator/internal/core"
	"github.com/Lontor/article-validator/internal/parser"
)

func main() {
	p := parser.New()
	client := apis.New("https://api.semanticscholar.org/graph/v1/paper/search/match", 3)
	core := core.New(p, []core.APIClient{client})

	cli := app.New(core)
	cli.Run()

}
