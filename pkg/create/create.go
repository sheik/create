package create

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func command(cmdline string) error {
	cmd := exec.Command("/bin/bash", "-c", cmdline)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Output(cmdline string) []byte {
	out, _ := exec.Command("/bin/bash", "-c", cmdline).Output()
	return out
}

func (s Steps) Execute(name string) error {
	if step, ok := s[name]; ok {
		if step.Check != "" {
			err := command(step.Check)
			if err == nil {
				fmt.Println("skipping", name)
				step.executed = true
				s[name] = step
				return nil
			}
		}
		for _, stepName := range s[name].Depends {
			if !s[name].executed {
				s.Execute(stepName)
			}
		}
		if !step.executed {
			fmt.Println("executing", name)
			err := command(step.Command)
			if err != nil {
				return err
			}
			step.executed = true
			s[name] = step
		}

		return nil
	}
	return fmt.Errorf("build target \"%s\" not found", name)
}

type Step struct {
	Command  string
	Check    string
	Depends  []string
	Default  bool
	executed bool
}

type Steps map[string]Step

func Complete(args ...string) []string {
	return args
}

func Plan(steps Steps) {
	flag.Parse()
	if len(flag.Args()) > 0 {
		steps.Execute(flag.Arg(0))
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
