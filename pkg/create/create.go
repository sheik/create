package create

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func command(cmdline string) error {
	cmd := exec.Command("/bin/bash", "-c", cmdline)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Output(cmdline string) string {
	outBytes, _ := exec.Command("/bin/bash", "-c", cmdline).Output()
	return strings.TrimSuffix(string(outBytes), "\n")
}

func (s Steps) Execute(name string) (err error) {
	if step, ok := s[name]; ok {
		if step.Check != "" && !step.executed {
			err = command(step.Check)
			if err == nil {
				fmt.Println("skipping", name)
				step.executed = true
				s[name] = step
				return
			}
		}
		for _, stepName := range s[name].Depends {
			if !s[name].executed {
				s.Execute(stepName)
			}
		}
		if !step.executed {
			fmt.Println("executing", name)
			fmt.Println(step.Command)
			if step.Interactive {
				err = InteractiveCommand(step.Command)
			} else {
				err = command(step.Command)
			}
			if err != nil {
				return
			}
			step.executed = true
			s[name] = step
		}

		return
	}
	return fmt.Errorf("build target \"%s\" not found", name)
}

type Step struct {
	Command     string
	Check       string
	Depends     []string
	Default     bool
	Interactive bool
	executed    bool
}

type Steps map[string]Step

func Complete(args ...string) []string {
	return args
}

func Plan(steps Steps) {
	flag.Parse()
	if len(flag.Args()) > 0 {
		err := steps.Execute(flag.Arg(0))
		if err != nil {
			fmt.Printf("error running target \"%s\": %s\n", flag.Arg(0), err)
		}
		return
	}
	for name, step := range steps {
		if step.Default {
			err := steps.Execute(name)
			if err != nil {
				fmt.Printf("error running target \"%s\": %s\n", name, err)
			}
		}
	}
}
