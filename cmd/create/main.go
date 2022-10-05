package main

import (
	"github.com/sheik/create/pkg/plan"
	"github.com/sheik/create/pkg/shell"
	"os"
	"strings"
)

func main() {
	var args string
	if len(os.Args) > 1 {
		args = strings.Join(os.Args[1:], " ")
		if os.Args[1] == "update" {
			plan.Run(plan.Steps{"update": plan.UpdateStep})
			return
		}
	}
	shell.InteractiveCommand("go run ./cmd/createfile " + args)
}
