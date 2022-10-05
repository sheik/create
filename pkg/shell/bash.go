package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Output(cmdline string) string {
	outBytes, _ := exec.Command("/bin/bash", "-c", cmdline).Output()
	return strings.TrimSuffix(string(outBytes), "\n")
}

func Exec(cmdline string) error {
	var cmd *exec.Cmd
	if strings.Contains(cmdline, "\n") {
		file, err := os.CreateTemp("/tmp", "create-script")
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())
		defer file.Close()
		if _, err := file.WriteString(cmdline); err != nil {
			return err
		}
		cmd = exec.Command("/bin/bash", file.Name())
	} else {
		cmd = exec.Command("/bin/bash", "-c", cmdline)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Bash(command string) bool {
	cmd := fmt.Sprintf("bash -c \"%s\"", command)
	return Exec(cmd) == nil
}
