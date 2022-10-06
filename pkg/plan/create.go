package plan

import (
	"flag"
	"fmt"
	"github.com/sheik/create/pkg/color"
	"github.com/sheik/create/pkg/shell"
	"os"
	"sort"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "more verbose output")
)

func (steps Steps) Execute(name string) (err error) {
	if step, ok := steps[name]; ok {
		if step.Gate != nil {
			err = fmt.Errorf("target \"%s\" did not pass gate: %s",
				name,
				step.Gate.Error(),
			)
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
			if step.Function != nil {

			}
			if step.Function != nil {
				step.Function.(func(...interface{}) error)()
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
	Function     interface{}
	Precondition string
	Check        bool
	Gate         error
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
		go install github.com/sheik/create/cmd/create@latest
		touch go.sum
		rm -f go.sum
		go clean -modcache
		go mod tidy
		go get -u github.com/sheik/create/pkg/build
		go get -u github.com/sheik/create/pkg/docker
		go get -u github.com/sheik/create/pkg/git
		go get -u github.com/sheik/create/pkg/plan
		go get -u github.com/sheik/create/pkg/shell
		go get -u github.com/sheik/create/pkg/util
		go mod vendor
		`,
	Help: "update create",
}

var HelpStep = Step{
	Help: "print help message for createfile",
}

func (steps Steps) PrintHelp(args ...interface{}) error {
	var items []string
	for name, _ := range steps {
		items = append(items, name)
	}
	sort.Strings(items)
	for _, item := range items {
		if steps[item].Help != "" {
			fmt.Printf("%30s : %s\n", color.Green(item), steps[item].Help)
		}
	}
	return nil
}

func Run(steps Steps) {
	var err error
	flag.Parse()

	// populate Steps map with auto targets
	steps["update"] = UpdateStep
	HelpStep.Function = steps.PrintHelp
	steps["help"] = HelpStep

	target := flag.Arg(0)
	if target == "" {
		if target, err = steps.DefaultTarget(); err != nil {
			color.Error(err.Error())
			os.Exit(3)
		}
	}
	steps.ProcessTarget(target)
	os.Exit(0)
}

func (steps Steps) DefaultTarget() (string, error) {
	for target, step := range steps {
		if step.Default {
			return target, nil
		}
	}
	return "", fmt.Errorf("no default target found in createfile")
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
			fmt.Printf(color.Red("[!] failed precondition for %s\n"), name)
			os.Exit(1)
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
			os.Exit(2)
		}
	}
}
