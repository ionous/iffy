package composer

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/ionous/errutil"
	"golang.org/x/net/context"
)

// write the contents to the passed filename but only if the file doesnt already exist
func ensureFile(fullPath string, contents map[string]interface{}) (err error) {
	if n, e := os.Stat(fullPath); os.IsNotExist(e) {
		dir := path.Dir(fullPath)
		if e := os.MkdirAll(dir, 0755); e != nil {
			err = e
		} else if b, e := json.MarshalIndent(contents, "", "  "); e != nil {
			err = e
		} else {
			err = ioutil.WriteFile(fullPath, b, 0644)
		}
	} else if e != nil {
		err = e
	} else if n.IsDir() {
		err = errutil.New("unexpected conflicting directory")
	}
	return
}

// uses the command line tool versions for now....
func tempTest(ctx context.Context, file string, in io.Reader) (err error) {
	cfg := ctx.Value(configKey).(*Config)
	base := cfg.PathTo("stories")
	if !strings.HasPrefix(file, base) {
		err = errutil.New("unexpected path", file, "from", base)
	} else {
		// note: .Split keeps a trailing slash, .Dir does not.
		dir, _ := path.Split(file[len(base)+1:])
		const shared = "shared/"
		const stories = "stories/"
		// we'll always include the shared files in our build
		src := cfg.PathTo(stories, shared)
		if strings.HasPrefix(dir, shared) {
			dir = shared
		} else {
			// get the first part of the name -- that's the project name
			i := strings.Index(dir, "/")
			dir = dir[0:i] // the project relative dir
			src += "," + cfg.PathTo(stories, dir)
		}
		// src is now one or two absolute paths to project directories
		// dir is a relative dir
		if ephFile, e := runImport(ctx, cfg, src, dir); e != nil {
			log.Println(tab, "Import error", cfg.Import, exitError(e))
			err = e
		} else if playFile, e := runAsm(ctx, cfg, ephFile, dir); e != nil {
			log.Println(tab, "Assembly error", cfg.Assemble, exitError(e))
			err = e
		} else if e := runCheck(ctx, cfg, playFile); e != nil {
			log.Println(tab, "Check error", cfg.Check, exitError(e))
			err = e
		}
	}
	return
}

const tab = '\t'

// note: for now, these read from CombinedOutput to grab any log/println traces...
func runImport(ctx context.Context, cfg *Config, inFile, path string) (ret string, err error) {
	log.Println("Importing", inFile+"...")
	ephFile := cfg.Scratch(path, "ephemera.db")
	log.Println(">", cfg.Import, "-in", inFile, "-out", ephFile)
	imported, e := exec.CommandContext(ctx, cfg.Import, "-in", inFile, "-out", ephFile).CombinedOutput()
	if e != nil {
		err = e
	} else {
		ret = ephFile
	}
	logBytes(imported)
	return
}

func runAsm(ctx context.Context, cfg *Config, ephFile, path string) (ret string, err error) {
	log.Println("Assembling", ephFile+"...")
	inFile, playFile := ephFile, cfg.Scratch(path, "play.db")
	log.Println(">", cfg.Assemble, "-in", inFile, "-out", playFile)
	assembled, e := exec.CommandContext(ctx, cfg.Assemble, "-in", inFile, "-out", playFile).CombinedOutput()
	if e != nil {
		err = e
	} else {
		ret = playFile
	}
	logBytes(assembled)
	return

}
func runCheck(ctx context.Context, cfg *Config, playFile string) (err error) {
	log.Println("Checking", playFile+"...")
	log.Println(">", cfg.Check, "-in", playFile)
	checked, e := exec.CommandContext(ctx, cfg.Check, "-in", playFile).CombinedOutput()
	if e != nil {
		err = e
	}
	logBytes(checked)
	return
}

func logBytes(b []byte) {
	if s := strings.Trim(string(b), "\n"); len(s) > 0 {
		log.Println(s)
	}
}

func exitError(e error) (ret string) {
	if x, ok := e.(*exec.ExitError); ok {
		ret = string(x.Stderr)
	} else {
		ret = "unknown."
	}
	return
}
