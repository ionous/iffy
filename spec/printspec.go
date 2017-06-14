package spec

import (
	"fmt"
)

func PrintArray(cmds []*Command) {
	printArg(&(cmds), "", "")
}

func printCmd(cmd *Command, space, header string) {
	if len(cmd.name) > 0 {
		fmt.Println(space, header, cmd.name)
	}
	space = space + " "
	for _, arg := range cmd.args {
		dash := fmt.Sprint(space, "-")
		printArg(arg, space, dash)
	}
	for k, arg := range cmd.keys {
		dash := fmt.Sprint(space, "- ", k)
		printArg(arg, space, dash)
	}
}

func printArg(arg interface{}, space, header string) {
	switch arg := arg.(type) {
	case *Command:
		printCmd(arg, space, header)
	case *[]*Command:
		fmt.Println(space, header, "{")
		indent := space + " "
		for i, cmd := range *arg {
			printCmd(cmd, indent, fmt.Sprint(" ", i, ":"))
		}
		fmt.Println(space, "}")
	default:
		fmt.Println(space, header, arg)
	}
}
