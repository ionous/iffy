package composer

import (
	"encoding/json"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/web"
)

// Compose starts the composer server, this function doesnt return.
func Compose(cfg *Config) {
	// configure server
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/compose/index.html", http.StatusMovedPermanently)
	})
	http.Handle("/compose/", http.StripPrefix("/compose/", http.FileServer(http.Dir("./www"))))
	http.HandleFunc("/story/", web.HandleResource(StoryApi(cfg)))

	log.Println("Composer using", cfg.Documents)
	log.Println("Listening on port", strconv.Itoa(cfg.Port)+"...")
	if e := http.ListenAndServe(":3000", nil); e != nil {
		log.Fatal(e)
	}

}

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

func StoryApi(cfg *Config) web.Resource {
	return &web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			switch name {
			case "story":
				ret = &web.Wrapper{
					Finds: func(name string) (ret web.Resource) {
						switch name {
						case "check":
							ret = &web.Wrapper{
								Puts: func(ctx context.Context, in io.Reader, out http.ResponseWriter) (err error) {
									if e := tempTest(ctx, cfg, in); e != nil {
										err = e
									}
									return
								},
							}
						}
						return
					},
				}
			}
			return
		},
	}
}

// uses the command line tool versions
func tempTest(ctx context.Context, cfg *Config, in io.Reader) (err error) {
	hash := fnv.New64a()
	d := make(map[string]interface{})
	// all reads from "in" (via tee) are written to hash.
	tee := io.TeeReader(in, hash)
	dec := json.NewDecoder(tee)
	if e := dec.Decode(&d); e != nil && e != io.EOF {
		err = e
	} else {
		const tab = "\t"
		hashed := strconv.FormatUint(hash.Sum64(), 36)
		inFile := path.Join(cfg.Documents, hashed, "story.js")
		log.Println("Saving", inFile+"...")
		if e := ensureFile(inFile, d); e != nil {
			log.Println(tab, "Save error", e)
			err = e
		} else if ephFile, e := runImport(ctx, cfg, inFile, hashed); e != nil {
			log.Println(tab, "Import error", cfg.Import, exitError(e))
			err = e
		} else if playFile, e := runAsm(ctx, cfg, ephFile, hashed); e != nil {
			log.Println(tab, "Assembly error", cfg.Assemble, exitError(e))
			err = e
		} else if e := runCheck(ctx, cfg, playFile); e != nil {
			log.Println(tab, "Check error", cfg.Check, exitError(e))
			err = e
		}
	}
	return
}

// note: for now, these read from CombinedOutput to grab any log/println traces...
func runImport(ctx context.Context, cfg *Config, inFile, hashed string) (ret string, err error) {
	log.Println("Importing", inFile+"...")
	ephFile := path.Join(cfg.Documents, hashed, "ephemera.db")
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

func runAsm(ctx context.Context, cfg *Config, ephFile, hashed string) (ret string, err error) {
	log.Println("Assembling", ephFile+"...")
	inFile, playFile := ephFile, path.Join(cfg.Documents, hashed, "play.db")
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
