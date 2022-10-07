package main

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/sheik/create/pkg/parser"
	"os"
)

type CreatefileListener struct {
	*parser.BasecreatefileListener
	Steps map[string]parser.IBlockbodyContext
}

func (l *CreatefileListener) EnterFile_(ctx *parser.File_Context) {
	fmt.Println("Entered file!")
	l.Steps = make(map[string]parser.IBlockbodyContext)
}

func (l *CreatefileListener) EnterCommand(ctx *parser.CommandContext) {
	fmt.Println("Executing command:", ctx.STRING())
}

func (l *CreatefileListener) EnterStep(ctx *parser.StepContext) {
	step := ctx.Stepname().GetText()
	body := ctx.Blockbody()
	l.Steps[step] = body
	fmt.Printf("Defined \"%s\"\n", step)
}

func (l *CreatefileListener) EnterDependencies(ctx *parser.DependenciesContext) {
	fmt.Println("Execute dependency first:")
	for _, dep := range ctx.AllIDENTIFIER() {
		fmt.Println("execute", dep.GetText())
	}
}

func main() {
	contents, err := os.ReadFile("createfile")
	if err != nil {
		panic(err)
	}
	is := antlr.NewInputStream(string(contents))

	// Create lexer
	lexer := parser.NewcreatefileLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create parser
	p := parser.NewcreatefileParser(stream)
	var listener = &CreatefileListener{}
	antlr.ParseTreeWalkerDefault.Walk(listener, p.File_())
}
