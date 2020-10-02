package composer

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ionous/iffy/web"
	"golang.org/x/net/context"
)

// Compose starts the composer server, this function doesnt return.
func Compose(cfg *Config) {
	// configure server
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/compose/index.html", http.StatusMovedPermanently)
	})
	http.Handle("/compose/", http.StripPrefix("/compose/", http.FileServer(http.Dir("./www"))))
	http.HandleFunc("/stories/", web.HandleResourceWithContext(FilesApi(cfg), func(ctx context.Context) context.Context {
		return context.WithValue(ctx, configKey, cfg)
	}))

	log.Println("Composer using", cfg.Root)
	log.Println("Listening on port", strconv.Itoa(cfg.Port)+"...")
	if e := http.ListenAndServe(":3000", nil); e != nil {
		log.Fatal(e)
	}
}

type key int

var configKey key
