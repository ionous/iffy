package main

import (
	"flag"
	"go/build"

	"github.com/ionous/iffy/composer"
	"github.com/ionous/iffy/web/support"
)

// ex. go run compose.go -open -dir /Users/ionous/Documents/Iffy
// needs a subdirectory "stories"
func main() {
	var dir string
	var open bool
	flag.StringVar(&dir, "dir", "", "directory for processing iffy files.")
	flag.BoolVar(&open, "open", false, "open a new browser window.")
	flag.Parse()
	//
	cfg := composer.DevConfig(build.Default.GOPATH)
	if len(dir) > 0 {
		cfg.Root = dir
	}
	if open {
		support.OpenBrowser("http://localhost:3000/compose/")
	}
	// by design, this never returns.
	composer.Compose(cfg)
}
