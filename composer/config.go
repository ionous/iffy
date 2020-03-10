package composer

import (
	"go/build"
	"path"
)

// Config contains paths to the standalone console utils.
// Rather than creating one big app, for now, iffy is split into a bunch of separate commands.
type Config struct {
	Composer  string
	Importer  string
	Assembler string
	Player    string
}

func DevConfig() Config {
	base := build.Default.GOPATH
	repo := "github.com/ionous/iffy/cmd"
	return Config{
		Composer:  path.Join(base, repo, "compose"),
		Importer:  path.Join(base, repo, "import"),
		Assembler: path.Join(base, repo, "asm"),
		Player:    path.Join(base, repo, "play"),
	}
}

//

// to create a temp dir.
// then i can call apps to create

// 	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
// 	defer cancel()

// 	if err := exec.CommandContext(ctx, "sleep", "5").Run(); err != nil {
// 		// This will fail after 100 milliseconds. The 5 second sleep
// 		// will be interrupted.
// 	}

// you'll have to have the commands listed;
// default paths based maybe on gopath or current exe path by default

// fmt.Println(os.Getenv("GOPATH"))
// gopath = build.Default.GOPATH
