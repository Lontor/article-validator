package cli

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
	response, err := c.core.Validate(ref)
	if err != nil {
		fmt.Fprintf(c.out, "[PARSER ERROR] %s: %v\n", ref, err)
		return
	}

	status := "INVALID"
	if response.IsValid {
		status = "VALID"
	}
	fmt.Fprintf(c.out, "[%s] %s\n", status, ref)

	for _, res := range response.Results {
		status := "❌"
		if res.Valid {
			status = "✅"
		}

		errMsg := ""
		if res.Error != nil {
			errMsg = fmt.Sprintf(" | Error: %v", res.Error)
		}

		fmt.Fprintf(c.out, "  %s %s%s\n", status, res.APIName, errMsg)
	}
}
