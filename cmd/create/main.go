package main

import (
	"github.com/sheik/create/pkg/shell"
	"os"
	"strings"
)

func main() {
	shell.InteractiveCommand("go build -o Createfile ./cmd/createfile")
	var args string
	if len(os.Args) > 1 {
		args = strings.Join(os.Args[1:], " ")
	}
	shell.InteractiveCommand("./Createfile " + args)
}
