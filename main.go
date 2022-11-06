package main

import (
	"flag"
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type WebhookPushData struct {
	Pusher    string   `json:"pusher"`
	PushedAt  int64    `json:"pushed_at"`
	Tag       string   `json:"tag"`
	Images    []string `json:"images"`
	MediaType string   `json:"media_type"`
}

type WebhookRepository struct {
	Status          string `json:"status"`
	Namespace       string `json:"namespace"`
	Name            string `json:"name"`
	RepoName        string `json:"repo_name"`
	RepoUrl         string `json:"repo_url"`
	Description     string `json:"description"`
	FullDescription string `json:"full_description"`
	StarCount       int64  `json:"star_count"`
	Dockerfile      string `json:"dockerfile"`
	IsPrivate       bool   `json:"is_private"`
	IsTrusted       bool   `json:"is_trusted"`
	IsOfficial      bool   `json:"is_official"`
	Owner           string `json:"owner"`
	DateCreated     int64  `json:"date_created"`
}

type WebhookBody struct {
	CallbackUrl string            `json:"callback_url"`
	PushData    WebhookPushData   `json:"push_data"`
	Repository  WebhookRepository `json:"repository"`
}

func main() {
	listenPort := os.Getenv("NOMAD_PORT_http")
	if len(listenPort) == 0 {
		listenPort = "3000"
	}

	listenAddr := flag.String("listen-addr", fmt.Sprintf(":%s", listenPort),
		"Address on which to listen for HTTP requests")
	saveDir := flag.String("save-dir", "", "Root directory to save request bodies to")
	flag.Parse()

	if *saveDir == "" {
		*saveDir = os.Getenv("NOMAD_TASK_DIR")
		if len(*saveDir) == 0 {
			wd, _ := os.Getwd()
			*saveDir = fmt.Sprintf("%s/requests/", wd)
		}
	}

	os.MkdirAll(*saveDir, os.ModePerm)

	log.Printf("Saving to %s, listening on %s", *saveDir, *listenAddr)
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Printf("Received %s %s", r.Method, r.URL)
		if r.Method == "POST" {
			p := path.Clean(path.Join(*saveDir, fmt.Sprintf("%s.txt", r.URL.Path[1:])))
			log.Printf("Writing to %s", p)
			if strings.HasPrefix(p, *saveDir) {
				f, err := os.Create(p)
				if err == nil {
					defer f.Close()

					var body WebhookBody
					err = json.NewDecoder(r.Body).Decode(&body)
					if err != nil {
						f.Write([]byte(fmt.Sprintf("Error: %s", err.Error())))
						log.Printf("Error: %s", err.Error())
					} else {
						f.Write([]byte(fmt.Sprintf("%s:%s", body.Repository.RepoName, body.PushData.Tag)))
						log.Printf("Written: %s:%s", body.Repository.RepoName, body.PushData.Tag)
					}
				}
				rw.WriteHeader(204)
				return
			}
		}
		rw.WriteHeader(403)
	})

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
