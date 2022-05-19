package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/gorilla/mux"

	"github.com/agilestacks/stack-disposer/disposer/api"
	"github.com/agilestacks/stack-disposer/disposer/config"
)

func init() {
	defaultPort, exists := os.LookupEnv("PORT")
	if !exists || len(defaultPort) <= 0 {
		defaultPort = "8080"
	}

	flag.StringVar(&config.Port, "port", defaultPort, "port where to listen")
	flag.StringVar(&config.GitUrl, "gitUrl", "https://github.com/agilestacks/google-stacks.git", "Git URL with stacks")
	flag.StringVar(&config.GitDir, "gitDir", "/tmp/stacks", "directory where clone stacks to")
	flag.BoolVar(&config.Verbose, "verbose", false, "verbose logging")
	flag.DurationVar(&config.Timeout, "timeout", 1*time.Hour, "request timeout")

	flag.Parse()

	if config.Verbose {
		log.Print("VERBOSE logging enabled")
	}
}

func init() {
	log.Println("Prepare sandboxes repository", config.GitUrl)

	var progress io.Writer
	if config.Verbose {
		progress = log.Writer()
	}

	_, err := git.PlainClone(config.GitDir, false, &git.CloneOptions{
		URL:      config.GitUrl,
		Progress: progress,
	})
	// If exists - checkout HEAD with overwrite of local changes
	if err != nil {
		repo, err := git.PlainOpen(config.GitDir)
		if err != nil {
			log.Println("Unable to open git repo:", err)
			return
		}

		ref, err := repo.Head()
		if err != nil {
			log.Println("Unable to retrive ref to HEAD of git repo:", err)
			return
		}

		wt, err := repo.Worktree()
		if err != nil {
			log.Println("Unable to read work tree of git repo:", err)
			return
		}

		err = wt.Checkout(&git.CheckoutOptions{
			Hash:  ref.Hash(),
			Force: true,
			Keep:  false,
		})
		if err != nil {
			log.Println("Unable to checkout git repo:", err)
			return
		}
	}
}

func handleRequests() {
	router := mux.NewRouter()
	router.HandleFunc("/{sandboxId}/{stackId}", api.UndeployStackHandler).Methods(http.MethodDelete)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", config.Port),
		WriteTimeout: config.Timeout,
		ReadTimeout:  config.Timeout,
	}

	log.Println("Server listens on port", config.Port)
	log.Fatal(srv.ListenAndServe())
}

func main() {
	handleRequests()
}
