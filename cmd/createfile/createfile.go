package main

import (
	"fmt"
	"github.com/sheik/create/pkg/create"
	"github.com/sheik/create/pkg/docker"
	"github.com/sheik/create/pkg/git"
)

var (
	project           = "create"
	buildVersion      = create.Output("grep VERSION builder/Dockerfile | cut -d'=' -f2")
	imageName         = "builder:" + buildVersion
	dockerRun         = "docker run -h builder --rm -v $PWD:/code " + imageName
	dockerInteractive = "docker run -h builder --rm -v $PWD:/code -it " + imageName
	version           = create.Output("git describe --tags | sed 's/-/_/g'")
	newVersion        = git.IncrementMinorVersion(version)
	rpm               = fmt.Sprintf("%s-%s-1.x86_64.rpm", project, version)
)

var steps = create.Steps{
	"clean": create.Step{
		Command: "rm -rf create *.rpm usr Createfile",
		Help:    "clean build artifacts from repo",
	},
	"pull_build_image": create.Step{
		Command: fmt.Sprintf("docker pull %s", imageName),
		Check:   docker.ImageExists(imageName),
		Fail:    "build_image",
	},
	"build_image": create.Step{
		Command: fmt.Sprintf("docker build . -f builder/Dockerfile --tag %s", imageName),
		Check:   docker.ImageExists(imageName),
		Help:    "create the docker container used for building",
	},
	"build": create.Step{
		Command: dockerRun + " go build ./cmd/create",
		Gate:    git.RepoClean,
		Check:   create.Bash("stat create &>/dev/null"),
		Depends: create.Complete("pull_build_image"),
		Help:    "build the go binary",
	},
	"tag": create.Step{
		Command: fmt.Sprintf("git tag %s ; git push origin %s", newVersion, newVersion),
		Help:    "create a new minor tag and push",
	},
	"pre-package": create.Step{
		Command: "rm -rf usr && mkdir -p usr/local/bin && cp create usr/local/bin",
		Help:    "prepare dir structure for packaging",
	},
	"package": create.Step{
		Command: fmt.Sprintf("%s fpm --vendor CREATE -v %s -s dir -t rpm -n create usr", dockerRun, version),
		Check:   create.Bash(fmt.Sprintf("stat %s &>/dev/null", rpm)),
		Depends: create.Complete("pull_build_image", "build", "pre-package"),
		Help:    "create rpm",
		Default: true,
	},
	"commit": create.Step{
		Command: "git commit -a -m \":INPUT:\"",
	},
	"publish": create.Step{
		Command: "git push " + newVersion,
		Depends: create.Complete("commit", "tag"),
	},
	"shell": create.Step{
		Command:     dockerInteractive + " /bin/bash",
		Interactive: true,
		Depends:     create.Complete("pull_build_image"),
		Help:        "open a shell in the build container",
	},
}

func main() {
	create.Plan(steps)
}
