package main

import (
	"fmt"
	"github.com/sheik/create/pkg/build"
	"github.com/sheik/create/pkg/docker"
	"github.com/sheik/create/pkg/git"
	"github.com/sheik/create/pkg/plan"
	"github.com/sheik/create/pkg/shell"
)

const project = "create"

var (
	imageName  = "builder:" + shell.Output("grep VERSION builder/Dockerfile | cut -d'=' -f2")
	version    = git.GetVersion()
	newVersion = git.IncrementMinorVersion(version)
	rpm        = build.RPMFilename(project, version)
	container  = docker.Container(imageName).Mount("$PWD", "/code")
)

var steps = plan.Steps{
	"clean": plan.Step{
		Command: "rm -rf create *.rpm usr createfile.go",
		Help:    "clean build artifacts from repo",
	},
	"pull_build_image": plan.Step{
		Command: fmt.Sprintf("docker pull %s", imageName),
		Check:   docker.ImageExists(imageName),
		Fail:    "build_image",
		Help:    "pull build image from docker registry",
	},
	"build_image": plan.Step{
		Command: fmt.Sprintf("docker build . -f builder/Dockerfile --tag %s", imageName),
		Check:   docker.ImageExists(imageName),
		Help:    "create the docker container used for building",
	},
	"parser": plan.Step{
		Command: "peg -noast -switch -inline -strict -output cmd/create/parser.go grammar/createfile.peg",
		Help:    "generate createfile parser from grammar file",
	},
	"build": plan.Step{
		Command: container.Run("go build ./cmd/create"),
		Check:   shell.ReturnZero("stat create &>/dev/null"),
		Depends: plan.Complete("pull_build_image", "parser"),
		Help:    "build the go binary",
	},
	"pre_package": plan.Step{
		Command: "rm -rf usr && mkdir -p usr/local/bin && cp create usr/local/bin",
		Depends: plan.Complete("build"),
	},
	"package": plan.Step{
		Command: container.Run("fpm --vendor CREATE -v %s -s dir -t rpm -n create usr", version),
		Check:   shell.ReturnZero(fmt.Sprintf("stat %s &>/dev/null", rpm)),
		Depends: plan.Complete("pre_package"),
		Help:    "create rpm",
		Default: true,
	},
	"commit": plan.Step{
		Command: "git commit -a -m \":INPUT:\"",
		Help:    "create a git commit",
	},
	"publish": plan.Step{
		Command: fmt.Sprintf("git tag %s ; git push ; git push origin %s", newVersion, newVersion),
		Gate:    git.RepoClean(),
		Help:    "commit, tag, and push code to repo",
	},
	"shell": plan.Step{
		Command:     container.Interactive().Run("/bin/bash"),
		Interactive: true,
		Depends:     plan.Complete("pull_build_image"),
		Help:        "open a shell in the build container",
	},
}

func main() {
	plan.Run(steps)
}
