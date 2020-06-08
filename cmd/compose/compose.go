package main

import (
	"flag"
	"go/build"

	"github.com/ionous/iffy/composer"
)

// ex. -dir /Users/ionous/Documents/Iffy
func main() {
	var dir string
	flag.StringVar(&dir, "dir", "", "directory for processing iffy files.")
	flag.Parse()
	//
	cfg := composer.DevConfig(build.Default.GOPATH)
	if len(dir) > 0 {
		cfg.Documents = dir
	}
	composer.Compose(cfg)
}
