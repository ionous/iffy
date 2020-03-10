package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/", fs)

	log.Println("Listening on :3000...")
	if e := http.ListenAndServe(":3000", nil); e != nil {
		log.Fatal(e)
	}

	// handle POST "play" button
	// save that to a file -- tmp
	//dir, err := ioutil.TempFile("", "iffy")
	// "example.*.txt"
	//
	// invoke import -- import needs to take a command line input, output
	// same for asm and play..

	// maybe? it would be possible to open/close pipes for play to get input

	// flag

	// func init() {
	var in string
	flag.StringVar(&in, "in", "", "input file name")
	flag.Parse()

}
