package main

import (
	"flag"
	"fmt"

	"github.com/reiver/go-porterstemmer"
)

// command-line util for trying out stemmed input
// ex. "contain/s/ed/ing" to "contain"
func main() {
	flag.Parse()
	try := flag.Arg(0)
	fmt.Println("stemming:", try)
	fmt.Println(porterstemmer.StemString(try))
}
