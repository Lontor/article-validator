package app

import (
	"flag"
	"io"
	"os"
)

type Core interface {
	Validate(string) (bool, error)
}

type cli struct {
	core Core
	args []string
	out  io.Writer
	fs   *flag.FlagSet
}

func New(core Core) *cli {
	return &cli{core: core, out: os.Stdout, args: os.Args[1:]}
}

func (c *cli) SetOutput(w io.Writer) {
	c.out = w
}

func (c *cli) SetArgs(args []string) {
	c.args = args
}
