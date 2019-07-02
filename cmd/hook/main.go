package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.com/pongsatt/githook/pkg/gogshook"
	"gitlab.com/pongsatt/githook/pkg/tekton"
	"gopkg.in/go-playground/webhooks.v5/gogs"
)

const (
	// Environment variable containing the HTTP port
	envPort = "PORT"

	// EnvSecret environment variable containing gogs secret token
	envSecret = "GOGS_SECRET_TOKEN"
)

func main() {
	namespace := flag.String("namespace", "default", "namespace to create pipelinerun")
	name := flag.String("name", "", "name of the pipelinerun")
	runSpecJSON := flag.String("runSpecJSON", "", "pipelinerun spec in json format")

	flag.Parse()

	if runSpecJSON == nil || *runSpecJSON == "" {
		log.Fatalf("No runSpecJSON given")
	}

	port := os.Getenv(envPort)
	if port == "" {
		port = "8080"
	}

	secretToken := os.Getenv(envSecret)
	if secretToken == "" {
		log.Fatalf("No secret token given")
	}

	log.Printf("runSpecJSON is: %q", *runSpecJSON)

	tektonClient, err := tekton.New()

	if err != nil {
		log.Fatalf("cannot create tekton client: %s", err)
	}

	ra := &gogshook.GogsReceiveAdapter{
		TektonClient: tektonClient,
		Namespace:    *namespace,
		Name:         *name,
		RunSpecJSON:  *runSpecJSON,
	}

	hook, err := gogs.New(gogs.Options.Secret(secretToken))

	if err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf(":%s", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r,
			gogs.CreateEvent,
			gogs.DeleteEvent,
			gogs.ForkEvent,
			gogs.PushEvent,
			gogs.IssuesEvent,
			gogs.IssueCommentEvent,
			gogs.PullRequestEvent,
			gogs.ReleaseEvent)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
		}

		ra.HandleEvent(payload, r.Header)
	})
	http.ListenAndServe(addr, nil)
}
