package main

// import (
// 	"flag"
// 	"log"
// 	// "regexp"
// )

// const usage = `spectator func_name "struct_name param_name param_name"

// spectator generates implementations of iffy commands based on the existence of the named func.

// Example:

// spectator NumWorder "NumWord Num"
// // or, to avoid ast parsing: spectator NumWorder "NumWord:TextEval Num:NumEval"

// When used with the function func NumWorder(x float64) string {...}
// spectator generates a struct "NumWord" with the member "Num rt.NumberEval".
// Because the function returns a string, the struct implements rt.TextEval.
// The implementation converts the Num to a float, checks for errors, and finally, calls NumWorder.

// Spectator can handle functions that return just one value, or a value and error pair.
// `

// //
// func Split(s string) (names []string, err error) {
// 	// im not smart enough for regexp, so we'll do it more manually.
// 	parts := strings.Split(s, `"`)
// 	if len(parts) != 3 {
// 		err = errutil.New("couldnt parse", s)
// 	} else {
// 		names = append(names, strings.TrimSpace(parts[0]))
// 		names = append(names, strings.Fields(parts[1])...)
// 	}
// 	return
// }

func main() {
}

// 	flag.Parse()

// 	// https://github.com/vektra/mockery
// 	// https://github.com/josharian/impl/blob/master/impl.go
// 	// https://github.com/ernesto-jimenez/gogen/blob/master/cmd/goautomock/main.go

// 	spec := flag.Arg(0)
// 	if parts, e := Split(spec); e != nil {
// 		log.Fatalln(e)
// 	}
// 	// validLine := regexp.MustCompile(ValidLine)

// 	// take a function, generate a spec + implementation around unpacking evals to outputs

// }

// // desired input:
// // go:generate go run spectator "Add A B"
// // func Add(a, b float64) (ret float64) {
// //     return a + b
// // }

// // desired output:
// // type Add struct {
// // 	A, B rt.NumberEval
// // }

// // func (op *Add) GetNumber(run rt.Runtime) (ret float64, err error) {
// // 	if vA, e := op.A.GetNumber(run); e != nil {
// // 		err = errutil.New("Add.A", e)
// // 	} else if vB, e := op.B.GetNumber(run); e != nil {
// // 		err = errutil.New("Add.B", e)
// // 	} else {
// // 		ret = Add(vA,Vb) // generates ret,err if spec demands it
// // 	}
// // 	return
// // }

// // 	if a, b, e := Pair(*cmd).Get(run); e != nil {
// // 		err = errutil.New("Add", e)
// // 	} else {
// // 		ret = a + b
// // 	}
// // 	return
// // }
