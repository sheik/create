package create

import (
	"flag"
	"fmt"
	"github.com/sheik/create/pkg/color"
	"os"
	"os/exec"
	"path"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "more verbose output")
)

func Command(cmdline string) error {
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

func Output(cmdline string) string {
	flag.Parse()
	if *verbose {
		fmt.Println("evaluating: " + cmdline)
	}
	outBytes, _ := exec.Command("/bin/bash", "-c", cmdline).Output()
	return strings.TrimSuffix(string(outBytes), "\n")
}

func (steps Steps) Execute(name string) (err error) {
	if step, ok := steps[name]; ok {
		if step.Gate != nil && !step.Gate() {
			err = fmt.Errorf("target \"%s\" did not pass gate: %s", name, path.Base(runtime.FuncForPC(reflect.ValueOf(step.Gate).Pointer()).Name()))
			return
		}
		if step.Check && !step.executed {
			fmt.Println(color.Purple("[-] skipping ", name))
			step.executed = true
			steps[name] = step
			return
		}
		for _, stepName := range steps[name].Depends {
			if !steps[name].executed {
				steps.ProcessTarget(stepName)
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
			steps[name] = step
		}

		return
	}
	return fmt.Errorf("build target \"%s\" not found", name)
}

type Step struct {
	Command      string
	Precondition string
	Check        bool
	Gate         func() bool
	Fail         string
	Help         string
	Depends      []string
	Default      bool
	Interactive  bool
	executed     bool
}

type Steps map[string]Step

func Complete(args ...string) []string {
	return args
}

var UpdateStep = Step{
	Command: `
		export GOPRIVATE=github.com
		touch go.sum
		rm -f go.sum
		go clean -modcache
		sed -i "s/^.*github.com\/sheik\/create.*$//g" go.mod
		go mod tidy
		go mod vendor
		go install github.com/sheik/create/cmd/create@latest
		`,
}

func (steps Steps) PrintHelp() {
	var items []string
	for name, _ := range steps {
		items = append(items, name)
	}
	sort.Strings(items)
	for _, item := range items {
		fmt.Println(color.Green(item), ":", steps[item].Help)
	}
}

func Plan(steps Steps) {
	flag.Parse()
	steps["update"] = UpdateStep
	if len(flag.Args()) > 0 {
		target := flag.Arg(0)
		if target == "help" {
			steps.PrintHelp()
			return
		}
		steps.ProcessTarget(target)
		return
	}
	for target, step := range steps {
		if step.Default {
			steps.ProcessTarget(target)
		}
	}
}

func (steps Steps) ProcessTarget(name string) {
	var err error
	preconditionFailed := false
	if steps[name].Precondition != "" {
		err = Command(steps[name].Precondition)
		if err != nil {
			preconditionFailed = true
			fmt.Printf(color.Teal("[X] failed precondition for %s\n"), name)
		}
	}
	if !preconditionFailed {
		err = steps.Execute(name)
	}
	if err != nil || preconditionFailed {
		if steps[name].Fail != "" {
			fmt.Printf(color.Teal("[X] error running target \"%s\": failing over to %s\n"), name, steps[name].Fail)
			err = steps.Execute(steps[name].Fail)
		}
		if err != nil {
			fmt.Printf(color.Red("[!] error running target \"%s\": %s\n"), flag.Arg(0), err)
		}
	}
}
