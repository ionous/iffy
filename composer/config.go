package composer

import (
	"io/ioutil"
	"log"
	"path"
)

// Config contains paths to the standalone console utils.
// Rather than creating one big app, for now, iffy is split into a bunch of separate commands.
type Config struct {
	Import    string
	Assemble  string
	Check     string
	Play      string
	Documents string
	Port      int
}

// DevConfig creates a reasonable(?) config based on the developer go path.
func DevConfig(base string) *Config {
	repo := "src/github.com/ionous/iffy/cmd"
	var dir string
	if temp, e := ioutil.TempDir("", "iffy"); e != nil {
		log.Fatal(e)
	} else {
		dir = temp
	}
	i, a, c, p := "import", "asm", "check", "play"
	return &Config{
		Import:    path.Join(base, repo, i, i),
		Assemble:  path.Join(base, repo, a, a),
		Check:     path.Join(base, repo, c, c),
		Play:      path.Join(base, repo, p, p),
		Documents: dir,
		Port:      3000,
	}
}
