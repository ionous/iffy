package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/ionous/iffy/web"
	"github.com/kr/pretty"
)

func main() {
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/compose/index.html", http.StatusMovedPermanently)
	})
	http.Handle("/compose/", http.StripPrefix("/compose/", http.FileServer(http.Dir("./www"))))
	http.HandleFunc("/story/", web.HandleResource(StoryApi()))

	log.Println("Listening on :3000...")
	if e := http.ListenAndServe(":3000", nil); e != nil {
		log.Fatal(e)
	}

	var in string
	flag.StringVar(&in, "in", "", "input file name")
	flag.Parse()
}

func StoryApi() web.Resource {
	return &web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			switch name {
			case "story":
				ret = &web.Wrapper{
					Finds: func(name string) (ret web.Resource) {
						switch name {
						case "test":
							ret = &web.Wrapper{
								Puts: func(in io.Reader, out http.ResponseWriter) (err error) {
									d := make(map[string]interface{})
									dec := json.NewDecoder(in)
									if e := dec.Decode(&d); e != nil && e != io.EOF {
										err = e
									} else {
										log.Println(pretty.Sprint(d))
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
