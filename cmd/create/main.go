package main

import (
	"flag"
	"github.com/sheik/create/pkg/create"
	"strings"
)

func main() {
	flag.Parse()
	create.InteractiveCommand("go build -o Createfile ./cmd/create")
	create.InteractiveCommand("./Createfile " + strings.Join(flag.Args(), " "))
}
