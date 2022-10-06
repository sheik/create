package log

import "fmt"

var (
	Error = Message(Red, '!')
	Warn  = Message(Teal, 'X')
	Info  = Message(Green, '*')
)

var (
	Black   = color("\033[1;30m%s\033[0m")
	Red     = color("\033[1;31m%s\033[0m")
	Green   = color("\033[1;32m%s\033[0m")
	Yellow  = color("\033[1;33m%s\033[0m")
	Purple  = color("\033[1;34m%s\033[0m")
	Magenta = color("\033[1;35m%s\033[0m")
	Teal    = color("\033[1;36m%s\033[0m")
	White   = color("\033[1;37m%s\033[0m")
)

type Func func(args ...interface{}) string

func color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func Message(colorFunc func(...interface{}) string, symbol rune) func(string, ...interface{}) {
	return func(message string, args ...interface{}) {
		fmt.Printf(colorFunc("[%c] %s\n"), symbol, fmt.Sprintf(message, args...))
	}
}
