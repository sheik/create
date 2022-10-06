package docker

import (
	"fmt"
	"github.com/sheik/create/pkg/shell"
)

type ContainerObj struct {
	Name  string
	Flags string
}

func Container(name string) *ContainerObj {
	return &ContainerObj{
		Name: name,
	}
}

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

func (image *ContainerObj) Interactive() *ContainerObj {
	image.Flags += " -it "
	return image
}

func (image *ContainerObj) Mount(source, dest string) *ContainerObj {
	image.Flags += fmt.Sprintf("-v %s:%s", source, dest)
	return image
}

func (image *ContainerObj) Run(formatString string, args ...interface{}) string {
	dockerCommand := fmt.Sprintf(formatString, args...)
	return fmt.Sprintf("docker run %s --rm -v $PWD:/code %s %s", image.Flags, image.Name, dockerCommand)
}
