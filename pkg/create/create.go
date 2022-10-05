package create

import (
	"flag"
	"fmt"
	"github.com/sheik/create/pkg/color"
	"github.com/sheik/create/pkg/shell"
	"os"
	"path"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "more verbose output")
)

func (steps Steps) Execute(name string) (err error) {
	if step, ok := steps[name]; ok {
		if step.Gate != nil && !step.Gate() {
			err = fmt.Errorf("target \"%s\" did not pass gate: %s", name, path.Base(runtime.FuncForPC(reflect.ValueOf(step.Gate).Pointer()).Name()))
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
				err = shell.InteractiveCommand(step.Command)
			} else {
				err = shell.Exec(step.Command)
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
	Help: "update create",
}

func (steps Steps) PrintHelp() {
	var items []string
	for name, _ := range steps {
		items = append(items, name)
	}
	sort.Strings(items)
	for _, item := range items {
		fmt.Printf("%30s : %s\n", color.Green(item), steps[item].Help)
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
	step := steps[name]
	if strings.Contains(step.Command, ":INPUT:") {
		step.Command = strings.ReplaceAll(step.Command, ":INPUT:", strings.Join(os.Args[2:], " "))
		steps[name] = step
	}
	if step.Check && !step.executed {
		fmt.Println(color.Purple("[-] skipping ", name))
		step.executed = true
		steps[name] = step
		return
	}
	if step.Precondition != "" {
		err = shell.Exec(step.Precondition)
		if err != nil {
			preconditionFailed = true
			fmt.Printf(color.Teal("[X] failed precondition for %s\n"), name)
		}
	}
	if !preconditionFailed {
		err = steps.Execute(name)
	}
	if err != nil || preconditionFailed {
		if step.Fail != "" {
			fmt.Printf(color.Teal("[X] error running target \"%s\": failing over to %s\n"), name, step.Fail)
			err = steps.Execute(step.Fail)
		}
		if err != nil {
			fmt.Printf(color.Red("[!] error running target \"%s\": %s\n"), name, err)
		}
	}
}
