package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	listenPort := os.Getenv("NOMAD_PORT_http")
	if len(listenPort) == 0 {
		listenPort = "3000"
	}

	listenAddr := flag.String("listen-addr", fmt.Sprintf(":%s", listenPort),
		"Address on which to listen for HTTP requests")
	saveDir := flag.String("save-dir", "./requests/", "Root directory to save request bodies to")
	flag.Parse()

	os.MkdirAll(*saveDir, os.ModePerm)

	log.Printf("Saving to %s, listening on %s", *saveDir, *listenAddr)
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		f, err := os.Create(path.Join(*saveDir, fmt.Sprintf("%s.txt", r.URL.Path[1:])))
		if err == nil {
			r.Write(f)
		}
		defer f.Close()
		rw.WriteHeader(204)
	})

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
