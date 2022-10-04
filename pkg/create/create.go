package create

import (
	"flag"
	"fmt"
	"github.com/sheik/create/pkg/color"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "more verbose output")
)

func Command(cmdline string) error {
	cmd := exec.Command("/bin/bash", "-c", cmdline)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Output(cmdline string) string {
	if *verbose {
		fmt.Println("evaluating: " + cmdline)
	}
	outBytes, _ := exec.Command("/bin/bash", "-c", cmdline).Output()
	return strings.TrimSuffix(string(outBytes), "\n")
}

func (s Steps) Execute(name string) (err error) {
	if step, ok := s[name]; ok {
		if step.Check != "" && !step.executed {
			err = Command(step.Check)
			if err == nil {
				fmt.Println(color.Purple("[-] skipping ", name))
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
			fmt.Println(color.Green("[*] executing ", name))
			if *verbose {
				fmt.Println(step.Command)
			}
			if step.Interactive {
				err = InteractiveCommand(step.Command)
			} else {
				err = Command(step.Command)
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
	Help        string
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
		if flag.Arg(0) == "help" {
			var items []string
			for name, _ := range steps {
				items = append(items, name)
			}
			sort.Strings(items)
			for _, item := range items {
				fmt.Println(item)
			}
			return
		}
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
