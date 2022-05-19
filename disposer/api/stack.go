package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/agilestacks/stack-disposer/disposer/config"
	"github.com/gorilla/mux"
)

func UndeployStackHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	args := make([]string, 0)
	verbose := len(query.Get("verbose")) > 0
	if verbose {
		args = append(args, "--verbose")
	}

	vars := mux.Vars(r)
	sandboxId := vars["sandboxId"]
	stackId := vars["stackId"]

	stackDir := filepath.Join(config.GitDir, sandboxId)

	_, err := os.Stat(stackDir)
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
	log.Println("Undeploying", stackId)

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
