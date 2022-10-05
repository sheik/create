package main

import (
	"fmt"
	"github.com/sheik/create/pkg/docker"
	"github.com/sheik/create/pkg/git"
	"github.com/sheik/create/pkg/plan"
	"github.com/sheik/create/pkg/shell"
)

var (
	project           = "plan"
	buildVersion      = shell.Output("grep VERSION builder/Dockerfile | cut -d'=' -f2")
	imageName         = "builder:" + buildVersion
	dockerRun         = "docker run -h builder --rm -v $PWD:/code " + imageName
	dockerInteractive = "docker run -h builder --rm -v $PWD:/code -it " + imageName
	version           = shell.Output("git describe --tags | sed 's/-/_/g'")
	newVersion        = git.IncrementMinorVersion(version)
	rpm               = fmt.Sprintf("%s-%s-1.x86_64.rpm", project, version)
)

var steps = plan.Steps{
	"clean": plan.Step{
		Command: "rm -rf plan *.rpm usr Createfile",
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
		Help:    "plan the docker container used for building",
	},
	"build": plan.Step{
		Command: dockerRun + " go build ./cmd/plan",
		Gate:    git.RepoClean,
		Check:   shell.Bash("stat plan &>/dev/null"),
		Depends: plan.Complete("pull_build_image"),
		Help:    "build the go binary",
	},
	"pre_package": plan.Step{
		Command: "rm -rf usr && mkdir -p usr/local/bin && cp plan usr/local/bin",
		Help:    "prepare dir structure for packaging",
	},
	"package": plan.Step{
		Command: fmt.Sprintf("%s fpm --vendor CREATE -v %s -s dir -t rpm -n plan usr", dockerRun, version),
		Check:   shell.Bash(fmt.Sprintf("stat %s &>/dev/null", rpm)),
		Depends: plan.Complete("pull_build_image", "build", "pre_package"),
		Help:    "plan rpm",
		Default: true,
	},
	"commit": plan.Step{
		Command: "git commit -a -m \":INPUT:\"",
		Help:    "plan a git commit",
	},
	"publish": plan.Step{
		Command: fmt.Sprintf("git tag %s ; git push ; git push origin %s", newVersion, newVersion),
		Depends: plan.Complete("commit"),
		Help:    "commit, tag, and push code to repo",
	},
	"shell": plan.Step{
		Command:     dockerInteractive + " /bin/bash",
		Interactive: true,
		Depends:     plan.Complete("pull_build_image"),
		Help:        "open a shell in the build container",
	},
}

func main() {
	plan.Run(steps)
}
