package main

import (
	"flag"
	"github.com/sheik/create/pkg/create"
	"strings"
)

func main() {
	flag.Parse()
	create.InteractiveCommand("go run Createfile " + strings.Join(flag.Args(), " "))
}
