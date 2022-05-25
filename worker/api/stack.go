package api

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/agilestacks/stack-disposer/disposer/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gorilla/mux"
)

func UndeployStackHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	args := make([]string, 0)
	verbose := false
	verbose, _ = strconv.ParseBool(query.Get("verbose"))
	if verbose {
		args = append(args, "--verbose")
	}
	commit := query.Get("commit")

	vars := mux.Vars(r)
	sandboxId := vars["sandboxId"]
	stackId := vars["stackId"]

	log.Println("Undeploying", stackId)

	dir := filepath.Join(config.GitDir, stackId)
	err := checkout(dir, commit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stackDir := filepath.Join(dir, sandboxId)

	_, err = os.Stat(stackDir)
	if os.IsNotExist(err) {
		log.Printf("Sandbox type '%s' not found", sandboxId)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprint(err)))
		return
	}

	_, err = os.Stat(filepath.Join(stackDir, "hub.yaml"))
	if os.IsNotExist(err) {
		log.Printf("File hub.yaml for sandbox type '%s' not found", sandboxId)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = undeploy(stackDir, stackId, args...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func undeploy(stackDir string, stackId string, args ...string) error {

	subCommands := make([]string, 0)
	subCommands = append(subCommands, "stack", "init", stackId, "--force")
	subCommands = append(subCommands, args...)
	err := hubCommand(stackDir, subCommands...)
	if err != nil {
		log.Println(err)
		log.Println("Failed to initialize stack", stackId)
		return err
	}

	err = hubCommand(stackDir, "stack", "undeploy")
	if err != nil {
		log.Println(err)
		log.Println("Failed to undeploy stack", stackId)
		return err
	}

	log.Printf("Stack %s is undeployed", stackId)

	return nil
}

func hubCommand(dir string, args ...string) error {
	log.Printf("Execute: hub %s at %s", strings.Join(args, " "), dir)
	cmd := exec.Command("hub", args...)

	cmd.Dir = dir

	if config.Verbose {
		cmd.Stdout = log.Writer()
		cmd.Stderr = log.Writer()
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func checkout(dir string, commit string) error {
	log.Println("Prepare sandboxes repository", dir, commit)

	if !plumbing.IsHash(commit) {
		return errors.New(fmt.Sprint("invalid commit hash", commit))
	}

	var progress io.Writer
	if config.Verbose {
		progress = log.Writer()
	}

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      config.GitUrl,
		Progress: progress,
	})

	if err != nil {
		repo, err = git.PlainOpen(dir)
		if err != nil {
			log.Println("Unable to open git repo:", err)
			return err
		}
	}

	refHash := plumbing.NewHash(commit)
	if len(commit) <= 0 {
		ref, err := repo.Head()
		if err != nil {
			log.Println("Unable to retrive ref to HEAD of git repo:", err)
			return err
		}

		refHash = ref.Hash()
	}

	wt, err := repo.Worktree()
	if err != nil {
		log.Println("Unable to read work tree of git repo:", err)
		return err
	}

	log.Println("Checkout", refHash)
	err = wt.Checkout(&git.CheckoutOptions{
		Hash:  refHash,
		Force: true,
		Keep:  false,
	})
	if err != nil {
		log.Println("Unable to checkout git repo:", err)
		return err
	}

	return nil
}
