package app

import (
	"flag"
	"fmt"
)

func (c *cli) parseFlags() {
	fs := flag.NewFlagSet("validator", flag.ContinueOnError)
	fs.SetOutput(c.out)
	if err := fs.Parse(c.args); err != nil {
		fmt.Fprintf(c.out, "Argument parsing error")
		return
	}
	c.fs = fs
}

func (c *cli) Run() {
	c.parseFlags()
	args := c.fs.Args()
	if len(args) == 0 {
		fmt.Fprintf(c.out, "No references provided")
		return
	}

	for _, arg := range args {
		c.processReference(arg)
	}
}

func (c *cli) processReference(ref string) {
	valid, err := c.core.Validate(ref)
	if err != nil {
		fmt.Fprintf(c.out, "Validation failed for \"%s\":\nError: %s", ref, err)
		return
	}

	fmt.Fprintf(c.out, "[%t] %s\n", valid, ref)
}
