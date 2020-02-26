// Package main for 'stem'.
// A command-line util for exploring stemmed input.
// ex. "contain/s/ed/ing" to "contain"
package main

import (
	"flag"
	"fmt"

	"github.com/reiver/go-porterstemmer"
)

func main() {
	flag.Parse()
	try := flag.Arg(0)
	fmt.Println("stemming:", try)
	fmt.Println(porterstemmer.StemString(try))
}
