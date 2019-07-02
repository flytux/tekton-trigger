package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.com/pongsatt/githook/pkg/githook"
	"gitlab.com/pongsatt/githook/pkg/server"
	"gitlab.com/pongsatt/githook/pkg/tekton"
)

const (
	// Environment variable containing the HTTP port
	envPort = "PORT"

	// EnvSecret environment variable containing gogs secret token
	envSecret = "SECRET_TOKEN"
)

func main() {
	gitprovider := flag.String("gitprovider", "", "git provider ex. gitlab github")
	namespace := flag.String("namespace", "default", "namespace to create pipelinerun")
	name := flag.String("name", "", "name of the pipelinerun")
	runSpecJSON := flag.String("runSpecJSON", "", "pipelinerun spec in json format")

	flag.Parse()

	if gitprovider == nil || *gitprovider == "" {
		log.Fatalf("No gitprovider given")
	}

	if name == nil || *name == "" {
		log.Fatalf("No name given")
	}

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

	hook, err := buildProvider(*gitprovider, secretToken)

	if err != nil {
		log.Fatal(err)
	}

	ra := &githook.ReceiveAdapter{
		TektonClient: tektonClient,
		HookServer:   hook,
		Namespace:    *namespace,
		Name:         *name,
		RunSpecJSON:  *runSpecJSON,
	}

	addr := fmt.Sprintf(":%s", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ra.HandleRequest(w, r)
	})
	http.ListenAndServe(addr, nil)
}

func buildProvider(gitprovider, secretToken string) (githook.HookServer, error) {
	switch gitprovider {
	case "gogs":
		return server.NewGogsServer(secretToken)
	}

	return nil, fmt.Errorf("provider %s not supported", gitprovider)
}
