package docker

import (
	"fmt"
	"github.com/sheik/create/pkg/shell"
)

// ImageExists checks to see if a docker image exists in the local registry
// takes one argument, the image name: docker.ImageExists("NameOfImage")
func ImageExists(image string) bool {
	command := fmt.Sprintf(`bash -c "if [[ \"$(docker images -q %s)\" == \"\" ]]; then exit 1; else exit 0; fi"`, image)
	return shell.Exec(command) == nil
}

func Pull(image string) bool {
	command := "docker pull " + image
	return shell.Exec(command) == nil
}
