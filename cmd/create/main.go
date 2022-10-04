package main

import (
	"fmt"
	"github.com/sheik/create/pkg/create"
)

var (
	docker            = "docker run -v $PWD:/code builder:1"
	dockerInteractive = "docker run -v $PWD:/code -it builder:1"
	version           = create.Output("git describe --tags")
)

var steps = create.Steps{
	"clean": create.Step{
		Command: "rm -rf create *.rpm usr",
	},
	"build_container": create.Step{
		Command: "docker build . --tag builder:1",
		Check:   `bash -c "if [[ \"$(docker images -q builder:1)\" == \"\" ]]; then exit 1; else exit 0; fi"`,
	},
	"build": create.Step{
		Command: "go build ./cmd/create",
		Check:   "stat create &>/dev/null",
		Depends: create.Complete("build_container"),
	},
	"pre-package": create.Step{
		Command: "rm -rf usr && mkdir -p usr/local/bin && cp create usr/local/bin",
	},
	"package": create.Step{
		Command: fmt.Sprintf("%s fpm --vendor CREATE -v %s -s dir -t rpm -n create usr", docker, version),
		Check:   "stat *.rpm &>/dev/null",
		Depends: create.Complete("build_container", "build", "pre-package"),
		Default: true,
	},
	"shell": create.Step{
		Command:     fmt.Sprintf("%s /bin/bash", dockerInteractive),
		Interactive: true,
	},
}

func main() {
	create.Plan(steps)
}
