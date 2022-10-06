package docker

import (
	"fmt"
	"github.com/sheik/create/pkg/shell"
)

type ImageObj struct {
	Name  string
	Flags string
}

func Image(name string) *ImageObj {
	return &ImageObj{
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

func (image *ImageObj) Interactive() *ImageObj {
	image.Flags += " -it "
	return image
}

func (image *ImageObj) Mount(source, dest string) *ImageObj {
	image.Flags += fmt.Sprintf("-v %s:%s", source, dest)
	return image
}

func (image *ImageObj) Run(formatString string, args ...interface{}) string {
	dockerCommand := fmt.Sprintf(formatString, args...)
	return fmt.Sprintf("docker run %s --rm -v $PWD:/code %s %s", image.Flags, image.Name, dockerCommand)
}
