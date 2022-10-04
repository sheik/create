package create

import "fmt"

func Bash(command string) bool {
	cmd := fmt.Sprintf("bash -c \"%s\"", command)
	return Command(cmd) == nil
}
