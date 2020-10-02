package composer

import (
	"io/ioutil"
	"log"
	"path"
)

// Config contains paths to the standalone console utils.
// Rather than creating one big app, for now, iffy is split into a bunch of separate commands.
type Config struct {
	Import   string
	Assemble string
	Check    string
	Play     string
	Root     string
	Port     int
}

func (c *Config) Scratch(parts ...string) string {
	return path.Join(append([]string{c.Root, "scratch"}, parts...)...)
}

func (c *Config) PathTo(parts ...string) string {
	return path.Join(append([]string{c.Root}, parts...)...)
}

// DevConfig creates a reasonable(?) config based on the developer go path.
func DevConfig(base string) *Config {
	bin := "bin"
	var dir string // echo $TMPDIR
	if temp, e := ioutil.TempDir("", "iffy"); e != nil {
		log.Fatal(e)
	} else {
		dir = temp
	}
	i, a, c, p := "import", "asm", "check", "play"
	return &Config{
		Import:   path.Join(base, bin, i),
		Assemble: path.Join(base, bin, a),
		Check:    path.Join(base, bin, c),
		Play:     path.Join(base, bin, p),
		Root:     dir,
		Port:     3000,
	}
}
