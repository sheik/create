package main

import (
	"fmt"
	"github.com/sheik/create/pkg/create"
)

var (
	buildVersion      = create.Output("grep VERSION Dockerfile | cut -d'=' -f2")
	imageName         = "builder:" + buildVersion
	docker            = "docker run -h builder --rm -v $PWD:/code " + imageName
	dockerInteractive = "docker run -h builder --rm -v $PWD:/code -it " + imageName
	version           = create.Output("git describe --tags")
)

var steps = create.Steps{
	"clean": create.Step{
		Command: "rm -rf create *.rpm usr",
		Help:    "clean build artifacts from repo",
	},
	"build_container": create.Step{
		Command: fmt.Sprintf("docker build . --tag %s", imageName),
		Check:   fmt.Sprintf(`bash -c "if [[ \"$(docker images -q %s)\" == \"\" ]]; then exit 1; else exit 0; fi"`, imageName),
		Help:    "create the docker container used for building",
	},
	"build": create.Step{
		Command: docker + " go build ./cmd/create",
		Check:   "stat create &>/dev/null",
		Depends: create.Complete("build_container"),
		Help:    "build the go binary",
	},
	"pre-package": create.Step{
		Command: "rm -rf usr && mkdir -p usr/local/bin && cp create usr/local/bin",
		Help:    "prepare dir structure for packaging",
	},
	"package": create.Step{
		Command: fmt.Sprintf("%s fpm --vendor CREATE -v %s -s dir -t rpm -n create usr", docker, version),
		Check:   "stat *.rpm &>/dev/null",
		Depends: create.Complete("build_container", "build", "pre-package"),
		Default: true,
		Help:    "create rpm",
	},
	"shell": create.Step{
		Command:     dockerInteractive + " /bin/bash",
		Interactive: true,
		Depends:     create.Complete("build_container"),
		Help:        "open a shell in the build container",
	},
}

func main() {
	create.Plan(steps)
}
