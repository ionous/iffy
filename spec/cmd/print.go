package cmd

import "fmt"

// PrintSpec writes the passed command to stdout.
func Print(cmd *Command) {
	printCmd(cmd, "", "")
}

func printCmd(cmd *Command, space, header string) {
	if len(cmd.Name) > 0 {
		fmt.Println(space, header, cmd.Name)
	}
	space = space + " "
	for _, arg := range cmd.Args {
		dash := fmt.Sprint(space, "-")
		printArg(arg, space, dash)
	}
	for k, arg := range cmd.Keys {
		dash := fmt.Sprint(space, "- ", k)
		printArg(arg, space, dash)
	}
}

func printArg(arg interface{}, space, header string) {
	switch arg := arg.(type) {
	case *Command:
		printCmd(arg, space, header)
	case *Commands:
		fmt.Println(space, header, "{")
		indent := space + " "
		for i, cmd := range arg.Els {
			printCmd(cmd, indent, fmt.Sprint(" ", i, ":"))
		}
		fmt.Println(space, "}")
	default:
		fmt.Println(space, header, arg)
	}
}
