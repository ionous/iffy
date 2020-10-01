package composer

import (
	"path/filepath"

	"github.com/ionous/iffy/web"
)

func FilesApi(cfg *Config) web.Resource {
	return &web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			switch name {
			case "stories":
				// by adding a trailing slash, walk will follow a symlink.
				ret = storyFolder(filepath.Join(cfg.Root, "stories") + "/")
			}
			return
		},
	}
}
