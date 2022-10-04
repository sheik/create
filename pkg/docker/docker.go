package docker

import (
	"fmt"
	"github.com/sheik/create/pkg/create"
)

func ImageExists(image string) bool {
	command := fmt.Sprintf(`bash -c "if [[ \"$(docker images -q %s)\" == \"\" ]]; then exit 1; else exit 0; fi"`, image)
	return create.Command(command) == nil
}
