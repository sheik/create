package main

import (
	"github.com/sheik/create/pkg/create"
	"os"
	"strings"
)

func main() {
	create.InteractiveCommand("go build -o Createfile ./cmd/createfile")
	var args string
	if len(os.Args) > 1 {
		args = strings.Join(os.Args[1:], " ")
	}
	create.InteractiveCommand("./Createfile " + args)
}
