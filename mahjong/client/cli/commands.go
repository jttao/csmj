package cli

import (
	"github.com/codegangsta/cli"
)

var (
	commands = make([]cli.Command, 0)
)

func appendCmd(cmd cli.Command) {
	commands = append(commands, cmd)
}
