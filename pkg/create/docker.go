package create

import "fmt"

func DockerImageExists(image string) bool {
	command := fmt.Sprintf(`bash -c "if [[ \"$(docker images -q %s)\" == \"\" ]]; then exit 1; else exit 0; fi"`, image)
	return Command(command) == nil
}
