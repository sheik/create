package main

import (
	"github.com/sheik/create/pkg/create"
	"github.com/sheik/create/pkg/shell"
	"os"
	"strings"
)

func main() {
	var args string
	if len(os.Args) > 1 {
		args = strings.Join(os.Args[1:], " ")
		if os.Args[1] == "update" {
			create.Plan(create.Steps{"update": create.UpdateStep})
			return
		}
	}
	shell.InteractiveCommand("go run ./cmd/createfile " + args)
}
